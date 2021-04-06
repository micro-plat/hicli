package ui

const vueConfigJS = `

module.exports = {
    runtimeCompiler: true,

    chainWebpack: config => {
        // 移除 prefetch 插件
        config.plugins.delete('prefetch')
    },

    publicPath: process.env.NODE_ENV === 'production' ? '/' : '/',

    // outputDir: 在npm run build时 生成文件的目录 type:string, default:'dist'
    outputDir: 'dist',

    //   lintOnSave：{ type:Boolean default:true } 问你是否使用eslint
    lintOnSave: false,

    // productionSourceMap：{ type:Bollean,default:true } 生产源映射
    // 如果您不需要生产时的源映射，那么将此设置为false可以加速生产构建
    productionSourceMap: false,

    // devServer:{type:Object} 3个属性host,port,https
    // 它支持webPack-dev-server的所有选项

    devServer: {
        port: 8080, // 端口号
        host: '0.0.0.0',
        https: false, // https:{type:Boolean}
        // open: true, //配置自动启动浏览器
    }
}`
