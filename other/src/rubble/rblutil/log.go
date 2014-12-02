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

package rblutil

import "math/rand"
import "time"
import "fmt"
import "os"
import "io"

// Logger writes to rubble.log and os.Stdout by default, you can add more locations if you wish.
type Logger struct {
	io.Writer
}

// NewLogger creates a new logger.
func NewLogger(other ...io.Writer) (error, *Logger) {
	file, err := os.Create("./rubble.log")
	if err != nil {
		return err, nil
	}
	// Ouch. That is quite the hack...
	return nil, &Logger{io.MultiWriter(append([]io.Writer{file, os.Stdout}, other...)...)}
}

func (log *Logger) Printf(format string, msg ...interface{}) {
	fmt.Fprintf(log, format, msg...)
}

func (log *Logger) Println(msg ...interface{}) {
	fmt.Fprintln(log, msg...)
}

func (log *Logger) Print(msg ...interface{}) {
	fmt.Fprint(log, msg...)
}

// Separator writes a section separator to the log.
// Use for consistency.
func (log *Logger) Separator() {
	fmt.Fprint(log, "==============================================\n")
}

// Header writes the standard header to the log.
// Use for consistency.
func (log *Logger) Header(version string) {
	rand.Seed(time.Now().Unix())
	fmt.Fprint(log, "Rubble v"+version+"\n")
	fmt.Fprint(log, startupLines[rand.Int()%(len(startupLines)-1)]+"\n")
	log.Separator()
}

var startupLines = []string{
	"After Blast comes Rubble.",
	"Modding made easy!",
	"Scriptable!",
	"Templates!",
	"Now with random startup lines!",
	"Why did I add this feature?",
	"Coming up with cool startup lines is tough!",
	"Rubblize it!",
	"Now with a web UI!",
	"Use the command line!",
	"Configurable!",
	"Now with more addons you will never use!",
	"More bugs than a termite mound!",
	"It'll abort, you know it will.",
	"Please report any problems.",
	"Feedback is greatly appreciated!",
	"Unintentionally Ironic!",
	"Why do these all end with exclamation points!",
	"Free exclamation points!",
	"Guaranteed 50% bug free!",
	"There better not be an error log!",
	"Run it again, this line might change.",
	"Over 100 addons!",
	"Now with meta data!",
	"Scripting by Rex!",
	"Need help porting your mod? Just ask Milo!",
	"Have a cool idea for something to go here? Submit it!",
	"Read the documentation!",
	"Runs natively on Windows, Linux, and OSX!",
	"Under continuous development since June 2013!",
	"Done in less than 4 seconds or your money back!\n(disclaimer: Rubble is free.)",
	"Open source!",
	"Supports DFHack!",
}
