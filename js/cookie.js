layui.define(function (exports) {
    var obj={
        setCookie:function(name,value,time)
        {
            var strsec = obj.getsec(time);
            var exp = new Date();
            exp.setTime(exp.getTime() + strsec*1);
            document.cookie = name + "="+ value + ";expires=" + exp.toGMTString();
        },
        getsec:function(str)
        {
            var str1=str.substring(1,str.length)*1;
            var str2=str.substring(0,1);
            if (str2=="s")
            {
                return str1*1000;
            }
            else if (str2=="h")
            {
                return str1*60*60*1000;
            }
            else if (str2=="d")
            {
                return str1*24*60*60*1000;
            }
        },
        getCookie:function(name)
        {
            var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");
            if(arr=document.cookie.match(reg))
                return unescape(arr[2]);
            else
                return null;
        },
        delCookie:function(name)
        {
            var exp = new Date();
            exp.setTime(exp.getTime() - 1);
            var cval=obj.getCookie(name);
            if(cval!=null)
                document.cookie= name + "="+cval+";expires="+exp.toGMTString();
        }
    }

    exports('cookie',obj);
})