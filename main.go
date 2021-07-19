package main

import (
	"SimBlock/config"
	"net/http"

	"SimBlock/simulations"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(config.CORS)
	r.StaticFile("/", config.MAIN_PAGE)
	r.StaticFS("/static", http.Dir(config.STATIC_PATH))
	simulations.RegisterSimnet(r)
	r.Run(":8000")
}
