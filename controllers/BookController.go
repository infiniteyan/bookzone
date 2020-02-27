package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/mq"
	"bookzone/util"
	"bookzone/util/graphics"
	"bookzone/util/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/mvc"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"ziyoubiancheng/mbook/utils"
)

type BookController struct {
	BaseController
}

func (this *BookController) BeforeActivation(a mvc.BeforeActivation) {
	log.Infof("BookController BeforeActivation")
	a.Handle("GET", "/", "Index")
	a.Handle("GET", "/{id:int}", "Score")
	a.Handle("GET", "/collect/{id:int}", "Collection")
	a.Handle("GET", "/{key:string}/setting", "Setting")
	a.Handle("POST", "/comment/{id:int}", "Comment")
	a.Handle("POST", "/create", "Create")
	a.Handle("POST", "/uploadcover", "UploadCover")
	a.Handle("POST", "/createtoken", "CreateToken")
	a.Handle("POST", "/savebook", "SaveBook")
	a.Handle("POST", "/release/{key:string}", "Release")
}

func (this *BookController) Index() mvc.Result {
	data := make(map[string]interface{})
	data["SITE_NAME"] = "BOOKZONE"
	data["SettingBook"] = true

	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if ok {
		data["Member"] = member
	} else {
		data["Member"] = models.Member{}
	}

	pageIndex, err := this.Ctx.URLParamInt("page")
	if err != nil {
		pageIndex = 1
	}
	private, err := this.Ctx.URLParamInt("private")
	if err != nil {
		private = 1
	}
	books, totalCount, err := models.NewBook().SelectPage(pageIndex, common.PageSize, member.MemberId, private)
	if err != nil {
		log.Errorf("BookController.Index => ", err)
		return mvc.Response{
			Code: 404,
		}
	}
	if totalCount > 0 {
		data["PageHtml"] = util.NewPaginations(common.RollPage, totalCount, common.PageSize, pageIndex, "/book", fmt.Sprintf("&private=%v", private))
	} else {
		data["PageHtml"] = ""
	}

	for idx, book := range books {
		book.Cover = book.Cover
		books[idx] = book
	}
	b, err := json.Marshal(books)
	if err != nil || len(books) <= 0 {
		data["Result"] = template.JS("[]")
	} else {
		data["Result"] = template.JS(string(b))
	}
	data["Private"] = private

	return mvc.View{
		Name: "book/index.html",
		Data: data,
	}
}

func (this *BookController) Create() {
	log.Infof("BookController.Create")
	identify := strings.TrimSpace(this.Ctx.FormValue("identify"))
	bookName := strings.TrimSpace(this.Ctx.FormValue("book_name"))
	author := strings.TrimSpace(this.Ctx.FormValue("author"))
	authorURL := strings.TrimSpace(this.Ctx.FormValue("author_url"))
	privatelyOwned, _ := strconv.Atoi(this.Ctx.FormValue("privately_owned"))
	description := strings.TrimSpace(this.Ctx.FormValue("description"))

	if identify == "" || len(identify) > 50 {
		this.JsonResult(common.HttpCodeErrorParameter, "请正确填写图书标识，不能超过50字")
		return
	}
	if bookName == "" {
		this.JsonResult(common.HttpCodeErrorParameter, "请填图书名称")
		return
	}

	if len(description) > 500 {
		this.JsonResult(common.HttpCodeErrorParameter, "图书描述需小于500字")
		return
	}

	if privatelyOwned != 0 && privatelyOwned != 1 {
		privatelyOwned = 1
	}

	if book, _ := models.NewBook().Select("identify", identify); book != nil && book.BookId > 0 {
		log.Errorf("book with identify %s already exists", identify)
		this.JsonResult(common.HttpCodeErrorParameter, "identify冲突")
		return
	}

	session := this.getSession()
	member, _ := session.Get(common.MemberSessionName).(models.Member)

	book := models.NewBook()
	book.BookName = bookName
	book.Identify = identify
	book.Description = description
	book.CommentCount = 0
	book.PrivatelyOwned = privatelyOwned
	book.Cover = "/static/images/book.jpg"
	book.DocCount = 0
	book.MemberId = member.MemberId
	book.CommentCount = 0
	book.Editor = "markdown"
	book.ReleaseTime = time.Now()
	book.CreateTime = time.Now()
	book.ModifyTime = time.Now()
	book.Score = 40
	book.Author = author
	book.AuthorUrl = authorURL

	if err := book.Insert(); err != nil {
		log.Errorf(err.Error())
		this.JsonResult(common.HttpCodeErrorDatabase, "数据库错误")
		return
	}

	bookResult, err := models.NewBookData().SelectByIdentify(book.Identify, member.MemberId)
	if err != nil {
		log.Errorf(err.Error())
		this.JsonResult(common.HttpCodeErrorDatabase, "数据库错误")
		return
	}

	this.JsonResult(common.HttpCodeSuccess, "ok", bookResult)
}

func (this *BookController) SaveBook() {
	log.Infof("BookController.SaveBook")

	bookResult, err := this.hasPermission()
	if err != nil {
		this.JsonResult(common.HttpCodeErrorPermissionDeny, err.Error())
		return
	}

	book, err := models.NewBook().Select("book_id", bookResult.BookId)
	if err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorDatabase, err.Error())
		return
	}

	bookName := strings.TrimSpace(this.Ctx.FormValue("book_name"))
	description := strings.TrimSpace(this.Ctx.FormValue("description"))
	editor := strings.TrimSpace(this.Ctx.FormValue("editor"))

	if len(description) > 500 {
		this.JsonResult(common.HttpCodeErrorParameter, "描述需小于500字")
		return
	}

	if editor != "markdown" && editor != "html" {
		editor = "markdown"
	}

	book.BookName = bookName
	book.Description = description
	book.Editor = editor
	book.Author = this.Ctx.FormValue("author")
	book.AuthorUrl = this.Ctx.FormValue("author_url")

	if err := book.Update(); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "保存失败")
		return
	}
	bookResult.BookName = bookName
	bookResult.Description = description

	cidValues := this.Ctx.FormValues()["cid"]

	if len(cidValues) != 0 {
		models.NewBookCategory().SetBookCates(book.BookId, cidValues)
	}

	this.JsonResult(common.HttpCodeSuccess, "ok", bookResult)
}

func (this *BookController) Setting() mvc.Result {
	log.Infof("BookController Setting")
	key := this.Ctx.Params().Get("key")

	if key == "" {
		return mvc.Response{
			Code: 404,
		}
	}

	dataMap := make(map[string]interface{})
	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if ok {
		dataMap["Member"] = member
	} else {
		dataMap["Member"] = models.Member{}
	}
	book, err := models.NewBookData().SelectByIdentify(key, member.MemberId)
	if err != nil {
		return mvc.Response{
			Code: 404,
		}
	}

	if book.RoleId != common.BookFounder && book.RoleId != common.BookAdmin {
		return mvc.Response{
			Code: 404,
		}
	}

	if selectedCates, rows, _ := models.NewBookCategory().SelectByBookId(book.BookId); rows > 0 {
		var maps = make(map[int]bool)
		for _, cate := range selectedCates {
			maps[cate.Id] = true
		}
		dataMap["Maps"] = maps
	}

	dataMap["SITE_NAME"] = "BOOKZONE"
	dataMap["Cates"], _ = new(models.Category).GetCates(-1, 1)
	dataMap["Model"] = book

	return mvc.View{
		Name: "book/setting.html",
		Data: dataMap,
	}
}

func (this *BookController) Collection() {
	log.Infof("BookController Collection")

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

	bookId, err := this.Ctx.Params().GetInt("id")
	if err != nil {
		this.JsonResult(common.HttpCodeErrorParameter, "参数错误")
		return
	}

	if bookId <= 0 {
		this.JsonResult(common.HttpCodeErrorParameter, "收藏失败，图书不存在")
		return
	}

	cancel, err := models.NewCollection().Collection(member.MemberId, bookId)
	data := map[string]bool{"IsCancel": cancel}
	if err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorDatabase, err.Error())
		return
	}

	if cancel {
		this.JsonResult(common.HttpCodeSuccess, "取消收藏成功", data)
		return
	}
	this.JsonResult(common.HttpCodeSuccess, "添加收藏成功", data)
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
		msgBody := &mq.CommentEntity{
			BookId: bookId,
			MemberId: member.MemberId,
			Content: content,
		}
		dataStr, _ := json.Marshal(&msgBody)
		msgEntity := &mq.MsgEntity{
			Type: mq.MSG_TYPE_COMMENT,
			Data: string(dataStr),
		}
		ret, _ := json.Marshal(&msgEntity)
		mq.GlobalSynWorker.Push(string(ret))
		this.JsonResult(common.HttpCodeSuccess, "评论成功")
	} else {
		this.JsonResult(common.HttpCodeErrorParameter, "文档图书不存在")
	}
}

func (this *BookController) Score() {
	log.Infof("BookController Score")
	bookId, err := this.Ctx.Params().GetInt("id")
	if err != nil {
		this.JsonResult(common.HttpCodeErrorParameter, "参数错误")
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

func (this *BookController) hasPermission() (*models.BookData, error) {
	identify := this.Ctx.FormValue("identify")
	session := this.getSession()
	member, _ := session.Get(common.MemberSessionName).(models.Member)
	book, err := models.NewBookData().SelectByIdentify(identify, member.MemberId)
	if err != nil {
		return book, err
	}

	if book.RoleId != common.BookAdmin && book.RoleId != common.BookFounder {
		return book, errors.New("权限不足")
	}
	return book, nil
}

func (this *BookController) UploadCover() {
	log.Infof("BookController.UploadCover")

	bookResult, err := this.hasPermission()
	if err != nil {
		this.JsonResult(common.HttpCodeErrorPermissionDeny, err.Error())
		return
	}

	book, err := models.NewBook().Select("book_id", bookResult.BookId)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, err.Error())
		return
	}

	file, header, err := this.Ctx.FormFile("image-file")
	if err != nil {
		log.Infof(err.Error())
		this.JsonResult(common.HttpCodeErrorFile, "读取文件异常")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)

	if !strings.EqualFold(ext, ".png") && !strings.EqualFold(ext, ".jpg") && !strings.EqualFold(ext, ".gif") && !strings.EqualFold(ext, ".jpeg") {
		this.JsonResult(common.HttpCodeErrorFile, "不支持图片格式")
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

	fileName := strconv.FormatInt(time.Now().UnixNano(), 16)
	filePath := filepath.Join("uploads", time.Now().Format("200601"), fileName + ext)

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

	filePath = filepath.Join(common.WorkingDirectory, "uploads", time.Now().Format("200601"), fileName + ext)
	err = graphics.ImageResizeSaveFile(subImg, 175, 230, filePath)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorInternal, "保存文件失败")
		return
	}

	url := "/" + strings.Replace(strings.TrimPrefix(filePath, common.WorkingDirectory), "\\", "/", -1)
	if strings.HasPrefix(url, "//") {
		url = string(url[1:])
	}
	book.Cover = url

	if err := book.Update(); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "保存图片失败")
		return
	}

	this.JsonResult(common.HttpCodeSuccess, "ok", book.Cover)
}

func (this *BookController) CreateToken() {
	log.Infof("BookController.CreateToken")
	action := this.Ctx.FormValue("action")
	bookResult, err := this.hasPermission()
	if err != nil {
		this.JsonResult(common.HttpCodeErrorPermissionDeny, err.Error())
		return
	}

	log.Infof("bookid: %d", bookResult.BookId)

	book := &models.Book{}
	if book, err = models.NewBook().Select("book_id", bookResult.BookId); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "图书不存在")
		return
	}

	if action == "create" {
		if bookResult.PrivatelyOwned == 0 {
			this.JsonResult(common.HttpCodeErrorInternal, "公开图书不能创建令牌")
			return
		}
		book.PrivateToken = string(util.Krand(12, util.KC_RAND_KIND_ALL))
		if err := book.Update(); err != nil {
			this.JsonResult(common.HttpCodeErrorDatabase, "生成阅读失败")
			return
		}
		this.JsonResult(common.HttpCodeSuccess, "ok", book.PrivateToken)
		return
	} else {
		book.PrivateToken = ""
		if err := book.ResetPrivateToken(); err != nil {
			log.Errorf(err.Error())
			this.JsonResult(common.HttpCodeErrorDatabase, "删除令牌失败")
			return
		}
		this.JsonResult(common.HttpCodeSuccess, "ok", "")
	}
}

func (this *BookController) Release() {
	log.Infof("BookController.Release")
	identify := this.Ctx.Params().Get("key")

	bookId := 0
	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if !ok {
		log.Infof("please login first")
		this.JsonResult(common.HttpCodeErrorPermissionDeny, "请先登录")
		return
	}

	if member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			this.JsonResult(common.HttpCodeErrorDatabase, "内部错误")
			return
		}
		bookId = book.BookId
	} else {
		book, err := models.NewBookData().SelectByIdentify(identify, member.MemberId)
		if err != nil {
			this.JsonResult(common.HttpCodeErrorDatabase, "内部错误")
			return
		}
		if book.RoleId != common.BookAdmin && book.RoleId != common.BookFounder && book.RoleId != common.BookEditor {
			this.JsonResult(common.HttpCodeErrorPermissionDeny, "鉴权失败")
			return
		}
		bookId = book.BookId
	}

	if exist := utils.BooksRelease.Exist(bookId); exist {
		this.JsonResult(common.HttpCodeErrorBookRelease, "正在发布中，请稍后操作")
		return
	}

	go func() {
		models.NewDocument().ReleaseContent(bookId)
	}()

	this.JsonResult(common.HttpCodeSuccess, "已发布")
}