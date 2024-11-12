package writer

import (
	"io"
	"rent_service/logger"
	"rent_service/logger/implementations/writer"
)

type WriterBuilder struct {
	writer io.Writer
	host   *string
}

func New() *WriterBuilder {
	return &WriterBuilder{}
}

func (self *WriterBuilder) WithWriter(writer io.Writer) *WriterBuilder {
	self.writer = writer
	return self
}

func (self *WriterBuilder) WithHost(host *string) *WriterBuilder {
	if nil == host {
		self.host = nil
	} else {
		self.host = new(string)
		*self.host = *host
	}

	return self
}

func (self *WriterBuilder) Build() (logger.ILogger, error) {
	return writer.New(self.writer, self.host)
}

