package utils

import (
	"fmt"
	"io"
	"os"
	"regexp"
)

func StripNewLines(str string) string {
	nlre := regexp.MustCompile(`\r?\n`)
	return nlre.ReplaceAllString(str, "")
}

func GetFileSize(filePath string) int64 {
	fileStats, err := os.Stat(filePath)
	if err != nil {
		panic(err)
	}

	return fileStats.Size()
}

func GetFileSeekPos(filePath string) int64 {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	seekPos, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}

	return seekPos
}

func ReadFileAsString(filePath string, stripNewLines bool, seekPos int64) (string, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return "", err
	}
	defer f.Close()
	content := ""
	buf := make([]byte, 5)

	if seekPos >= 0 {
		f.Seek(seekPos, 1)
	}
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
		if n > 0 {
			if stripNewLines {
				content += StripNewLines(string(buf[:n]))
			} else {
				content += string(buf[:n])
			}

		}
	}

	return content, nil
}
