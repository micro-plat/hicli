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
				{
					path: '{{$v.Name|rmhd|rpath}}',
					name: '{{$v.Name|rmhd|varName}}',
					component: () => import('../pages/{{$v.Name|rmhd|rpath|parentPath}}/{{$v.Name|rmhd|l2d}}.list.vue')
				},
				{{- if $v.HasDetail }}
				{
					path: '{{$v.Name|rmhd|rpath}}/detail',
					name: '{{$v.Name|rmhd|varName}}Detail',
					component: () => import('../pages/{{$v.Name|rmhd|rpath|parentPath}}/{{$v.Name|rmhd|l2d}}.detail.vue')
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
