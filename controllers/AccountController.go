package controllers

import "github.com/kataras/iris/mvc"

type AccountController struct {}

func (this *AccountController) GetLogin() mvc.Result {
	return mvc.View{
		Name: "account/login.html",
		Data:map[string]interface{} {
			"SITE_NAME":"BOOKZONE",
		},
	}
}

func (this *AccountController) GetRegist() mvc.Result {
	return mvc.View{
		Name: "account/register.html",
	}
}

func (this *AccountController) Logout() mvc.Result {
	return mvc.Response{
		Content: []byte("logout"),
	}
}

func (this *AccountController) PostRegist() mvc.Result {
	return mvc.Response{
		Content: []byte("postregist"),
	}
}