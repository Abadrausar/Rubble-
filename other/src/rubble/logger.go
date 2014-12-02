/*
Copyright 2013 by Milo Christiansen

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

import "os"
import "io"
import "fmt"

var logFile io.Writer

func InitLogging() {
	file, err := os.Create("./rubble.log")
	if err != nil {
		fmt.Println("FATAL ERROR: Unable to open log file:", err)
		fmt.Println("This should NEVER, NEVER happen! Please report this message ASAP!")
		os.Exit(1)
	}
	logFile = io.MultiWriter(file, os.Stdout)
}

func LogPrintf(format string, msg ...interface{}) {
	fmt.Fprintf(logFile, format, msg...)
}

func LogPrintln(msg ...interface{}) {
	fmt.Fprintln(logFile, msg...)
}

func LogPrint(msg ...interface{}) {
	fmt.Fprint(logFile, msg...)
}
