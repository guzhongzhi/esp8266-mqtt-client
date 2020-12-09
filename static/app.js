const APP_ID = "camera360";

jQuery(document).ready(function () {

    var content = '<div>当前应用: <span data-bind="text:appId"></span></div>\
    <div>当前客户端列表:</div> \
    <ul data-bind="foreach:users">\
    <li><span data-bind="text:mac"></span> <span data-bind="text:mac"></span> <span data-bind="text:relay"></span></li>\
    </ul>\
    <select data-bind="value:currentDevice,options: userMacs,optionsText:\'label\',optionsValue:\'value\'"></select>\
    <div>操作</div>\
    <div><a href="javascript:void(0)" class="on-btn">电源开</a> <a href="javascript:void(0)" class="off-btn">电源关</a></div>\
    <ul data-bind="foreach: devices" class="devices"> \
    <li class="device"> \
        <b data-bind="text: name"></b> \
        <ul data-bind="foreach: commands" class="commands"> \
            <li> \
            <a href="javascript:void(0)" class="commands-item" data-bind="attr: {data:value},text:label"> </a> \
            </li> \
        </ul> \
    </li> \
</ul>'
    jQuery('#content').append(content);
    jQuery('#loading').hide();

    let users = [
        {
            mac: "EEEE",
            relay: "off",
        },
        {
            mac: "FFFFF",
            relay: "off",
        }
    ];
    let userMacs = [
        {
            value: "",
            label: "所有设备",
        },
        {
            value: "EEEE",
            label: "B3",
        }
    ];
    let model = {
        devices: devices,
        appId: APP_ID,
        users: users,
        userMacs: userMacs,
        currentDevice: "",
    }
    ko.applyBindings(model, document.getElementById("content"));

    let getUsers = function () {
        jQuery.get("/" + APP_ID + "/users", function (res) {
            console.log(res);
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
            jQuery('#loading').hide();
        })
    })
    let sendCmd = function (cmd) {
        let url = "/" + APP_ID + "/message?cmd=off";
        if (model.currentDevice != "") {
            url = "/" + APP_ID + "/" + model.currentDevice + "/message?cmd=off";
        }
        console.log(url)
        jQuery('#loading').show();
        jQuery.get(url, function (res) {
            jQuery('#loading').hide();
        })
    }
    jQuery(".on-btn").click(function () {
        sendCmd("on")
    })
    jQuery(".off-btn").click(function () {
        sendCmd("off")
    })

})