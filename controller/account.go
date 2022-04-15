package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"teamt5interview/service"
	"teamt5interview/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Account struct {
	Username string
	Password string
}

type AccountController interface {
}

type IAccountController struct {
	engine      *gin.RouterGroup
	fileService service.FileService
	prefix      string
}

const fileSavePathDir = "file/account/"

func CreateAccountController(engine *gin.Engine) AccountController {
	group := engine.Group("/v1/account")
	c := &IAccountController{
		prefix:      "/v1/account",
		engine:      group,
		fileService: service.CreateFileService(),
	}
	c.Registry()
	return c
}

func (controller *IAccountController) Registry() {
	controller.engine.POST("/registry", controller.RegistryAccount)
	controller.engine.POST("/login", controller.loginAccount)
	controller.engine.GET("/logout", controller.logoutAccount)
}

func (controller *IAccountController) RegistryAccount(c *gin.Context) {
	var account Account
	err := c.Bind(&account)
	if err != nil {
		return
	}
	fmt.Printf("%v\n", account)
	fileName := ""
	if account.Username != "" && account.Password != "" {
		fileName = b64.StdEncoding.EncodeToString([]byte(account.Username))
		filepath := fileSavePathDir + fileName + ".json"
		if isExist, err := utils.Exists(filepath); isExist {
			goto Exist
		} else if err != nil {
			goto Error
		}
		jsonEntity, err := json.MarshalIndent(account, "", " ")
		if err != nil {
			goto Error
		}
		err = controller.fileService.Write(filepath, jsonEntity)
		if err != nil {
			goto Error
		}
	}
	c.JSON(200, "success")
	return
Error:
	c.JSON(400, "error parameter")
	return

Exist:
	c.JSON(http.StatusForbidden, "error create account")
	return
}

func (controller *IAccountController) loginAccount(c *gin.Context) {
	session := sessions.Default(c)
	var account Account
	err := c.Bind(&account)
	if err != nil {
		return
	}
	if session.Get(utils.LOGIN_STATUSKEY) == true {
		c.Status(http.StatusNoContent)
		c.Done()
		return
	} else {
		fileName := b64.StdEncoding.EncodeToString([]byte(account.Username))
		filepath := fileSavePathDir + fileName + ".json"
		bData, err := controller.fileService.Read(filepath)
		readAccount := &Account{}
		if err != nil {
			fmt.Println(err)
			goto Error
		}
		err = json.Unmarshal(bData, readAccount)
		if err != nil {
			fmt.Println(err)
			goto Error
		}
		if readAccount.Password == account.Password {
			session.Set(utils.LOGIN_STATUSKEY, true)
			session.Set(utils.LOGIN_USERNAMEKEY, account.Username)
			session.Save()
			c.JSON(200, "login success")
			return
		} else {
			goto Error
		}
	}
Error:
	c.JSON(400, "error parameter")
	return
}

func (controller *IAccountController) logoutAccount(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get(utils.LOGIN_STATUSKEY) == true {
		session.Delete(utils.LOGIN_STATUSKEY)
		session.Delete(utils.LOGIN_USERNAMEKEY)
		session.Save()
		c.Status(200)
		c.Done()
		return
	} else {
		goto Error
	}
Error:
	c.JSON(400, "error parameter")
	return
}
