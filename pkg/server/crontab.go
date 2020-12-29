package server

import (
	"camera360.com/tv/pkg/remotecontrol"
	"camera360.com/tv/pkg/tv"
	"code.aliyun.com/MIG-server/micro-base/config"
	"code.aliyun.com/MIG-server/micro-base/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pinguo/pgo2/util"
	"github.com/robfig/cron/v3"
	"log"
	"strings"
	"time"
)

var cronIns *cron.Cron
var cronJobs = make(map[string]cron.EntryID)

type CrontabConfig struct {
	OnTime    int64  `json:"ontime"`    //定时开机时间, 时间戳只计算时间部分
	OnDevice  string `json:"ondevice"`  //定时开机设备mac，多个都好分隔
	OnCodes   string `json:"oncodes"`   //开机所执行的红外命令,多个逗号分隔
	OffTime   int64  `json:"offtime"`   //定时关机时间
	OffDevice string `json:"offdevice"` //定时关机设备
	OffCodes  string `json:"offcodes"`  //关机所执行的红外命令, 多个逗号分隔
}

type Crontab struct {
	Second  string `json:"second"`
	Minute  string `json:"minute"`
	Hour    string `json:"hour"`
	Day     string `json:"day"`
	Month   string `json:"month"`
	Week    string `json:"week"`
	Command string `json:"command"`
	Devices string `json:"devices"`
}

func (s *Crontab) Time() string {
	return fmt.Sprintf("%v %v %v %v %v %v", s.Second, s.Minute, s.Hour, s.Day, s.Month, s.Week)
}

func (s *Crontab) GetKey() string {
	v := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v", s.Second, s.Minute, s.Hour, s.Day, s.Month, s.Week, s.Command, s.Devices)
	return util.Md5String(v)
}

func (s *Crontab) sendMessage(device *tv.Device) {
	commands := strings.Split(s.Command, ",")
	app := tv.NewApp("", tv.NewAppNameOption(device.GetPlainObject().AppName))
	var cmd *tv.Command
	modeIds := device.GetPlainObject().ModeId
	mode, _ := remotecontrol.NewModel(context.Background())
	for _, command := range commands {
		logger.Default().Debug("execute command:", command)
		command = strings.TrimSpace(command)
		switch command {
		case "none", "":
			time.Sleep(time.Second * 5)
			continue
		case "on":
			if device.GetPlainObject().RelayTriggeredByLowLevel {
				cmd = tv.NewOffCommand()
			} else {
				cmd = tv.NewOnCommand()
			}
			app.SendMessageToUser(device.GetPlainObject().Mac, cmd)
		case "off":
			if device.GetPlainObject().RelayTriggeredByLowLevel {
				cmd = tv.NewOnCommand()
			} else {
				cmd = tv.NewOffCommand()
			}
			app.SendMessageToUser(device.GetPlainObject().Mac, cmd)
		default:
			for _, modeId := range modeIds {
				mode.Load(modeId)
				btn := mode.GetButtonByCode(command)
				if btn == nil || btn.IrCode == "" {
					continue
				}
				logger.Default().Info("send ir command:", btn.Name, btn.Code, btn.IrCode)
				app.SendMessageToUser(device.GetPlainObject().Mac, tv.NewIrSendCommand(btn.IrCode))
				time.Sleep(time.Second * 2)
			}
		}
	}
}

func (s *Crontab) Run() {
	macs := strings.Split(s.Devices, ",")
	for _, mac := range macs {
		logger.Default().Info("execute cronjob for mac:", mac)
		d, _ := tv.NewDevice(context.Background())
		d.LoadByMac(mac)
		if ! d.HasId() {
			logger.Default().Error("invalid device mac: ", mac)
			continue
		}
		s.sendMessage(d)
	}

}

func RunCronTab() {
	cronIns = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	go cronIns.Run()
	for {
		configJobs := make([]*Crontab, 0)
		newKeys := make(map[string]bool)
		cfg := config.Config().GetString("params.crontab", "[]")
		fmt.Println("cfg", cfg)
		json.Unmarshal([]byte(cfg), &configJobs)
		for _, configJob := range configJobs {
			//已经存在的直接跳到下一个
			if _, ok := cronJobs[configJob.GetKey()]; ok {
				newKeys[configJob.GetKey()] = true
				continue
			}
			entityId, err := cronIns.AddJob(configJob.Time(), configJob)
			if err != nil {
				log.Println("add cronjob failure", err.Error())
				continue
			}
			cronJobs[configJob.GetKey()] = entityId
			newKeys[configJob.GetKey()] = true
		}
		//判断已经运行的cronjob是否在新配置里
		for k, entityId := range cronJobs {
			if _, ok := newKeys[k]; ok {
				continue
			}
			cronIns.Remove(entityId)
		}
		time.Sleep(30 * time.Second)
	}
}
