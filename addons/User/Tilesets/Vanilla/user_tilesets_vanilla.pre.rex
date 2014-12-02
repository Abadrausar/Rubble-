
(if (exists [rubble:raws] "user_tilesets_vanilla.pre.rex") {
	(axis:write [rubble:fs] "out:prep/user_tilesets_vanilla.prep.rex" [rubble:raws "user_tilesets_vanilla.pre.rex"])
})

(rubble:d_init_to_defaults)

(rubble:set_fullscreen_font_graphics "curses_square_16x16.png")
(rubble:set_windowed_font_graphics "curses_square_16x16.png")
(rubble:set_fullscreen_font "curses_800x600.png")
(rubble:set_windowed_font "curses_640x300.png")
