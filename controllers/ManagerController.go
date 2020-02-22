package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util"
	"bookzone/util/log"
	"fmt"
	"github.com/kataras/iris/mvc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ManagerController struct {
	BaseController
}

func (this *ManagerController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("GET", "/category", "Category")
	a.Handle("GET", "/delcategory", "DelCategory")
	a.Handle("GET", "/updatecategory", "UpdateCategory")
	a.Handle("POST", "/addcategory", "AddCategory")
	a.Handle("POST", "/updateicon", "UpdateCateIcon")
}

func (this *ManagerController) Category() mvc.Result {
	log.Infof("ManagerController.Category")

	cates, err := models.NewCategory().GetCates(-1, -1)
	if err != nil {
		return mvc.Response{
			Code: 404,
		}
	}

	var parents []models.Category
	for idx, item := range cates {
		if strings.TrimSpace(item.Icon) == "" {
			item.Icon = "/static/images/icon.png"
		} else {
			item.Icon = item.Icon
		}
		if item.Pid == 0 {
			parents = append(parents, item)
		}
		cates[idx] = item
	}

	dataMap := make(map[string]interface{})
	dataMap["Parents"] = parents
	dataMap["Cates"] = cates
	dataMap["IsCategory"] = true
	dataMap["SITE_NAME"] = "BOOKZONE"

	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if ok {
		dataMap["Member"] = member
	} else {
		dataMap["Member"] = models.Member{}
	}

	return mvc.View{
		Name: "manager/category.html",
		Data: dataMap,
	}
}

func (this *ManagerController) AddCategory() {
	log.Infof("ManagerController.AddCategory")

	pidStr := this.Ctx.FormValue("pid")
	pid, _ := strconv.Atoi(pidStr)
	if err := models.NewCategory().InsertMulti(pid, this.Ctx.FormValue("cates")); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "新增失败：" + err.Error())
	}
	this.JsonResult(common.HttpCodeSuccess, "新增成功")
}

func (this *ManagerController) DelCategory() {
	log.Infof("ManagerController.DelCategory")
	var err error
	var id int
	if id, err = this.Ctx.URLParamInt("id"); err != nil {
		this.JsonResult(common.HttpCodeErrorParameter, "参数错误")
		return
	}
	err = models.NewCategory().Delete(id)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, err.Error())
		return
	}
	this.JsonResult(common.HttpCodeSuccess, "删除成功")
}

func (this *ManagerController) UpdateCategory() {
	log.Infof("ManagerController.UpdateCategory")
	field := this.Ctx.URLParam("field")
	val := this.Ctx.URLParam("value")
	id, _ := strconv.Atoi(this.Ctx.URLParam("id"))
	if err := models.NewCategory().UpdateField(id, field, val); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "更新失败：" + err.Error())
		return
	}
	this.JsonResult(common.HttpCodeSuccess, "更新成功")
}

func (this *ManagerController) UpdateCateIcon() {
	log.Infof("ManagerController.UpdateCateIcon")

	var err error
	idStr := this.Ctx.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	if id == 0 {
		this.JsonResult(common.HttpCodeErrorParameter, "参数不正确")
		return
	}
	category := models.NewCategory()
	if cate, err := category.Find(id); err == nil && cate.Id > 0 {
		cate.Icon = strings.TrimLeft(cate.Icon, "/")
		f, h, err1 := this.Ctx.FormFile("icon")
		if err1 != nil {
			this.JsonResult(common.HttpCodeErrorFile, "文件上传错误")
			return
		}
		defer f.Close()

		tmpFile := fmt.Sprintf("uploads/icons/%v%v" + filepath.Ext(h.Filename), id, time.Now().Unix())
		os.MkdirAll(filepath.Dir(tmpFile), os.ModePerm)
		if err = this.SaveToFile("icon", tmpFile); err == nil {
			util.DeleteLocalFiles(cate.Icon)
			err = category.UpdateField(cate.Id, "icon", "/" + tmpFile)
		}
	}

	if err != nil {
		this.JsonResult(common.HttpCodeErrorFile, err.Error())
		return
	}
	this.JsonResult(common.HttpCodeSuccess, "更新成功")
}