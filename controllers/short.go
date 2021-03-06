package controllers

import (
	"shorturl/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
)

var (
	urlcache cache.Cache
)

func init() {
	urlcache, _ = cache.NewCache("memory", `{"interval:0"}`)
}

type ShortResult struct {
	Code int
	Msg  string
	Data struct {
		UrlShort string
		UrlLong  string
	}
}

type ShortController struct {
	beego.Controller
}

func (_this *ShortController) ToShort() {
	var result ShortResult
	var err error
	longurl := _this.Input().Get("longurl")
	beego.Info(longurl)
	result.Data.UrlLong = longurl
	md5Url := models.GetMD5(longurl)
	beego.Info(md5Url)

	if urlcache.IsExist(md5Url) {
		if value, ok := urlcache.Get(md5Url).(string); ok {
			result.Data.UrlShort = value
		}
	} else {
		result.Data.UrlShort = models.Generate()
		err = urlcache.Put(md5Url, result.Data.UrlShort, 0)
		if err != nil {
			beego.Info(err)
		}
		err = urlcache.Put(result.Data.UrlShort, longurl, 0)
		if err != nil {
			beego.Info(err)
		}
	}
	_this.Data["json"] = result
	_this.ServeJSON()
}
