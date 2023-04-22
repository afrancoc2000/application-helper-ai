package cli

import (
	"fmt"
	"os"

	"github.com/afrancoc2000/application-helper-ai/internal/models"
)

type FileFactory interface {
	CreateFiles(files []models.AppFile) error
}

type fileFactory struct {
}

func NewFileFactory() FileFactory {
	return &fileFactory{}
}

func (f *fileFactory) CreateFiles(files []models.AppFile) error {

	for _, file := range files {
		err := f.saveFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fileFactory) saveFile(factoryFile models.AppFile) error {
	err := os.MkdirAll(factoryFile.Path, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s%s", factoryFile.Path, factoryFile.Name))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(factoryFile.Content)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
