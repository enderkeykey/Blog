package main

import "github.com/gin-gonic/gin"

func main() {
	DataInit() //sqlite3 init
	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.POST("/upload", AuthMiddleWare(), Photo)
	r.POST("/comment", Comment)
	r.POST("/response", Response)
	r.POST("/showComment", ShowComment)
	r.Run(":8008")
}
