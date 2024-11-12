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
}

