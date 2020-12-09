const APP_ID = "camera360";

jQuery(document).ready(function () {

    var content = '<div>当前应用: <span data-bind="text:appId"></span></div>\
    <div>当前客户端列表:</div> \
    <ul class="users" data-bind="foreach:users">\
    <li><span data-bind="text:username"></span> <span data-bind="text:wifi"></span> <span data-bind="text:ip"></span> <span data-bind="text:mac"></span> <span data-bind="text:relay"></span> <span data-bind="text:$parent.timeformat(heartbeat_at)"></span></li>\
    </ul>\
    <select data-bind="value:currentDevice,options: userMacs,optionsText:\'label\',optionsValue:\'value\'"></select>\
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

    let userMacs = [
        {
            value: "",
            label: "所有设备",
        }
    ];
    let model = {
        devices: devices,
        appId: APP_ID,
        users: ko.observableArray([]),
        userMacs: ko.observableArray(userMacs),
        currentDevice: "",
        timeformat:function (v) {
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
        jQuery.get("/" + APP_ID + "/users", function (res) {
            console.log(res);
            model.users.splice(0,1000);
            model.userMacs.splice(1,1000);
            console.log(model.users());
            res.map((user)=>{
                console.log(user);
                model.users.push(user);
                model.userMacs.push({
                    value:user.mac,
                    label:user.username,
                })
            })
        })
    }
    getUsers();
    setInterval(getUsers, 5000);

    jQuery(".commands-item").click(function () {
        let url = "/" + APP_ID + "/ir?code=" + jQuery(this).attr("data");
        if (model.currentDevice != "") {
            url = "/" + APP_ID + "/" + model.currentDevice + "/ir?code=" + jQuery(this).attr("data");
        }
        console.log(url)
        jQuery('#loading').show();
        jQuery.get(url, function (res) {
            setTimeout(function () {
                jQuery('#loading').hide();
            }, 500)
        })
    })
    let sendCmd = function (cmd) {
        let url = "/" + APP_ID + "/message?cmd=" + cmd;
        if (model.currentDevice != "") {
            url = "/" + APP_ID + "/" + model.currentDevice + "/message?cmd=" + cmd;
        }
        console.log(url)
        jQuery('#loading').show();
        jQuery.get(url, function (res) {
            getUsers();
            setTimeout(function () {
                jQuery('#loading').hide();
            }, 500)
        })
    }
    jQuery(".on-btn").click(function () {
        sendCmd("on")
    })
    jQuery(".off-btn").click(function () {
        sendCmd("off")
    })

})