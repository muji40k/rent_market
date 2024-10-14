package photoregistry

import (
	"errors"
	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	photo_provider "rent_service/internal/repository/context/providers/photo"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"

	"github.com/google/uuid"
)

type repoproviders struct {
	photo photo_provider.IProvider
	temp  photo_provider.ITempProvider
}

type Registry struct {
	repos   repoproviders
	storage IStorage
}

type IStorage interface {
	// Returns path to temp data
	WriteTempData(content []byte) (string, error)
	// Move temp data to persistent storage and return new path
	SaveTempData(tempPath string) (string, error)
	// Convert path to href
	ConvertPath(path string) string
}

func New(
	photo photo_provider.IProvider,
	temp photo_provider.ITempProvider,
	storage IStorage,
) *Registry {
	return &Registry{repoproviders{photo, temp}, storage}
}

func mapTempPhoto(value *models.TempPhoto) models.Photo {
	return models.Photo{
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
	}
}

func (self *Registry) SaveTempPhoto(tempId uuid.UUID, content []byte) error {
	repo := self.repos.temp.GetTempPhotoRepository()
	photo, err := repo.GetById(tempId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err && nil != photo.Path {
		err = cmnerrors.AlreadyExists("photo_data")
	}

	if nil == err {
		var path string
		path, err = self.storage.WriteTempData(content)

		if nil != err {
			err = cmnerrors.Internal(err)
		} else {
			photo.Path = new(string)
			*photo.Path = path
		}
	}

	if nil == err {
		err = repo.Update(photo)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return err
}

// Move temp photo to persistent storage and return new models.Photo
// identifier
func (self *Registry) MoveFromTemp(tempId uuid.UUID) (uuid.UUID, error) {
	var photo models.Photo
	trepo := self.repos.temp.GetTempPhotoRepository()
	temp, err := trepo.GetById(tempId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	if nil == err && nil == temp.Path {
		err = cmnerrors.Unknown("photo_data")
	}

	if nil == err {
		photo = mapTempPhoto(&temp)
		photo.Path, err = self.storage.SaveTempData(*temp.Path)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	if nil == err {
		repo := self.repos.photo.GetPhotoRepository()
		photo, err = repo.Create(photo)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.AlreadyExists(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return photo.Id, err
}

func (self *Registry) MoveFromTemps(tempIds ...uuid.UUID) ([]uuid.UUID, error) {
	var err error
	out := make([]uuid.UUID, 0, len(tempIds))

	for i := 0; nil == err && len(tempIds) > i; i++ {
		id, cerr := self.MoveFromTemp(tempIds[i])

		if nil == cerr {
			out = append(out, id)
		} else {
			err = cerr
		}
	}

	return out, err
}

func (self *Registry) ConvertPath(path string) string {
	return self.storage.ConvertPath(path)
}

