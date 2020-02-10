package controllers

import (
	"bookzone/models"
	"errors"
	"github.com/kataras/iris/mvc"
	"log"
	"strings"
)

type DocumentController struct {
	BaseController
}

func (this *DocumentController) BeforeActivation(a mvc.BeforeActivation) {
	a.Handle("GET", "/{key:string}", "Index")
	a.Handle("GET", "/{key:string}/{id:string}", "Read")
	a.Handle("POST", "/{key:string}/search", "Search")
}

func (this *DocumentController) getBookData(identify, token string) (*models.BookData, error) {
	book, err := models.NewBook().SelectByIdentify(identify)
	if err != nil {
		log.Println(err)
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
	if this.Member != nil {
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
	tab := strings.ToLower(this.Ctx.URLParam("tab"))
	dataMap["Tab"] = tab
	dataMap["Book"] = bookResult
	if this.Member != nil {
		dataMap["MyScore"] = new(models.Score).BookScoreByUid(this.Member.MemberId, bookResult.BookId)
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

func (this *DocumentController) Read() mvc.Result {
	return mvc.View{
	}
}

func (this *DocumentController) Search() mvc.Result {
	return mvc.View{
	}
}