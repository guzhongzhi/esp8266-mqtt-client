package remotecontrol

import (
	"camera360.com/tv/pkg/controller"
	"code.aliyun.com/MIG-server/micro-base/model"
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"context"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func NewControl() *Control {
	ctl := &Control{}
	ctl.ActionIndex(4)
	return ctl
}

type Control struct {
	controller.Controller
	AppName string
}

func (c *Control) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c.AppName = vars["appId"]
	c.SetRequest(r)
	c.SetResponse(w)
	c.Dispatch(c)
}

func (c *Control) ButtonSave() error {
	btn, err := NewButton(c.Request().Context())
	if err != nil {
		return err
	}
	id := c.String("id")
	modeId := c.String("modeId")
	code := c.String("code")
	name := c.String("name")
	btn.Load(id)
	btn.GetPO().ModeId = modeId
	btn.GetPO().Code = code
	btn.GetPO().Name = name
	btn.Save()
	return nil
}

func (c *Control) ButtonDelete() error {
	btn, err := NewButton(c.Request().Context())
	if err != nil {
		return err
	}
	id := c.String("id")
	btn.Load(id)
	btn.Delete()
	c.WriteStatusData(true, http.StatusOK, "OK")
	return nil
}

func (c *Control) ButtonList() error {
	btn, err := NewButton(c.Request().Context())
	if err != nil {
		return err
	}
	modeId := c.String("modeId")
	log.Println("modeId", modeId)
	var pager *model.Pager
	if modeId != "" {
		pager, err = btn.GetCollection().Where(mongo.M{
			"appName": c.AppName,
			"modeId":  modeId,
		}).GetPager(1, 10000)
	} else {
		pager, err = btn.GetCollection().Where(mongo.M{
			"appName": c.AppName,
			"modeId":  "",
		}).GetPager(1, 10000)
	}
	if err != nil {
		return err
	}
	c.WriteStatusData(pager, http.StatusOK, "OK")
	return nil
}

func (c *Control) Info() error {
	modeId := c.String("modeId")
	mod, err := NewModel(context.Background())
	if err != nil {
		log.Println("err", err)
	}
	mod.Load(modeId)
	return c.WriteJSON(mongo.M{
		"buttons": mod.GetButtons(),
		"mode":    mod.Data,
	})
}

func (c *Control) List() error {
	mod, err := NewModel(context.Background())
	if err != nil {
		log.Println("err", err)
	}
	pager, err := mod.GetCollection().Where(mongo.M{
		"appName": c.AppName,
	}).GetPager(1, 100)
	if err != nil {
		log.Println("err", err)
	}
	c.WriteStatusData(pager, http.StatusOK, "OK")
	return nil
}

func (c *Control) Delete() error {
	id := c.String("id")
	mod, _ := NewModel(c.Request().Context())
	mod.Load(id)
	mod.Delete()
	return nil
}

func (c *Control) Save() error {
	name := c.String("name")
	id := c.String("id")
	if name == "" {
		return errors.New("name can not be null")
	}
	mod, err := NewModel(c.Request().Context())
	if id != "" {
		mod.Load(id)
	}
	if err != nil {
		return err
	}

	if !mod.HasId() {
		mod.LoadByName(name, c.AppName)
	}
	if !mod.HasId() {
		mod.GetPO().CreatedAt = time.Now().Unix()
	}
	mod.GetPO().Name = name
	mod.GetPO().AppName = c.AppName
	mod.GetPO().UpdatedAt = time.Now().Unix()
	mod.Save()
	c.WriteStatusData(mod.GetPO(), http.StatusOK, "OK")
	return nil
}
