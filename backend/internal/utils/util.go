package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	envFile, err := FindFile(".env")
	if err != nil {
		return fmt.Errorf("error finding .env file")
	}

	err = godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("error loading .env file")
	}

	return nil
}

// FindFolder searches for a folder by name starting from the Go project root.
func FindFolder(folderName string) (string, error) {
	rootDir, err := FindProjectRoot()
	if err != nil {
		return "", err
	}

	var foundPath string

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == folderName {
			foundPath = path
			return filepath.SkipDir // Stop searching after the first match.
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("folder %s not found in %s", folderName, rootDir)
	}

	return foundPath, nil
}

// FindFile searches for a file by name starting from the Go project root.
func FindFile(fileName string) (string, error) {
	rootDir, err := FindProjectRoot()
	if err != nil {
		return "", err
	}

	var foundPath string

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == fileName {
			foundPath = path
			return filepath.SkipDir // Stop searching after the first match
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("file %s not found in %s", fileName, rootDir)
	}

	return foundPath, nil
}

// FindFilesWithExtension searches for all files with a specified extension starting from the Go project root.
func FindFilesWithExtension(extension string) ([]string, error) {
	rootDir, err := FindProjectRoot()
	if err != nil {
		return nil, err
	}

	var filesWithExtension []string

	err = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(info.Name()) == extension {
			filesWithExtension = append(filesWithExtension, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filesWithExtension, nil
}

// FindProjectRoot attempts to find the root directory of the Go project.
// It first checks for a go.mod file, and if not found, returns the executable's directory.
func FindProjectRoot() (string, error) {
	// Start from the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up the directory tree until we find a go.mod file or reach the root
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil // Found go.mod, return this directory
		}

		// Reached the root
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}

	// If go.mod is not found, use the directory of the executable
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	execDir := filepath.Dir(execPath)

	return execDir, nil
}

func Filter[T any](ss []T, test func(T) bool) (ret []T) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
