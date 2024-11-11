package models

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math/rand/v2"
	modelsb "rent_service/builders/domain/models"
	"rent_service/builders/misc/uuidgen"
	"rent_service/misc/nullable"
	"time"
)

func PhotoRandomId() *modelsb.PhotoBuilder {
	return modelsb.NewPhoto().
		WithId(uuidgen.Generate())
}

func PhotoExample(
	prefix string,
	date *nullable.Nullable[time.Time],
) *modelsb.PhotoBuilder {
	return PhotoRandomId().
		WithPath("/path/to/" + prefix).
		WithMime("image/png").
		WithPlaceholder(prefix).
		WithDescription("Description for " + prefix).
		WithDate(nullable.GetOrFunc(date, time.Now))
}

func TempPhotoRandomId() *modelsb.TempPhotoBuilder {
	return modelsb.NewTempPhoto().
		WithId(uuidgen.Generate())
}

func TempPhotoExampleEmpty(
	prefix string,
	date *nullable.Nullable[time.Time],
) *modelsb.TempPhotoBuilder {
	return TempPhotoRandomId().
		WithMime("image/png").
		WithPlaceholder(prefix).
		WithDescription("Description for " + prefix).
		WithCreate(nullable.GetOrFunc(date, time.Now))
}

func TempPhotoExampleUploaded(
	prefix string,
	date *nullable.Nullable[time.Time],
) *modelsb.TempPhotoBuilder {
	return TempPhotoExampleEmpty(prefix, date).
		WithPath(nullable.Some("/path/to/" + prefix))
}

func getRandomColor() color.NRGBA {
	return color.NRGBA{
		uint8(rand.Uint32() % 256),
		uint8(rand.Uint32() % 256),
		uint8(rand.Uint32() % 256),
		255,
	}
}

func ImagePNGContent(size *nullable.Nullable[int]) []byte {
	var rsize = nullable.GetOr(size, 200)
	img := image.NewRGBA(image.Rect(0, 0, rsize, rsize))

	for i := range rsize {
		for j := range rsize {
			img.Set(i, j, getRandomColor())
		}
	}

	var buf bytes.Buffer
	var out []byte
	err := png.Encode(&buf, img)

	if nil == err {
		out = make([]byte, buf.Len())
		_, err = buf.Read(out)
	}

	if nil != err {
		panic(err)
	}

	return out
}

