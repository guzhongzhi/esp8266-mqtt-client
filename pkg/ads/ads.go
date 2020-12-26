package ads

import (
	"code.aliyun.com/MIG-server/micro-base/orm/mongo"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AdTypeText  = "text"
	AdTypeAlbum = "album"
	AdTypeVideo = "video"
)

type TextInfoConfig struct {
	BackgroundColor string  `json:"backgroundColor" bson:"backgroundColor"` //ARGB
	BackgroundImage string  `json:"backgroundImage" bson:"backgroundImage"` //背景图
	ContentColor    string  `json:"contentColor" bson:"contentColor"`       //ARGB
	ContentSize     float32 `json:"contentSize" bson:"contentSize"`
	Duration        int     `json:"duration" bson:"duration"`
	SignColor       string  `json:"signColor" bson:"signColor"` //ARGB
	SignSize        float32 `json:"signSize" bson:"signSize"`
	TitleColor      string  `json:"titleColor" bson:"titleColor"` //ARGB
	TitleSize       float32 `json:"titleSize" bson:"titleSize"`
}
type TextInfo struct {
	Content string          `json:"content" bson:"content"`
	Sign    string          `json:"sign" bson:"sign"`
	Title   string          `json:"title" bson:"title"`
	Config  *TextInfoConfig `json:"config" bson:"config"`
}

type AlbumInfoConfig struct {
	DisplayMode  string `json:"displayMode" bson:"displayMode"`
	DisplayOrder string `json:"displayOrder" bson:"displayOrder"`
	Duration     int    `json:"duration" bson:"duration"`
}

type AlbumInfo struct {
	Images []string         `json:"images" bson:"images"`
	Config *AlbumInfoConfig `json:"config" bson:"config"`
}

type VideoInfo struct {
	VideoTitle string `json:"videoTitle" bson:"videoTitle"` //标题
	VideoUrl   string `json:"videoUrl" bson:"videoUrl"`     //url地址
}

type AdPO struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Sort      int                `json:"sort" bson:"sort"` //排序
	Title     string             `json:"title" bson:"title"`
	Tvs       []string           `json:"tvs" bson:"tvs"`         //在哪些电视上放,电视Id
	Type      string             `json:"type" bson:"type"`       //广告类型: text,album,video
	Status    string             `json:"status" bson:"status"`   //状态: enabled,disabled
	StartAt   int64              `json:"startAt" bson:"startAt"` //开始时间
	EndAt     int64              `json:"endAt" bson:"endAt"`     //结束时间
	AppName   string             `json:"appName" bson:"appName"`
	TextInfo  *TextInfo          `json:"textInfo" bson:"textInfo"`
	AlbumInfo *AlbumInfo         `json:"albumInfo" bson:"albumInfo"`
	VideoInfo *VideoInfo         `json:"videoInfo" bson:"videoInfo"`
	CreatedAt int64              `json:"createdAt" bson:"createdAt"`
	UpdatedAt int64              `json:"updatedAt" bson:"updatedAt"`
}

func (*AdPO) ValidateMessages() map[string]string {
	panic("implement me")
}

type Ad struct {
	mongo.Object
}

func (s *Ad) PlainObject() *AdPO {
	return s.Data.(*AdPO)
}

func NewAd() *Ad {
	ad := &Ad{}
	ad.initialize(context.Background())
	return ad
}

func (s *Ad) initialize(ctx context.Context) error {
	return s.InitializeWithContext(ctx, &AdPO{}, "Id", "ads", mongo.ConfigNodeNameOption("tvads"))
}
