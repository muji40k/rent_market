package models

import (
	"rent_service/builders/misc/nullcommon"
	"rent_service/internal/domain/models"
	"rent_service/misc/nullable"
	"time"

	"github.com/google/uuid"
)

type PhotoBuilder struct {
	id          uuid.UUID
	path        string
	mime        string
	placeholder string
	description string
	date        time.Time
}

func NewPhoto() *PhotoBuilder {
	return &PhotoBuilder{}
}

func (self *PhotoBuilder) WithId(id uuid.UUID) *PhotoBuilder {
	self.id = id
	return self
}

func (self *PhotoBuilder) WithPath(path string) *PhotoBuilder {
	self.path = path
	return self
}

func (self *PhotoBuilder) WithMime(mime string) *PhotoBuilder {
	self.mime = mime
	return self
}

func (self *PhotoBuilder) WithPlaceholder(placeholder string) *PhotoBuilder {
	self.placeholder = placeholder
	return self
}

func (self *PhotoBuilder) WithDescription(description string) *PhotoBuilder {
	self.description = description
	return self
}

func (self *PhotoBuilder) WithDate(date time.Time) *PhotoBuilder {
	self.date = date
	return self
}

func (self *PhotoBuilder) Build() models.Photo {
	return models.Photo{
		Id:          self.id,
		Path:        self.path,
		Mime:        self.mime,
		Placeholder: self.placeholder,
		Description: self.description,
		Date:        self.date,
	}
}

type TempPhotoBuilder struct {
	id          uuid.UUID
	path        *string
	mime        string
	placeholder string
	description string
	create      time.Time
}

func NewTempPhoto() *TempPhotoBuilder {
	return &TempPhotoBuilder{}
}

func (self *TempPhotoBuilder) WithId(id uuid.UUID) *TempPhotoBuilder {
	self.id = id
	return self
}

func (self *TempPhotoBuilder) WithPath(path *nullable.Nullable[string]) *TempPhotoBuilder {
	self.path = nullcommon.CopyPtrIfSome(path)
	return self
}

func (self *TempPhotoBuilder) WithMime(mime string) *TempPhotoBuilder {
	self.mime = mime
	return self
}

func (self *TempPhotoBuilder) WithPlaceholder(placeholder string) *TempPhotoBuilder {
	self.placeholder = placeholder
	return self
}

func (self *TempPhotoBuilder) WithDescription(description string) *TempPhotoBuilder {
	self.description = description
	return self
}

func (self *TempPhotoBuilder) WithCreate(create time.Time) *TempPhotoBuilder {
	self.create = create
	return self
}

func (self *TempPhotoBuilder) Build() models.TempPhoto {
	return models.TempPhoto{
		Id:          self.id,
		Path:        self.path,
		Mime:        self.mime,
		Placeholder: self.placeholder,
		Description: self.description,
		Create:      self.create,
	}
}

