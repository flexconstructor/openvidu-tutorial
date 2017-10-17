package main

import (
	"github.com/flexconstructor/openvidu-tutorial/route"
	"github.com/flexconstructor/openvidu-tutorial/service"
)

// Is a OpenViDu GoLang tutorial.
func main() {
	router := route.InitRouter(&service.Client{
		OpenViDuURL: "https://openvidu-server-kms:8443",
		Login:       "OPENVIDUAPP",
		Password:    "MY_SECRET",
	})
	router.LoadHTMLGlob("/resources/templates/*.tmpl")
	router.Static("/images", "resources/static/images")
	router.StaticFile("/style.css", "resources/static//style.css")
	router.StaticFile("/openvidu-browser-1.1.0.js",
		"resources/static/openvidu-browser-1.1.0.js")
	router.Run()
}
