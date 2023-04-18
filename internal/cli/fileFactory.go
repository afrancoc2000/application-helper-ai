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
		if openContent {
			currentFile.content = currentFile.content + instruction
		}

		if strings.Contains(instruction, fileQuotesLabel) {
			openContent = !openContent
			continue
		}

		if strings.Contains(instruction, filePathLabel) {
			re := regexp.MustCompile(filePathLabel + `\s+(\S+)`)
			currentFile.path = re.FindStringSubmatch(instruction)[1]
			continue
		}
		
		if strings.Contains(instruction, fileNameLabel) {
			currentFile := new(factoryFile)
			re := regexp.MustCompile(fileNameLabel + `\s+(\S+)`)
			currentFile.name = re.FindStringSubmatch(instruction)[1]
			files = append(files, *currentFile)
			continue
		}
	}
	return files
}

func (f *fileFactory) saveFile(factoryFile factoryFile) error {
	err := os.MkdirAll(factoryFile.path, os.ModePerm)
    if err != nil {
        return err
    }

	file, err := os.Create(fmt.Sprintf(".%s%s", factoryFile.path, factoryFile.name))
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