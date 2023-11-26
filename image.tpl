<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<link rel="shortcut icon" href="{{.Image}}.thumb.JPG">
<style>
body{
	background-color: black;
	text-align: center;
}
a{
	color: white;
}
.disabled{
	color: grey;
}
</style>
</head>
<body>
<p>{{if .Prev}}<a href="{{.Prev}}">prev</a>{{else}}<span class="disabled">prev</span>{{end}} | <a href=".">up</a> | {{if .Next}}<a href="{{.Next}}">next</a>{{else}}<span class="disabled">next</span>{{end}}</p>
<p><a href="{{.Image}}.full.JPG"><img src="{{.Image}}.big.JPG"/></a></p>
{{range .Tags}} <a href="#">#{{.}}</a>{{end}}
<form action="/api/tag" method="post">
</form>
</body>
</html>
