package controllers

import (
	"bookzone/cache"
	"bookzone/common"
	"bookzone/constant"
	"bookzone/models"
	"bookzone/util"
	"bookzone/util/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kataras/iris/mvc"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type DocumentController struct {
	BaseController
}

func (this *DocumentController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("GET", "/{key:string}", "Index")
	a.Handle("GET", "/{bookidentify:string}/{identify:string}", "Read")
	a.Handle("GET", "/{key:string}/edit", "Edit")
	a.Handle("GET", "/content/{key:string}/{id:string}", "Content")
	a.Handle("POST", "/{key:string}/search", "Search")
	a.Handle("POST", "/savecontent", "SaveContent")
}

func (this *DocumentController) getBookData(identify, token string) (*models.BookData, error) {
	book, err := models.NewBook().Select("identify", identify)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	if book.PrivatelyOwned == 1 && !this.Member.IsAdministrator() {
		isOk := false
		_, err := models.NewRelationship().SelectRoleId(book.BookId, this.Member.MemberId)
		if err == nil {
			isOk = true
		}
		if book.PrivateToken != "" && !isOk {
			if token != "" && strings.EqualFold(token, book.PrivateToken) {
				this.SetSession(identify, token)
			} else if token, ok := this.GetSession(identify).(string); !ok || !strings.EqualFold(token, book.PrivateToken) {
				return nil, errors.New("permission denied")
			}
		} else if !isOk {
			return nil, errors.New("permission denied")
		}
	}

	bookResult := book.ToBookData()
	if this.Member.MemberId > 0 {
		rsh, err := models.NewRelationship().Select(bookResult.BookId, this.Member.MemberId)
		if err == nil {
			bookResult.MemberId = rsh.MemberId
			bookResult.RoleId = rsh.RoleId
			bookResult.RelationshipId = rsh.RelationshipId
		}
	}
	return bookResult, nil
}

func (this *DocumentController) Index() mvc.Result {
	log.Infof("DocumentController.Index")
	token := this.Ctx.URLParam("token")
	identify := this.Ctx.Params().Get("key")
	if len(identify) == 0 {
		return mvc.View{
			Name: "error/error.html",
			Data:map[string]interface{}{
				"Info": "invalid request",
			},
		}
	}

	bookResult, err := this.getBookData(identify, token)
	if err != nil {
		return mvc.View{
			Name: "error/error.html",
			Data:map[string]interface{}{
				"Info": "fail to get data info",
			},
		}
	}
	if bookResult.BookId == 0 {
		return mvc.Response{
			Code: 302,
			Path: "/",
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
	tab := strings.ToLower(this.Ctx.URLParam("tab"))
	if tab == "" {
		tab = "default"
	}
	dataMap["SITE_NAME"] = "BOOKZONE"
	dataMap["Tab"] = tab
	dataMap["Book"] = bookResult
	if member.MemberId > 0 {
		dataMap["MyScore"] = new(models.Score).BookScoreByUid(this.Member.MemberId, bookResult.BookId)
	} else {
		dataMap["MyScore"] = 0
	}

	bookCacheKey := "dynamiccache_document_" + identify
	bookMenuCacheKey := bookCacheKey + "menu"
	bookCommentCacheKey := bookCacheKey + "comments"

	var dataMenu []*models.Document
	err = cache.ReadStruct(bookMenuCacheKey, &dataMenu)
	if err != nil {
		dataMenu, _ = new(models.Document).GetMenuTop(bookResult.BookId)
		cache.WriteStruct(bookMenuCacheKey, &dataMenu, util.RandomExpire(constant.MIN_REDIS_EXPIRE_SEC, constant.MAX_REDIS_EXPIRE_SEC))
	}
	dataMap["Menu"] = dataMenu

	var dataComments []*models.BookCommentsResult
	err = cache.ReadStruct(bookCommentCacheKey, &dataComments)
	if err != nil {
		dataComments, _ = new(models.Comments).BookComments(1, 30, bookResult.BookId)
		cache.WriteStruct(bookCommentCacheKey, &dataComments, util.RandomExpire(constant.MIN_REDIS_EXPIRE_SEC, constant.MAX_REDIS_EXPIRE_SEC))
	}
	dataMap["Comments"] = dataComments

	return mvc.View{
		Name: "document/intro.html",
		Data: dataMap,
	}
}

func (this *DocumentController) IsFromAjax() bool {
	if this.Ctx.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		return true
	} else {
		return false
	}
}

func (this *DocumentController) Read() mvc.Result {
	log.Infof("DocumentController.Read")
	bookIdentify := this.Ctx.Params().Get("bookidentify")
	identify := this.Ctx.Params().Get("identify")
	token := this.Ctx.URLParam("token")

	if bookIdentify == "" || identify == "" {
		return mvc.Response{
			Code: 404,
		}
	}

	bookData, err := this.getBookData(bookIdentify, token)
	if err != nil {
		return mvc.Response{
			Code: 404,
		}
	}

	document := models.NewDocument()
	doc, err := document.SelectByIdentify(bookData.BookId, identify)
	if err != nil {
		log.Infof(err.Error())
		return mvc.Response{
			Code: 404,
		}
	}

	if doc.BookId != bookData.BookId {
		return mvc.Response{
			Code: 404,
		}
	}
	if doc.Release != "" {
		//query, err := goquery.NewDocumentFromReader(bytes.NewBufferString(doc.Release))
		//if err != nil {
		//	beego.Error(err)
		//} else {
		//	ossdomain := strings.TrimRight(beego.AppConfig.String("oss_attach_domain"), "/")
		//	query.Find("img").Each(func(i int, contentSelection *goquery.Selection) {
		//		if src, ok := contentSelection.Attr("src"); ok {
		//			if !(strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "http://")) {
		//				src = ossdomain + "/" + strings.TrimLeft(src, "./")
		//				contentSelection.SetAttr("src", src)
		//				beego.Debug(src)
		//			}
		//		}
		//		if alt, _ := contentSelection.Attr("alt"); alt == "" {
		//			contentSelection.SetAttr("alt", doc.DocumentName+" - 图"+fmt.Sprint(i+1))
		//		}
		//	})
		//	html, err := query.Find("body").Html()
		//	if err != nil {
		//		beego.Error(err)
		//	} else {
		//		doc.Release = html
		//	}
		//}
	}

	attach, err := models.NewAttachment().SelectByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	//图书阅读人次+1
	if err := util.IncOrDec("md_books", "vcnt",
		fmt.Sprintf("book_id=%v", doc.BookId),
		true, 1,
	); err != nil {
		log.Infof(err.Error())
	}

	//文档阅读人次+1
	if err := util.IncOrDec("md_documents", "vcnt",
		fmt.Sprintf("document_id=%v", doc.DocumentId),
		true, 1,
	); err != nil {
		log.Infof(err.Error())
	}
	doc.Vcnt = doc.Vcnt + 1

	if this.IsFromAjax() {
		var data struct {
			Id        int    `json:"doc_id"`
			DocTitle  string `json:"doc_title"`
			Body      string `json:"body"`
			Title     string `json:"title"`
			View      int    `json:"view"`
			UpdatedAt string `json:"updated_at"`
		}
		data.DocTitle = doc.DocumentName
		data.Body = doc.Release
		data.Id = doc.DocumentId
		data.View = doc.Vcnt
		data.UpdatedAt = doc.ModifyTime.Format("2006-01-02 15:04:05")

		this.JsonResult(common.HttpCodeSuccess, "ok", data)
		return mvc.Response{
			Code: 200,
		}
	}

	tree, err := models.NewDocument().GetMenuHtml(bookData.BookId, doc.DocumentId)
	if err != nil {
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
	dataMap["SITE_NAME"] = "BOOKZONE"
	dataMap["Bookmark"] = false
	dataMap["Model"] = bookData
	dataMap["Book"] = bookData
	dataMap["Result"] = template.HTML(tree)
	dataMap["Title"] = doc.DocumentName
	dataMap["DocId"] = doc.DocumentId
	dataMap["Content"] = template.HTML(doc.Release)
	dataMap["View"] = doc.Vcnt
	dataMap["UpdatedAt"] = doc.ModifyTime.Format("2006-01-02 15:04:05")

	return mvc.View{
		Name: "document/default_read.html",
		Data: dataMap,
	}
}

func (this *DocumentController) Search() mvc.Result {
	log.Infof("DocumentController.Search")
	return mvc.View{
	}
}

func (this *DocumentController) Edit() mvc.Result {
	log.Infof("DocumentController.Edit")
	docId := 0

	identify := this.Ctx.Params().Get("key")
	if identify == "" {
		log.Infof("indentify empty")
		return mvc.Response{
			Code: 404,
		}
	}

	bookData := models.NewBookData()
	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if !ok {
		log.Infof("please login first")
		return mvc.Response{
			Code: 404,
		}
	}
	var err error
	if member.IsAdministrator() {
		book, err := models.NewBook().Select("identify", identify)
		if err != nil {
			log.Infof(err.Error())
			return mvc.Response{Code: 404}
		}
		bookData = book.ToBookData()
	} else {
		bookData, err = models.NewBookData().SelectByIdentify(identify, member.MemberId)
		if err != nil {
			log.Infof(err.Error())
			return mvc.Response{Code: 404}
		}

		if bookData.RoleId == common.BookGeneral {
			log.Infof(err.Error())
			return mvc.Response{Code: 404}
		}
	}

	dataMap := make(map[string]interface{})
	dataMap["Model"] = bookData
	r, _ := json.Marshal(bookData)
	dataMap["ModelResult"] = template.JS(string(r))

	dataMap["Result"] = template.JS("[]")

	//if id := c.GetString(":id"); id != "" {
	//	if num, _ := strconv.Atoi(id); num > 0 {
	//		docId = num
	//	} else { //字符串
	//		var doc = models.NewDocument()
	//		models.GetOrm("w").QueryTable(doc).Filter("identify", id).Filter("book_id", bookData.BookId).One(doc, "document_id")
	//		docId = doc.DocumentId
	//	}
	//}

	trees, err := models.NewDocument().GetMenu(bookData.BookId, docId, true)
	if err != nil {
		log.Infof(err.Error())
		return mvc.Response{Code: 404}
	} else {
		if len(trees) > 0 {
			if jsTree, err := json.Marshal(trees); err == nil {
				dataMap["Result"] = template.JS(string(jsTree))
			}
		} else {
			dataMap["Result"] = template.JS("[]")
		}
	}
	dataMap["BaiDuMapKey"] = "mapkey"

	return mvc.View{
		Name: "document/markdown_edit_template.html",
		Data: dataMap,
	}
}

func (this *DocumentController) Content() {
	log.Infof("DocumentController.Content")
	identify := this.Ctx.Params().Get("key")
	docId, _ := strconv.Atoi(this.Ctx.Params().Get("id"))

	session := this.getSession()
	member, ok := session.Get(common.MemberSessionName).(models.Member)
	if !ok {
		log.Infof("please login first")
		this.JsonResult(common.HttpCodeErrorPermissionDeny, "请先登录")
		return
	}
	if !member.IsAdministrator() {
		bookData, err := models.NewBookData().SelectByIdentify(identify, member.MemberId)
		if err != nil || bookData.RoleId == common.BookGeneral {
			this.JsonResult(common.HttpCodeErrorPermissionDeny, "鉴权失败")
			return
		}
	}

	if docId <= 0 {
		this.JsonResult(common.HttpCodeErrorParameter, "参数错误")
		return
	}

	documentStore := new(models.DocumentStore)
	doc, err := models.NewDocument().SelectByDocId(docId)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "文档不存在")
		return
	}
	attach, err := models.NewAttachment().SelectByDocumentId(doc.DocumentId)
	if err == nil {
		doc.AttachList = attach
	}

	doc.Release = ""
	doc.Markdown = documentStore.SelectField(doc.DocumentId, "markdown")
	this.JsonResult(common.HttpCodeSuccess, "ok", doc)
}

func (this *DocumentController) SaveContent() {
	log.Infof("DocumentController.SaveContent")
	identify := this.Ctx.FormValue("identify")
	docId, err := strconv.Atoi(this.Ctx.FormValue("doc_id"))

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
			this.JsonResult(common.HttpCodeErrorDatabase, "获取内容错误")
			return
		}
		bookId = book.BookId
	} else {
		bookData, err := models.NewBookData().SelectByIdentify(identify, member.MemberId)
		if err != nil || bookData.RoleId == common.BookGeneral {
			this.JsonResult(common.HttpCodeErrorPermissionDeny, "鉴权失败")
			return
		}
		bookId = bookData.BookId
	}

	if docId <= 0 {
		this.JsonResult(common.HttpCodeErrorParameter, "参数错误")
		return
	}

	documentStore := new(models.DocumentStore)

	markdown := strings.TrimSpace(this.Ctx.FormValue("markdown"))
	content := this.Ctx.FormValue("html")

	version, _ := strconv.Atoi(this.Ctx.FormValue("version"))
	isCover := this.Ctx.FormValue("cover")

	doc, err := models.NewDocument().SelectByDocId(docId)
	if err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "读取文档错误")
		return
	}
	if doc.BookId != bookId {
		this.JsonResult(common.HttpCodeErrorInternal, "内部错误")
		return
	}
	if int(doc.Version) != version && !strings.EqualFold(isCover, "yes") {
		this.JsonResult(common.HttpCodeErrorInternal, "文档将被覆盖")
		return
	}

	if markdown == "" && content != "" {
		documentStore.Markdown = content
	} else {
		documentStore.Markdown = markdown
	}
	documentStore.Content = content
	doc.Version = time.Now().Unix()
	if docId, err := doc.InsertOrUpdate(); err != nil {
		this.JsonResult(common.HttpCodeErrorDatabase, "保存失败")
		return
	} else {
		documentStore.DocumentId = int(docId)
		if err := documentStore.InsertOrUpdate("markdown", "content"); err != nil {
			this.JsonResult(common.HttpCodeErrorDatabase, "保存失败")
			return
		}
	}

	doc.Release = ""
	this.JsonResult(common.HttpCodeSuccess, "ok", doc)
}