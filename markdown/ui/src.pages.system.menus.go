package ui

const srcPagesSystemMenus = `
<template>
  <div id="app">
    <nav-menu
      :menus="menus"
      :copyright="copyright"
      :copyrightcode="copyrightcode"
      :themes="system.themes"
      :logo="system.logo"
      :systemName="system.systemName"
      :userinfo="userinfo"
      :items="items"
      :pwd="pwd"
      :signOut="signOutM"
      ref="NewTap"
    >
    </nav-menu>
  </div>
</template>

<script>
import navMenu from 'nav-menu'; // 引入
export default {
  name: 'app',
  data() {
    return {
      system: {
        logo: "",
        systemName: "",  //系统名称
        themes: "", //顶部左侧背景颜色,顶部右侧背景颜色,右边菜单背景颜色
      },
      copyright: (this.$env.conf.copyright.company || "") + "Copyright©" + new Date().getFullYear() + "版权所有",
      copyrightcode: this.$env.conf.copyright.code,
      menus: [{}],  //菜单数据
      userinfo: {},
      items: []
    }
  },
  components: { //注册插件
    navMenu
  },
  created() {

  },
  mounted() {
    this.getMenu();
    this.getSystemInfo();
    this.userinfo = this.$sys.getUserInfo()
  },
  methods: {
    pwd() {    
      this.$sys.changePwd()
    },
    signOutM() {
      this.$sys.logout();
    },
    getMenu() {
      this.$sys.getMenus().then(res => {
        this.menus = res;
        this.getUserOtherSys();
        var cur = this.$sys.findMenuItem(res);
        this.$refs.NewTap.open(cur.name, cur.path);
      });
    },
    //获取系统的相关数据
    getSystemInfo() {
      let that = this
      this.$sys.getSystemInfo().then(res => {
        this.system = res;
        document.title = that.$sys.getTitle(window.location.path)
      })
    },
    //用户可用的其他系统
    getUserOtherSys() {
      this.$sys.getSystemList().then(res => {
        this.items = res;
      })
    },
  }
}
</script>
`
