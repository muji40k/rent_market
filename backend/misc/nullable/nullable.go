package nullable

type Nullable[T any] struct {
	Value T
	Valid bool
}

func Some[T any](value T) *Nullable[T] {
	return &Nullable[T]{value, true}
}

func None[T any]() *Nullable[T] {
	return &Nullable[T]{}
}

func And[T any, F any](self *Nullable[T], value *Nullable[F]) *Nullable[F] {
	if nil == self {
		return nil
	}

	if self.Valid {
		return value
	} else {
		return &Nullable[F]{}
	}
}

func AndFunc[T any, F any](self *Nullable[T], f func(*T) *Nullable[F]) *Nullable[F] {
	if nil == self {
		return nil
	}

	if self.Valid {
		return f(&self.Value)
	} else {
		return &Nullable[F]{}
	}
}

func Or[T any](self *Nullable[T], value *Nullable[T]) *Nullable[T] {
	if nil == self {
		return nil
	}

	if !self.Valid {
		return value
	} else {
		return self
	}
}

func OrFunc[T any](self *Nullable[T], f func() *Nullable[T]) *Nullable[T] {
	if nil == self {
		return nil
	}

	if !self.Valid {
		return f()
	} else {
		return self
	}
}

func Map[T any, F any](self *Nullable[T], f func(*T) F) *Nullable[F] {
	if nil == self {
		return nil
	}

	var out = &Nullable[F]{}

	if self.Valid {
		out.Value = f(&self.Value)
		out.Valid = true
	}

	return out
}

func IfSome[T any](self *Nullable[T], f func(*T)) {
	if nil != self && self.Valid {
		f(&self.Value)
	}
}

func IfNone[T any](self *Nullable[T], f func()) {
	if nil != self && !self.Valid {
		f()
	}
}

func GetOr[T any](self *Nullable[T], def T) T {
	if nil != self && self.Valid {
		return self.Value
	} else {
		return def
	}
}

func GetOrFunc[T any](self *Nullable[T], f func() T) T {
	if nil != self && self.Valid {
		return self.Value
	} else {
		return f()
	}
}

