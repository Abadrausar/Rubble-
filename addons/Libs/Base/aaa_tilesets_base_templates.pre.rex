
var rubble:tile_data = <map>
var rubble:color_data = <map>

# Used by the "Dev/Extract Tileset" addon
var rubble:tileset = <array>

(rubble:template "SET_TILE" block id tile {
	(if (exists [rubble:tile_data] [id]){
	}{
		[rubble:tile_data [id] = [tile]]
	})
	(ret "")
})
(rubble:template "#TILE" block id default {
	[rubble:tileset append = (str:add "{SET_TILE;" [id] ";" [default] "}")]
	(if (exists [rubble:tile_data] [id]){
		(ret [rubble:tile_data [id]])
	}{
		(ret [default])
	})
})
(rubble:template "SET_COLOR" block id color {
	(if (exists [rubble:color_data] [id]){
	}{
		[rubble:color_data [id] = [color]]
	})
	(ret "")
})
(rubble:template "#COLOR" block id default {
	[rubble:tileset append = (str:add "{SET_COLOR;" [id] ";" [default] "}")]
	(if (exists [rubble:color_data] [id]){
		(ret [rubble:color_data [id]])
	}{
		(ret [default])
	})
})

(rubble:template "INSTALL_TILESHEET" block tilesheet {
	(axis:write [rubble:fs] (str:add "df:data/art/" [tilesheet]) [rubble:raws [tilesheet]])
	(ret "")
})
(rubble:template "SET_FULLSCREEN_TILESHEET" block tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FULLFONT:[^]]+\\]" [init] (str:add "[FULLFONT:" [tilesheet] "]"))]
	[init = (regex:replace "\\[GRAPHICS_FULLFONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FULLFONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
	(ret "")
})
(rubble:template "SET_WINDOWED_TILESHEET" block tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FONT:[^]]+\\]" [init] (str:add "[FONT:" [tilesheet] "]"))]
	[init = (regex:replace "\\[GRAPHICS_FONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
	(ret "")
})
(rubble:template "SET_FULLSCREEN_FONT_GRAPHICS" block tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[GRAPHICS_FULLFONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FULLFONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
	(ret "")
})
(rubble:template "SET_WINDOWED_FONT_GRAPHICS" block tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[GRAPHICS_FONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
	(ret "")
})
(rubble:template "SET_FULLSCREEN_FONT" block tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FULLFONT:[^]]+\\]" [init] (str:add "[FULLFONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
	(ret "")
})
(rubble:template "SET_WINDOWED_FONT" block tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FONT:[^]]+\\]" [init] (str:add "[FONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
	(ret "")
})

var rubble:d_init = ""

var rubble:init_settings = <map>

var rubble:init_defaults = <map
	"SKY"="178:3:0:0"
	"CHASM"="250:0:0:1"
	"PILLAR_TILE"="'O'"
	"TRACK_N"="208I"
	"TRACK_S"="210I"
	"TRACK_E"="198I"
	"TRACK_W"="181I"
	"TRACK_NS"="186I"
	"TRACK_NE"="200I"
	"TRACK_NW"="188I"
	"TRACK_SE"="201I"
	"TRACK_SW"="187I"
	"TRACK_EW"="205I"
	"TRACK_NSE"="204I"
	"TRACK_NSW"="185I"
	"TRACK_NEW"="202I"
	"TRACK_SEW"="203I"
	"TRACK_NSEW"="206I"
	"TRACK_RAMP_N"="30I"
	"TRACK_RAMP_S"="30I"
	"TRACK_RAMP_E"="30I"
	"TRACK_RAMP_W"="30I"
	"TRACK_RAMP_NS"="30I"
	"TRACK_RAMP_NE"="30I"
	"TRACK_RAMP_NW"="30I"
	"TRACK_RAMP_SE"="30I"
	"TRACK_RAMP_SW"="30I"
	"TRACK_RAMP_EW"="30I"
	"TRACK_RAMP_NSE"="30I"
	"TRACK_RAMP_NSW"="30I"
	"TRACK_RAMP_NEW"="30I"
	"TRACK_RAMP_SEW"="30I"
	"TRACK_RAMP_NSEW"="30I"
>

var rubble:init_regexes = <map
	"SKY"="\\[SKY:[0-9]+:[0-9]+:[0-9]+:[0-9]+\\]"
	"CHASM"="\\[CHASM:[0-9]+:[0-9]+:[0-9]+:[0-9]+\\]"
	"PILLAR_TILE"="\\[PILLAR_TILE:(?:[0-9]+|'.')\\]"
	"TRACK_N"="\\[TRACK_N:(?:[0-9]+|'.')I?\\]"
	"TRACK_S"="\\[TRACK_S:(?:[0-9]+|'.')I?\\]"
	"TRACK_E"="\\[TRACK_E:(?:[0-9]+|'.')I?\\]"
	"TRACK_W"="\\[TRACK_W:(?:[0-9]+|'.')I?\\]"
	"TRACK_NS"="\\[TRACK_NS:(?:[0-9]+|'.')I?\\]"
	"TRACK_NE"="\\[TRACK_NE:(?:[0-9]+|'.')I?\\]"
	"TRACK_NW"="\\[TRACK_NW:(?:[0-9]+|'.')I?\\]"
	"TRACK_SE"="\\[TRACK_SE:(?:[0-9]+|'.')I?\\]"
	"TRACK_SW"="\\[TRACK_SW:(?:[0-9]+|'.')I?\\]"
	"TRACK_EW"="\\[TRACK_EW:(?:[0-9]+|'.')I?\\]"
	"TRACK_NSE"="\\[TRACK_NSE:(?:[0-9]+|'.')I?\\]"
	"TRACK_NSW"="\\[TRACK_NSW:(?:[0-9]+|'.')I?\\]"
	"TRACK_NEW"="\\[TRACK_NEW:(?:[0-9]+|'.')I?\\]"
	"TRACK_SEW"="\\[TRACK_SEW:(?:[0-9]+|'.')I?\\]"
	"TRACK_NSEW"="\\[TRACK_NSEW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_N"="\\[TRACK_RAMP_N:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_S"="\\[TRACK_RAMP_S:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_E"="\\[TRACK_RAMP_E:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_W"="\\[TRACK_RAMP_W:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NS"="\\[TRACK_RAMP_NS:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NE"="\\[TRACK_RAMP_NE:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NW"="\\[TRACK_RAMP_NW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_SE"="\\[TRACK_RAMP_SE:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_SW"="\\[TRACK_RAMP_SW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_EW"="\\[TRACK_RAMP_EW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NSE"="\\[TRACK_RAMP_NSE:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NSW"="\\[TRACK_RAMP_NSW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NEW"="\\[TRACK_RAMP_NEW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_SEW"="\\[TRACK_RAMP_SEW:(?:[0-9]+|'.')I?\\]"
	"TRACK_RAMP_NSEW"="\\[TRACK_RAMP_NSEW:(?:[0-9]+|'.')I?\\]"
>

(rubble:template "OPEN_D_INIT" {
	[rubble:d_init = (axis:read [rubble:fs] "df:data/init/d_init.txt")]
	(ret "")
})
(rubble:template "EDIT_D_INIT" block setting value {
	(if (exists [rubble:init_defaults] [setting]) {
		[rubble:init_settings [setting] = [value]]
	}{
		(rubble:abort "Error: Attempt to use invalid init setting, only tileset related settings allowed.")
	})
	
	(ret "")
})
(rubble:template "CLOSE_D_INIT" {
	(foreach [rubble:init_defaults] block key value {
		(if (exists [rubble:init_settings] [key]) {
			[rubble:d_init =
				(regex:replace [rubble:init_regexes [key]] [rubble:d_init] 
					(str:add "[" [key] ":" [rubble:init_settings [key]] "]"))
			]
		}{
			[rubble:d_init =
				(regex:replace [rubble:init_regexes [key]] [rubble:d_init] 
					(str:add "[" [key] ":" [rubble:init_defaults [key]] "]"))
			]
		})
		
		(break true)
	})
	(axis:write [rubble:fs] "df:data/init/d_init.txt" [rubble:d_init])
	(ret "")
})
(rubble:template "D_INIT_TO_DEFAULTS" {
	[rubble:d_init = (axis:read [rubble:fs] "df:/data/init/d_init.txt")]
	
	(foreach [rubble:init_defaults] block key value {
		[rubble:d_init =
			(regex:replace [rubble:init_regexes [key]] [rubble:d_init] 
				(str:add "[" [key] ":" [rubble:init_defaults [key]] "]"))
		]
		
		(break true)
	})
	(axis:write [rubble:fs] "df:data/init/d_init.txt" [rubble:d_init])
	(ret "")
})

(rubble:template "#INSTALL_GRAPHICS_FILE" block tilesheet {
	(axis:write [rubble:fs] (str:add "out:graphics/" [tilesheet]) [rubble:raws [tilesheet]])
	(ret "")
})

(rubble:template "@GRAPHICS_FILE" {
	(rubble:graphicsfile (rubble:currentfile))
	(ret "")
})
