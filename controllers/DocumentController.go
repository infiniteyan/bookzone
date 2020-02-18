package controllers

import (
	"bookzone/common"
	"bookzone/models"
	"bookzone/util"
	"bookzone/util/log"
	"errors"
	"fmt"
	"github.com/kataras/iris/mvc"
	"html/template"
	"strings"
)

type DocumentController struct {
	BaseController
}

func (this *DocumentController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("GET", "/{key:string}", "Index")
	a.Handle("GET", "/{bookidentify:string}/{identify:string}", "Read")
	a.Handle("POST", "/{key:string}/search", "Search")
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

	var dataMenu []*models.Document
	dataMenu, _ = new(models.Document).GetMenuTop(bookResult.BookId)
	dataMap["Menu"] = dataMenu

	var dataComments []*models.BookCommentsResult
	dataComments, _ = new(models.Comments).BookComments(1, 30, bookResult.BookId)
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
	return mvc.View{
	}
}