function WebSocketTest(model) {
    if ("WebSocket" in window) {

    } else {
        console.log("您的浏览器不支持 WebSocket!");
        return
    }

    // 打开一个 web socket
    let wsURL = ""
    if(location.protocol == "https:") {
        wsURL = "wss://" + location.host + "/ws?app=" + APP_ID;
    } else {
        wsURL = "ws://" + location.host + "/ws?app=" + APP_ID;
    }
    var ws = new WebSocket(wsURL);

    ws.onopen = function () {
        ws.send(JSON.stringify({
            operation: "users",
            data: APP_ID,
        }));
    };

    ws.onmessage = function (evt) {
        var data = evt.data;
        try {
            data = JSON.parse(data);
            switch (data.operation) {
                case "users":
                    console.log(data.data);
                    Object.values(data.data).map((user) => {
                        let hasUser = false;
                        model.users().map((u2)=>{
                            if(u2.mac == user.mac) {
                                hasUser = true;
                                model.users.replace(u2,user);
                            }
                        });
                        if(!hasUser) {
                            model.users.push(user);
                        }
                    })

                    break;
            }
        } catch (e) {
            console.log(e);
        }
    };
    ws.onclose = function () {
        // 关闭 websocket
        setTimeout(function (model) {
            WebSocketTest(model);
        }, 5000);
    };
}

jQuery(document).ready(function () {
    var content = '<div>当前应用: <span data-bind="text:appId"></span></div>\
    <div>当前客户端列表:</div> \
    <ul class="users" data-bind="foreach:users">\
    <li style="padding:10px 0px; border: solid 1px #e1e1e1;padding:10px;"><span><input data-bind="textInput: client_id" /></span> <span data-bind="text:wifi"></span> \
        <span data-bind="text:ip"></span> <span data-bind="text:mac"></span> <span data-bind="text:$parent.relayStatus($data)"></span> \
        <div style="padding: 10px;">Relay Pin: <input data-bind="textInput: relay_pin" />\
        自定义RelayPin: <input type="checkbox" data-bind="checked: has_custom_relay_pin" /> <input data-bind="textInput: custom_relay_pin" />\
        <span data-bind="text:$parent.timeformat(refreshed_at)"></span>\
        <span><a href="javascript:void(0)" data-bind="text:$parent.operationText($data), event: { click: $parent.operation}"></a></span>\
        <span><a href="javascript:void(0)" data-bind="event: { click: $parent.save}">保存</a></span>\
        <span><a href="javascript:void(0)" data-bind="event: { click: $parent.setCurrentDevice}">选择</a></span></div>\
        </li>\
    </ul>\
    <div class="currentDevice">当前设备: <span data-bind="text:currentDeviceName"></span><span data-bind="text:currentDevice"></span></div>\
    <div style="margin-top: 10px;" data-bind="if: currentDevice">\
    <ul data-bind="foreach: devices" class="devices"> \
    <li class="device"> \
    <div class="panel panel-default">\n' +
        '  <div class="panel-heading">\n' +
        '    <h3 class="panel-title" ><b data-bind="text: name"></b></h3>\n' +
        '  </div>' +
        '  <div class="panel-body">' + '\
        <span data-bind="foreach: commands" class="commands"> \
        <button  data-bind="attr: {data:value},text:label, click: $root.sendIR"> </button> \
        </span></div> \
        </div></li> \
        </ul></div>';
    jQuery('#content').append(content);
    jQuery('#loading').hide();
    let GlobalModes = {};


    let postJSON = function (url, data) {
        return new Promise((resolve => {
            jQuery.ajax({
                url: url,
                type: 'POST',
                data: typeof(data) == "string" ? data : JSON.stringify(data),
                contentType: 'application/json',
                dataType: 'json',
                success: function (data, status, xhr) {
                    resolve(data)
                },
                error: function (xhr, error, exception) {
                }
            });
        }))
    }

    let sendCmd = function (cmd, mac = null) {
        let url = "/app/" + appName + "/send-message?cmd=" + cmd;
        if (mac) {
            url = "/app/" + appName + "/user/send-message?mac=" + mac + "&cmd=" + cmd;
        }
        postJSON(url,{
            mac:mac,
            cmd:cmd,
        })

    }


    postJSON("/app/"+appName+"/mode/list",{}).then(res=>{
        res.data.items.map(item=>{
            GlobalModes[item.id] = item;
            GlobalModes[item.id].commands = [];
            postJSON("/app/"+appName+"/mode/button-list?modeId="+item.id,{}).then(res=>{
                res.data.items.map(btn=>{
                    GlobalModes[item.id].commands.push({
                        label:btn.name,
                        value:btn.irCode,
                    })
                })
            })
        })
    })
    let model = {
        devices: ko.observableArray(devices),
        appId: appName,
        users: ko.observableArray([]),
        currentDevice: ko.observable(""),
        currentDeviceName:ko.observable(""),
        operation: function (data) {
            let mac = data.mac;
            let relay = data.relay;
            if (relay == "off") {
                sendCmd("on", mac);
            } else {
                sendCmd("off", mac);
            }
        },
        setCurrentDevice(v) {
            console.log(this);
            model.currentDevice(this.mac);
            model.currentDeviceName(this.name);
            model.devices.splice(0,devices.length);
            this.modeId.map(modelId=>{
                if(!GlobalModes[modelId]) {
                    return;
                }
                model.devices.push({
                    name:GlobalModes[modelId].name,
                    commands:GlobalModes[modelId].commands,
                });
            })
        },
        save(v) {
            console.log(v);
            let data = {
                mac:this.mac,
                name:this.client_id,
            }
            postJSON("/app/"+appName+"/user/save",data).then(res=>{
                console.log(res);
            })
        },
        operationText(data) {
            if(data.relayTriggeredByLowLevel) {
                if(data.relay == "off") {
                    return "关闭";
                } else {
                    return "打开";
                }
            } else {
                return data.relay == "off" ? "打开" : "关闭";
            }
        },
        relayStatus(data) {
            let low = "(高电平)";

            if(data.relayTriggeredByLowLevel) {
                low = "(低电平)";
                if(data.relay == "off") {
                    return "on" + low ;
                } else {
                    return "off" + low;
                }
            } else {
                return data.relay + low;
            }
        },
        sendIR() {
            let url = "/app/" + appName + "/send-ir?code=" + this.value;
            if (model.currentDevice() != "") {
                url = "/app/" + APP_ID + "/device-send-ir?mac=" + model.currentDevice() + "&code=" + this.value;
            }
            console.log(url)
            jQuery('#loading').show();
            jQuery.get(url, function (res) {
                setTimeout(function () {
                    jQuery('#loading').hide();
                }, 500)
            })
        },
        timeformat: function (v) {
            let now = new Date(v * 1000);
            let y = now.getFullYear(),
                m = now.getMonth() + 1,
                d = now.getDate();
            return y + "-" + (m < 10 ? "0" + m : m) + "-" + (d < 10 ? "0" + d : d) + " " + now.toTimeString().substr(0, 8);

        }
    }
    ko.applyBindings(model, document.getElementById("content"));

    let getUsers = function () {
        jQuery.getJSON("/app/"+appName+"/users", function (res) {
            model.users.splice(0, 1000);
            res = Object.values(res);
            console.log(res);
            if(res && Array.isArray(res)) {
                res.map((user) => {
                    model.users.push(user);
                })
            } else {
                console.log(res);
            }
        })
    }
    getUsers();
    //setInterval(getUsers, 10000);
    WebSocketTest(model);
    /*
    jQuery(".commands-item").on("click",function () {

    })*/
    jQuery(".on-btn").click(function () {
        sendCmd("on")
    })
    jQuery(".off-btn").click(function () {
        sendCmd("off")
    })

})
