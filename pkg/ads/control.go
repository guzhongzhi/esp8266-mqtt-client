package ads

import (
	"camera360.com/tv/pkg/controller"
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Control struct {
	controller.Controller
	AppName string
}

func NewController() *Control {
	c := &Control{}
	return c
}

func (c *Control) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c.AppName = vars["appId"]
	c.ActionIndex(4)
	c.SetRequest(r)
	c.SetResponse(w)
	c.Dispatch(c)
}

func (c *Control) Save() error {
	ad := NewAd()
	request := c.loadRequest()
	if request.Id.Hex() != "" {
		ad.Load(request.Id)
	}
	if ad.HasId() {
		ad.PlainObject().VideoInfo = request.VideoInfo
		ad.PlainObject().TextInfo = request.TextInfo
		ad.PlainObject().AlbumInfo = request.AlbumInfo
		ad.PlainObject().Title = request.Title
		ad.PlainObject().StartAt = request.StartAt
		ad.PlainObject().EndAt = request.EndAt
		ad.PlainObject().Sort = request.Sort
		ad.PlainObject().Status = request.Status
	} else {
		ad.SetData(request)
		request.AppName = c.AppName
		request.CreatedAt = time.Now().Unix()
	}
	ad.PlainObject().UpdatedAt = time.Now().Unix()
	ad.Save()
	c.WriteStatusData(ad.GetData(), http.StatusOK, "OK")
	return nil
}

func (c *Control) SetStatus() error {
	id := c.String("id")
	status := c.String("status")
	ad := NewAd()
	ad.Load(id)
	ad.PlainObject().Status = status
	ad.PlainObject().UpdatedAt = time.Now().Unix()
	ad.Save()
	return c.WriteStatusData(ad.GetData(), http.StatusOK, "OK")
}

func (c *Control) Delete() error {
	id := c.String("id")
	ad := NewAd()
	ad.Load(id)
	ad.Delete()
	c.WriteStatusData(ad.GetData(), http.StatusOK, "OK")
	return nil
}

func (c *Control) Info() error {
	id := c.String("id")
	ad := NewAd()
	ad.Load(id)
	return c.WriteStatusData(ad.Data, http.StatusOK, "OK")
}

func (c *Control) List() error {
	ad := NewAd()
	collection := ad.GetCollection()
	collection.Sort(mongo.M{"createdAt": -1})
	pager, err := collection.Where(mongo.M{"appName": c.AppName}).GetPager(1, 1000)
	if err != nil {
		return err
	}
	c.WriteStatusData(pager, http.StatusOK, "OK")
	return nil
}

func (c *Control) loadRequest() *AdPO {
	r := &AdPO{}
	body, err := ioutil.ReadAll(c.Request().Body)
	log.Println("err", err, string(body))
	err = json.Unmarshal(body, r)
	log.Println("err", err, string(body))

	s, _ := json.Marshal(r)

	log.Println(string(s))
	return r
}
