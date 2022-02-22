package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(time.Now())
		fmt.Println(c.ClientIP())
		fmt.Println(c.Request.URL)
		fmt.Println(c.Request.UserAgent())
		fmt.Println(c.Request.Referer())
	}
}
