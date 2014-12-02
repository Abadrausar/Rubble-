
(if (exists [rubble:raws] "user_tilesets_mlc_pure_ascii.pre.rex") {
	(axis:write [rubble:fs] "out:prep/user_tilesets_mlc_pure_ascii.prep.rex" [rubble:raws "user_tilesets_mlc_pure_ascii.pre.rex"])
	
	(rubble:prepfile "MLC 16x16 - Font.png")
	(rubble:prepfile "MLC 10x10 - Font.png")
})

(rubble:d_init_to_defaults)
(rubble:install_tilesheet "MLC 16x16 - Font.png")
(rubble:install_tilesheet "MLC 10x10 - Font.png")
(rubble:set_fullscreen_tilesheet "MLC 16x16 - Font.png")
(rubble:set_windowed_tilesheet "MLC 10x10 - Font.png")
