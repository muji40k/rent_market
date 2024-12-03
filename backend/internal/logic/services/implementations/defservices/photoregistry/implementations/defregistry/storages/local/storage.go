package local

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"rent_service/internal/logic/services/errors/cmnerrors"
	"rent_service/internal/logic/services/implementations/defservices/photoregistry/implementations/defregistry"

	"github.com/google/uuid"
)

type storage struct {
	tempPath       string
	persistentPath string
	converter      func(string) string
}

func New(
	tempPath string,
	persistentPath string,
	converter func(string) string,
) defregistry.IStorage {
	return &storage{
		filepath.FromSlash(tempPath),
		filepath.FromSlash(persistentPath),
		converter,
	}
}

func mkdir(path string) error {
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	return err
}

func checkFile(path string) (bool, error) {
	if _, cerr := os.Stat(path); errors.Is(cerr, os.ErrNotExist) {
		return false, nil
	} else if nil != cerr {
		return false, cmnerrors.Internal(cerr)
	}

	return true, nil
}

func generatePath(base string) (string, error) {
	var path string
	var err error
	var id uuid.UUID

	for exists := true; nil == err && exists; {
		id, err = uuid.NewRandom()

		if nil != err {
			err = cmnerrors.Internal(err)
		} else {
			path = filepath.Join(base, id.String())
			exists, err = checkFile(path)
		}
	}

	return path, err
}

// Returns path to temp data
func (self *storage) WriteTempData(content []byte) (string, error) {
	var file *os.File
	var path string
	err := mkdir(self.tempPath)

	if nil == err {
		path, err = generatePath(self.tempPath)
	}

	if nil == err {
		file, err = os.Create(path)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	if nil == err {
		_, err = file.Write(content)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	if nil != file {
		file.Close()
	}

	return path, err
}

// Move temp data to persistent storage and return new path
func (self *storage) SaveTempData(tempPath string) (string, error) {
	var file *os.File
	var path string
	const BUFSIZE int64 = 4096
	var chunk [BUFSIZE]byte
	var tempFile *os.File

	err := mkdir(self.persistentPath)

	if nil == err {
		err = mkdir(self.tempPath)
	}

	if nil == err {
		tempFile, err = os.Open(tempPath)
	}

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	if nil == err {
		path, err = generatePath(self.persistentPath)
	}

	if nil == err {
		file, err = os.Create(path)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	for finish, read := false, 0; nil == err && !finish; {
		read, err = tempFile.Read(chunk[:])

		if errors.Is(err, io.EOF) {
			err = nil
			finish = true
		} else if nil != err {
			err = cmnerrors.Internal(err)
		}

		if nil == err && 0 != read {
			_, err = file.Write(chunk[:read])

			if nil != err {
				err = cmnerrors.Internal(err)
			}
		}
	}

	if nil != tempFile {
		tempFile.Close()
	}

	if nil != file {
		file.Close()
	}

	if nil == err {
		err = os.Remove(tempPath)

		if nil != err {
			err = cmnerrors.Internal(err)
		}
	}

	return path, err
}

// Convert path to href
func (self *storage) ConvertPath(path string) string {
	return self.converter(path)
}

