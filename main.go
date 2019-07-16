package main

import (
	"github.com/gin-gonic/gin"
	"ws-server/app"
)

func main() {
	app.StartServer(gin.DebugMode,"/ws")
}
