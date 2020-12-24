package ads

import (
	"camera360.com/tv/pkg/controller"
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hash/crc32"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var lastChanged int64
var lastIds = ""

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

func Crc32IEEE(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func (c *Api) Index() error {
	ad := NewAd()
	version := c.Int64("version")
	ad.SetData(&RespItem{})
	now := time.Now().Unix()
	col := ad.GetCollection()

	where := mongo.M{
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
	}
	e, _ := json.Marshal(where)
	fmt.Println("eeeeeee", string(e))
	col.Where(where)
	sort := primitive.D{}
	sort = append(sort, primitive.E{"sort", 1})
	sort = append(sort, primitive.E{"createdAt", -1})
	col.Sort(sort)
	//mongo.M{"sort": 1, "createdAt": -1}
	pager, err := col.GetPager(1, 1000)

	ids := ""

	for _, item := range pager.Items.([]*RespItem) {
		ids += fmt.Sprintf("%s-%v", item.Id.Hex(), item.UpdatedAt)
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
	if err != nil {
		return err
	}

	v := Crc32IEEE([]byte(ids))
	lastChanged = int64(v)
	lastIds = ids

	if false && lastChanged == version && lastChanged > 0 {
		c.Response().WriteHeader(http.StatusNotModified)
		c.Response().Write([]byte("OK"))
		return nil
	}

	return c.WriteStatusData(mongo.M{
		"version": v,
		"items":   pager.Items,
	}, http.StatusOK, "OK")
}
