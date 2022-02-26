package main

import (
	"Gallery/webdav"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

const folderBaseUrl = "/folder"
const fileBaseUrl = "/file"
const previewBaseUrl = "/preview"
const detailsBaseUrl = "/details"

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(cors.Default())
	//r.Use(Logger())

	r.LoadHTMLGlob("templates/gin/*")
	r.StaticFile("/favicon.ico", "public/favicon.ico")
	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	sites := readConfig()
	for baseUrl, baseDirectory := range sites {
		hGroup := handlersGroup{GroupName: baseUrl, BaseDirectory: baseDirectory}
		r.Static(path.Join(baseUrl, "/assets"), "public/assets")
		r.StaticFile(path.Join(baseUrl, "/404.html"), "public/404.html")
		group := r.Group(baseUrl)
		{
			group.GET(path.Join(folderBaseUrl, "*action"), hGroup.folderHandler)
			group.GET(path.Join(fileBaseUrl, "*action"), hGroup.fileHandler)
			group.GET(path.Join(previewBaseUrl, "*action"), hGroup.previewHandler)
			group.GET(path.Join(detailsBaseUrl, "*action"), hGroup.detailsHandler)
			group.GET("/", func(context *gin.Context) {
				context.Redirect(http.StatusMovedPermanently, path.Join(baseUrl, folderBaseUrl))
			})
		}

		webdav.Add(baseUrl, baseDirectory)
	}

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080

	go webdav.Run()
	err := r.Run(":80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
