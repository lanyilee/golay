layui.define(['tree'],function (exports) {
    var tree = layui.tree;
    //渲染
    tree.render({
        elem: '#testTree'  //绑定元素
        ,data: [{
            title: '江西' //一级菜单
            ,children: [{
                title: '南昌' //二级菜单
                ,children: [{
                    title: '高新区' //三级菜单
                    //…… //以此类推，可无限层级
                }]
            }]
        },{
            title: '陕西' //一级菜单
            ,children: [{
                title: '西安' //二级菜单
            }]
        }]
    });
    exports('privileges',{});
})