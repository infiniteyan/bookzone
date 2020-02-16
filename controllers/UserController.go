package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util/log"
	"github.com/kataras/iris/mvc"
	"ziyoubiancheng/mbook/utils"
)

type UserController struct {
	BaseController
	UcenterMember 	models.Member
}

func (this *UserController) BeforeActivation(a mvc.BeforeActivation) {
	log.Infof("BookController BeforeActivation")
	a.Handle("GET","/{username:string}", "Index")
}

func (this *UserController) Index() mvc.Result {
	username := this.Ctx.Params().Get("username")
	log.Infof("UserController Index, username:%s", username)
	if username == "" {
		return mvc.Response{
			Code: 404,
		}
	}
	m, err := models.NewMember().GetByUsername(username)
	if err != nil || m == nil {
		return mvc.Response{
			Code: 404,
		}
	}
	if m.MemberId == 0 {
		return mvc.Response{
			Code: 404,
		}
	}

	dataMap := make(map[string]interface{})
	session := this.getSession()
	member, _ := session.Get(common.MemberSessionName).(models.Member)
	this.UcenterMember = *m
	dataMap["IsSelf"] = this.UcenterMember.MemberId == member.MemberId
	dataMap["User"] = this.UcenterMember
	dataMap["Tab"] = "share"

	dataMap["SITE_NAME"] = "BOOKZONE"
	dataMap["Member"] = member

	page, _ := this.Ctx.URLParamInt("page")
	pageSize := 10
	if page < 1 {
		page = 1
	}
	var books []*models.BookData
	var totalCount int
	books, totalCount, _ = models.NewBook().SelectPage(page, pageSize, this.UcenterMember.MemberId, 0)
	dataMap["Books"] = books

	if totalCount > 0 {
		html := utils.NewPaginations(common.RollPage, totalCount, pageSize, page, "/taafzdfxv", "")
		dataMap["PageHtml"] = html
	} else {
		dataMap["PageHtml"] = ""
	}
	dataMap["Total"] = totalCount

	return mvc.View{
		Name: "user/index.html",
		Data: dataMap,
	}
}