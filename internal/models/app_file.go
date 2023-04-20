package models

import "encoding/json"

type AppFile struct {
	Name    string
	Path    string
	Content string
}

func AppFileFromString(text string) ([]AppFile, error) {
	files := []AppFile{}
	err := json.Unmarshal([]byte(text), &files)

	return files, err
}
