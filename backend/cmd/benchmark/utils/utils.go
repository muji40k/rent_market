package utils

import "testing"

type Benchmark interface {
	Before()
	Func(b *testing.B)
	After()
}

type Reporter interface {
	Report(res *testing.BenchmarkResult)
}

func RunBenchmark(bench Benchmark, times uint, reporter Reporter) {
	bench.Before()

	for range times {
		res := testing.Benchmark(bench.Func)
		if nil != reporter {
			reporter.Report(&res)
		}
	}

	bench.After()
}

type BenchmarkFunc func(b *testing.B)

func (self BenchmarkFunc) Before()           {}
func (self BenchmarkFunc) Func(b *testing.B) { self(b) }
func (self BenchmarkFunc) After()            {}

type ReporterFunc func(*testing.BenchmarkResult)

func (self ReporterFunc) Report(res *testing.BenchmarkResult) { self(res) }

type MetaWrap struct {
	Name   string
	Target string
	Extra  map[string]string
	bench  Benchmark
}

func NewMeta(name string, target string, bench Benchmark) *MetaWrap {
	return &MetaWrap{
		name, target,
		make(map[string]string),
		bench,
	}
}

func (self *MetaWrap) Before()           { self.bench.Before() }
func (self *MetaWrap) Func(b *testing.B) { self.bench.Func(b) }
func (self *MetaWrap) After()            { self.bench.After() }

