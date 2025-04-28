package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SearchFonts() []string {
	fd := filepath.Join(workingDir(), "fonts")

	var fontFiles []string

	// Walk through the directory
	err := filepath.Walk(fd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if file is .ttf or .otf
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".ttf") {
			fontFiles = append(fontFiles, info.Name()) // Add file name
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the path:", err)
	}

	return fontFiles
}

func FontExists(font string) bool {
	for _, f := range SearchFonts() {
		if strings.EqualFold(f, font) {
			return true
		}
	}
	return false
}
