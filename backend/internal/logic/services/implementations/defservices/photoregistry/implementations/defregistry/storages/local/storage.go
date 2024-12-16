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

type fileWrap struct {
	file *os.File
}

func (self *fileWrap) Read(b []byte) (int, error) {
	if nil == self.file {
		return 0, cmnerrors.Internal(errors.New("No file wrapped"))
	}

	n, err := self.file.Read(b)

	if nil != err && !errors.Is(err, io.EOF) {
		err = cmnerrors.Internal(err)
	}

	return n, err
}

func (self *fileWrap) Write(b []byte) (int, error) {
	if nil == self.file {
		return 0, cmnerrors.Internal(errors.New("No file wrapped"))
	}

	n, err := self.file.Write(b)

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	return n, err
}

func (self *fileWrap) Close() error {
	if nil == self.file {
		return nil
	}

	err := self.file.Close()

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	return err
}

func openWrap(filename string) (*fileWrap, error) {
	file, err := os.Open(filename)

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	return &fileWrap{file}, err
}

func createWrap(filename string) (*fileWrap, error) {
	file, err := os.Create(filename)

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	return &fileWrap{file}, err
}

func removeWrap(filename string) error {
	err := os.Remove(filename)

	if nil != err {
		err = cmnerrors.Internal(err)
	}

	return err
}

// Returns path to temp data
func (self *storage) WriteTempData(content []byte) (string, error) {
	var file *fileWrap
	var path string
	err := mkdir(self.tempPath)

	if nil == err {
		path, err = generatePath(self.tempPath)
	}

	if nil == err {
		file, err = createWrap(path)
	}

	if nil == err {
		_, err = file.Write(content)
	}

	if nil != file {
		file.Close()
	}

	return path, err
}

func pipe(reader io.Reader, writer io.Writer) error {
	const BUFSIZE int64 = 4096
	var chunk [BUFSIZE]byte
	var err error

	for finish, read := false, 0; nil == err && !finish; {
		read, err = reader.Read(chunk[:])

		if errors.Is(err, io.EOF) {
			err = nil
			finish = true
		}

		if nil == err && 0 != read {
			_, err = writer.Write(chunk[:read])
		}
	}

	return err
}

// Move temp data to persistent storage and return new path
func (self *storage) SaveTempData(tempPath string) (string, error) {
	var path string
	var file *fileWrap
	var tempFile *fileWrap

	err := mkdir(self.persistentPath)

	if nil == err {
		err = mkdir(self.tempPath)
	}

	if nil == err {
		tempFile, err = openWrap(tempPath)
	}

	if nil == err {
		path, err = generatePath(self.persistentPath)
	}

	if nil == err {
		file, err = createWrap(path)
	}

	if nil == err {
		err = pipe(tempFile, file)
	}

	tempFile.Close()
	file.Close()

	if nil == err {
		err = removeWrap(tempPath)
	}

	return path, err
}

// Convert path to href
func (self *storage) ConvertPath(path string) string {
	return self.converter(path)
}

