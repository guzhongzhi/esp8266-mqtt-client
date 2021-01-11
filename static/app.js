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
    <li style="padding:10px 0px;"><span><input data-bind="textInput: name" /></span> <span data-bind="text:wifi"></span> \
        <span data-bind="text:ip"></span> <span data-bind="text:mac"></span> <span data-bind="text:relay"></span> \
        <input data-bind="textInput: relayPin" />\
        自定义RelayPin: <input type="checkbox" data-bind="checked: hasCustomRelayPin" /> <input data-bind="textInput: customRelayPin" />\
        <span data-bind="text:$parent.timeformat(heartbeatAt)"></span>\
        <span><a href="javascript:void(0)" data-bind="text:$parent.operationText(relay), event: { click: $parent.operation}"></a></span>\
        <span><a href="javascript:void(0)" data-bind="event: { click: $parent.save}">保存</a></span>\
        <span><a href="javascript:void(0)" data-bind="event: { click: $parent.setCurrentDevice}">选择</a></span>\
        </li>\
    </ul>\
    <div data-bind="text:currentDevice"></div>\
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

    let sendCmd = function (cmd, mac = null) {
        let url = "/app/" + APP_ID + "/send-message?cmd=" + cmd;
        if (mac) {
            url = "/app/" + APP_ID + "/device-send-message?mac=" + mac + "&cmd=" + cmd;
        }
        console.log(url)
        jQuery('#loading').show();
        jQuery.get(url, function (res) {
            setTimeout(() => {
                jQuery('#loading').hide();
            }, 500)
        })
    }

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

    postJSON("/app/"+APP_ID+"/mode/list",{}).then(res=>{
        res.data.items.map(item=>{
            GlobalModes[item.id] = item;
            GlobalModes[item.id].commands = [];
            postJSON("/app/"+APP_ID+"/mode/button-list?modeId="+item.id,{}).then(res=>{
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
        appId: APP_ID,
        users: ko.observableArray([]),
        currentDevice: ko.observable(""),
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
            model.currentDevice(this.mac,this.modeId);
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
            console.log(this);
            postJSON("/app/guz/device-save",this).then(res=>{
                console.log(res);
            })
        },
        operationText(v) {
            return v == "off" ? "打开" : "关闭";
        },
        sendIR() {
            let url = "/app/" + APP_ID + "/send-ir?code=" + this.value;
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
        jQuery.getJSON("/app/"+APP_ID+"/device-list", function (res) {
            model.users.splice(0, 1000);
            res = res.data.items;
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
