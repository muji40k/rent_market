package logger

import "rent_service/logger"

type IBuilder interface {
	Build() (logger.ILogger, error)
}

