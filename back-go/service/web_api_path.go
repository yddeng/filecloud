package service

type pathHandler struct {
}

func (*pathHandler) mkdir(wait *WaitConn, req struct {
	Path string `json:"path"`
}) {
	logger.Infof("%s %v", wait.GetRoute(), req)

	defer func() { wait.Done() }()

	if req.Path == "" {
		wait.SetResult("创建路径错误", nil)
		return
	}

	_, err := filePtr.findPath(req.Path, false)
	if err != nil {
		wait.SetResult("文件夹名错误，可能与文件名相同", nil)
		return
	}

}
