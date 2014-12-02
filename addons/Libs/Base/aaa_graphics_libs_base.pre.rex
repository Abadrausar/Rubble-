
command rubble:install_graphics_file tilesheet {
	(axis:write [rubble:fs] (str:add "out:graphics/" [tilesheet]) [rubble:raws [tilesheet]])
}
(rubble:template "@INSTALL_GRAPHICS_FILE" block tilesheet {
	(rubble:install_graphics_file [tilesheet])
	(ret "")
})

command rubble:install_images_as_graphics addon {
	var addonRef = [rubble:addonstbl [addon]]
	(if (isnil [addonRef]) {
		(rubble:abort (str:add "Attempt to install PNG images from non-existent addon: \"" [addon] "\" As graphics."))
	})
	(if [addonRef Active] {}{
		(rubble:abort (str:add "Attempt to install PNG images from inactive addon: \"" [addon] "\" As graphics."))
	})
	
	(foreach [addonRef Files] block _ file {
		(if (str:cmp (str:right [file Name] 4) ".png") {
			(rubble:install_graphics_file [file Name])
		})
	})
}
(rubble:template "@INSTALL_IMAGES_AS_GRAPHICS" block addon {
	(rubble:install_images_as_graphics [addon])
	(ret "")
})

(rubble:template "@GRAPHICS_FILE" {
	(rubble:filetag (rubble:currentfile) "GraphicsFile" true)
	(ret "")
})
