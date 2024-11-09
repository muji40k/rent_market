package codegen

//go:generate mockgen -source=interface.go -destination=mock/mock.go

type IGenerator interface {
	Generate() string
}

