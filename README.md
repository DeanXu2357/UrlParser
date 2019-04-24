# UrlParser
Nginx log 檔案解析器

## 緣由
公司要重構舊專案，由於沒有頭緒起頭，所以寫了這個分析現形 nginx 上的 access log ，藉此分析哪些 api 有流量。

* access log 格式 
```
xxx.xxx.xxx.xxx - - [19/Mar/2019:00:00:04 +0800] "POST /event/points.php HTTP/1.1" 200 147 -
xxx.xxx.xxx.xxx - - [19/Mar/2019:00:00:04 +0800] "GET /checkout-1.php?pid=220345&site=m&sp_id=684700 HTTP/1.1" 200 14094 -
xxx.xxx.xxx.xxx - - [19/Mar/2019:00:00:04 +0800] "GET /js/jquery.sticky.js?1540469307 HTTP/1.1" 200 2601 -
xxx.xxx.xxx.xxx - - [19/Mar/2019:00:00:04 +0800] "GET /css/new/responsive.css?1542685352 HTTP/1.1" 200 3455 -
xxx.xxx.xxx.xxx - - [19/Mar/2019:00:00:04 +0800] "GET /img/new/etc/m_gomaji_logo.png HTTP/1.1" 200 4449 -
```
