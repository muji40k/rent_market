package writer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"rent_service/logger"
	"time"
)

func getStatus(status logger.Status) string {
	switch status {
	case logger.ERROR:
		return "E"
	case logger.WARNING:
		return "W"
	case logger.INFO:
		return "I"
	case logger.DEBUG:
		return "D"
	default:
		return "X"
	}
}

type writerLogger struct {
	host   string
	writer io.Writer
}

func New(writer io.Writer, host *string) (logger.ILogger, error) {
	var hostname string
	var err error

	if nil == writer {
		return nil, errors.New("No writer provided")
	}

	if nil != host {
		hostname = *host
	} else {
		hostname, err = os.Hostname()

		if nil != err {
			err = errors.New("Can't retrieve hostname to setup logger")
		}
	}

	if nil == err {
		write(writer, hostname, logger.INFO, "Logger start")
		return &writerLogger{hostname, writer}, nil
	} else {
		write(writer, "unknown", logger.ERROR, err.Error())
		write(writer, "unknown", logger.INFO, "System will work without log")
		return nil, err
	}
}

func write(writer io.Writer, host string, status logger.Status, msg any) {
	fmt.Fprintf(writer, "[%v][%v][%v]: %v\n", host, time.Now(), getStatus(status), msg)
}

func (self *writerLogger) Log(status logger.Status, msg any) {
	write(self.writer, self.host, status, msg)
}

func (self *writerLogger) Logf(status logger.Status, format string, args ...any) {
	write(self.writer, self.host, status, fmt.Sprintf(format, args...))
}

