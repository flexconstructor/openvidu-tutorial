package main

import (
	"github.com/flexconstructor/openvidu-tutorial/route"
)

func main() {
	router := route.InitRouter()
	router.Run()
}
