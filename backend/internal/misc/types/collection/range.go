package collection

type rng struct {
	start int
	end   int
	inc   int
}

func def() rng {
	return rng{0, 0, 1}
}

func (self *rng) check() {
	if 0 < self.inc && self.end < self.start ||
		0 > self.inc && self.start < self.end {
		self.start, self.end = self.end, self.start
	} else if 0 == self.inc {
		if self.end < self.start {
			self.inc = -1
		} else if self.end > self.start {
			self.inc = 1
		}
	}
}

type rangeCollection struct {
	rng
}

type rangeIterator struct {
	rng
}

func (self *rangeIterator) Next() (int, bool) {
	if self.start == self.end {
		return self.start, false
	}

	out := self.start
	self.start += self.inc

	return out, true
}

func (self *rangeIterator) Skip() bool {
	if self.start == self.end {
		return false
	}

	self.start += self.inc

	return true
}

func (self *rangeCollection) Iter() Iterator[int] {
	return &rangeIterator{self.rng}
}

func RangeStart(start int) func(*rng) {
	return func(col *rng) {
		col.start = start
	}
}

func RangeEnd(end int) func(*rng) {
	return func(col *rng) {
		col.end = end
	}
}

func RangeStep(step int) func(*rng) {
	return func(col *rng) {
		if 0 != step {
			col.inc = step
		}
	}
}

func setupRng(setters []func(*rng)) rng {
	rng := def()

	for _, f := range setters {
		f(&rng)
	}

	rng.check()

	return rng
}

func RangeCollection(setters ...func(*rng)) Collection[int] {
	return &rangeCollection{setupRng(setters)}
}

func RangeIterator(setters ...func(*rng)) Iterator[int] {
	return &rangeIterator{setupRng(setters)}
}

