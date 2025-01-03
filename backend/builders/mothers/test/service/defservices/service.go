package defservices

import (
	"errors"
	"fmt"
	"os"
	"rent_service/builders/misc/uuidgen"
	"rent_service/builders/service/factory/v1/defservices"
	delivery_dummy "rent_service/internal/logic/delivery/implementations/dummy"
	"rent_service/internal/logic/services/implementations/defservices/codegen/simple"
	checker_dummy "rent_service/internal/logic/services/implementations/defservices/paymentcheckers/dummy"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry/storages/local"
	rv1 "rent_service/internal/repository/context/v1"
	"strings"
)

const (
	DEFAULT_MEDIA string = "/server/media"
	DEFAULT_TEMP  string = "/server/temp"
	DEFAULT_HREF  string = "/img"
)

func DefaultPathConverter(path string) string {
	return strings.Replace(
		path,
		DEFAULT_MEDIA,
		DEFAULT_HREF,
		-1,
	)
}

func DefaultServiceFactory(factories rv1.Factories) *defservices.Builder {
	return defservices.New().
		WithRepositoryContext(rv1.New(factories)).
		WithCodegen(simple.New(6)).
		WithDeliveryCreator(delivery_dummy.New()).
		WithPaymentChecker(checker_dummy.New()).
		WithPhotoStorage(local.New(
			DEFAULT_TEMP, DEFAULT_MEDIA,
			DefaultPathConverter,
		))
}

type PhotoRegistry struct{}

func mkdir(path string) {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	if nil != err {
		panic(err)
	}
}

func NewPhotoRegistry() *PhotoRegistry {
	mkdir(DEFAULT_MEDIA)
	mkdir(DEFAULT_TEMP)
	return &PhotoRegistry{}
}

func getPath(base string, file any) string {
	return fmt.Sprintf("%v/%v", base, file)
}

func SavePhoto(base string, content []byte) string {
	id := uuidgen.Generate()

	if _, err := os.Stat(getPath(base, id)); !errors.Is(err, os.ErrNotExist) {
		if nil != err {
			panic(err)
		} else {
			id = uuidgen.Generate()
		}
	}

	path := getPath(base, id)
	file, err := os.Create(path)

	if nil == err {
		_, err = file.Write(content)
	}

	if nil == err {
		file.Close()
	} else {
		panic(err)
	}

	return path
}

func (*PhotoRegistry) SavePhoto(content []byte) string {
	return SavePhoto(DEFAULT_MEDIA, content)
}

func (*PhotoRegistry) SaveTempPhoto(content []byte) string {
	return SavePhoto(DEFAULT_TEMP, content)
}

func Clear(dirname string) {
	entires, err := os.ReadDir(dirname)

	if nil != err {
		panic(err)
	}

	for _, entry := range entires {
		name := dirname + "/" + entry.Name()

		if entry.Type().IsDir() {
			os.RemoveAll(name)
		} else {
			os.Remove(name)
		}
	}
}

func (*PhotoRegistry) Clear() {
	Clear(DEFAULT_MEDIA)
	Clear(DEFAULT_TEMP)
}

