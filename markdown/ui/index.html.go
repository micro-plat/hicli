package ui

const publicIndexHTML = `
<!DOCTYPE html>
<html lang="">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1,minimum-scale=1">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
  </head>
  <body>
    <noscript>
      <strong>We're sorry but <%= htmlWebpackPlugin.options.title %> doesn't work properly without JavaScript enabled. Please enable it to continue.</strong>
    </noscript>
    <div id="app"></div>
    <!-- built files will be auto injected -->
  </body>
</html>

<style>
  .app-content-body{
    adding-bottom: 0px !important;
  }
  .panel{
    margin-bottom: 0px !important;
  }
</style>

<script type="text/javascript">

    var changeViewPort = function () {
      var viewportmeta = document.querySelector('meta[name="viewport"]')
      if (!viewportmeta) {
        return
      }
      var winW = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth
      var scale = winW / 1680
      scale = Math.round(scale * 100) / 100
      if (scale <= 1) {
        viewportmeta.content = 'width=device-width,initial-scale=' + scale + ',maximum-scale=' + scale + ',minimum-scale=' + scale;
      }
    }

    //调整页面尺寸
    changeViewPort()   
  </script>
`
