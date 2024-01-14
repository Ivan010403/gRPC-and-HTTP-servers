package imagefiles

import (
	"fmt"
	"io"
	"os"
)

type ImageFile struct {
	Name     string
	Filetype string
}

// TODO: validation of data (name and TYPE!)
func (i *ImageFile) Read() ([]byte, error) {
	var file *os.File
	path := i.Name + "." + i.Filetype

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("file or directory doesn't exist")
	} else {
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (i *ImageFile) Write(data []byte) (string, error) {
	var file *os.File
	path := i.Name + "." + i.Filetype

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err = os.Create(i.Name + "." + i.Filetype)
		if err != nil {
			return "", err
		}
		defer file.Close()
	} else {
		file, err = os.Open(path)
		if err != nil {
			return "", err
		}
		defer file.Close()
	}

	_, err := file.Write(data)
	if err != nil {
		return "", err
	}
	return i.Name + "." + i.Filetype, nil
}
