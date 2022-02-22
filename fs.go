package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type fileStruct struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	IsDir   bool      `json:"is_dir"`
}

//const baseDirectory = "/storage"

//const baseDirectory = "./storage"

//const baseDirectory = "C:/Users/Yan_X/Workspace/test"

//const baseDirectory = "/mnt/c/Users/Yan_X/Workspace/test"

func modifiedTime(name string) time.Time {
	stat, err := os.Stat(name)
	if err != nil {
		fmt.Println(err)
	}
	return stat.ModTime()
}

func readFile(name string) []byte {
	file, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return file
}

func filePath(name string) string {
	return name
}

func readFileBytes(name string, bytes int) []byte {
	buf := make([]byte, bytes)
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	_, err = file.Read(buf)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buf
}

func fileStream(name string) *os.File {
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return file
}

func contentType(name string) string {
	buf := make([]byte, 512)
	file, err := os.Open(name)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	_, err = file.Read(buf)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return http.DetectContentType(buf)
}

func Files(name string) ([]fileStruct, []fileStruct, error) {
	list, err := ioutil.ReadDir(name)
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, fmt.Errorf("Not Found\n")
	}
	var files []fileStruct
	var folders []fileStruct
	for _, item := range list {
		//fmt.Println(item.Name(), item.IsDir())
		if item.IsDir() {
			folders = append(folders, fileStruct{
				Name:    item.Name(),
				Size:    item.Size(),
				ModTime: item.ModTime(),
				IsDir:   item.IsDir(),
			})
		} else {
			files = append(files, fileStruct{
				Name:    item.Name(),
				Size:    item.Size(),
				ModTime: item.ModTime(),
				IsDir:   item.IsDir(),
			})
		}
	}
	return files, folders, nil
}
