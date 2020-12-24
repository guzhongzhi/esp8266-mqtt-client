package server

import (
	"camera360.com/tv/pkg/remotecontrol"
	"camera360.com/tv/pkg/tv"
	"code.aliyun.com/MIG-server/micro-base/config"
	"code.aliyun.com/MIG-server/micro-base/logger"
	"code.aliyun.com/MIG-server/micro-base/runtime"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CrontabConfig struct {
	OnTime    int64  `json:"ontime"`    //定时开机时间, 时间戳只计算时间部分
	OnDevice  string `json:"ondevice"`  //定时开机设备mac，多个都好分隔
	OffTime   int64  `json:"offtime"`   //定时关机时间
	OffDevice string `json:"offdevice"` //定时关机设备
}

func RunCronTab() {
	for {
		runCrontab()
		time.Sleep(33 * time.Second)
	}
}

func runCrontab() error {
	cfg := config.Config().GetString("params.timedtask", "")
	fmt.Println("cfg", cfg)
	if cfg == "" {
		return errors.New("invalid config")
	}
	config := &CrontabConfig{}
	err := json.Unmarshal([]byte(cfg), config)
	if err != nil {
		logger.Default().Error("定时任务配置解析失败")
		return errors.New("invalid config")
	}

	if config.OnTime > 0 {
		logger.Default().Info("自动开")
		runOnDevice(config.OnTime, config.OnDevice, tv.NewOnCommand())
	}
	if config.OffTime > 0 {
		logger.Default().Info("自动关")
		runOnDevice(config.OffTime, config.OffDevice, tv.NewOffCommand())
	}
	return nil
}

func runOnDevice(timestamp int64, device string, cmd *tv.Command) {
	now := time.Now()
	t := time.Unix(timestamp/1000, 0)
	t.In(time.UTC)
	if runtime.IsDebug() {
		logger.Default().Info("now.Hour(),now.Minute(),now.Second()", now.Hour(), now.Minute(), now.Second())
		logger.Default().Info("t.Hour(),t.Minute(),t.Second()", t.Month(), "-", t.Day(), " ", t.Hour(), t.Minute(), t.Second())
	}
	if now.Hour() != t.Hour() || now.Minute() != t.Minute() {
		return
	}
	//如果运行了，拖到下一分钟再结束
	defer time.Sleep(time.Second * 60)

	second := t.Second() - now.Second()
	for second > 0 {
		if runtime.IsDebug() {
			logger.Default().Info("second", second)
		}
		time.Sleep(time.Second)
		second -= 1
	}

	macs := strings.Split(device, ",")
	for _, mac := range macs {
		logger.Default().Info("mac", mac)
		d, _ := tv.NewDevice(context.Background())
		d.LoadByMac(mac)
		if ! d.HasId() {
			logger.Default().Error("invalid device mac: ", mac)
			continue
		}
		app := tv.NewApp("", tv.NewAppNameOption(d.GetPlainObject().AppName))
		modeIds := d.GetPlainObject().ModeId
		mode, _ := remotecontrol.NewModel(context.Background())

		app.SendMessageToUser(mac, cmd)
		//开电后延迟5秒再发送红外信号
		time.Sleep(time.Second * 5)

		//如果有继电器,关电后不用再按遥控板
		//@TODO
		//if cmd.IsTurnOff() && d.GetPlainObject().Relay == "on" {
		//continue
		//}

		for _, modeId := range modeIds {
			mode.Load(modeId)
			for _, btn := range mode.GetButtons() {
				if !btn.IsPower() {
					continue
				}
				logger.Default().Info("send ir command:", btn.Name, btn.Code, btn.IrCode)
				app.SendMessageToUser(mac, tv.NewIrSendCommand(btn.IrCode))
			}
		}
		logger.Default().Info("execute crontab", cmd.ToString())
	}
}
