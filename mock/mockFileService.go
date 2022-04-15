package service

import (
	"io/fs"
	"io/ioutil"
	"teamt5interview/service"
)

type IMockFileService struct {
	fileContext []byte
	readDirPath string
}

func CreateMockFileService(fileContext []byte, readDirPath string) service.FileService {
	return &IMockFileService{
		fileContext: fileContext,
		readDirPath: readDirPath,
	}
}

func (mockService *IMockFileService) Read(filepath string) ([]byte, error) {
	return mockService.fileContext, nil
}

func (mockService *IMockFileService) Write(filepath string, writeBytes []byte) error {
	mockService.fileContext = writeBytes
	return nil
}

func (mockService *IMockFileService) MakeDir(fileDir string) (bool, error) {
	return true, nil
}

func (mockService *IMockFileService) ReadDir(fileDir string) ([]fs.FileInfo, error) {
	return ioutil.ReadDir(mockService.readDirPath)
}

func (mockService *IMockFileService) Delete(fileDir string) error {
	mockService.fileContext = nil
	return nil
}
