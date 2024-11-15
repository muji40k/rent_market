package photo

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"

	"rent_service/internal/domain/models"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/authenticator"
	"rent_service/internal/logic/services/implementations/defservices/emptymathcer"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry"
	"rent_service/internal/logic/services/interfaces/photo"
	"rent_service/internal/logic/services/types/date"
	"rent_service/internal/logic/services/types/token"
	photo_providers "rent_service/internal/repository/context/providers/photo"
	repo_errors "rent_service/internal/repository/errors/cmnerrors"
)

var ALLOWED_MIMES []string = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
	"image/gif",
}

type repoproviders struct {
	photo photo_providers.IProvider
	temp  photo_providers.ITempProvider
}

type service struct {
	repos         repoproviders
	authenticator authenticator.IAuthenticator
	registry      photoregistry.IRegistry
}

func New(
	authenticator authenticator.IAuthenticator,
	registry photoregistry.IRegistry,
	photo photo_providers.IProvider,
	temp photo_providers.ITempProvider,
) photo.IService {
	return &service{repoproviders{photo, temp}, authenticator, registry}
}

func checkMime(mime string) error {
	if !slices.Contains(ALLOWED_MIMES, mime) {
		return photo.ErrorUnsupportedMime{Mime: mime}
	}

	return nil
}

func (self *service) CreateTempPhoto(
	token token.Token,
	photo photo.Description,
) (uuid.UUID, error) {
	created := models.TempPhoto{
		Mime:        photo.Mime,
		Placeholder: photo.Placeholder,
		Description: photo.Description,
		Create:      time.Now(),
	}

	err := emptymathcer.Match(
		emptymathcer.Item("mime", photo.Mime),
		emptymathcer.Item("placeholder", photo.Placeholder),
	)

	if nil == err {
		err = checkMime(photo.Mime)
	}

	if nil == err {
		_, err = self.authenticator.LoginWithToken(token)
	}

	if nil == err {
		repo := self.repos.temp.GetTempPhotoRepository()
		created, err = repo.Create(created)

		if cerr := (repo_errors.ErrorDuplicate{}); errors.As(err, &cerr) {
			err = cmnerrors.AlreadyExists(cerr.What...)
		} else if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return created.Id, err
}

func (self *service) UploadTempPhoto(
	token token.Token,
	photoId uuid.UUID,
	content []byte,
) error {
	_, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		err = self.registry.SaveTempPhoto(photoId, content)
	}

	return err
}

func (self *service) mapTempPhoto(value *models.TempPhoto) photo.TempPhoto {
	return photo.TempPhoto{
		Id:          value.Id,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Date:        date.New(value.Create),
	}
}

func (self *service) GetTempPhoto(
	token token.Token,
	photoId uuid.UUID,
) (photo.TempPhoto, error) {
	var photo models.TempPhoto
	_, err := self.authenticator.LoginWithToken(token)

	if nil == err {
		repo := self.repos.temp.GetTempPhotoRepository()
		photo, err = repo.GetById(photoId)

		if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
			err = cmnerrors.NotFound(cerr.What...)
		} else if nil != err {
			err = cmnerrors.Internal(cmnerrors.DataAccess(err))
		}
	}

	return self.mapTempPhoto(&photo), err
}

func (self *service) mapPhoto(value *models.Photo) photo.Photo {
	return photo.Photo{
		Id:          value.Id,
		Mime:        value.Mime,
		Placeholder: value.Placeholder,
		Description: value.Description,
		Href:        self.registry.ConvertPath(value.Path),
		Date:        date.New(value.Date),
	}
}

func (self *service) GetPhoto(photoId uuid.UUID) (photo.Photo, error) {
	var photo models.Photo

	repo := self.repos.photo.GetPhotoRepository()
	photo, err := repo.GetById(photoId)

	if cerr := (repo_errors.ErrorNotFound{}); errors.As(err, &cerr) {
		err = cmnerrors.NotFound(cerr.What...)
	} else if nil != err {
		err = cmnerrors.Internal(cmnerrors.DataAccess(err))
	}

	return self.mapPhoto(&photo), err
}

