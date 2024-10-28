package rqactions

import (
	"fmt"
	"rent_service/internal/logic/services/types/token"

	"github.com/google/uuid"
)

type Form interface {
	Action() string
}

type Action func(token.Token) error
type Getter[C any, F Form] func(controller *C, id uuid.UUID, form F) (Action, error)

type Actions[C any, F Form] struct {
	getters map[string]Getter[C, F]
	keys    []string
}

func New[C any, F Form](getters map[string]Getter[C, F]) Actions[C, F] {
	keys := make([]string, len(getters))
	i := 0

	for k := range getters {
		keys[i] = k
		i++
	}

	return Actions[C, F]{getters, keys}
}

func (self *Actions[C, F]) GetAction(
	controller *C,
	id uuid.UUID,
	form F,
) (Action, error) {
	if getter, ok := self.getters[form.Action()]; ok {
		return getter(controller, id, form)
	} else {
		return nil, fmt.Errorf(
			"Action '%v' isn't in list of allowed operations %v",
			form.Action(), self.keys,
		)
	}
}

