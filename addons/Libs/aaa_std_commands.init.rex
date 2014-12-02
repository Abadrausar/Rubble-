
# Load some helpful commands.

# Rubble specific commands

# Placeholder template generator.
command rubble:placeholder name {
	(rubble:template [name] block ... {
		"Library template addon disabled."
	})
}

command rubble:skipfile name {
	(rubble:filetag [name] Skip true)
}

command rubble:nowritefile name {
	(rubble:filetag [name] NoWrite true)
}

command rubble:graphicsfile name {
	(rubble:filetag [name] GraphicsFile true)
}

command rubble:prepfile name {
	(rubble:filetag [name] PrepFile true)
}
