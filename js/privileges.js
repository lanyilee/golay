layui.define(['tree'],function (exports) {
    var tree = layui.tree,
        $ = layui.$;
    // 请求数据
    $.ajax({
        type: "post",
        data: "",
        async:false,
        timeout:300000,
        url: "/GetConfigPrivileges.do",
        dataType: "json",
        beforeSend:function(xhr){
        },
        success: function (data, textStatus) {
            var treeDatas = [];
            function SetTree(PrivilegeData) {
                if(PrivilegeData==null){
                    return null
                }
                var tPrivilege = new Object();
                tPrivilege.id = PrivilegeData.Id;
                tPrivilege.title = PrivilegeData.Name;
                tPrivilege.href = PrivilegeData.Selector;
                if(PrivilegeData.Privilege!=null){
                    tPrivilege.children = [];
                    for(let tchild of PrivilegeData.Privilege){
                        var cdata = SetTree(tchild);
                        tPrivilege.children.push(cdata);
                    }
                }
                return tPrivilege;
            }
            if(data==null||data.StatusCode!=200){
                layui.msg("获取出错");
                return
            }else{
                for(let tchild of data.Data){
                    var cdata = SetTree(tchild);
                    treeDatas.push(cdata)
                }
            }
            //渲染
            tree.render({
                //绑定元素
                elem: '#testTree',
                data:treeDatas,
                showCheckbox: true
            });

        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            layer.msg("你没有相关权限");
            return false;
        }
    });
    exports('privileges',{});

})