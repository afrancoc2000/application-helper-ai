package models

import (
	"encoding/json"
	"fmt"
)

type AppFile struct {
	Name    string `required:"true" json:"fileName"`
	Path    string `required:"true" json:"filePath"`
	Content string `required:"true" json:"fileContent"`
}

func AppFileFromString(text string) ([]AppFile, error) {
	files := []AppFile{}
	err := json.Unmarshal([]byte(text), &files)
	if err != nil {
		return nil, fmt.Errorf("Sorry, Couldn't parse OpenAI response: %s", err)
	}

	return files, nil
}
