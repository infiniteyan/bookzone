package util

import (
	"bookzone/sysinit"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func NewPaginations(rollPage, totalRows, listRows, currentPage int, urlPrefix string, urlSuffix string, urlParams ...interface{}) template.HTML {
	var (
		htmlPage, path string
		pages          []int
		params         []string
	)
	//总页数
	totalPage := totalRows / listRows
	if totalRows%listRows > 0 {
		totalPage += 1
	}
	//只有1页的时候，不分页
	if totalPage < 2 {
		return ""
	}
	paramsLen := len(urlParams)
	if paramsLen > 0 {
		if paramsLen%2 > 0 {
			paramsLen = paramsLen - 1
		}
		for i := 0; i < paramsLen; {
			key := strings.TrimSpace(fmt.Sprintf("%v", urlParams[i]))
			val := strings.TrimSpace(fmt.Sprintf("%v", urlParams[i+1]))
			//键存在，同时值不为0也不为空
			if len(key) > 0 && len(val) > 0 && val != "0" {
				params = append(params, key, val)
			}
			i = i + 2
		}
	}

	path = strings.Trim(urlPrefix, "/")
	if len(params) > 0 {
		path = path + "/" + strings.Trim(strings.Join(params, "/"), "/")
	}
	//最后再处理一次“/”，是为了防止urlPrifix参数为空时，出现多余的“/”
	path = "/" + strings.Trim(path, "/")

	if currentPage > totalPage {
		currentPage = totalPage
	}
	if currentPage < 1 {
		currentPage = 1
	}
	index := 0
	rp := rollPage * 2
	for i := rp; i > 0; i-- {
		p := currentPage + rollPage - i
		if p > 0 && p <= totalPage {

			pages = append(pages, p)
		}
	}
	for k, v := range pages {
		if v == currentPage {
			index = k
		}
	}
	pages_len := len(pages)
	if currentPage > 1 {
		htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`?page=1%v">1..</a></li><li><a class="num" href="`+path+`?page=%d%v">«</a></li>`, urlSuffix, currentPage-1, urlSuffix)
	}
	if pages_len <= rollPage {
		for _, v := range pages {
			if v == currentPage {
				htmlPage += fmt.Sprintf(`<li class="active"><a href="javascript:void(0);">%d</a></li>`, v)
			} else {
				htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`?page=%d%v">%d</a></li>`, v, urlSuffix, v)
			}
		}

	} else {
		var pageSlice []int
		indexMin := index - rollPage/2
		indexMax := index + rollPage/2
		if indexMin > 0 && indexMax < pages_len { //切片索引未越界
			pageSlice = pages[indexMin:indexMax]
		} else {
			if indexMin < 0 {
				pageSlice = pages[0:rollPage]
			} else if indexMax > pages_len {
				pageSlice = pages[(pages_len - rollPage):pages_len]
			} else {
				pageSlice = pages[indexMin:indexMax]
			}

		}

		for _, v := range pageSlice {
			if v == currentPage {
				htmlPage += fmt.Sprintf(`<li class="active"><a href="javascript:void(0);">%d</a></li>`, v)
			} else {
				htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`?page=%d%v">%d</a></li>`, v, urlSuffix, v)
			}
		}

	}
	if currentPage < totalPage {
		htmlPage += fmt.Sprintf(`<li><a class="num" href="`+path+`?page=%v%v">»</a></li><li><a class="num" href="`+path+`?page=%v%v">..%d</a></li>`, currentPage+1, urlSuffix, totalPage, urlSuffix, totalPage)
	}

	return template.HTML(`<ul class="pagination">` + htmlPage + `</ul>`)
}

func ScoreFloat(score int) string {
	return fmt.Sprintf("%1.1f", float32(score)/10.0)
}

func InMap(maps map[int]bool, key int) (ret bool) {
	if _, ok := maps[key]; ok {
		return true
	}
	return
}

func IncOrDec(table string, field string, condition string, incre bool, step ...int) error {
	mark := "-"
	if incre {
		mark = "+"
	}

	s := 1
	if len(step) > 0 {
		s = step[0]
	}

	sql := fmt.Sprintf("update %v set %v = %v %v %v where %v", table, field, field, mark, s, condition)
	_, err := sysinit.DatabaseEngine.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func Map2struct(m map[string]string, model interface{}) {
	modelType := reflect.TypeOf(model).Elem()
	modelValue := reflect.ValueOf(model).Elem()
	for i := 0; i < modelType.NumField(); i++ {
		jsonTag, ok := modelType.Field(i).Tag.Lookup("json")

		if !ok || jsonTag == "" {
			continue
		}

		if v, ok := m[jsonTag]; ok {
			kind := modelType.Field(i).Type.Kind()
			switch kind {
			case reflect.String:
				modelValue.Field(i).Set(reflect.ValueOf(v))
			case reflect.Int:
				intVal, err :=strconv.Atoi(v)
				if err == nil {
					modelValue.Field(i).Set(reflect.ValueOf(intVal))
				}
			case reflect.Int32:
				intVal, err :=strconv.Atoi(v)
				if err == nil {
					modelValue.Field(i).Set(reflect.ValueOf(intVal))
				}
			case reflect.Float32:
				floatVal, err := strconv.ParseFloat(v, 32)
				if err == nil {
					modelValue.Field(i).Set(reflect.ValueOf(floatVal))
				}
			case reflect.Float64:
				floatVal, err := strconv.ParseFloat(v, 64)
				if err == nil {
					modelValue.Field(i).Set(reflect.ValueOf(floatVal))
				}
			case reflect.Bool:
				boolVal, err := strconv.ParseBool(v)
				if err == nil {
					modelValue.Field(i).Set(reflect.ValueOf(boolVal))
				}
			case reflect.Struct:
				name := modelType.Field(i).Type.Name()
				if name == "Time" {
					timeVal, err := time.ParseInLocation("20060102 15:04:05", v, time.Local)
					if err != nil {
						continue
					}
					modelValue.Field(i).Set(reflect.ValueOf(timeVal))
				}

			default:
			}
		}
	}
}

func DeleteLocalFiles(object ...string) error {
	for _, file := range object {
		os.Remove(strings.TrimLeft(file, "/"))
	}
	return nil
}

func RandomExpire(min int64, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max - min) + min
}