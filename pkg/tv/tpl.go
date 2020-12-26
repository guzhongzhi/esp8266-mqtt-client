package tv

const indexHTML = `<html><head>
    <meta charset="utf-8" />
    <title>Daemon</title> 
 <meta name="robots" content="index,follow" />
<script language="javascript">window.APP_ID = "%s"</script>
              <script src='https://libs.baidu.com/jquery/2.0.0/jquery.min.js'></script>
              <script src='/static/knockout.js'></script>
              <script src='/static/config.js'></script>
<!-- 最新版本的 Bootstrap 核心 CSS 文件 -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<!-- 可选的 Bootstrap 主题文件（一般不用引入） -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

<!-- 最新的 Bootstrap 核心 JavaScript 文件 -->
<script src="https://cdn.jsdelivr.net/npm/bootstrap@3.3.7/dist/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
              <link rel='stylesheet' type='text/css' href='/static/main.css'>
              <script src='/static/app.js'></script>
 <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1" />
</head><body>

              <div class='title'><h1>控制中心<span id='loading'>Loading</span></h1></div>
              <div id='content'></div>
</body></html>`
