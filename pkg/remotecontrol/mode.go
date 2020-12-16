package remotecontrol

import (
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//遥控器型号
type ModePO struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	AppName   string             `json:"appName" bson:"appName"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`
}

func (s *ModePO) ValidateMessages() map[string]string {
	return nil
}

func NewModel(ctx context.Context) (*Mode, error) {
	m := &Mode{}
	m.initialize(ctx)
	return m, nil
}

type Mode struct {
	mongo.Object
}

func (s *Mode) initialize(ctx context.Context) error {
	return s.InitializeWithContext(ctx, &ModePO{}, "Id", "mode", mongo.ConfigNodeNameOption("tvads"))
}
func (s *Mode) GetPO() *ModePO {
	return s.Data.(*ModePO)
}
