package remotecontrol

import (
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ButtonPO struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"` //键名
	AppName   string             `json:"appName" bson:"appName"`
	Code      string             `json:"code" bson:"code"`     //键盘code   多设备一起操作时，根据按钮code查找红外码
	IrCode    string             `json:"irCode" bson:"irCode"` //红外码
	ModeId    string             `json:"modeId" bson:"modeId"` //遥控器
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`
}

func (s *ButtonPO) ValidateMessages() map[string]string {
	return nil
}

func NewButton(ctx context.Context) (*Button, error) {
	b := &Button{}
	b.initialize(ctx)
	return b, nil
}

type Button struct {
	mongo.Object
}

func (s *Button) initialize(ctx context.Context) error {
	return s.InitializeWithContext(ctx, &ButtonPO{}, "Id", "button", mongo.ConfigNodeNameOption("tvads"))
}

func (s *Button) GetPO() *ButtonPO {
	return s.Data.(*ButtonPO)
}
