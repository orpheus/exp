package api

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
}
