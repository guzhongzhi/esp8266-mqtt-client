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