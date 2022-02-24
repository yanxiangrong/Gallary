package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	"golang.org/x/sync/semaphore"
	"html/template"
	"net/http"
	"path"
	"runtime/debug"
	"time"
)

const (
	maxThread = 4
)

var (
	sem *semaphore.Weighted
)

func init() {
	sem = semaphore.NewWeighted(maxThread)
}

type handlersGroup struct {
	GroupName     string
	BaseDirectory string
}

func (group handlersGroup) folderHandler(c *gin.Context) {
	action := c.Param("action")

	if action[len(action)-1] != '/' {
		c.Redirect(http.StatusMovedPermanently, path.Join(group.GroupName, folderBaseUrl, action)+"/")
	}
	directory := path.Join(group.BaseDirectory, action)

	files, folders, err := Files(directory)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	folderWaterfall := folderWaterfall(group.GroupName, path.Join(group.GroupName, folderBaseUrl, action), folders)
	pictureWaterfall := pictureWaterfall(
		path.Join(group.GroupName, fileBaseUrl, action), path.Join(group.GroupName, previewBaseUrl, action),
		files)

	c.HTML(http.StatusOK, "waterfall.html", gin.H{
		"details":          path.Join(group.GroupName, detailsBaseUrl, action) + "/",
		"title":            path.Join(group.GroupName, action) + "/",
		"folderWaterfall":  template.HTML(folderWaterfall),
		"pictureWaterfall": template.HTML(pictureWaterfall),
		"num":              len(files) + len(folders),
		"baseUrl":          group.GroupName,
	})
}

func (group handlersGroup) previewHandler(c *gin.Context) {
	action := c.Param("action")
	if action[len(action)-1] == '/' {
		c.String(http.StatusNotFound, "Not Found")
		return
	}
	directory := path.Join(group.BaseDirectory, action)

	file := readFileBytes(directory, 8192)

	fileIsImage := filetype.IsImage(file)
	fileIsVideo := filetype.IsVideo(file)

	if !(fileIsImage || fileIsVideo) {
		c.Redirect(http.StatusFound, path.Join(group.GroupName, "/assets/file.webp"))
		return
	}

	fileModifiedTime := modifiedTime(directory)

	c.Header("Last-Modified", fileModifiedTime.Format(http.TimeFormat))
	str := c.GetHeader("If-Modified-Since")
	if str != "" {
		sinceTime, err := time.Parse(http.TimeFormat, str)
		if err != nil {
			fmt.Println(err)
		} else {
			if sinceTime.Unix() > fileModifiedTime.Unix() {
				c.Status(http.StatusNotModified)
				return
			}
		}
	}

	err := sem.Acquire(c, 1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer func() {
		debug.FreeOSMemory()
		sem.Release(1)
	}()

	switch http.DetectContentType(file) {
	case "image/gif":
		c.Data(http.StatusOK, "image/gif", gifPreview(directory))
		return
	}
	if fileIsImage {
		c.Data(http.StatusOK, "image/jpeg", imagePreview(directory))
		return
	}
	if fileIsVideo {
		c.Data(http.StatusOK, "image/gif", videoPreview(directory))
		return
	}
	c.Redirect(http.StatusFound, path.Join(group.GroupName, "/assets/file.webp"))
}

func (group handlersGroup) fileHandler(c *gin.Context) {
	action := c.Param("action")
	if action[len(action)-1] == '/' {
		c.Redirect(http.StatusMovedPermanently, path.Join(group.GroupName, fileBaseUrl, action))
	}
	directory := path.Join(group.BaseDirectory, action)

	fileModifiedTime := modifiedTime(directory)

	c.Header("Last-Modified", fileModifiedTime.Format(http.TimeFormat))

	str := c.GetHeader("If-Modified-Since")
	if str != "" {
		sinceTime, err := time.Parse(http.TimeFormat, str)
		if err != nil {
			fmt.Println(err)
		} else {
			if sinceTime.Unix() > fileModifiedTime.Unix() {
				c.Status(http.StatusNotModified)
				return
			}
		}
	}

	c.File(filePath(directory))
}

func (group handlersGroup) detailsHandler(c *gin.Context) {
	action := c.Param("action")

	if action[len(action)-1] != '/' {
		c.Redirect(http.StatusMovedPermanently, path.Join(group.GroupName, folderBaseUrl, action)+"/")
	}
	directory := path.Join(group.BaseDirectory, action)

	files, folders, err := Files(directory)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	table := detailsList(group.GroupName, action, path.Join(group.BaseDirectory, action), folders, files)

	c.HTML(http.StatusOK, "details.html", gin.H{
		"waterfall": path.Join(group.GroupName, folderBaseUrl, action) + "/",
		"title":     path.Join(group.GroupName, action) + "/",
		"tbody":     template.HTML(table),
		"num":       len(files) + len(folders),
		"baseUrl":   group.GroupName,
	})
}
