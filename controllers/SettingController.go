package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util/graphics"
	"bookzone/util/log"
	"bookzone/util/store"
	"github.com/kataras/iris/mvc"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type SettingController struct {
	BaseController
}

func (this *SettingController) BeforeActivation(a mvc.BeforeActivation) {
	log.Infof("SettingController BeforeActivation")
	a.Handle("GET", "/", "Index")
	a.Handle("POST", "/", "Update")
	a.Handle("POST", "/upload", "Upload")
}

func (this *SettingController) Update() {
	log.Infof("SettingController Update")

	email := strings.TrimSpace(this.Ctx.FormValue("email"))
	phone := strings.TrimSpace(this.Ctx.FormValue("phone"))
	description := strings.TrimSpace(this.Ctx.FormValue("description"))
	if email == "" {
		this.JsonResult(common.HttpCodeErrorParameter, "邮箱不能为空")
		return
	}

	session := this.getSession()
	member, _ := session.Get(common.MemberSessionName).(models.Member)
	member.Email = email
	member.Phone = phone
	member.Description = description
	if err := member.Update("email", "phone", "description"); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "提交信息错误")
		return
	}
	this.SetMemberSession(member)
	this.JsonResult(common.HttpCodeSuccess, "ok")
}

func (this *SettingController) Index() mvc.Result {
	log.Infof("SettingController Index")

	dataMap := make(map[string]interface{})
	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if ok {
		dataMap["Member"] = member
	} else {
		dataMap["Member"] = models.Member{}
	}
	dataMap["SettingBasic"] = true

	return mvc.View{
		Name: "setting/index.html",
		Data: dataMap,
	}
}

func (this *SettingController) Upload() {
	var err error
	file, fileHeader, err := this.Ctx.FormFile("image-file")
	if err != nil {
		this.JsonResult(common.HttpCodeErrorFile, "文件异常")
		return
	}
	defer file.Close()
	ext := filepath.Ext(fileHeader.Filename)
	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") &&
		!strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		this.JsonResult(common.HttpCodeErrorFile, "图片格式异常")
		return
	}

	x1, err := strconv.ParseFloat(this.Ctx.FormValue("x"), 32)
	if err != nil {
		x1 = 10
	}
	y1, err := strconv.ParseFloat(this.Ctx.FormValue("y"), 32)
	if err != nil {
		y1 = 10
	}
	w1, err := strconv.ParseFloat(this.Ctx.FormValue("width"), 32)
	if err != nil {
		w1 = 10
	}
	h1, err := strconv.ParseFloat(this.Ctx.FormValue("height"), 32)
	if err != nil {
		h1 = 10
	}

	x := int(x1)
	y := int(y1)
	width := int(w1)
	height := int(h1)

	randomName := strconv.FormatInt(time.Now().UnixNano(), 16)
	filePath := filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("202011"), randomName + ext)
	path := filepath.Dir(filePath)
	os.MkdirAll(path, os.ModePerm)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Infof("fail to save file")
		this.JsonResult(common.HttpCodeErrorInternal, "保存失败")
		return
	}
	_, err = io.Copy(f, file)

	if err != nil {
		log.Infof("fail to save file")
		this.JsonResult(common.HttpCodeErrorInternal, "保存失败")
		return
	}

	subImg, err := graphics.ImageCopyFromFile(filePath, x, y, width, height)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorInternal, "剪切失败")
		return
	}
	os.Remove(filePath)

	filePath = filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("202011"), randomName + ext)
	graphics.ImageResizeSaveFile(subImg, 120, 120, filePath)
	err = graphics.SaveImage(filePath, subImg)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorInternal, "保存文件失败")
		return
	}

	url := "/" + strings.Replace(strings.TrimPrefix(filePath, common.WorkingDirectory), "\\", "/", -1)
	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}

	session := this.getSession()
	member, _ := session.Get(common.MemberSessionName).(models.Member)
	if member, err := models.NewMember().Find(member.MemberId); err == nil {
		avatar := member.Avatar
		member.Avatar = url
		err = member.Update()
		if err != nil {
			this.JsonResult(common.HttpCodeErrorInternal, "保存信息失败")
			return
		}
		if strings.HasPrefix(avatar, "/uploads/") {
			oldPath := filepath.Join(common.WorkingDirectory, avatar)
			os.Remove(oldPath)
		}
		this.SetMemberSession(*member)
	}

	if err := store.SaveToLocal("."+url, strings.TrimLeft(url, "./")); err != nil {
		log.Infof(err.Error())
	} else {
		url = "/" + strings.TrimLeft(url, "./")
	}

	this.JsonResult(common.HttpCodeSuccess, "ok", url)
}