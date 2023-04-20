package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/afrancoc2000/application-helper-ai/internal/models"
)

const (
	fileNameLabel    = "File Name:"
	filePathLabel    = "File Path:"
	fileContentLabel = "File Content:"
	fileQuotesLabel  = "```"
)

type FileFactory interface {
	BuildProject(openAIResult string) error
}

type fileFactory struct {
}

func NewFileFactory() fileFactory {
	return fileFactory{}
}

func (f *fileFactory) BuildProject(openAIResult string) error {

	instructions := strings.Split(openAIResult, "\n")
	files := f.generateFiles(instructions)

	for _, file := range files {
		err := f.saveFile(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *fileFactory) generateFiles(instructions []string) []models.AppFile {
	files := []models.AppFile{}
	openContent := false

	var currentFile models.AppFile
	for _, instruction := range instructions {

		if strings.Contains(instruction, fileQuotesLabel) {
			openContent = !openContent
			if !openContent {
				files = append(files, currentFile)
			}
			continue

		} else if openContent {
			currentFile.Content = currentFile.Content + "\n" + instruction
			continue

		} else if strings.Contains(instruction, filePathLabel) {
			re := regexp.MustCompile(filePathLabel + `\s+(\S+)`)
			currentFile.Path = re.FindStringSubmatch(instruction)[1]
			continue

		} else if strings.Contains(instruction, fileNameLabel) {
			currentFile = models.AppFile{}
			re := regexp.MustCompile(fileNameLabel + `\s+(\S+)`)
			currentFile.Name = re.FindStringSubmatch(instruction)[1]
			continue
		}
	}

	fmt.Printf("Files: %v", len(files))
	for _, file := range files {
		fmt.Printf("{name: %s, path: %s, content: %v}\n", file.Name, file.Path, len(file.Content))
	}

	return files
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
