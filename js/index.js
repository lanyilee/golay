/**
 项目JS主入口
 以依赖layui的layer和form模块
 **/

layui.define(['layer','cookie'],function (exports) {
    var $ = layui.$,
        cookie = layui.cookie,
        layer = layui.layer;
    $.ajax({
        type:"POST",
        async:false,
        data:null,
        dataType:"json",
        url:"/GetMenu.do",
        beforeSend:function(xhr){
            var golayToken = cookie.getCookie("golay_token");
            if(golayToken==null){
                return false;
            }
            xhr.setRequestHeader("GolayToken", golayToken);
        },
        success: function (data, textStatus) {
            debugger
            if (data.StatusCode!=200){
                layer.alert(data.Message)
            }else{
                layer.alert(123)
            }
            //returnHtml = data;
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            layer.msg("获取菜单列表失败");
        }


    });
    exports('index',{});
})