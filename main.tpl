<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<meta http-equiv="content-type" content="text/html; charset=utf-8">
<style>
body{
	background-color: black;
}
a{
	color: white;
}
div{
	float: left;
	text-align: center;
	padding: 0.25cm;
	width: 200px;
	height: 200px;
}
img{
	width: 160px;
	height: 160px;
}
</style>
</head>
<body>
<p>&nbsp;</p>
{{range .Years}}
<div><a href="{{.}}/index.html"><img src="{{.}}/montage.jpg"/><p>{{.}}</p></a></div>
{{end}}
{{range .Albums}}
<div><a href="{{.}}/index.html"><img src="{{.}}/montage.jpg"/><p>{{.}}</p></a></div>
{{end}}
<hr style="clear: both;" />
{{range .Tags}}
<a href="/tags/{{.}}">#{{.}}</a>
{{end}}
</body>
</html>
