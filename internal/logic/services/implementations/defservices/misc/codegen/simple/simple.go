package simple

import (
	"bytes"
	"math/rand"
	"rent_service/internal/logic/services/implementations/defservices/misc/codegen"
)

type generator struct {
	len uint
}

func New(lenght uint) codegen.IGenerator {
	return &generator{lenght}
}

func getNumber() rune {
	switch rand.Intn(10) {
	case 0:
		return '0'
	case 1:
		return '1'
	case 2:
		return '2'
	case 3:
		return '3'
	case 4:
		return '4'
	case 5:
		return '5'
	case 6:
		return '6'
	case 7:
		return '7'
	case 8:
		return '8'
	case 9:
		return '9'
	default:
		return '0'
	}

}

func (self *generator) Generate() string {
	var buffer bytes.Buffer

	for i := uint(0); self.len > i; i++ {
		buffer.WriteRune(getNumber())
	}

	return buffer.String()
}

