package job

import "go.uber.org/zap"

func (s *JobScheduler) hello(a string, b int) {
	logger := zap.L().Sugar()
	defer logger.Sync()
	logger.Infof("hello %v %v", a, b)
}
