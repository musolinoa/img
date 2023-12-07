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
<script type="text/javascript">
document.onkeydown = function (e) {
	e = e || window.event
	switch(e.keyCode){
	case 37:
		document.getElementById("prev").click()
		break
	case 38:
		document.getElementById("up").click()
		break
	case 39:
		document.getElementById("next").click()
		break
	}
}
</script>
</head>
<body>
<p>
{{if .Prev}}
<a id="prev" href="../{{.Prev}}/index.html">prev</a>
{{else}}
<span class="disabled">prev</span>
{{end}}
 | <a id="up" href="../index.html">index</a> | 
{{if .Next}}
<a id="next" href="../{{.Next}}/index.html">prev</a>
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
