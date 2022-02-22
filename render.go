package main

import (
	"bytes"
	"fmt"
	"github.com/docker/go-units"
	"github.com/h2non/filetype"
	"html/template"
	"path"
)

const previewHeight = 330

func pictureWaterfall(fullPath string, previewPath string, files []fileStruct) string {
	tp, err := template.ParseFiles("templates/img.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	buf := new(bytes.Buffer)
	for _, file := range files {
		data := map[string]interface{}{
			"alt":  file.Name,
			"src":  path.Join(previewPath, file.Name),
			"href": path.Join(fullPath, file.Name),
		}

		err = tp.Execute(buf, data)

		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return buf.String()
}

func folderWaterfall(baseUrl string, fullPath string, folders []fileStruct) string {
	tp, err := template.ParseFiles("templates/folder-img.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	buf := new(bytes.Buffer)
	for _, folder := range folders {
		data := map[string]interface{}{
			"alt":  folder.Name,
			"src":  path.Join(baseUrl, "/assets/folder.webp"),
			"href": path.Join(fullPath, folder.Name+"/"),
		}

		err = tp.Execute(buf, data)

		if err != nil {
			fmt.Println(err.Error())
		}
	}
	return buf.String()
}

func detailsList(groupName, action string, fileDirectory string, folders []fileStruct, files []fileStruct) string {
	tp, err := template.ParseFiles("templates/details-item.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	buf := new(bytes.Buffer)
	for _, item := range append(folders, files...) {
		var fileType, href string
		if item.IsDir {
			fileType = "文件夹"
			href = path.Join(groupName, detailsBaseUrl, action, item.Name) + "/"
		} else {
			buf := readFileBytes(path.Join(fileDirectory, item.Name), 8192)
			if filetype.IsVideo(buf) {
				fileType = "视频"
			} else if filetype.IsDocument(buf) {
				fileType = "文档"
			} else if filetype.IsApplication(buf) {
				fileType = "应用程序"
			} else if filetype.IsArchive(buf) {
				fileType = "压缩包"
			} else if filetype.IsAudio(buf) {
				fileType = "音频"
			} else if filetype.IsFont(buf) {
				fileType = "字体"
			} else if filetype.IsImage(buf) {
				fileType = "图像"
			} else {
				fileType = "未知"
			}

			href = path.Join(groupName, fileBaseUrl, action, item.Name)
		}

		data := map[string]interface{}{
			"name":    item.Name,
			"modTime": item.ModTime.Format("2006-01-02 15:04:05"),
			"type":    fileType,
			"size":    units.HumanSize(float64(item.Size)),
			"href":    href,
		}

		err = tp.Execute(buf, data)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return buf.String()
}
