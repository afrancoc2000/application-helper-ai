package models

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type AppFile struct {
	Name    string `required:"true" json:"fileName"`
	Path    string `required:"true" json:"filePath"`
	Content string `required:"true" json:"fileContent"`
}

const parseError = "Sorry, Couldn't parse OpenAI response: %s"

func AppFileFromString(text string) ([]AppFile, error) {
	files := []AppFile{}
	err := json.Unmarshal([]byte(text), &files)
	if err != nil {
		return nil, fmt.Errorf(parseError, err)
	}

	return files, nil
}

func removeFormatting(data string) (string, error) {
	var v interface{}
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func removeTextBeforeAndAfterBrackets(data string) string {
	re := regexp.MustCompile(`(?s)^.*?(\[.*\]).*$`)
	return re.FindStringSubmatch(data)[1]
}
