package main

import (
	"net/http"
	"teamt5interview/controller"
	"teamt5interview/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"

	"github.com/gin-gonic/gin"
)

//http://www.manongjc.com/detail/29-skskxedphvkaucv.html
func main() {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("t5InterviewSession", store))
	controller.CreateAccountController(router)
	controller.CreateNoteController(router)
	router.LoadHTMLGlob("templates/**/*")
	router.GET("/", GetIndex)
	router.Static("/static", "static")
	router.Run(":8080")
}

func GetIndex(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get(utils.LOGIN_STATUSKEY) == true {
		ctx.HTML(http.StatusOK, "main.html", nil)
		return
	} else {
		ctx.HTML(http.StatusOK, "login.html", nil)
	}
}
