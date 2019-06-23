layui.define(function(exports){
    var obj ={
        //获取当前页面的相对路径前缀
        getCurPageRelativePathPrefix:function () {
            var curWwwPath = window.document.location.href;
            var pathName = window.document.location.pathname;
            var pos = curWwwPath.indexOf(pathName);
            var localhostPath = curWwwPath.substring(0, pos);
            return localhostPath;
        }
    }
    exports('jqUtils',obj);
})












