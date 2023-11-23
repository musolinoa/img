<!DOCTYPE html>
<html>
<head>
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
<p>
{{if .Prev}}
<a href="{{.Prev}}">prev</a>
{{else}}
<span class="disabled">prev</span>
{{end}}
 | <a href=".">up</a> | 
{{if .Next}}
<a href="{{.Next}}">next</a>
{{else}}
<span class="disabled">next</span>
{{end}}
</p>
<p><a href="{{.Image}}.full.JPG"><img src="{{.Image}}.big.JPG"/></a></p>
</body>
</html>
