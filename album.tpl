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
<a id="prev" href="{{.Prev}}">prev</a>
{{else}}
<span class="disabled">prev</span>
{{end}}
 | <a id="up" href="{{.UpLink}}">{{.UpText}}</a> | 
{{if .Next}}
<a id="next" href="{{.Next}}">next</a>
{{else}}
<span class="disabled">next</span>
{{end}}
</p>
{{range .Images}}<a href="{{.ID}}.html"><img src="{{.Prefix}}{{.ID}}.thumb.JPG"/></a>
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
