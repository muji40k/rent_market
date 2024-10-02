package deffactory

type DefaultFactory struct {
}

func New() DefaultFactory {
	return DefaultFactory{}
}

