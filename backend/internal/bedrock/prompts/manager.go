package prompts

import (
	"embed"
	"io/fs"
)

//go:embed extractor.txt
var promptFiles embed.FS

func LoadPrompt(filename string) (string, error) {
	content, err := fs.ReadFile(promptFiles, filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
