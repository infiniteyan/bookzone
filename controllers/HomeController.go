package controllers

import "github.com/kataras/iris/mvc"

type HomeController struct {}

func (this *HomeController) Get() mvc.Result {
	return mvc.View{
		Name: "home/list.html",
	}
}