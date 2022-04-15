package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	service "teamt5interview/mock"
	"teamt5interview/utils"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestRegistryAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Body: io.NopCloser(strings.NewReader(`{
			"username": "testRegistry",
			"password": "testpassword",
		}`)),
	}
	controller := &IAccountController{}
	controller.RegistryAccount(c)
	assert.Equal(t, c.Writer.Status(), 200)
}

func TestLoginAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	body := &Account{
		Username: "test",
		Password: "test",
	}
	fileContext, _ := json.Marshal(body)
	fileService := service.CreateMockFileService(fileContext, "")
	controller := &IAccountController{
		fileService: fileService,
	}
	r.POST("/login", controller.loginAccount)
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	req, _ := http.NewRequest("POST", "/login", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "\"login success\"", w.Body.String())
}

func TestLogoutAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	controller := &IAccountController{}
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(utils.LOGIN_STATUSKEY, true)
		session.Save()
		c.Next()
	})
	r.GET("/logout", controller.logoutAccount)
	logoutw := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	logoutReq, _ := http.NewRequest("GET", "/logout", nil)
	r.ServeHTTP(logoutw, logoutReq)
	fmt.Println(logoutw.Code, logoutw.Body.String())
	assert.Equal(t, http.StatusOK, logoutw.Code)
}
