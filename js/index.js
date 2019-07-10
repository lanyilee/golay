/**
 项目JS主入口
 以依赖layui的layer和form模块
 **/

layui.define(['layer','cookie','jqUtils'],function (exports) {
    var $ = layui.$,
        cookie = layui.cookie,
        jqUtils = layui.jqUtils,
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
            if (data.StatusCode == 200){
                var items = data.Data;
                var appendHtml = '';
                $.each(items,function (index,item) {
                    appendHtml =  appendHtml + '<li class="layui-nav-item layui-nav-itemed"><a class="" href="javascript:;">'+item.Name+' <span class="layui-nav-more"></span></a>';
                    if (item.LeftMenu.length>0){
                        appendHtml = appendHtml+ '<dl class="layui-nav-child">';
                        $.each(item.LeftMenu,function (num,menu) {
                            appendHtml = appendHtml+ '<dd><a href="javascript:;" lay-href="' + menu.Redirecturl + '">' + menu.Name + '</a></dd>';
                        })
                        appendHtml = appendHtml+ '</dl></li>';
                    }
                })
                $("#LeftMenuUL").html(appendHtml);
                $("#LoginRealName").text(data.Message);
                layui.use('element', function() {
                    var element = layui.element;
                    element.init();
                });
            }else{
                layer.alert("登录已过期",{icon: 3, title:'提示'},function () {
                    location.href = jqUtils.getCurPageRelativePathPrefix() +"/html/login.html";
                });
            }
            //returnHtml = data;
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            layer.msg("获取菜单列表失败");
        }


    });
    exports('index',{});
})