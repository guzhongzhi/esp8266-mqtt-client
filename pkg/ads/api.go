package ads

import (
	"camera360.com/tv/pkg/controller"
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var lastChanged int64

func CssRBGAToAndroidARGB(v string) string {
	reg := regexp.MustCompile("\\s+")
	v = reg.ReplaceAllString(v, "")
	if v == "" {
		return ""
	}
	v = v[5:]
	v = v[:len(v)-1]

	temp := strings.Split(v, ",")
	if len(temp) != 4 {
		return ""
	}
	var r, g, b int
	var a float64

	r, _ = strconv.Atoi(temp[0])
	g, _ = strconv.Atoi(temp[1])
	b, _ = strconv.Atoi(temp[2])
	a, _ = strconv.ParseFloat(temp[3], 32)
	ai := int(a * 255)

	fmHex := func(v int) string {
		s := fmt.Sprintf("%X", v)
		if len(s) == 1 {
			s = "0" + s
		}
		return s
	}
	hex := fmt.Sprintf("#%s%s%s%s", fmHex(ai), fmHex(r), fmHex(g), fmHex(b))
	return hex
}

type RespItem struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Sort      int                `json:"sort" bson:"sort"` //排序
	Title     string             `json:"title" bson:"title"`
	Type      string             `json:"type" bson:"type"` //广告类型: text,album,video
	TextInfo  *TextInfo          `json:"textInfo,omitempty" bson:"textInfo"`
	AlbumInfo *AlbumInfo         `json:"albumInfo,omitempty" bson:"albumInfo"`
	VideoInfo *VideoInfo         `json:"videoInfo,omitempty" bson:"videoInfo"`
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`
}

func NewApi() *Api {
	c := &Api{}
	return c
}

type Api struct {
	controller.Controller
}

func (c *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.ActionIndex(2)
	c.SetRequest(r)
	c.SetResponse(w)
	c.Dispatch(c)
}

func (c *Api) Index() error {
	ad := NewAd()
	version := c.Int64("version")
	ad.SetData(&RespItem{})
	now := time.Now().Unix()
	col := ad.GetCollection()

	col.Where(mongo.M{
		"status": "enabled",
		"$and": []interface{}{
			mongo.M{"$or": []interface{}{
				mongo.M{"startAt": 0},
				mongo.M{"startAt": mongo.M{"$gt": now}},
			}},
			mongo.M{"$or": []interface{}{
				mongo.M{"endAt": 0},
				mongo.M{"endAt": mongo.M{"$gt": now}},
			}},
		},
	})
	col.Sort(mongo.M{"sort": 1, "createdAt": -1})
	pager, err := col.GetPager(1, 1000)

	for _, item := range pager.Items.([]*RespItem) {
		if item.UpdatedAt > lastChanged {
			lastChanged = item.UpdatedAt
		}
		switch item.Type {
		case AdTypeAlbum:
			item.TextInfo = nil
			item.VideoInfo = nil
		case AdTypeText:
			item.VideoInfo = nil
			item.AlbumInfo = nil
			item.TextInfo.Config.BackgroundColor = CssRBGAToAndroidARGB(item.TextInfo.Config.BackgroundColor)
			item.TextInfo.Config.ContentColor = CssRBGAToAndroidARGB(item.TextInfo.Config.ContentColor)
			item.TextInfo.Config.SignColor = CssRBGAToAndroidARGB(item.TextInfo.Config.SignColor)
			item.TextInfo.Config.TitleColor = CssRBGAToAndroidARGB(item.TextInfo.Config.TitleColor)
		case AdTypeVideo:
			item.AlbumInfo = nil
			item.TextInfo = nil
		}
	}

	if version == lastChanged {
		c.Response().WriteHeader(http.StatusNotModified)
		return nil
	}

	if err != nil {
		return err
	}
	return c.WriteStatusData(mongo.M{
		"version": lastChanged,
		"items":   pager.Items,
	}, http.StatusOK, "OK")
}
