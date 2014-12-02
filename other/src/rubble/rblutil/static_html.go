/*
Copyright 2013-2014 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package rblutil

import "html/template"
import "dctech/axis"

func LoadHTMLTemplatesStatic(fs axis.DataSource) (*template.Template, error) {
	tmpl := template.New("static")
	
	_, err := tmpl.New("mainpage").Parse(LoadOr(fs, "mainpage_static.html", `
<html>
<head>
	<title>Rubble Generation Report: Menu</title>
	
	<style>
		body
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		table
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a
		{
			color: DarkGray;
			text-decoration: underline;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a:hover 
		{
			color: Black;
			background: DarkGray;
			text-decoration: none;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		p.mono
		{
			font-family: monospace;
			font-size: 10pt;
		}
		
		pre
		{
			border-style: none;
			font-family: monospace;
			font-size: 10pt;
			padding: 5px 5px 5px 5px;
			margin: 5px 5px 5px 5px;
		}
		
		h2
		{
			color: Black;
			background: DarkGray;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 15pt;
		}
	</style>
</head>
<body>
	<h2>Welcome to the Rubble Generation Report!</h1>
	<p><a href="./Docs/Addon List.html">Active Addon List</a>
	<p><a href="./Docs/Configuration.html">Configuration Variables</a>
	<p><a href="./Docs/About.html">About Rubble</a>
</body>
</html>
`))
	if err != nil {
		return nil, err
	}
	
	_, err = tmpl.New("about").Parse(LoadOr(fs, "about_static.html", `
<html>
<head>
	<title>Rubble Generation Report: About</title>
	
	<style>
		body
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a
		{
			color: DarkGray;
			text-decoration: underline;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a:hover 
		{
			color: Black;
			background: DarkGray;
			text-decoration: none;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		pre
		{
			border-style: dashed;
			border-width: 1px;
			font-family: monospace;
			font-size: 10pt;
			padding: 5px 5px 5px 5px;
			margin: 5px 5px 5px 5px;
		}
		
	</style>
</head>
<body>
	<p>Rubble is Copyright 2013-2014 by Milo Christiansen
	
	<pre>
This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
	</pre>
	
	<p>Written in <a href="http://golang.org">go</a>.
</body>
</html>
`))
	if err != nil {
		return nil, err
	}
	
	_, err = tmpl.New("config").Parse(LoadOr(fs, "config_static.html", `
<html>
<head>
	<title>Rubble Generation Report: Configuration Variables</title>
	
	<style>
		body
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		table
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		h2
		{
			color: Black;
			background: DarkGray;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 15pt;
		}
	</style>
</head>
<body>
	<h2>Rubble Configuration Variable Report</h2>
	<table border="0">
	<script language="JavaScript"> 
		var items = {{.}}
		for (var i in items) {
			document.write("<tr><td>")
			document.write(items[i].Name)
			document.write("</td><td>")
			document.write(items[i].Meta.Name)
			document.write("</td><td>")
			document.write(items[i].Value)
			document.write("</td></tr>")
		}
	</script>
	</table>
</body>
</html>
`))
	if err != nil {
		return nil, err
	}

	_, err = tmpl.New("addondata").Parse(LoadOr(fs, "addondata_static.html", `
<html>
<head>
	<title>Rubble Generation Report: Addon Information for "{{.Name}}"</title>
	
	<style>
		body
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		table
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a
		{
			color: DarkGray;
			text-decoration: underline;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a:hover 
		{
			color: Black;
			background: DarkGray;
			text-decoration: none;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		p.mono
		{
			font-family: monospace;
			font-size: 10pt;
		}
		
		pre
		{
			border-style: none;
			font-family: monospace;
			font-size: 10pt;
			padding: 5px 5px 5px 5px;
			margin: 5px 5px 5px 5px;
		}
		
		h2
		{
			color: Black;
			background: DarkGray;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 15pt;
		}
	</style>
</head>
<body>
	<h2>{{.Name}}</h2>
	<script language="JavaScript"> 
		var meta = {{.Meta}}
		if (meta.Lib == true) {
			document.write("(An automatically managed library)")
		}
		if (meta.Header != "") {
			document.write("<p class=\"mono\">" + meta.Header + "</p>")
		}
		
		if (meta.Description != "") {
			document.write("<pre>" + meta.Description + "</pre>")
		}
		
		document.write("<h4>Dependencies (automatically activated):</h4>")
		for (var i in meta.Activates) {
			document.write("<li>" + meta.Activates[i] + "</li>")
		}
		
	</script>
</body>
</html>
`))
	if err != nil {
		return nil, err
	}

	_, err = tmpl.New("addonlist").Parse(LoadOr(fs, "addonlist_static.html", `
<html>
<head>
	<title>Rubble Generation Report: Addon List</title>
	
	<style>
		body
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		table
		{
			color: DarkGray;
			background-color: Black;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a
		{
			color: DarkGray;
			text-decoration: underline;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		a:hover 
		{
			color: Black;
			background: DarkGray;
			text-decoration: none;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 10pt;
		}
		
		h2
		{
			color: Black;
			background: DarkGray;
			font-family: Verdana, Arial, Helvetica, sans-serif;
			font-size: 15pt;
		}
	</style>
</head>
<body>
	<h2>Rubble Active Addon List</h2>
	<table border="0">
	<script language="JavaScript"> 
		var items = {{.}}
		for (var i in items) {
			document.write("<tr><td>")
			document.write("<a href=\"./Addons/" + items[i].Name + ".html\">" + items[i].Name + "</a>")
			document.write("</td><td>")
			document.write(items[i].Meta.Header)
			document.write("</td></tr>")
		}
	</script>
	</table>
</body>
</html>
`))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}