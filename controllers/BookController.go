package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util/log"
	"github.com/kataras/iris/mvc"
	"strconv"
)

type BookController struct {
	BaseController
}

func (this *BookController) BeforeActivation(a mvc.BeforeActivation) {
	log.Infof("BookController BeforeActivation")
	a.Handle("POST", "/comment/{id:int}", "Comment")
	a.Handle("GET", "/{id:int}", "Score")
}

func (this *BookController) Comment() {
	log.Infof("BookController Post")
	content := this.Ctx.FormValue("content")
	bookId, err := this.Ctx.Params().GetInt("id")
	if err != nil {
		this.JsonResult(common.HttpCodeErrorParameter, "请求参数错误")
		return
	}

	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if !ok {
		this.JsonResult(common.HttpCodeErrorLoginFirst, "评论失败，请先登录再操作")
		return
	}
	if member.MemberId == 0 {
		this.JsonResult(common.HttpCodeErrorLoginFirst, "评论失败，请先登录再操作")
		return
	}

	log.Infof("BookController, comment:%s", content)

	if bookId > 0 {
		if err := models.NewComments().AddComments(member.MemberId, bookId, content); err != nil {
			this.JsonResult(common.HttpCodeErrorDatabase, "评论失败")
		}
		this.JsonResult(common.HttpCodeSuccess, "评论成功")
	} else {
		this.JsonResult(common.HttpCodeErrorParameter, "文档图书不存在")
	}
}

func (this *BookController) Score() {
	log.Infof("BookController Score")
	bookId, err := this.Ctx.Params().GetInt("id")
	if err != nil {
		this.JsonResult(common.HttpCodeErrorParameter, "请求参数错误")
		return
	}

	scoreStr := this.Ctx.URLParam("score")
	score, _ :=strconv.Atoi(scoreStr)
	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if !ok {
		this.JsonResult(common.HttpCodeErrorLoginFirst, "评论失败，请先登录再操作")
		return
	}
	if member.MemberId == 0 {
		this.JsonResult(common.HttpCodeErrorLoginFirst, "评论失败，请先登录再操作")
		return
	}
	log.Infof("BookController, book id:%d score:%d", bookId, score)

	if err := models.NewScore().AddScore(member.MemberId, bookId, score); err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorDatabase, "您已经为文档评分")
		return
	}

	this.JsonResult(common.HttpCodeSuccess, "感谢您为当前文档打分")
}