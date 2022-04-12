package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
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
	engine *gin.Engine
	prefix string
}

const notePathDir = "file/note/"

func CreateNoteController(engine *gin.Engine) AccountController {
	c := &INoteController{
		prefix: "/v1/note",
		engine: engine,
	}
	c.Registry()
	return c
}

func (controller *INoteController) Registry() {
	controller.engine.Use(controller.Autherization())
	controller.engine.GET(controller.prefix+"/fileId/:fileId", controller.GetNote)
	controller.engine.GET(controller.prefix+"/", controller.GetAllNote)
	controller.engine.POST(controller.prefix+"/", controller.CreateNote)
	controller.engine.PUT(controller.prefix+"/fileId/:fileId", controller.UpdateNote)
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
		if _, err := utils.MakeDir(userDirPath); err != nil {
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
		bData, err := ioutil.ReadFile(noteFilePath)
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
	subitems, _ := ioutil.ReadDir(userDirPath)
	retData := make([]*Note, 0)
	for _, subitem := range subitems {
		bData, err := ioutil.ReadFile(userDirPath + "/" + subitem.Name())
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
		bData, err := ioutil.ReadFile(noteFilePath)
		if err != nil {
			goto Error
		}
		currentnote := &Note{}
		err = json.Unmarshal(bData, currentnote)
		if err != nil {
			goto Error
		}
		updatenote.Id = currentnote.Id
		if updatenote.Message == "" {
			goto Empty
		}
		if updatenote.Time == nil {
			updatenote.Time = currentnote.Time
		}
		jsonEntity, err := json.MarshalIndent(updatenote, "", " ")
		err = ioutil.WriteFile(noteFilePath, jsonEntity, 0644)
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
	err = ioutil.WriteFile(noteFilePath, jsonEntity, 0644)
	if err != nil {
		goto Error
	}
	c.JSON(200, "success")
	return
Error:
	c.JSON(400, "error parameter")
	return
}
