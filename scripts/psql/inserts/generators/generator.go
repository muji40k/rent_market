package main

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/loremipsum.v1"
)

type Content interface {
	Width() uint
	At(i uint) string
	Next() bool
}

type Generator struct {
	Table   string
	Schema  string
	Content Content
}

func (self *Generator) Generate() string {
	if 0 == self.Content.Width() {
		return ""
	}

	lines := make([]string, 0)

	for self.Content.Next() {
		line := "(" + self.Content.At(0)

		for i := uint(1); self.Content.Width() > i; i++ {
			line += ", " + self.Content.At(i)
		}

		lines = append(lines, line+")")
	}

	return fmt.Sprintf(
		"INSERT INTO %v\n(%v)\nVALUES\n%v;",
		self.Table, self.Schema, strings.Join(lines, ",\n"),
	)
}

func quoteWrap(f func() string) func() string {
	return func() string {
		return fmt.Sprintf("'%v'", f())
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

func NewSampler(n uint, table string, schema string, creators ...func() string) Sampler {
	return Sampler{table, schema, n, creators}
}

type Sampler struct {
	table    string
	schema   string
	samples  uint
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

func (self *Sampler) toContent() Content {
	return &SIterator{self.samples + 1, self.creators}
}

func (self *Sampler) ToGenerator() Generator {
	return Generator{
		Table:   self.table,
		Schema:  self.schema,
		Content: self.toContent(),
	}
}

// const PRODUCT_TABLE string = "products.products"
// const PRODUCT_SCHEMA string = `id, "name", category_id, description, modification_date, modification_source`
//
// type ProductContent struct {
// }

var CATEGORY_IDS = [15]string{
	"2c1a5afd-043f-42c3-8563-1ff24f0adea2",
	"e19dfa5d-902a-484c-95fb-4bf0818d4420",
	"80c172b1-ae1e-4baa-9cd3-a47dd2047ca5",
	"e08c9037-8fc9-424a-837e-eacfbf1c637d",
	"bebe188c-da95-4dfd-a431-bed64a204e7f",
	"d9114bdc-e10c-4b60-8939-5d10c833daa8",
	"923297f7-3dc3-4d02-a4ef-f82dd4ce1fb9",
	"a12b94b6-9d01-406f-868c-fc2e8efb1524",
	"03d35033-c447-41fe-90a1-bdfc80536991",
	"3973e21f-a423-4720-bd91-ec5ff9b90bd8",
	"4fde0ac7-0d81-421e-97c9-d6d0957c2f7d",
	"e373fc89-1023-4237-bb4f-4f8f3b653b00",
	"ff406572-79c2-4040-b2a8-e0605ab9cf57",
	"32e76612-1cbd-443a-a54a-7be810f87c1d",
	"9b8e6e20-0d57-4415-a2df-9504e85492ac",
}

func product() Sampler {
	return NewSampler(
		100,
		"products.products",
		`id, "name", category_id, description, modification_date, modification_source`,
		uuidWrap(idGenerator), quoteWrap(sampleNamer("Product")),
		uuidWrap(randomFromList(CATEGORY_IDS[:], func(s *string) string { return *s })),
		quoteWrap(func() string {
			return loremipsum.New().Sentences((rand.Int() % 4) + 1)
		}),
		constant("now()"),
		constant("'preset'"),
	)
}

var PRODUCT_IDS = [100]string{
	"f97fff60-b293-4499-be56-171fef99e227",
	"2ee9a40b-b442-4666-b9ea-2ec7892689bc",
	"3c33893d-6b75-4765-8d0c-00457775e7a6",
	"4237efc8-de40-4e71-b2b6-7db4fbb22f59",
	"b44fa60f-9ce0-43f6-89dd-5a6860aeb0a3",
	"6e9d66d9-4249-47f7-9478-8e6409094563",
	"0aaabacf-8bd4-45ee-a5a1-2330f80774d3",
	"41ab1bfa-718b-4f97-bcaf-0ca34d1fc0a2",
	"1104b325-ef31-4ea2-bf83-e39101f8ec97",
	"3ea61213-db2d-42a6-9e1e-e3c4ce4e09f1",
	"9d535d89-8fc6-47f6-b9cc-4d241eff697b",
	"442c8e21-6d88-400a-b1a0-b8b96a4b2da1",
	"8abebf07-64a9-43ba-b359-1aa99a03b104",
	"9d6aff17-93a4-47a1-a366-11fdc49013f8",
	"a02358c6-97e4-44cd-8547-a4ab6bd9cd66",
	"c802e290-f497-422c-8f71-20679d3f1a40",
	"62751cfe-00ff-4817-bd98-38970ea7a4e5",
	"517c445c-ea26-4d85-ac94-feeb48b3cd27",
	"4f907358-d2ad-4d67-8b27-8ec44cd6153a",
	"09c95223-6068-4ecc-b203-8245f960e508",
	"95db5232-0c09-4114-b32c-8a0782ca0dc8",
	"f23d7936-9c2d-4f57-a00e-fa3f891a822e",
	"681aef16-c2d9-48d3-971c-87a9656530ed",
	"ee203ee3-4dbb-4af7-8962-870a73221d8b",
	"ae227d80-01d7-4f2e-a02e-266d6c8306d8",
	"d7bd127e-394d-48c5-a713-6ffb838d820d",
	"43917cf4-8ac5-4581-a7fd-3abb7e9b5126",
	"a56c1374-9086-4276-932e-6ec6b616c9e5",
	"d26ab980-0d2a-46fc-88d2-508e000c9cb0",
	"fab12310-7e9b-4097-bd1e-d4b238cc628b",
	"6b5f78cc-1da6-4a59-8359-a6364fca5f10",
	"3c3a1f00-7975-4c9c-a90b-9400b23de20e",
	"20f41fe4-369b-4ce0-b00a-2ee55a423882",
	"a7675ae0-1e82-464d-acb6-5d8d6c9576fd",
	"977a8cd3-a6ab-42d9-ac2e-2e3b4b53b39d",
	"23436910-a969-46d5-915a-8dd4bde0653d",
	"7528ef0b-73b9-4a2d-ad06-b03bc3b648b2",
	"91eab63d-2dad-4148-adb6-1923d26dd3a5",
	"6d4e7818-ff4c-4f95-82c0-811aff3dbbd7",
	"dbe1d455-eb9b-4a73-805a-93166e542835",
	"68d43a14-8a7f-4456-b0f6-90cf29f0fa79",
	"ed5be431-a82f-495a-bbcd-a9e791aaa29e",
	"18439cdc-f9e6-40ef-8aba-45cb0ea92307",
	"695e14be-9a34-4c55-bd1e-c41b5ccd30ec",
	"460d919e-6df0-4dc1-bcce-5619af804681",
	"60a28565-1328-446e-ba1b-0a2b207d765a",
	"c7864316-7e01-466d-b2d3-774c6ea5d454",
	"2d83dfdf-fb29-47f2-8c9f-da786b019f29",
	"4059ccfa-2934-4e3c-abb3-4f0211978772",
	"db600bea-355a-47dd-803d-a86dce9034cc",
	"f22f81e0-f458-4ddc-ac8a-e99bfec16593",
	"7325dc68-864e-43a7-9a04-16a9d6de5767",
	"13285d17-1b23-4478-af27-c8c676b5dea9",
	"b03af559-663e-4916-bd55-a33278ec632d",
	"2c2f1f5e-2033-41ef-9c4e-50aff6b427cf",
	"d3daae2d-f6a8-4cde-9a09-74aee626a983",
	"465c43b3-31a3-4475-9b4b-91f4b175d278",
	"65b66a9e-f103-476b-86a3-066da7186421",
	"6ba80b9b-d167-4458-a30b-ff5662f8c69e",
	"ba50245a-dfc5-4e98-abb4-245ebfaa6fe4",
	"ac4a81d2-3a23-46aa-9345-41c2bd42dd4c",
	"7f12893a-4226-4ffa-9e63-9ca21d647c92",
	"54ce0b3e-6072-4735-9eae-c3fc937a2879",
	"997e3205-af88-4148-b523-a3e27d0c88a1",
	"ac8a92d1-39c3-494e-978e-be70be6c9c3e",
	"c04b1ffa-d498-4a1c-8e90-19a7cabacf31",
	"30c10def-8813-4538-bc2e-6f9de66f7abb",
	"a9e93af7-026d-476c-aa0f-9c23400ef095",
	"aca136b8-d5d2-4cf1-abcd-7fc69ce43d24",
	"e5acd67c-b54d-4617-a50a-ef5b95f03f9a",
	"78c35734-3452-4894-87d0-c235f43da56a",
	"1360860e-20cd-4dd5-9340-5f81d80c21fe",
	"736cf167-539c-48e3-a4f2-150b547e1ea4",
	"c4190781-84dc-4492-adaf-a3f35e31a3e1",
	"699055ea-25f3-4c7d-a8d7-2d9e8be1c9ab",
	"292f579e-7de4-4618-9140-c32355636eab",
	"17982a80-b0bd-458e-8f31-190285844022",
	"45eb2af7-fc8c-42db-9e64-795147c914aa",
	"414985d4-4d84-4376-ac2f-8f9894c3f4d6",
	"b00d3dc9-18f4-40da-ab46-d86c07e048f0",
	"8e149dd0-c052-45fa-be64-9f7b3e3424f5",
	"fc4d28f5-ae66-48ae-8ec7-d05435abcdb1",
	"832125e6-9c08-46bc-ac7f-1fa551570ad4",
	"46a45e8d-bee8-4514-b700-4ee161864f71",
	"dda2d66c-3aad-444f-b37a-050dd2b893f1",
	"7f96bb72-e785-427c-9d6c-8045a90eecf8",
	"68a5fcaf-5a25-4c23-8e4a-0386e66830e4",
	"881a2782-84de-4ccd-abb9-7c68b5fb1775",
	"40594748-314a-46b7-8f0f-99ca08852807",
	"1bde65ad-b65e-4a74-970d-224cadcae4cd",
	"131a4cae-d711-4552-8d97-d61b739c84ea",
	"0f4a95df-d83c-4db2-95b7-45ca6051b0ff",
	"1d86a7a0-8023-40eb-8f89-f83dcd38de56",
	"decfabb1-7605-4c4d-8316-e276bf43e0d7",
	"4b6f8ade-9179-47ac-bdd0-697998e40daa",
	"309ae9ef-b2cc-4670-9474-8ecf5c4668cd",
	"c892c6a9-93f8-4989-800c-a7d1f21fe48c",
	"ffd03d36-b898-43f0-bf5d-80e53fbc92fe",
	"ada0275b-5f4e-47bf-b9a3-b4701901b991",
	"f675e5f9-ce4b-44e9-affb-3cfc3defc375",
}

func product_chars() Sampler {
	return NewSampler(
		200,
		`products."characteristics"`,
		`id, product_id, "name", value, modification_date, modification_source`,
		uuidWrap(idGenerator),
		duplicator(2, uuidWrap(cycleList(PRODUCT_IDS[:], func(s *string) string { return *s }))),
		quoteWrap(cycleList([]string{"param_1", "param_2"}, func(s *string) string { return *s })),
		cycle(quoteWrap(sampleNamer("Value")), func() string { return fmt.Sprint(10 * rand.Float64()) }),
		constant("now()"),
		constant("'preset'"),
	)
}

func main() {
	s := product_chars()
	g := s.ToGenerator()

	fmt.Println(g.Generate())
}

