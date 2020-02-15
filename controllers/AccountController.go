package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util/log"
	"errors"
	"github.com/kataras/iris/mvc"
	"regexp"
	"time"
)

type AccountController struct {
	BaseController
}

func (this *AccountController) BeforeActivation(a mvc.BeforeActivation) {
	log.Infof("AccountController BeforeActivation")
	a.Handle("GET", "/logout", "Logout")
}

func (this *AccountController) GetLogin() mvc.Result {
	log.Infof("AccountController, get login")
	return mvc.View{
		Name: "account/login.html",
		Data:map[string]interface{} {
			"SITE_NAME":"BOOKZONE",
		},
	}
}

func (this *AccountController) PostLogin() {
	log.Infof("AccountController, post login")
	account := this.Ctx.FormValue("account")
	password := this.Ctx.FormValue("password")

	if account == "" || password == "" {
		this.JsonResult(common.HttpCodeErrorParameter, "参数错误")
	}
	log.Infof("%s   %s", account, password)
	member, err := models.NewMember().Login(account, password)
	if err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorPassword, "密码错误，请重新输入")
		return
	}
	member.LastLoginTime = time.Now()
	member.Update()
	this.SetMemberSession(*member)
	this.JsonResult(common.HttpCodeSuccess, "ok")
	log.Infof("AccountController, login success")
	this.Ctx.Redirect("/")
}

func (this *AccountController) GetRegist() mvc.Result {
	return mvc.View{
		Name: "account/register.html",
	}
}

func (this *AccountController) Logout() {
	log.Infof("AccountController, logout")
	this.SetMemberSession(models.Member{MemberId: -1})
	this.Ctx.Redirect("/")
}

func (this *AccountController) PostRegist() {
	account := this.Ctx.FormValue("account")
	password1 := this.Ctx.FormValue("password1")
	password2 := this.Ctx.FormValue("password2")
	nickname := this.Ctx.FormValue("nickname")
	email := this.Ctx.FormValue("email")

	log.Infof("AccountController, post register:%s %s %s %s %s", account, password1, password2, nickname, email)

	if password1 != password2 {
		this.JsonResult(common.HttpCodeErrorPassword, "登录密码与确认密码不一致")
		return
	}

	if l := len(password1); password1 == "" || l < 6 || l > 20 {
		this.JsonResult(common.HttpCodeErrorPassword, "密码必须在6-20个字符之间")
		return
	}

	if ok, err := regexp.MatchString(common.RegexpEmail, email); !ok || err != nil {
		this.JsonResult(common.HttpCodeErrorEmail, "邮箱格式错误")
		return
	}

	if l := len(nickname); l < 2 || l > 20 {
		this.JsonResult(common.HttpCodeErrorNickname, "用户昵称限制在2-20个字符")
		return
	}

	member := models.NewMember()
	member.Account = account
	member.Nickname = nickname
	member.Password = password1
	member.Email = email
	member.Status = 0
	if account == "admin" || account == "administrator" {
		member.Role = common.MemberSuperRole
	} else {
		member.Role = common.MemberGeneralRole
	}
	member.Avatar = common.DefaultAvatar()
	member.CreateAt = 0
	member.CreateTime = time.Now()

	if err := member.Add(); err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorRegisterFail, "注册失败")
		return
	}

	log.Infof("register user success")

	if err := this.login(member.MemberId); err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorLoginFail, "注册后登录失败")
		return
	}

	this.Ctx.Redirect("/")
}

func (this *AccountController) login(memberId int) error {
	member, err := models.NewMember().Find(memberId)
	if err != nil || member == nil || member.MemberId == 0 {
		return errors.New("用户不存在")
	}

	member.LastLoginTime = time.Now()
	member.Update()
	this.SetMemberSession(*member)
	var remember CookieRemember
	remember.MemberId = member.MemberId
	remember.Account = member.Account
	remember.Time = time.Now()
	return nil
}