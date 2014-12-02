
(if (exists [rubble:raws] "user_tilesets_mlc_normal.pre.rex") {
	(axis:write [rubble:fs] "out:prep/user_tilesets_mlc_normal.prep.rex" [rubble:raws "user_tilesets_mlc_normal.pre.rex"])
	
	(rubble:prepfile "MLC 16x16.png")
	(rubble:prepfile "MLC 10x10.png")
})

(rubble:d_init_to_defaults)
(rubble:install_tilesheet "MLC 16x16.png")
(rubble:install_tilesheet "MLC 10x10.png")
(rubble:set_fullscreen_tilesheet "MLC 16x16.png")
(rubble:set_windowed_tilesheet "MLC 10x10.png")
