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

package main

import "html/template"
import "io/ioutil"

func loadOr(name, file string) string {
	content, err := ioutil.ReadFile("./other/webUI/" + name)
	if err != nil {
		ioutil.WriteFile("./other/webUI/"+name, []byte(file), 04755)
		return file
	}
	return string(content)
}

// Main Menu and Common Pages

var menuPage = loadOr("menu.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Menu</title>
</head>
<body>
	<h1>Welcome to the Rubble Web UI!</h1>
	<p><a href="./genaddons">Generate raws</a>
	<p><a href="./prep">Prep a region</a>
	<p><a href="./regen">Regenerate a region</a>
	<p><a href="./addons">List of all addons</a>
	<p><a href="./log">View Rubble log</a>
	<p><a href="./about">About</a>
	<p><a href="./kill">Quit</a>
	<p>WARNING! Be careful with the browser back button!<br>
	Every time a new page loads Rubble's internal state changes and some of these changes cannot be undone!<br>
	Returning to the main menu will reset the state, making it safe to start a new operation.
</body>
</html>
`)

var logPage = template.Must(template.New("log").Parse(loadOr("log.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Log</title>
</head>
<body>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
	<pre>{{.}}</pre>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

var killPage = loadOr("kill.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Exit</title>
</head>
<body>
	<h2>Thank you for using Rubble Web UI!</h2>
	<p>The server has shutdown.
	<p>Please post any suggestions or bugs to the Rubble thread on the forum!
</body>
</html>
`)

var aboutPage = loadOr("about.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: About</title>
</head>
<body>
	<p>Rubble Web UI
	<p>Rubble is Copyright 2013-2014 by Milo Christiansen
	
	<pre class="border">
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
	
	<p><a href="./">Back to Menu</a>
</body>
</html>
`)

var addondataPage = template.Must(template.New("addondata").Parse(loadOr("addondata.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Addon Information for "{{.Name}}"</title>
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
			document.write("<pre>​" + meta.Description + "</pre>")
		}
		
		document.write("<h4>Dependencies (automatically activated):</h4>")
		for (var i in meta.Activates) {
			document.write("<li><a href=\"./addondata?addon=" + meta.Activates[i] + "\">" + meta.Activates[i] + "</a>")
		}
		
		var files = {{.Files}}
		document.write("<h4>Addon Files:</h4>")
		for (var i in files) {
			document.write("<li><a href=\"./addonfile?addon={{.Name}}&file=" + files[i].Name + "\">" + files[i].Name + "</a>")
		}
		
	</script>
</body>
</html>
`)))

var addonfilePage = template.Must(template.New("addonfile").Parse(loadOr("addonfile.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: File "{{.Name}}"</title>
</head>
<body>
	<h2>{{.Name}}</h2>
	<p>{{.Source}}/{{.Name}}
	<pre class="border">
{{printf "%s" .Content}}
	</pre>
</body>
</html>
`)))

// Master Addon List

var addonsPage = template.Must(template.New("genaddons").Parse(loadOr("addons.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Addon List</title>
</head>
<body>
	<table border="0">
	<script language="JavaScript"> 
		var items = {{.}}
		for (var i in items) {
			document.write("<tr><td>")
			document.write("<a href=\"./addondata?addon=" + items[i].Name + "\">" + items[i].Name + "</a>")
			document.write("</td><td>")
			document.write(items[i].Meta.Header)
			document.write("</td></tr>")
		}
	</script>
	</table>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

// Normal Generation

var genaddonsPage = template.Must(template.New("genaddons").Parse(loadOr("genaddons.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Select Addons</title>
</head>
<body>
	<form action="./genbranch">
		<table border="0">
		<script language="JavaScript"> 
			var items = {{.}}
			for (var i in items) {
				if (items[i].Meta.Lib) {
					continue
				}
				
				document.write("<tr><td nowrap>")
				if (items[i].Active) {
					document.write("<input type=\"checkbox\" value=\"true\" checked name=\"" + items[i].Name + "\">")
				} else {
					document.write("<input type=\"checkbox\" value=\"true\" name=\"" + items[i].Name + "\">")
				}
				document.write("&nbsp;<a href=\"./addondata?addon=" + items[i].Name + "\">" + items[i].Name.replace(" ", "&nbsp;") + "</a>")
				document.write("</td><td>")
				document.write(items[i].Meta.Header)
				document.write("</td></tr>")
			}
		</script>
		
		<tr><td><input type="submit" name="__EditConfig__" value="Edit Configuration Variables"/></td></tr>
		<tr><td><input type="submit" name="__Generate__" value="Go Directly To Generation"/></td></tr>
		</table>
	</form>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

var genvarsPage = template.Must(template.New("genvars").Parse(loadOr("genvars.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Edit Configuration Variables</title>
</head>
<body>
	<form action="./genrun">
		<script language="JavaScript"> 
			var items = {{.Addons}}
			for (var i in items) {
				document.write("<input type=\"hidden\" value=\"true\" name=\"" + items[i] + "\">")
			}
		</script>
		
		<table border="0">
		<script language="JavaScript"> 
			var items = {{.Vars}}
			var noitems = true
			
			for (var i in items) {
				noitems = false
				
				var name = i
				if (items[i].Name != "") {
					name = items[i].Name
				}
				
				document.write("<tr><td>" + name + "</td>")
				document.write("<td>")
				if (items[i].Choices != null) {
					document.write("<select name=\"__CONFIG_VAR_" + i + "\">")
					for (var j in items[i].Choices) {
						if (items[i].Choices[j] == items[i].Val) {
							document.write("<option selected>" + items[i].Choices[j])
						} else {
							document.write("<option>" + items[i].Choices[j])
						}
					}
					document.write("</select>")
				} else {
					document.write("<input type=\"text\" value=\"" + items[i].Val + "\" name=\"__CONFIG_VAR_" + i + "\">")
				}
				document.write("</td></tr>")
			}
			if (noitems) {
				document.write("<tr><td>No Variables listed in Meta Data</td></tr>")
			}
		</script>
		
		<tr><td><input type="submit" value="Generate Raws"/></td></tr>
		</table>
	</form>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

// Prep

var prepPage = template.Must(template.New("prep").Parse(loadOr("prep.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Select Region</title>
</head>
<body>
	<form action="./preprun">
		<table border="0">
		<tr><td><input type="radio" value="raw" checked name="region"> raw</td></tr>
		<script language="JavaScript"> 
			var items = {{.}}
			for(var i in items) 
			{
				document.write("<tr><td>")
				document.write("<input type=\"radio\" value=\"" + items[i] + "\" name=\"region\">")
				document.write(" " + items[i])
				document.write("</td></tr>")
			}
		</script>
		
		<tr><td><input type="submit" value="Prep Selected Region"/></td></tr>
		</table>
	</form>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

// Regen

var regenPage = loadOr("regen.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Regen Warning</title>
</head>
<body>
	<p><h2>WARNING!</h2>
	<p>Regenerating a region's raws is dangerous! If you make any mistakes the region could become unplayable!<br>
	It is highly recommended that this be only used for switching tilesets and other actions that do not significantly edit the raws.
	<p>NEVER use regen to update from one version of an addon to another unless you are absolutely SURE that the new version is save compatible with the old version!
	
	<p><a href="./">Back to Menu</a>
	<p><a href="./regenregion">Continue to Regeneration</a>
</body>
</html>
`)

var regenregionPage = template.Must(template.New("regenregion").Parse(loadOr("regenregion.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Select Region</title>
</head>
<body>
	<form action="./regenaddons">
		<table border="0">
		<script language="JavaScript"> 
			var items = {{.}}
			var once = true
			for(var i in items) {
				if (once) {
					once = false
					document.write("<tr><td>")
					document.write("<input type=\"radio\" value=\"" + items[i] + "\" checked name=\"region\">")
					document.write(" " + items[i])
					document.write("</td></tr>")
					continue
				}
				
				document.write("<tr><td>")
				document.write("<input type=\"radio\" value=\"" + items[i] + "\" name=\"region\">")
				document.write(" " + items[i])
				document.write("</td></tr>")
			}
			if (once) {
				once = false
				document.write("<tr><td>")
				document.write("No Regions available to regenerate.")
				document.write("</td></tr>")
			} else {
				document.write("<tr><td><input type=\"submit\" value=\"Choose Addons For Selected Region\"/></td></tr>")
			}
		</script>
		
		
		</table>
	</form>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

var regenaddonsPage = template.Must(template.New("regenaddons").Parse(loadOr("regenaddons.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Select Addons</title>
</head>
<body>
	<form action="./regenbranch">
		<table border="0">
		<script language="JavaScript"> 
			var items = {{.Addons}}
			for (var i in items) {
				if (items[i].Meta.Lib) {
					continue
				}
				
				document.write("<tr><td nowrap>")
				if (items[i].Active) {
					document.write("<input type=\"checkbox\" value=\"true\" checked name=\"" + items[i].Name + "\">")
				} else {
					document.write("<input type=\"checkbox\" value=\"true\" name=\"" + items[i].Name + "\">")
				}
				document.write("&nbsp;<a href=\"./addondata?addon=" + items[i].Name + "\">" + items[i].Name + "</a>")
				document.write("</td><td>")
				document.write(items[i].Meta.Header)
				document.write("</td></tr>")
			}
		</script>
		
		<tr><td><input type="submit" name="__EditConfig__" value="Edit Configuration Variables"/></td></tr>
		<tr><td><input type="submit" name="__Generate__" value="Go Directly To Regeneration"/></td></tr>
		</table>
	</form>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))

var regenvarsPage = template.Must(template.New("regenvars").Parse(loadOr("regenvars.html", `
<html>
<head>
	<link rel="stylesheet" type="text/css" href="./css"/>
	<title>Rubble Web UI: Edit Configuration Variables</title>
</head>
<body>
	<form action="./regenrun">
		<script language="JavaScript"> 
			var items = {{.Addons}}
			for (var i in items) {
				document.write("<input type=\"hidden\" value=\"true\" name=\"" + items[i] + "\">")
			}
		</script>
		
		<table border="0">
		<script language="JavaScript"> 
			var items = {{.Vars}}
			var noitems = true
			
			for (var i in items) {
				noitems = false
				
				var name = i
				if (items[i].Name != "") {
					name = items[i].Name
				}
				
				document.write("<tr><td>" + name + "</td>")
				document.write("<td>")
				if (items[i].Choices != null) {
					document.write("<select name=\"__CONFIG_VAR_" + i + "\">")
					for (var j in items[i].Choices) {
						if (items[i].Choices[j] == items[i].Val) {
							document.write("<option selected>" + items[i].Choices[j])
						} else {
							document.write("<option>" + items[i].Choices[j])
						}
					}
					document.write("</select>")
				} else {
					document.write("<input type=\"text\" value=\"" + items[i].Val + "\" name=\"__CONFIG_VAR_" + i + "\">")
				}
				document.write("</td></tr>")
			}
			if (noitems) {
				document.write("<tr><td>No Variables listed in Meta Data</td></tr>")
			}
		</script>
		
		<tr><td><input type="submit" value="Regenerate Raws"/></td></tr>
		</table>
	</form>
	<p><a href="./">Back to Menu</a>
	<p><a href="./kill">Quit</a>
</body>
</html>
`)))
