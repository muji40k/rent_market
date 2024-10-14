package states

type State uint

const (
	STATE_STORAGE State = iota
	STATE_RENT_AWAIT_DELIVERY_START
	STATE_RENT_AWAIT_DELIVERY
	STATE_RENT_AWAIT
	STATE_RENT
	STATE_RETURN_AWAIT
	STATE_RETURN_FORCED_ISSUE
	STATE_RETURN_FORCED_AWAIT
	STATE_REVOKE_AWAIT_DELIVERY_START
	STATE_REVOKE_AWAIT_DELIVERY
	STATE_REVOKE_AWAIT
	STATE_DELIVERY_AWAIT_START
	STATE_DELIVERY_AWAIT
	STATE_REVOKED
	STATE_MAX
)

func getAll() uint {
	return (1 << STATE_MAX) - 1
}

func getInvertMask(mask uint) uint {
	return getAll() & ^mask
}

func getMask(states ...State) uint {
	out := uint(0)

	for _, state := range states {
		out |= 1 << state
	}

	return out
}

func getStates(state uint) []State {
	out := make([]State, 0, STATE_MAX)

	for i := uint(0); uint(STATE_MAX) > i && 0 != state; i++ {
		if 1 == 1&state {
			out = append(out, State(i))
		}

		state >>= 1
	}

	return out
}

