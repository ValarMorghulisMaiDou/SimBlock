package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

var API_PREFIX = "/api"
var SIMNET_PREFIX = API_PREFIX + "/simnet"

var CORS = cors.New(cors.Config{
	//准许跨域请求网站,多个使用,分开,限制使用*
	AllowOrigins: []string{"*"},
	//准许使用的请求方式
	AllowMethods: []string{"PUT", "PATCH", "POST", "GET", "DELETE"},
	//准许使用的请求表头
	AllowHeaders: []string{"Origin", "Authorization", "Content-Type"},
	//显示的请求表头
	ExposeHeaders: []string{"Content-Type"},
	//凭证共享,确定共享
	AllowCredentials: true,
	//容许跨域的原点网站
	AllowOriginFunc: func(origin string) bool {
		return true
	},
	//超时时间设定
	MaxAge: 24 * time.Hour,
})

var STATIC_BASE_PATH = "./SimBlock/dist"
var MAIN_PAGE = STATIC_BASE_PATH + "/index.html"
var STATIC_PATH = STATIC_BASE_PATH + "/static"
