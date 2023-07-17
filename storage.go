package edge_tts_go

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type IStorage interface {
	Write([]byte) error
	Read(string) (io.Reader, error)
	Exist(string) (bool, error)
	Create(string) error
	Close() error
}

var _ IStorage = new(OssStorage)

type OssStorage struct {
	*oss.Client
	bucket  string
	nextPos int64
	objName string
}

func NewOssStorage(endpoint, ak, sk string) (*OssStorage, error) {

	client, err := oss.New(endpoint, ak, sk)
	if err != nil {
		return nil, err
	}
	return &OssStorage{
		Client: client,
	}, nil
}

func (storage *OssStorage) Create(objName string) error {
	storage.objName = objName
	return nil
}

func (storage *OssStorage) Exist(objNmae string) (bool, error) {
	bucket, err := storage.Client.Bucket(storage.bucket)
	if err != nil {
		return false, err
	}
	isExist, err := bucket.IsObjectExist(objNmae)
	if err != nil {
		return false, err
	}
	return isExist, nil
}

func (storage *OssStorage) Write(data []byte) error {
	return storage.upload(storage.bucket, storage.objName, bytes.NewReader(data))
}

func (storage *OssStorage) Read(objectName string) (io.Reader, error) {
	data, err := storage.download(storage.bucket, objectName)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}

func (storage *OssStorage) Close() error {
	storage.nextPos = 0
	return nil
}

func (storage *OssStorage) upload(bucketName string, objName string, objValue io.Reader) error {
	bucket, err := storage.Client.Bucket(bucketName)
	if err != nil {
		return err
	}
	{
		var err error
		storage.nextPos, err = bucket.AppendObject(objName, objValue, storage.nextPos)
		if err != nil {
			return err
		}

	}
	storageType := oss.ObjectStorageClass(oss.StorageStandard)
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)
	return bucket.PutObject(objName, objValue, storageType, objectAcl)
}

func (storage *OssStorage) download(bucketName, objName string) ([]byte, error) {
	bucket, err := storage.Client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	body, err := bucket.GetObject(objName)
	if err != nil {
		return nil, err
	}

	defer body.Close()

	data, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}
	return data, nil

}

type FileStorage struct {
	file *os.File
}

func (storage *FileStorage) Write(data []byte) error {
	_, err := storage.file.Write(data)
	return err
}

func (storage *FileStorage) Read(string) (io.Reader, error) {
	return nil, nil
}
func (storage *FileStorage) Exist(fileName string) (bool, error) {
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	file.Close()
	return true, nil
}

func (storage *FileStorage) Create(fileName string) error {
	file, err := os.Open(fileName)
	if err == nil {
		return nil
	} else {
		file, err = os.Create(fileName)
		if err != nil {
			return err
		}
	}
	storage.file = file
	return nil
}

func (storage *FileStorage) Close() error {
	return storage.file.Close()
}
