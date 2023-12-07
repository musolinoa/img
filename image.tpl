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
<script type="text/javascript">
document.onkeydown = function (e) {
	if (document.activeElement == document.getElementById("tag-list"))
		return
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
<p>{{if .Prev}}<a id="prev" href="{{.Prev}}">prev</a>{{else}}<span class="disabled">prev</span>{{end}} | <a id="up" href=".">up</a> | {{if .Next}}<a id="next" href="{{.Next}}">next</a>{{else}}<span class="disabled">next</span>{{end}}</p>
<p><a href="{{.Image}}.full.JPG"><img src="{{.Image}}.big.JPG"/></a></p>
{{range .Tags}} <a href="#">#{{.}}</a>{{end}}
<form action="/api/tag" method="post">
</form>
</body>
</html>
