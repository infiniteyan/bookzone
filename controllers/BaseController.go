package controllers

import (
	"bookzone/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"time"
)

var Global_session *sessions.Sessions

func init() {
	Global_session = sessions.New(sessions.Config{Cookie:"sid"})
}

type BaseController struct {
	Ctx 				iris.Context
	Member 				*models.Member
	Option 				map[string]string
	EnableAnonymous		bool
}

func (this *BaseController) SetSession(key string, value string) {
	session := Global_session.Start(this.Ctx)
	session.Set(key, value)
}

func (this *BaseController) GetSession(key string) interface{} {
	session := Global_session.Start(this.Ctx)
	return session.Get(key)
}

type CookieRemember struct {
	MemberId 	int
	Account 	string
	Time 		time.Time
}