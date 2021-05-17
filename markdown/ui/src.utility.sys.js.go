package ui

const srcUtilitySysJS = `var __user_info__="__user_info__"
export function Sys(Vue) {
    Sys.prototype.Vue = Vue
}

//checkAuthCode 向服务器发送请求，验证auth code
Sys.prototype.checkAuthCode = function (router, url){
    //检查请求参数中是否有code
    if (!router.query.code){
        return
    }
    
    //从服务器拉取数据
    let that = Sys.prototype.Vue.prototype;

    //检查verify地址
    let sys = that.$env.conf.system || {}
    var verifyURL = url || sys.verifyURL || "/login/verify";

    var userInfo = that.$http.xget(verifyURL, {code: router.query.code})
    if (Object.getOwnPropertyNames(userInfo).length == 0){
        throw new Error("userInfo数据为空");
    }
    //保存用户信息
    window.localStorage.setItem(__user_info__, JSON.stringify(userInfo));
}

//lognout 退出登录
Sys.prototype.logout = function(url){
    let that = Sys.prototype.Vue.prototype;

    //清除http认证头信息及cookie
    that.$http.clearAuthorization();

    //清除cookie getSystemInfo
    let destURL = url || that.$env.conf.system.logoutDestURL;

    var redirctURL= "?returnurl=" + encodeURIComponent(window.location.href);
    if(destURL){
        redirctURL = "?desturl=" + encodeURIComponent(destURL);
    }
    //检查logoutURL是否配置
    if (that.$env.conf.system.logoutURL){
        window.location = that.$env.conf.system.logoutURL+redirctURL;
    }

}

//changePwd 修改密码
Sys.prototype.changePwd = function(url){
    let that = Sys.prototype.Vue.prototype;
    //跳转到修改密码页面
    let changepwdURL = url || that.$env.conf.system.changepwdURL;

    if (changepwdURL){
        window.location.href = changepwdURL
    }  
}

//getUserInfo 获取用户信息
Sys.prototype.getUserInfo = function(){
   let userInfo = window.localStorage.getItem(__user_info__)  
   if (!userInfo){
       return {}
   }
   return JSON.parse(userInfo)
}

//根据路由获取标题
Sys.prototype.getTitle = function(path){
    
    //获取本地配置的菜单
    let that = Sys.prototype.Vue.prototype  
    var menus = that.$env.conf.menus
   
    //根据路径查找名称
    var cur = Sys.prototype.findMenuItem(menus, path)
    if (that.$env.conf.system && that.$env.conf.system.name){
        return cur ? cur.name + " - " + that.$env.conf.system.name : "";
    } 
    return cur ? cur.name : "";
}

//递归查找父级菜单
Sys.prototype.findMenuItem = function(menus, path){   
    path = path || window.location.pathname;

    if(path != "/" && path[path.length - 1] == '/'){
        path = path.substring(0, path.length - 1);
    }
    //递归查找父级菜单
    var cur = getMenuItem(menus, path);
    if (!cur){                
        path = path.substring(0, path.lastIndexOf('/'));
        if (!path){
            return getMenuItem(menus, '/')
        }
        cur = getMenuItem(menus, path)
    }
    return cur || getMenuItem(menus, "/")
}

//getMenus获取菜单数据
Sys.prototype.getMenus = function(url){  
    return new Promise((resolve, reject) => {
        let that = Sys.prototype.Vue.prototype  
        //获取本地配置的菜单
        let menus = that.$env.conf.menus
        if (menus){
            let data = typeof menus == "string" ? that.$http.xget(menus) || [] : menus
            resolve(data)
            return
        }       

        //远程获取菜单
        let menuURL = url || "/member/menus/get"
        that.$http.get(menuURL)
        .then(res => {             
            //保存菜单信息
            Object.assign(that.$env.conf, { menus: res || [] }) 
            resolve(res);
        })
        .catch(err => {
            reject(err)
        })
    });
}

//getSystemInfo获取系统信息
Sys.prototype.getSystemInfo = function(url ){   
    return new Promise((resolve, reject) => {   
        let that = Sys.prototype.Vue.prototype 
        //获取本地配置的系统信息
        if (that.$env.conf.system){
            resolve(that.$env.conf.system)
            return
        }  

        //获取远程系统信息
        let systemInfoURL = url || "/system/info/get"
        that.$http.get(systemInfoURL)
        .then(res => {             
            //保存系统信息
            Object.assign(that.$env.conf, { system: res || {} }) 
            resolve(res);
        })
        .catch(err => {
            reject(err)
        })
    });
}

//getSystemList获取用户系统列表
Sys.prototype.getSystemList = function(url ){      
    return new Promise((resolve, reject) => {
        let that = Sys.prototype.Vue.prototype  
        //获取本地配置的菜单
        if (that.$env.conf.sysList){
            resolve(that.$env.conf.sysList)
            return
        }  

        //远程获取其它系统
        let systemsListURL = url || "/member/systems/get"
        that.$http.get(systemsListURL)
        .then(res => {             
            Object.assign(that.$env.conf, { sysList: res || [] }) 
            resolve(res);
        })
        .catch(err => {
            reject(err)
        })
    });
}

//获取路由name
function getMenuItem(menus, path){    
    if(!menus || !menus.length){
        return null;
    }
    for (var i in menus){
        var cur = menus[i];
        if(cur.path == path){
            return cur;
        }
        if(path == "/" && cur.path && cur.path != "-"){
            return cur;
        }        
        var res = getMenuItem(cur.children || [], path);
        if(res){
            return res;
        }
    }
    return null;
}
`
