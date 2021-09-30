package ui

const srcRouterIndexJS = `
import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

const VueRouterPush = Router.prototype.push
Router.prototype.push = function push (to) {
  return VueRouterPush.call(this, to).catch(err => err)
}

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'menus',
      component: () => import('../pages/system/menus.vue'),
      children:[
        // {
        // path: 'index',
        // name: 'index',
        // component: () => import('../pages/system/index.vue'),
        // meta: { title: "首页" }
        // },
      ]
    }
  ]
})

`

const SnippetSrcRouterIndexJS = `
{{- $router:=.router -}}
{{- $ext:=.ext -}}
import Vue from 'vue';
import Router from 'vue-router';

Vue.use(Router);

const VueRouterPush = Router.prototype.push
Router.prototype.push = function push (to) {
  return VueRouterPush.call(this, to).catch(err => err)
}

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'menus',
      component: () => import('../pages/system/menus.vue'),
      children:[
        {{- range $i,$v:=$router}}
        {{- if $v.HasList }}
				{
					path: '{{$v|pathPrefix|rmhd|rpath}}',
					name: '{{$v|pathPrefix|rmhd|varName}}',
					component: () => import('../pages/{{$v|pathPrefix|rmhd|rpath|parentPath}}/{{$v|pathPrefix|rmhd|l2d}}.list.vue')
				},{{- end}}
				{{- if $v.HasDetail }}
				{
					path: '{{$v|pathPrefix|rmhd|rpath}}/detail',
					name: '{{$v|pathPrefix|rmhd|varName}}Detail',
					component: () => import('../pages/{{$v|pathPrefix|rmhd|rpath|parentPath}}/{{$v|pathPrefix|rmhd|l2d}}.detail.vue')
				},{{- end}}
        {{- end}}
        {{- range $i,$v:=$ext}}
        {{- if not $v.Independent}}
				{
					path: '{{$v.Path}}',
					name: '{{$v.Name}}',
					component: () => import('../pages/{{$v.Component}}.vue')
				},
				{{- if $v.HasDetail }}
				{
					path: '{{$v.Path}}/detail',
					name: '{{$v.Name}}Detail',
					component: () => import('../pages/{{$v.Component|trimlist}}.detail.vue')
				},
        {{- end}}
        {{- end}}
        {{- end}}
      ]
    }
    {{- range $i,$v:=$ext}}
    {{- if $v.Independent}},
    {
      path: '{{$v.Path}}',
      name: '{{$v.Name}}',
      component: () => import('../pages/{{$v.Component}}.vue')
    }
    {{- end}}
    {{- end}}
  ]
})

`
