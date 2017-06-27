package controllers

import (
	"ignite/models"
	"ignite/ss"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (router *MainRouter) ResetAccountHandler(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
		c.JSON(http.StatusOK, resp)
		return
	}

	user := new(models.User)
	user.PackageUsed = 0

	router.db.Id(uid).Cols("package_used").Update(user)
	resp := models.Response{Success: true, Message: "success"}
	c.JSON(http.StatusOK, resp)
}

func (router *MainRouter) DestroyAccountHandler(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
		c.JSON(http.StatusOK, resp)
		return
	}

	user := new(models.User)
	router.db.Id(uid).Get(user)

	//1. Destroy user's container
	err = ss.RemoveContainer(user.ServiceId)

	if err != nil {
		resp := models.Response{Success: false, Message: "终止用户容器失败!"}
		c.JSON(http.StatusOK, resp)
		return
	} else {
		//2. Delete user's account
		router.db.Id(uid).Delete(user)
	}

	resp := models.Response{Success: true, Message: "success"}
	c.JSON(http.StatusOK, resp)
}

func (router *MainRouter) StopServiceHandler(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
		c.JSON(http.StatusOK, resp)
		return
	}

	user := new(models.User)
	router.db.Id(uid).Get(user)

	//1. Stop user's container
	ss.StopContainer(user.ServiceId)

	//2. Update service status
	user.Status = 2
	router.db.Id(uid).Cols("status").Update(user)
	resp := models.Response{Success: true, Message: "success"}
	c.JSON(http.StatusOK, resp)
}

func (router *MainRouter) StartServiceHandler(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		resp := models.Response{Success: false, Message: "用户ID参数不正确"}
		c.JSON(http.StatusOK, resp)
		return
	}

	user := new(models.User)
	router.db.Id(uid).Get(user)

	//1. Stop user's container
	ss.StartContainer(user.ServiceId)

	//2. Update service status
	user.Status = 2
	router.db.Id(uid).Cols("status").Update(user)
	resp := models.Response{Success: true, Message: "success"}
	c.JSON(http.StatusOK, resp)
}
