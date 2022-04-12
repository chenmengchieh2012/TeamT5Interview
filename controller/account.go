package controller

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
	engine *gin.Engine
	prefix string
}

const fileSavePathDir = "file/account/"

func CreateAccountController(engine *gin.Engine) AccountController {
	c := &IAccountController{
		prefix: "/v1/account",
		engine: engine,
	}
	c.Registry()
	return c
}

func (controller *IAccountController) Registry() {
	controller.engine.POST(controller.prefix+"/registry", controller.RegistryAccount)
	controller.engine.POST(controller.prefix+"/login", controller.loginAccount)
	controller.engine.POST(controller.prefix+"/logout", controller.logoutAccount)
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
		err = ioutil.WriteFile(filepath, jsonEntity, 0644)
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
		bData, err := ioutil.ReadFile(filepath)
		readAccount := &Account{}
		if err != nil {
			goto Error
		}
		err = json.Unmarshal(bData, readAccount)
		if err != nil {
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
	var account Account
	err := c.Bind(&account)
	if err != nil {
		return
	}
	if session.Get(account.Username) == true {
		session.Delete(utils.LOGIN_STATUSKEY)
		session.Delete(utils.LOGIN_USERNAMEKEY)
		session.Save()
		return
	} else {
		goto Error
	}
Error:
	c.JSON(400, "error parameter")
	return
}
