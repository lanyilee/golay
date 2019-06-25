layui.define(function(exports){
    var $ = layui.$;
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
        }
    }
    exports('jqUtils',obj);
})












