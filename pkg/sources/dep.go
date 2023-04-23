package sources

type logger interface {
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}
