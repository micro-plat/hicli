package ui

const packageJSON = `
{
	"name": "web",
	"version": "0.1.0",
	"private": true,
	"scripts": {
    "dev": "vue-cli-service serve --mode dev",
    "serve": "vue-cli-service serve --mode dev",
    "build": "vue-cli-service build --mode prod",
		"prod": "vue-cli-service build --mode prod",
		"test": "vue-cli-service build --mode test",
		"fat": "vue-cli-service build --mode fat",
		"lint": "vue-cli-service lint"
	},
	"dependencies": {
		"axios": "^0.18.0",
		"bootstrap": "^4.1.3",
		"core-js": "^3.16.2",
		"element-ui": "^2.15.5",
		"font-awesome": "^4.7.0",
		"jquery": "^3.3.1",
		"nav-menu": "^1.3.50",
		"popper.js": "^1.14.3",
		"qxnw-utility": "^1.0.9",
		"vue": "^2.6.11",
		"vue-cookies": "^1.5.7",
		"vue-router": "^3.0.1",
		"vuex": "^3.0.1",
		"xlsx": "^0.16.9"
	},
	"devDependencies": {
		"@vue/cli-plugin-babel": "^4.5.0",
		"@vue/cli-plugin-eslint": "^4.5.0",
		"@vue/cli-service": "^4.5.0",
		"babel-eslint": "^10.1.0",
		"eslint": "^6.7.2",
		"eslint-plugin-vue": "^6.2.2",
		"postcss-px-to-viewport": "^1.1.1",
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
	],
	"postcss": {
		"plugins": {
			"autoprefixer": {},
			"postcss-px-to-viewport": {
				"unitToConvert": "px",
				"viewportWidth": 1920,
				"unitPrecision": 6,
				"propList": [
					"*"
				],
				"viewportUnit": "vw",
				"fontViewportUnit": "vw",
				"selectorBlackList": [
					"wrap"
				],
				"minPixelValue": 1,
				"mediaQuery": true,
				"replace": true,
				"landscape": false
			}
		}
	}
}
`
