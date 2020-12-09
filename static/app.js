const APP_ID = "camera360";

jQuery(document).ready(function () {

    var content = '<div>当前应用<span data-bind="text:appId"></span></div>\
    <div>当前客户端</div> \
    <ul data-bind="foreach:users">\
    <li><span data-bind="text:mac"></span> <span data-bind="text:mac"></span> <span data-bind="text:relay"></span></li>\
    </ul>\
    <select data-bind="options: userMacs,optionsText:\'label\',optionsValue:\'value\'"></select>\
    <div>操作</div>\
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
        },
        {
            mac: "FFFFF",
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
    ko.applyBindings({
        devices: devices,
        appId: APP_ID,
        users: users,
        userMacs: userMacs,
    }, document.getElementById("content"));


    jQuery(".commands-item").click(function () {
        jQuery('#loading').show();
        jQuery.get("/" + APP_ID + "/ir?code=" + jQuery(this).attr("data"), function (res) {
            console.log(res);
            jQuery('#loading').hide();
        })

    })
})