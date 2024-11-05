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

func Then[T any, F any](self *Nullable[T], f func(*T) Nullable[F]) *Nullable[F] {
	if nil == self {
		return nil
	}

	var out = &Nullable[F]{}

	if self.Valid {
		*out = f(&self.Value)
	}

	return out
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

func GerOrInsert[T any](self *Nullable[T], def T) T {
	if nil != self && self.Valid {
		return self.Value
	} else {
		return def
	}
}

func GerOrInsertFunc[T any](self *Nullable[T], f func() T) T {
	if nil != self && self.Valid {
		return self.Value
	} else {
		return f()
	}
}

