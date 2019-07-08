layui.define(['cookie'],function(exports){
    var $ = layui.$
        ,cookie = layui.cookie;
    var obj ={
        //获取当前页面的相对路径前缀
        getCurPageRelativePathPrefix:function () {
            var curWwwPath = window.document.location.href;
            var pathName = window.document.location.pathname;
            var pos = curWwwPath.indexOf(pathName);
            var localhostPath = curWwwPath.substring(0, pos);
            return localhostPath;
        },
        GetHtml:function(url){
            var returnHtml="";
            $.ajax({
                type: "GET",
                url: url,
                async:false,
                dataType: "html",
                success: function (data, textStatus) {
                    returnHtml = data;
                },
                error: function (XMLHttpRequest, textStatus, errorThrown) {
                    layer.msg("你没有相关权限");
                }
            });
            return returnHtml;
        },
        IsEnoughPrivilege:function(url){
            var flagObj = new Object();
            $.ajax({
                type: "post",
                data: url,
                async:false,
                timeout:300000,
                url: "/Privilege.do",
                dataType: "json",
                beforeSend:function(xhr){
                    var golayToken = cookie.getCookie("golay_token");
                    if(golayToken==null){
                        flagObj.flag = false;
                        return flagObj;
                    }
                    xhr.setRequestHeader("GolayToken", golayToken);
                },
                success: function (data, textStatus) {
                    if(data!=null && data.StatusCode==200){
                        flagObj.flag = data.Data;
                    }else if(data!=null && data.StatusCode==402){
                        flagObj.message = "登录已过期"
                        flagObj.href = obj.getCurPageRelativePathPrefix() +"/html/login.html";
                    }
                    else {
                        flagObj.message = "你没有相关权限"
                        flagObj.href = obj.getCurPageRelativePathPrefix() +"/html/index.html";
                    }
                },
                error: function (XMLHttpRequest, textStatus, errorThrown) {
                    layer.alert("你没有权限",{icon: 3, title:'提示'},function () {
                        location.href = obj.getCurPageRelativePathPrefix() +"/html/login.html";
                    });
                }
            });
            return flagObj;
        }
    }
    exports('jqUtils',obj);
})












