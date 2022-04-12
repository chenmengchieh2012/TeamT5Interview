package main

import (
	"teamt5interview/controller"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("t5InterviewSession", store))
	controller.CreateAccountController(router)
	controller.CreateNoteController(router)
	router.LoadHTMLGlob("public/*")
	router.Run(":8080")
}
