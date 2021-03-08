package ui

const packageJSON = `
{
	"name": "web",
	"version": "0.1.0",
	"private": true,
	"scripts": {
	  "serve": "vue-cli-service serve",
	  "build": "vue-cli-service build",
	  "lint": "vue-cli-service lint"
	},
	"dependencies": {
	"axios": "^0.18.0",
	"popper.js": "^1.14.3",
	"jquery": "^3.3.1",
	"bootstrap": "^4.1.3",
	"element-ui": "^2.4.5",
	"nav-menu": "^1.3.50",
	"font-awesome": "^4.7.0",
	  "core-js": "^3.6.5",
	  "vue": "^2.6.11",
	  "vue-cookies": "^1.5.7",
	  "vue-router": "^3.0.1",
    "vuex": "^3.0.1"
	},
	"devDependencies": {
	  "@vue/cli-plugin-babel": "^4.5.0",
	  "@vue/cli-plugin-eslint": "^4.5.0",
	  "@vue/cli-service": "^4.5.0",
	  "babel-eslint": "^10.1.0",
	  "eslint": "^6.7.2",
	  "eslint-plugin-vue": "^6.2.2",
	  "vue-template-compiler": "^2.6.11"
	},
	"eslintConfig": {
	  "root": true,
	  "env": {
		"node": true
	  },
	  "extends": [
		"plugin:vue/essential",
		"eslint:recommended"
	  ],
	  "parserOptions": {
		"parser": "babel-eslint"
	  },
	  "rules": {}
	},
	"browserslist": [
	  "> 1%",
	  "last 2 versions",
	  "not dead"
	]
  }  
`
