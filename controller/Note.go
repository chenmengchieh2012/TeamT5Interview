package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"teamt5interview/service"
	"teamt5interview/utils"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 項目內容、時間、日期
type Note struct {
	Id      string
	Message string
	Time    *utils.MyTime
}

type NoteController interface {
}

type INoteController struct {
	engine      *gin.RouterGroup
	fileService service.FileService
	prefix      string
}

const notePathDir = "file/note/"

func CreateNoteController(engine *gin.Engine) NoteController {
	group := engine.Group("/v1/note")
	c := &INoteController{
		prefix:      "/v1/note",
		engine:      group,
		fileService: service.CreateFileService(),
	}
	c.Registry()
	return c
}

func (controller *INoteController) Registry() {
	controller.engine.Use(controller.Autherization())
	controller.engine.GET("/fileId/:fileId", controller.GetNote)
	controller.engine.DELETE("/fileId/:fileId", controller.DeleteNote)
	controller.engine.GET("/", controller.GetAllNote)
	controller.engine.POST("/", controller.CreateNote)
	controller.engine.PUT("/fileId/:fileId", controller.UpdateNote)
}

func (controller *INoteController) Autherization() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(utils.LOGIN_STATUSKEY) != true {
			c.Abort()
			c.Status(http.StatusUnauthorized)
			return
		}
		username := session.Get(utils.LOGIN_USERNAMEKEY)
		if username == nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, "not found in session")
			return
		}
		userDirName := b64.StdEncoding.EncodeToString([]byte(username.(string)))
		userDirPath := notePathDir + userDirName
		if _, err := controller.fileService.MakeDir(userDirPath); err != nil {
			c.Abort()
			c.JSON(http.StatusInternalServerError, "create filenote path error")
			return
		}
		c.Next()
	}
}

func (controller *INoteController) GetNote(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(utils.LOGIN_USERNAMEKEY)
	userDirName := b64.StdEncoding.EncodeToString([]byte(username.(string)))
	userDirPath := notePathDir + userDirName
	if fileId, ok := c.Params.Get("fileId"); !ok {
		goto Error
	} else {
		noteFilePath := userDirPath + "/" + fileId
		bData, err := controller.fileService.Read(noteFilePath)
		if err != nil {
			goto Error
		}
		note := &Note{}
		err = json.Unmarshal(bData, note)
		if err != nil {
			goto Error
		}
		c.JSON(http.StatusOK, note)
		return
	}
Error:
	c.JSON(400, "error parameter")
	return
}

func (controller *INoteController) GetAllNote(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get(utils.LOGIN_USERNAMEKEY)
	userDirName := b64.StdEncoding.EncodeToString([]byte(username.(string)))
	userDirPath := notePathDir + userDirName
	subitems, _ := controller.fileService.ReadDir(userDirPath)
	retData := make([]*Note, 0)
	for _, subitem := range subitems {
		bData, err := controller.fileService.Read(userDirPath + "/" + subitem.Name())
		if err != nil {
			goto Error
		}
		note := &Note{}
		err = json.Unmarshal(bData, note)
		if err != nil {
			goto Error
		}
		retData = append(retData, note)
	}
	c.JSON(http.StatusOK, retData)
	return
Error:
	c.JSON(400, "error parameter")
	return
}

func (controller *INoteController) UpdateNote(c *gin.Context) {
	var updatenote Note
	err := c.Bind(&updatenote)
	if err != nil {
		fmt.Println("updateNote", err)
		return
	}
	session := sessions.Default(c)
	username := session.Get(utils.LOGIN_USERNAMEKEY)
	userDirName := b64.StdEncoding.EncodeToString([]byte(username.(string)))
	userDirPath := notePathDir + userDirName
	if fileId, ok := c.Params.Get("fileId"); !ok {
		goto Error
	} else {
		noteFilePath := userDirPath + "/" + fileId
		bData, err := controller.fileService.Read(noteFilePath)
		if err != nil {
			fmt.Println("err", err)
			goto Error
		}
		currentnote := &Note{}
		err = json.Unmarshal(bData, currentnote)
		if err != nil {
			fmt.Println("err", err)
			goto Error
		}
		updatenote.Id = currentnote.Id
		if updatenote.Message == "" {
			fmt.Println("err", err)
			goto Empty
		}
		if updatenote.Time == nil {
			updatenote.Time = currentnote.Time
		}
		jsonEntity, err := json.MarshalIndent(updatenote, "", " ")
		err = controller.fileService.Write(noteFilePath, jsonEntity)
		if err != nil {
			goto Error
		}
		c.JSON(http.StatusOK, "update OK")
		return
	}
Empty:
	c.JSON(http.StatusNotAcceptable, "no message")
	return
Error:
	c.JSON(400, "error parameter")
	return
}

func (controller *INoteController) CreateNote(c *gin.Context) {
	var note Note
	err := c.Bind(&note)
	if err != nil {
		return
	}

	session := sessions.Default(c)
	username := session.Get(utils.LOGIN_USERNAMEKEY)
	userDirName := b64.StdEncoding.EncodeToString([]byte(username.(string)))
	userDirPath := notePathDir + userDirName
	noteFileName := uuid.New().String()
	noteFilePath := userDirPath + "/" + noteFileName
	note.Id = noteFileName
	if note.Time == nil {
		now := utils.MyTime(time.Now())
		note.Time = &now
	}
	jsonEntity, err := json.MarshalIndent(note, "", " ")
	if err != nil {
		goto Error
	}
	err = controller.fileService.Write(noteFilePath, jsonEntity)
	if err != nil {
		goto Error
	}
	c.JSON(200, "success")
	return
Error:
	c.JSON(400, "error parameter")
	return
}

func (controller *INoteController) DeleteNote(c *gin.Context) {
	var note Note
	err := c.Bind(&note)
	if err != nil {
		return
	}
	session := sessions.Default(c)
	username := session.Get(utils.LOGIN_USERNAMEKEY)
	userDirName := b64.StdEncoding.EncodeToString([]byte(username.(string)))
	userDirPath := notePathDir + userDirName
	if fileId, ok := c.Params.Get("fileId"); !ok {
		goto Error
	} else {
		noteFilePath := userDirPath + "/" + fileId
		err := controller.fileService.Delete(noteFilePath)
		if err != nil {
			goto Error
		}
	}
	c.JSON(200, "success")
	return
Error:
	c.JSON(400, "error parameter")
	return
}
