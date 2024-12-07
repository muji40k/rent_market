package main

import (
	"bufio"
	"fmt"
	"generators/static"
	"io"
	"math/rand/v2"
	"os"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/loremipsum.v1"
)

type Generator interface {
	Generate(writer io.Writer)
}

type Content interface {
	Width() uint
	At(i uint) string
	Next() bool
}

type SQLGenerator struct {
	Table   string
	Schema  string
	Content Content
}

func (self *SQLGenerator) Generate(writer io.Writer) {
	if 0 == self.Content.Width() {
		return
	}

	var line []string

	_, err := writer.Write([]byte(fmt.Sprintf(
		"INSERT INTO %v\n(%v)\nVALUES",
		self.Table, self.Schema,
	)))

	if nil == err {
		line = make([]string, self.Content.Width())
	}

	if nil == err && self.Content.Next() {
		for i := uint(0); self.Content.Width() > i; i++ {
			line[i] = self.Content.At(i)
		}

		_, err = writer.Write([]byte("\n(" + strings.Join(line, ",") + ")"))
	}

	for nil == err && self.Content.Next() {
		for i := uint(0); self.Content.Width() > i; i++ {
			line[i] = self.Content.At(i)
		}

		_, err = writer.Write([]byte(",\n(" + strings.Join(line, ",") + ")"))
	}

	if nil == err {
		_, err = writer.Write([]byte(";"))
	}

	if nil != err {
		panic(err)
	}
}

type CSVGenerator struct {
	Content Content
}

func (self *CSVGenerator) Generate(writer io.Writer) {
	if 0 == self.Content.Width() {
		return
	}

	var err error
	line := make([]string, self.Content.Width())

	for nil == err && self.Content.Next() {
		for i := uint(0); self.Content.Width() > i; i++ {
			line[i] = self.Content.At(i)
		}

		_, err = writer.Write([]byte(strings.Join(line, ",") + "\n"))
	}

	if nil != err {
		panic(err)
	}
}

func quoteWrap(f func() string) func() string {
	return func() string {
		return fmt.Sprintf("'%v'", f())
	}
}

func doubleQuoteWrap(f func() string) func() string {
	return func() string {
		return fmt.Sprintf("\"%v\"", f())
	}
}

func castWrap(t string, f func() string) func() string {
	return func() string {
		return fmt.Sprintf("'%v'::%v", f(), t)
	}
}

func uuidWrap(f func() string) func() string {
	return castWrap("uuid", f)
}

func sampleNamer(entity string) func() string {
	counter := 0
	return func() string {
		counter++
		return fmt.Sprintf("Sample %v %v", entity, counter)
	}
}

func idGenerator() string {
	value, err := uuid.NewRandom()

	if nil != err {
		panic(err)
	}

	return value.String()
}

func randomFromList[T any](list []T, f func(*T) string) func() string {
	return func() string {
		return f(&list[rand.Uint()%uint(len(list))])
	}
}

func duplicator(n uint, f func() string) func() string {
	var last string
	i := uint(0)

	return func() string {
		if 0 == i {
			last = f()
		}

		i = (i + 1) % n

		return last
	}
}

func cycle(fs ...func() string) func() string {
	i := 0

	return func() string {
		j := i
		i = (i + 1) % len(fs)
		return fs[j]()
	}
}

func cycleList[T any](list []T, f func(*T) string) func() string {
	i := 0

	return func() string {
		j := i
		i = (i + 1) % len(list)
		return f(&list[j])
	}
}

func constant(value string) func() string {
	return func() string {
		return value
	}
}

func NewSampler(creators ...func() string) Sampler {
	return Sampler{creators}
}

type Sampler struct {
	creators []func() string
}

type SIterator struct {
	left     uint
	creators []func() string
}

func (self *SIterator) At(i uint) string {
	return self.creators[i%uint(len(self.creators))]()
}

func (self *SIterator) Next() bool {
	if 0 == self.left {
		return false
	}

	self.left--
	return 0 != self.left
}

func (self *SIterator) Width() uint {
	return uint(len(self.creators))
}

func (self Sampler) ToContent(samples uint) Content {
	return &SIterator{samples + 1, self.creators}
}

func SQLProductGenerator(content Content) Generator {
	return &SQLGenerator{
		Table:   "products.products",
		Schema:  `id, "name", category_id, description, modification_date, modification_source`,
		Content: content,
	}
}

func SQLProductSampler(sampler Sampler) Sampler {
	if 4 != len(sampler.creators) {
		panic("Wrong sampler")
	}

	cpy := Sampler{make([]func() string, 6)}

	cpy.creators[0] = uuidWrap(sampler.creators[0])
	cpy.creators[1] = quoteWrap(sampler.creators[1])
	cpy.creators[2] = uuidWrap(sampler.creators[2])
	cpy.creators[3] = quoteWrap(sampler.creators[3])
	cpy.creators[4] = constant("now()")
	cpy.creators[5] = constant("'preset'")

	return cpy
}

func DefaultProductSampler(categoryIds []string) Sampler {
	return NewSampler(
		idGenerator,
		sampleNamer("Product"),
		randomFromList(categoryIds, func(s *string) string { return *s }),
		func() string {
			return loremipsum.New().Sentences((rand.Int() % 4) + 1)
		},
	)
}

func BenchmarkProductSampler(category string) Sampler {
	return NewSampler(
		idGenerator,
		sampleNamer("Product"),
		constant(category),
		func() string {
			return loremipsum.New().Sentences((rand.Int() % 4) + 1)
		},
	)
}

func SQLProductCharsGenerator(content Content) Generator {
	return &SQLGenerator{
		Table:   `products."characteristics"`,
		Schema:  `id, product_id, "name", value, modification_date, modification_source`,
		Content: content,
	}
}

func SQLProductCharsSampler(sampler Sampler) Sampler {
	if 4 != len(sampler.creators) {
		panic("Wrong sampler")
	}

	cpy := Sampler{make([]func() string, 6)}

	cpy.creators[0] = uuidWrap(sampler.creators[0])
	cpy.creators[1] = uuidWrap(sampler.creators[1])
	cpy.creators[2] = quoteWrap(sampler.creators[2])
	cpy.creators[3] = quoteWrap(sampler.creators[3])
	cpy.creators[4] = constant("now()")
	cpy.creators[5] = constant("'preset'")

	return cpy
}

func DefaultProductCharsSampler(productIds []string) Sampler {
	return NewSampler(
		idGenerator,
		duplicator(2,
			cycleList(productIds, func(s *string) string { return *s }),
		),
		cycleList(
			[]string{"param_1", "param_2"},
			func(s *string) string { return *s }),
		cycle(
			sampleNamer("Value"),
			func() string { return fmt.Sprint(10 * rand.Float64()) },
		),
	)
}

func BenchmarkProductCharsSampler(productIds []string) Sampler {
	return NewSampler(
		idGenerator,
		duplicator(2,
			cycleList(productIds, func(s *string) string { return *s }),
		),
		cycleList(
			[]string{"key1", "key2"},
			func(s *string) string { return *s }),
		cycle(
			randomFromList(
				[]string{"value1", "value2", "value3", "value4"},
				func(s *string) string { return *s },
			),
			func() string { return fmt.Sprint(0.5 + 2*rand.Float64()) },
		),
	)
}

func GeneralCSVSampler(sampler Sampler) Sampler {
	cpy := Sampler{make([]func() string, len(sampler.creators)+1)}

	for i, f := range sampler.creators {
		cpy.creators[i] = doubleQuoteWrap(f)
	}

	cpy.creators[len(sampler.creators)] = constant("\"preset\"")

	return cpy
}

func main() {
	if 2 != len(os.Args) {
		panic(strings.Join(os.Args, ", "))
	}

	sampler := BenchmarkProductCharsSampler(static.BENCHMARK_PRODUCT_IDS[:])
	generator := CSVGenerator{
		Content: GeneralCSVSampler(sampler).ToContent(2000),
	}

	file, err := os.Create(os.Args[1])

	if nil != err {
		panic(err)
	}

	defer file.Close()
	writer := bufio.NewWriter(file)
	generator.Generate(writer)
	writer.Flush()
}

