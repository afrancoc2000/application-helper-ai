package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
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

type factoryFile struct {
	name    string
	path    string
	content string
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

func (f *fileFactory) generateFiles(instructions []string) []factoryFile {
	files := []factoryFile{}
	openContent := false

	var currentFile factoryFile
	for _, instruction := range instructions {

		if strings.Contains(instruction, fileQuotesLabel) {
			openContent = !openContent
			if !openContent {
				files = append(files, currentFile)
			}
			continue

		} else if openContent {
			currentFile.content = currentFile.content + "\n" + instruction
			continue

		} else if strings.Contains(instruction, filePathLabel) {
			re := regexp.MustCompile(filePathLabel + `\s+(\S+)`)
			currentFile.path = re.FindStringSubmatch(instruction)[1]
			continue

		} else if strings.Contains(instruction, fileNameLabel) {
			currentFile = factoryFile{}
			re := regexp.MustCompile(fileNameLabel + `\s+(\S+)`)
			currentFile.name = re.FindStringSubmatch(instruction)[1]
			continue
		}
	}

	fmt.Printf("Files: %v", len(files))		
	for _, file := range files {
		fmt.Printf("{name: %s, path: %s, content: %v}\n", file.name, file.path, len(file.content))		
	}

	return files
}

func (f *fileFactory) saveFile(factoryFile factoryFile) error {
	err := os.MkdirAll(factoryFile.path, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(fmt.Sprintf("%s%s", factoryFile.path, factoryFile.name))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(factoryFile.content)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
