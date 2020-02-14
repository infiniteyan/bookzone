package controllers

import (
	"bookzone/common"
	"bookzone/util/log"
	"github.com/kataras/iris/mvc"
)

type BookController struct {
	BaseController
}

func (this *BookController) BeforeActivation(a mvc.BeforeActivation) {
	log.Infof("BookController BeforeActivation")
	a.Handle("POST", "/", "Comment")
	a.Handle("GET", "/{id:int}", "Score")
}

func (this *BookController) Comment() {
	log.Infof("BookController Post")
	content := this.Ctx.Params().Get("content")
	log.Infof("book comment:%s", content)
	if this.Member.MemberId == 0 {
		this.JsonResult(common.HttpCodeErrorLoginFirst, "评论失败，请先登录再操作")
		return
	}
}

func (this *BookController) Score() {
	log.Infof("BookController Score")
	bookId, err := this.Ctx.Params().GetInt("id")
	if err != nil {
		this.JsonResult(common.HttpCodeErrorParameter, "请求参数错误")
		return
	}

	score := this.Ctx.URLParam("score")
	if uid := this.Member.MemberId; uid > 0 {
		//if err := new(models.Score).AddScore(uid, bookId, score); err != nil {
		//	c.JsonResult(1, err.Error())
		//}
		//c.JsonResult(0, "感谢您给当前文档打分")
	}
	log.Infof("book id:%d  score:%s", bookId, score)
	this.JsonResult(common.HttpCodeErrorLoginFirst, "给文档打分失败，请先登录再操作")
}