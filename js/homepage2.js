layui.define(['laydate'],function (exports) {
    var laydate = layui.laydate;
    laydate.render({
        elem:'#testTimePlugin',
        type:'datetime'
    });
    exports('homepage2',{});
})