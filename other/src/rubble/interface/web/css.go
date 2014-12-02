/*
Copyright 2013-2014 by Milo Christiansen

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
*/

package main

// Gray on black, kinda ugly, but easy on the eyes.
var css = loadOr("rbl.css", `
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

pre.border
{
	border-style: dashed;
	border-width: 1px;
	font-family: monospace;
	font-size: 10pt;
	padding: 5px 5px 5px 5px;
	margin: 5px 5px 5px 5px;
}

select
{
	border-color: LightGray;
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 10pt;
}

input
{
	border-color: LightGray;
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 10pt;
}

textarea
{
	border-color: LightGray;
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 10pt;
}

h1
{
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 20pt;
}

h2
{
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 15pt;
}

h3
{
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 12pt;
}

h4
{
	color: Black;
	background: DarkGray;
	font-family: Verdana, Arial, Helvetica, sans-serif;
	font-size: 10pt;
}
`)
