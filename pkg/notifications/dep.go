package notifications

type logger interface {
	Warn(args ...interface{})
	Debug(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}
