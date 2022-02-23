package webdav

import (
	"fmt"
	"golang.org/x/net/webdav"
	"net/http"
	"strings"
)

var fsList []*webdav.Handler

func Add(baseUrl string, dir string) {
	fs := &webdav.Handler{
		Prefix:     baseUrl,
		FileSystem: webdav.Dir(dir),
		LockSystem: webdav.NewMemLS(),
	}

	fsList = append(fsList, fs)
}

func Run() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		if !ok {
			writer.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		// 验证用户名/密码
		if username != "yan" || password != "world" {
			http.Error(writer, "WebDAV: need authorized!", http.StatusUnauthorized)
			return
		}

		for _, fs := range fsList {
			if strings.HasPrefix(request.RequestURI, fs.Prefix) {
				fs.ServeHTTP(writer, request)
				return
			}
		}
	})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
