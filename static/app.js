function WebSocketTest(model) {
    if ("WebSocket" in window) {

    } else {
        console.log("您的浏览器不支持 WebSocket!");
        return
    }

    // 打开一个 web socket
    var ws = new WebSocket("ws://" + location.host + "/ws?app=" + APP_ID);

    ws.onopen = function () {
        ws.send(JSON.stringify({
            operation: "users",
            data: APP_ID,
        }));
    };

    ws.onmessage = function (evt) {
        var data = evt.data;
        console.log(data)
        try {
            data = JSON.parse(data);
            switch (data.operation) {
                case "users":
                    console.log(data.data);
                    model.users.splice(0, 1000);
                    Object.values(data.data).map((user) => {
                        model.users.push(user);
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
    <li><span data-bind="text:name"></span> <span data-bind="text:wifi"></span> \
        <span data-bind="text:ip"></span> <span data-bind="text:mac"></span> <span data-bind="text:relay"></span> \
        <span data-bind="text:$parent.timeformat(heartbeatAt)"></span>\
        <span><a href="javascript:void(0)" data-bind="text:$parent.operationText(relay), event: { click: $parent.operation}"></a></span>\
        </li>\
    </ul>\
    <div style="margin-top: 10px;">操作</div>\
    <div style="padding:10px 0px;"><a href="javascript:void(0)" class="on-btn">电源开</a> <a href="javascript:void(0)" class="off-btn">电源关</a></div>\
    <ul data-bind="foreach: devices" class="devices"> \
    <li class="device"> \
        <div class="device-title"><b data-bind="text: name"></b></div> \
        <ul data-bind="foreach: commands" class="commands"> \
            <li> \
            <a href="javascript:void(0)" class="commands-item" data-bind="attr: {data:value},text:label"> </a> \
            </li> \
        </ul> \
        <div style="clear:both"></div>\
    </li> \
</ul>'
    jQuery('#content').append(content);
    jQuery('#loading').hide();

    let sendCmd = function (cmd, mac = null) {
        let url = "/app/" + APP_ID + "/message?cmd=" + cmd;
        if (mac) {
            url = "/app/" + APP_ID + "/" + mac + "/message?cmd=" + cmd;
        }
        console.log(url)
        jQuery('#loading').show();
        jQuery.get(url, function (res) {
            setTimeout(() => {
                jQuery('#loading').hide();
            }, 500)
        })
    }
    let model = {
        devices: devices,
        appId: APP_ID,
        users: ko.observableArray([]),
        currentDevice: "",
        operation: function (data) {
            let mac = data.mac;
            let relay = data.relay;
            if (relay == "off") {
                sendCmd("on", mac);
            } else {
                sendCmd("off", mac);
            }
        },
        operationText(v) {
            return v == "off" ? "打开" : "关闭";
        },
        timeformat: function (v) {
            let now = new Date(v * 1000);
            let
                y = now.getFullYear(),
                m = now.getMonth() + 1,
                d = now.getDate();
            return y + "-" + (m < 10 ? "0" + m : m) + "-" + (d < 10 ? "0" + d : d) + " " + now.toTimeString().substr(0, 8);

        }
    }
    ko.applyBindings(model, document.getElementById("content"));

    let getUsers = function () {
        jQuery.get("/app/" + APP_ID + "/users", function (res) {
            model.users.splice(0, 1000);
            res.map((user) => {
                model.users.push(user);
            })
        })
    }
    getUsers();
    //setInterval(getUsers, 10000);
    WebSocketTest(model);
    jQuery(".commands-item").click(function () {
        let url = "/app/" + APP_ID + "/send-ir?code=" + jQuery(this).attr("data");
        if (model.currentDevice != "") {
            url = "/app/" + APP_ID + "/" + model.currentDevice + "/device-send-ir?code=" + jQuery(this).attr("data");
        }
        console.log(url)
        jQuery('#loading').show();
        jQuery.get(url, function (res) {
            setTimeout(function () {
                jQuery('#loading').hide();
            }, 500)
        })
    })
    jQuery(".on-btn").click(function () {
        sendCmd("on")
    })
    jQuery(".off-btn").click(function () {
        sendCmd("off")
    })

})
