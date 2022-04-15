package service

import (
	"io/fs"
	"io/ioutil"
	"os"
	"teamt5interview/utils"
)

type FileService interface {
	MakeDir(fileDir string) (bool, error)
	ReadDir(fileDir string) ([]fs.FileInfo, error)
	Read(filepath string) ([]byte, error)
	Write(filepath string, writeBytes []byte) error
	Delete(fileDir string) error
}

type IFileService struct {
}

func CreateFileService() FileService {
	return &IFileService{}
}

func (service *IFileService) MakeDir(fileDir string) (bool, error) {
	return utils.MakeDir(fileDir)
}

func (service *IFileService) ReadDir(fileDir string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(fileDir)
}

func (service *IFileService) Read(filepath string) ([]byte, error) {
	bData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return bData, nil
}

func (service *IFileService) Write(filepath string, writeBytes []byte) error {
	err := ioutil.WriteFile(filepath, writeBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (service *IFileService) Delete(filepath string) error {
	return os.Remove(filepath)
}
