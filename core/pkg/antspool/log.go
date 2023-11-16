package antspool

import "xy3-proto/pkg/log"

type AntsLogger struct {
}

func (l *AntsLogger) Printf(format string, args ...interface{}) {
	log.Error(format, args...)
}
