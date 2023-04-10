package logging

import "go.uber.org/zap"

type loggerContainer struct {
	l *zap.Logger
}

var logger = loggerContainer{l: zap.NewExample()}

// callAdd tests an addition of the logging method with access path containing address taking operation
func callAdd(s2 string) string {
	const s1 = "42"
	(&logger).l.Error("error")
	return s1
}
