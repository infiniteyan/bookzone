package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util/log"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"time"
)

var GlobalSessions *sessions.Sessions

func init() {
	GlobalSessions = sessions.New(sessions.Config{Cookie:common.SessionId})
}

type BaseController struct {
	Ctx 				iris.Context
	Member 				models.Member
	Option 				map[string]string
	EnableAnonymous		bool
	baseSession 		*sessions.Session
}

func (this *BaseController) JsonResult(errCode common.HttpCode, errMsg string, data ...interface{}) {
	jsonData := make(map[string]interface{}, 3)
	jsonData["errcode"] = errCode
	jsonData["message"] = errMsg

	if len(data) > 0 && data[0] != nil {
		jsonData["data"] = data[0]
	}

	this.Ctx.JSON(jsonData)
}

func (this *BaseController) getSession() *sessions.Session {
	this.baseSession = GlobalSessions.Start(this.Ctx)
	return this.baseSession
}

func (this *BaseController) SetSession(key string, value interface{}) {
	session := this.getSession()
	session.Set(key, value)
}

func (this *BaseController) GetSession(key string) interface{} {
	session := this.getSession()
	return session.Get(key)
}

func (this *BaseController) DelSession(key string) {
	session := this.getSession()
	session.Delete(key)
}

func (this *BaseController) SetMemberSession(member models.Member) {
	if member.MemberId <= 0 {
		this.DelSession(common.MemberSessionName)
		this.DelSession(common.UidSessionName)
		GlobalSessions.Destroy(this.Ctx)
		log.Infof("destory session, member id:%d", member.MemberId)
	} else {
		this.SetSession(common.MemberSessionName, member)
		this.SetSession(common.UidSessionName, member.MemberId)
		log.Infof("set member session success, member id:%d", member.MemberId)
	}
}

type CookieRemember struct {
	MemberId 	int
	Account 	string
	Time 		time.Time
}