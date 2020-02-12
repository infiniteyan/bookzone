package controllers

import (
	"bookzone/models"
	"github.com/kataras/iris/mvc"
	"bookzone/util/log"
)

type HomeController struct {}

func (this *HomeController) Get() mvc.Result {
	data := make(map[string]interface{})
	data["SITE_NAME"] = "BOOKZONE"
	if cates, err := new(models.Category).GetCates(-1, 1); err != nil {
		log.Errorf("fail to get cateroty")
	} else {
		data["Cates"] = cates
	}
	return mvc.View{
		Name: "home/list.html",
		Data: data,
	}
}