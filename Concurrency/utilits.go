package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

const dirName = "test_directory"

func generateRunesArr(cnt int) []rune {
	runes := make([]rune, cnt)
	for i := range cnt {
		runes[i] = rune(rand.Intn(33) + 'а')
	}
	return runes
}

func generateBytesArr(cnt int) []byte {
	bytes := make([]byte, cnt)

	chunkSize := rand.Intn(20) + 3

	for i := range cnt {
		if (i % chunkSize) == 0 {
			bytes[i] = byte('\n')
		} else {
			bytes[i] = byte(rand.Intn('z'-'0') + '0')
		}
	}
	return bytes
}

func writeInFile(path string) error {
	cnt := rand.Intn(4000) + 1000
	bytes := generateBytesArr(cnt)

	err := os.WriteFile(path, bytes, 0644)
	return err
}

func checkDirIsExist(dirName string) error {
	if _, err := os.Stat(dirName); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.Mkdir(dirName, 0755); err != nil {
				return err
			}
		}
	}
	return nil
}

func createFiles(num int) ([]string, error) {
	var paths []string
	for i := range num {
		fileName := fmt.Sprintf("file_%3d.txt", i)
		path := filepath.Join(dirName, fileName)

		if err := writeInFile(path); err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func prepareFilesForProcessing(filesAmount int) ([]string, error) {
	if err := checkDirIsExist(dirName); err != nil {
		return nil, err
	}

	paths, err := createFiles(filesAmount)
	if err != nil {
		return nil, err
	}
	return paths, nil
}
