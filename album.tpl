<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<link rel="shortcut icon" href="montage.jpg">
<style>
body{
	background-color: black;
	text-align: center;
}
img{
	height: 96px;
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
<p>
{{if .Prev}}
<a href="{{.Prev}}">prev</a>
{{else}}
<span class="disabled">prev</span>
{{end}}
 | <a href="../index.html">index</a> | 
{{if .Next}}
<a href="{{.Next}}">next</a>
{{else}}
<span class="disabled">next</span>
{{end}}
</p>
{{range .Images}}<a href="{{.}}.html"><img src="{{.}}.thumb.JPG"/></a>
{{end}}
<p>
{{if .Prev}}
<a href="{{.Prev}}">prev</a>
{{else}}
<span class="disabled">prev</span>
{{end}}
 | <a href="../index.html">index</a> | 
{{if .Next}}
<a href="{{.Next}}">next</a>
{{else}}
<span class="disabled">next</span>
{{end}}
</p>
</body>
</html>
