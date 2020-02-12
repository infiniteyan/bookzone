package controllers

import (
	"bookzone/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"bookzone/util/log"
	"math"
	"strconv"
	"ziyoubiancheng/mbook/utils"
)

type ExploreController struct {
	Ctx iris.Context
}

func (this *ExploreController) Get() mvc.Result {
	urlPrefix := "/explore"
	dataMap := make(map[string]interface{})
	id, _ := strconv.Atoi(this.Ctx.URLParam("cid"))
	log.Infof("get items by cid:%d", id)
	category, _ := new(models.Category).Find(id)
	dataMap["Cid"] = id
	dataMap["Cate"] = category

	//pageStr := this.Ctx.Params().Get("page")
	pageIndex := 1
	pageSize := 24
	books, totalCount, err := new(models.Book).HomeData(pageIndex, pageSize, id)
	if err != nil {
		return mvc.View{
			Name: "error/error.html",
			Data:map[string]interface{}{
				"Info": "database query error",
			},
		}
	}

	if totalCount > 0 {
		urlSuffix := ""
		if id > 0 {
			urlSuffix = urlSuffix + "&cid=" + strconv.Itoa(id)
		}
		html := utils.NewPaginations(4, totalCount, pageSize, pageIndex, urlPrefix, urlSuffix)
		dataMap["PageHtml"] = html
	} else {
		dataMap["PageHtml"] = ""
	}
	dataMap["TotalPages"] = int(math.Ceil(float64(totalCount) / float64(pageSize)))
	dataMap["Lists"] = books

	return mvc.View{
		Name: "explore/index.html",
		Data: dataMap,
	}
}