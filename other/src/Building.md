
==============================================
Building Rubble from Source:
==============================================

To build Rubble you need to have Go (golang.org) installed.

One you have Go installed add <rubbledir>/other to your GOPATH or copy the contents of the src directory to a location on your GOPATH.

Now all you need to do is fire up a command prompt and run:
	go build rubble/interface/cli
Or:
	go build rubble/interface/web

To build the file encoder:
	go build rubble/encoder

When learning Rex it is very helpful to have the interactive shell available, to build the Rex shell:
	go build dctech/rex/rexsh
