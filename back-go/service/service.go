package service

import (
	"fmt"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/yddeng/utils/task"
	"net/http"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	app       *gin.Engine
	taskQueue *task.TaskPool
)

func Launch() {
	taskQueue = task.NewTaskPool(1, 1024)
	saveFileMultiple = config.SaveFileMultiple
	fileDiskTotal = config.FileDiskTotal * MB

	loadFilePath(config.FilePath)

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	// 跨域
	app.Use(func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type")
		ctx.Next()
	})

	// 前端
	if config.WebIndex != "" {
		app.Use(static.Serve("/", static.LocalFile(config.WebIndex, false)))
		app.NoRoute(func(ctx *gin.Context) {
			ctx.File(config.WebIndex + "/index.html")
		})
	}

	initHandler(app)

	port := strings.Split(config.WebAddr, ":")[1]
	webAddr := fmt.Sprintf("0.0.0.0:%s", port)

	logger.Infof("start web service on %s", config.WebAddr)

	if err := app.Run(webAddr); err != nil {
		logger.Error(err)
	}
}

func Stop() {

}

var (
	// 允许无token的路由
	allowTokenRoute = map[string]struct{}{
		"/auth/login":      {},
		"/shared/info":     {},
		"/shared/list":     {},
		"/shared/download": {},
	}
)

func checkToken(ctx *gin.Context, route string) bool {
	if _, ok := allowTokenRoute[route]; ok {
		return true
	}
	tkn := ctx.GetHeader("Access-Token")
	if tkn == "" {
		return false
	}

	if accessToken == "" || time.Now().After(accessTokenExpire) {
		accessToken = ""
		accessTokenExpire = time.Time{}
		return false
	}
	return tkn == accessToken
}

// 应答结构
type Result struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WaitConn struct {
	ctx      *gin.Context
	route    string
	result   Result
	done     chan struct{}
	doneOnce sync.Once
}

func newWaitConn(ctx *gin.Context, route string) *WaitConn {
	return &WaitConn{
		ctx:   ctx,
		route: route,
		done:  make(chan struct{}),
	}
}

func (this *WaitConn) Done() {
	this.doneOnce.Do(func() {
		if this.result.Message == "" {
			this.result.Success = true
		}
		close(this.done)
	})
}

func (this *WaitConn) GetRoute() string {
	return this.route
}

func (this *WaitConn) Context() *gin.Context {
	return this.ctx
}

func (this *WaitConn) SetResult(message string, data interface{}) {
	this.result.Message = message
	this.result.Data = data
}

func (this *WaitConn) Wait() {
	<-this.done
}

type webTask func()

func (t webTask) Do() {
	t()
}

func transBegin(ctx *gin.Context, fn interface{}, args ...reflect.Value) {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	if typ.NumIn() != len(args)+1 {
		panic("func argument error")
	}

	route := getCurrentRoute(ctx)
	wait := newWaitConn(ctx, route)
	if err := taskQueue.SubmitTask(webTask(func() {
		ok := checkToken(ctx, route)
		if !ok {
			wait.SetResult("Token验证失败", nil)
			wait.Done()
			return
		}
		val.Call(append([]reflect.Value{reflect.ValueOf(wait)}, args...))
	})); err != nil {
		wait.SetResult("访问人数过多", nil)
		wait.Done()
	}
	wait.Wait()

	ctx.JSON(http.StatusOK, wait.result)
}

func getCurrentRoute(ctx *gin.Context) string {
	return ctx.FullPath()
}

func getJsonBody(ctx *gin.Context, inType reflect.Type) (inValue reflect.Value, err error) {
	if inType.Kind() == reflect.Ptr {
		inValue = reflect.New(inType.Elem())
	} else {
		inValue = reflect.New(inType)
	}
	if err = ctx.ShouldBindJSON(inValue.Interface()); err != nil {
		return
	}
	if inType.Kind() != reflect.Ptr {
		inValue = inValue.Elem()
	}
	return
}

func WarpHandle(fn interface{}) gin.HandlerFunc {
	val := reflect.ValueOf(fn)
	if val.Kind() != reflect.Func {
		panic("value not func")
	}
	typ := val.Type()
	switch typ.NumIn() {
	case 1: // func(done *WaitConn)
		return func(ctx *gin.Context) {
			transBegin(ctx, fn)
		}
	case 2: // func(done *WaitConn, req struct)
		return func(ctx *gin.Context) {
			inValue, err := getJsonBody(ctx, typ.In(1))
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"message": "Json unmarshal failed!",
					"error":   err.Error(),
				})
				return
			}

			transBegin(ctx, fn, inValue)
		}
	default:
		panic("func symbol error")
	}
}

func initHandler(app *gin.Engine) {
	authHandle := new(authHandler)
	authGroup := app.Group("/auth")
	authGroup.POST("/login", WarpHandle(authHandle.login))

	fileHandle := new(fileHandler)
	fileGroup := app.Group("/file")
	fileGroup.POST("/mkdir", WarpHandle(fileHandle.mkdir))
	fileGroup.POST("/list", WarpHandle(fileHandle.list))
	fileGroup.POST("/remove", WarpHandle(fileHandle.remove))
	fileGroup.POST("/rename", WarpHandle(fileHandle.rename))
	fileGroup.POST("/mvcp", WarpHandle(fileHandle.mvcp))
	fileGroup.POST("/download", func(ctx *gin.Context) {
		var req *downloadArg
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Json unmarshal failed!",
				"error":   err.Error(),
			})
			return
		}
		download(ctx, req)
	})

	uploadHandle := new(uploadHandler)
	uploadGroup := app.Group("/upload")
	uploadGroup.POST("/check", WarpHandle(uploadHandle.check))
	uploadGroup.POST("/upload", WarpHandle(uploadHandle.upload))

	shareHandle := new(shareHandler)
	shareGroup := app.Group("/shared")
	shareGroup.POST("/create", WarpHandle(shareHandle.create))
	//shareGroup.POST("/s/:key", func(ctx *gin.Context) {
	//	// 动态路由
	//	var req *sharedCheckArg
	//	if err := ctx.ShouldBindJSON(&req); err != nil {
	//		ctx.JSON(http.StatusBadRequest, gin.H{
	//			"message": "Json unmarshal failed!",
	//			"error":   err.Error(),
	//		})
	//		return
	//	}
	//
	//	req.Key = ctx.Param("key")
	//
	//	transBegin(ctx, shareHandle.list, reflect.ValueOf(req))
	//})
	shareGroup.POST("/info", WarpHandle(shareHandle.info))
	shareGroup.POST("/list", WarpHandle(shareHandle.list))
	shareGroup.POST("/download", func(ctx *gin.Context) {
		var req *shareDownloadArg
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Json unmarshal failed!",
				"error":   err.Error(),
			})
			return
		}

		shared, err := shareHandle.checkShared(req.Key, req.SharedToken)
		if err != nil {
			ctx.Status(http.StatusBadRequest)
			return
		}

		if req.Path != shared.Path {
			children := false
			for _, name := range shared.Filename {
				if strings.Contains(req.Path, path.Join(shared.Path, name)) {
					children = true
					break
				}
			}

			if !children {
				ctx.Status(http.StatusBadRequest)
				return
			}
		}

		download(ctx, &downloadArg{
			Path:     req.Path,
			Filename: req.Filename,
		})
	})

}
