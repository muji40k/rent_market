package simple

import (
	"bytes"
	"math/rand"
	"rent_service/internal/logic/services/implementations/defservices/codegen"
)

type generator struct {
	len uint
}

func New(lenght uint) codegen.IGenerator {
	return &generator{lenght}
}

func getNumber() rune {
	return rune(rand.Intn(10) + '0')

}

func (self *generator) Generate() string {
	var buffer bytes.Buffer

	for i := uint(0); self.len > i; i++ {
		buffer.WriteRune(getNumber())
	}

	return buffer.String()
}

