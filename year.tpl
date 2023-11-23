<!DOCTYPE html>
<html>
<head>
<title>{{.CurrYear}}</title>
<style>
body{
	background-color: black;
	text-align: center;
}
div>p, a{
	color: white;
}
div{
	float: left;
	text-align: center;
	padding: 0.25cm;
	width: 200px;
	height: 200px;
}
div.bordered{
	float: none;
	display: inline-block;
	width: 160px;
	height: 160px;
	padding: 0px;
	border: 1px solid white;
}
img{
	width: 160px;
	height: 160px;
}
.disabled{
	color: grey;
}
</style>
</head>
<body>
<p>
{{if .PrevYear}}
<a href="../{{.PrevYear}}/index.html">prev</a>
{{else}}
<span class="disabled">prev</span>
{{end}}
 | <a href="../index.html">index</a> | 
{{if .NextYear}}
<a href="../{{.NextYear}}/index.html">prev</a>
{{else}}
<span class="disabled">next</span>
{{end}}
</p>
{{range .Months}}
	{{if .Empty}}
		<div><div class="bordered"></div><p>{{.Name}}</p></div>
	{{else}}
		<div><a href="{{.Number}}/index.html"><img src="{{.Number}}/montage.jpg"/><p>{{.Name}}</p></a></div>
	{{end}}
{{end}}
</body>
</html>
