#MQTT WEB API SERVER

##板子启动后需要发布 application-init 事件同步clientId初始化应用

事件值为: clientId=appId+"-"+随机值


## HTTP接口

#### 用户列表
/appId/users

#### 发送红外
/{appId}/ir

#### 发送红外到指定用户
/{appId}/mac/ir

#### 发送应用全局消息
/{appId}/message

#ESP8266 mqtt 消息

//executedAt 为消息id, 或执行时间，设备收到消息在执行后会在心跳中加上executedAt  以及data为最后一次执行成功的cmd
{"cmd":"on","executedAt":19939838} //on表示高电平, off表示低电平
{"cmd":"off","executedAt":19939838}
{"cmd":"setRelayPIN","executedAt":19939838,data:0}  //设置继电器引脚，设置后默认会是低电平
{"cmd":"setPinLow","executedAt":19939838,"data":12}
{"cmd":"serialSendHexStringArray","executedAt":19939838,"data":["A1","F1","0E","0E","0C"]}  //创维开机
{"cmd":"serialSendHexStringArray","executedAt":19939838,"data":["A1","F1","0E","0E","14"]}  //创维音量+
{"cmd":"serialSendHexStringArray","executedAt":19939838,"data":["A1","F1","0E","0E","15"]}  //创维音量-
{"cmd":"setPinHigh","executedAt":19939838,"data":12}
{"cmd":"irs","executedAt":19939838,"data":"9126,4344,668,476,640,504,640,504,752,388,640,500,740,406,738,1508,726,416,640,1604,726,1514,666,1576,640,1602,726,1510,642,1604,724,414,642,1602,726,418,614,1630,640,504,738,408,738,1506,640,502,736,412,726,414,738,1508,726,414,728,1512,726,1516,726,416,640,1602,666,1576,728,1512,726,40064,9126,2096,726"}