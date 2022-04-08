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

	_, err := filePtr.FileInfo.findDir(req.Path, true)
	if err != nil {
		wait.SetResult(err.Error(), nil)
		return
	}

}
