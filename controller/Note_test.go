package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	service "teamt5interview/mock"
	"teamt5interview/utils"
	"testing"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	now := utils.MyTime(time.Now())
	body := &Note{
		Id:      "testId",
		Message: "testMessage",
		Time:    &now,
	}
	fileContext, _ := json.Marshal(body)
	fileService := service.CreateMockFileService(fileContext, "")
	controller := &INoteController{
		fileService: fileService,
	}
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(utils.LOGIN_STATUSKEY, true)
		session.Set(utils.LOGIN_USERNAMEKEY, "test")
		session.Save()
		c.Next()
	})
	r.GET("/getNote/fileId/:fileId", controller.GetNote)
	logoutw := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	logoutReq, _ := http.NewRequest("GET", "/getNote/fileId/testId", nil)
	r.ServeHTTP(logoutw, logoutReq)
	fmt.Println(logoutw.Code, logoutw.Body.String())
	assert.Equal(t, http.StatusOK, logoutw.Code)
}

func TestGetAllNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	now := utils.MyTime(time.Now())
	body := &Note{
		Id:      "testId",
		Message: "testMessage",
		Time:    &now,
	}
	fileContext, _ := json.Marshal(body)
	fileService := service.CreateMockFileService(fileContext, "../mock/note")
	controller := &INoteController{
		fileService: fileService,
	}
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(utils.LOGIN_STATUSKEY, true)
		session.Set(utils.LOGIN_USERNAMEKEY, "test")
		session.Save()
		c.Next()
	})
	r.GET("/getAllNote", controller.GetAllNote)
	logoutw := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	logoutReq, _ := http.NewRequest("GET", "/getAllNote", nil)
	r.ServeHTTP(logoutw, logoutReq)
	fmt.Println(logoutw.Code, logoutw.Body.String())
	assert.Equal(t, http.StatusOK, logoutw.Code)
}

func TestCreateNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	now := utils.MyTime(time.Now())
	body := &Note{
		Id:      "testId",
		Message: "testMessage",
		Time:    &now,
	}
	fileService := service.CreateMockFileService(nil, "../mock/note")
	controller := &INoteController{
		fileService: fileService,
	}
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(utils.LOGIN_STATUSKEY, true)
		session.Set(utils.LOGIN_USERNAMEKEY, "test")
		session.Save()
		c.Next()
	})
	r.POST("/createNote", controller.CreateNote)
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(body)
	req, _ := http.NewRequest("POST", "/createNote", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())
	bData, _ := fileService.Read("")
	json.Unmarshal(bData, &body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, body.Message, "testMessage")
}

func TestUpdateNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	now := utils.MyTime(time.Now())
	body := &Note{
		Id:      "testId",
		Message: "testMessage",
		Time:    &now,
	}
	fileContext, _ := json.Marshal(body)
	fileService := service.CreateMockFileService(fileContext, "../mock/note")
	controller := &INoteController{
		fileService: fileService,
	}
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(utils.LOGIN_STATUSKEY, true)
		session.Set(utils.LOGIN_USERNAMEKEY, "test")
		session.Save()
		c.Next()
	})

	updateBody := &Note{
		Message: "updateMessage",
		Time:    &now,
	}
	r.PUT("/updateNote/fileId/:fileId", controller.UpdateNote)
	w := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	payloadBuf := new(bytes.Buffer)
	json.NewEncoder(payloadBuf).Encode(updateBody)
	req, _ := http.NewRequest("PUT", "/updateNote/fileId/testId", payloadBuf)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	fmt.Println(w.Code, w.Body.String())
	bData, _ := fileService.Read("")
	json.Unmarshal(bData, &body)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, body.Message, "updateMessage")
}

func TestDeleteNote(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("testT5InterviewSession", store))
	now := utils.MyTime(time.Now())
	body := &Note{
		Id:      "testId",
		Message: "testMessage",
		Time:    &now,
	}
	fileContext, _ := json.Marshal(body)
	fileService := service.CreateMockFileService(fileContext, "")
	controller := &INoteController{
		fileService: fileService,
	}
	r.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set(utils.LOGIN_STATUSKEY, true)
		session.Set(utils.LOGIN_USERNAMEKEY, "test")
		session.Save()
		c.Next()
	})
	r.DELETE("/deleteNote/fileId/:fileId", controller.DeleteNote)
	logoutw := httptest.NewRecorder() // 取得 ResponseRecorder 物件
	logoutReq, _ := http.NewRequest("DELETE", "/deleteNote/fileId/testId", nil)
	r.ServeHTTP(logoutw, logoutReq)
	fmt.Println(logoutw.Code, logoutw.Body.String())
	context, _ := fileService.Read("")
	assert.Equal(t, http.StatusOK, logoutw.Code)
	assert.Equal(t, nil, context)
}
