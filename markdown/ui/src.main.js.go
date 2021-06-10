package ui

const srcMainJS = `

import "jquery"
import "bootstrap"
 
import Vue from 'vue'
import App from './App'
import router from './router'

import VueCookies from 'vue-cookies'
Vue.use(VueCookies);

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
Vue.use(ElementUI);

import utility from 'qxnw-utility';
Vue.use(utility,false);

Vue.config.productionTip = false;

router.beforeEach((to, from, next) => {
  /* 路由发生变化修改页面title */
  Vue.prototype.$sys.checkAuthCode(to)
  if (to.path != "/") {
      document.title = Vue.prototype.$sys.getTitle(to.path)
  }
  next()
})

  /* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    components: {
        App
    },
    template: '<App/>'
});

`

const srcSSOMainJS = `

import "jquery"
import "bootstrap"
 
import Vue from 'vue'
import App from './App'
import router from './router'

import VueCookies from 'vue-cookies'
Vue.use(VueCookies);

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
Vue.use(ElementUI);

import utility from 'qxnw-utility';
Vue.use(utility);

Vue.config.productionTip = false;

router.beforeEach((to, from, next) => {
  /* 路由发生变化修改页面title */
  Vue.prototype.$sys.checkAuthCode(to)
  if (to.path != "/") {
      document.title = Vue.prototype.$sys.getTitle(to.path)
  }
  next()
})

  /* eslint-disable no-new */
new Vue({
    el: '#app',
    router,
    components: {
        App
    },
    template: '<App/>'
});

`
