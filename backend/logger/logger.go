package logger

type Status int

const (
	ERROR Status = iota
	WARNING
	INFO
	DEBUG
)

type ILogger interface {
	Log(status Status, msg any)
	Logf(status Status, format string, args ...any)
	Close()
}

func Log(log ILogger, status Status, msg any) {
	if nil != log {
		log.Log(status, msg)
	}
}

func Logf(log ILogger, status Status, format string, args ...any) {
	if nil != log {
		log.Logf(status, format, args)
	}
}

func Close(log ILogger) {
	if nil != log {
		log.Close()
	}
}

