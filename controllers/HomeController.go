package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util/log"
	"github.com/kataras/iris/mvc"
)

type HomeController struct {
	BaseController
}

func (this *HomeController) Get() mvc.Result {
	data := make(map[string]interface{})
	data["SITE_NAME"] = "BOOKZONE"
	if cates, err := new(models.Category).GetCates(-1, 1); err != nil {
		log.Errorf("fail to get cateroty")
	} else {
		data["Cates"] = cates
	}

	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if ok {
		data["Member"] = member
	} else {
		data["Member"] = models.Member{}
	}

	return mvc.View{
		Name: "home/list.html",
		Data: data,
	}
}