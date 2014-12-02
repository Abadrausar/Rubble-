
# this is the prep file version of the base templates
# only a small subset are included
# these are all copies of other templates/commands

command rubble:install_tilesheet tilesheet {
	(axis:write [rubble:fs] (str:add "df:data/art/" [tilesheet]) [rubble:raws [tilesheet]])
}
command rubble:set_fullscreen_tilesheet tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FULLFONT:[^]]+\\]" [init] (str:add "[FULLFONT:" [tilesheet] "]"))]
	[init = (regex:replace "\\[GRAPHICS_FULLFONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FULLFONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
}
command rubble:set_windowed_tilesheet tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FONT:[^]]+\\]" [init] (str:add "[FONT:" [tilesheet] "]"))]
	[init = (regex:replace "\\[GRAPHICS_FONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
}
command rubble:set_fullscreen_font_graphics tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[GRAPHICS_FULLFONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FULLFONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
}
command rubble:set_windowed_font_graphics tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[GRAPHICS_FONT:[^]]+\\]" [init] (str:add "[GRAPHICS_FONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
}
command rubble:set_fullscreen_font tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FULLFONT:[^]]+\\]" [init] (str:add "[FULLFONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
}
command rubble:set_windowed_font tilesheet {
	var init = (axis:read [rubble:fs] "df:data/init/init.txt")
	
	[init = (regex:replace "\\[FONT:[^]]+\\]" [init] (str:add "[FONT:" [tilesheet] "]"))]

	(axis:write [rubble:fs] "df:data/init/init.txt" [init])
}

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

command rubble:open_d_init {
	[rubble:d_init = (axis:read [rubble:fs] "df:data/init/d_init.txt")]
}
command rubble:edit_d_init setting value {
	(if (exists [rubble:init_defaults] [setting]) {
		[rubble:init_settings [setting] = [value]]
	}{
		(rubble:abort "Error: Attempt to use invalid init setting, only tileset related settings allowed.")
	})
}
command rubble:close_d_init {
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
}
command rubble:d_init_to_defaults {
	[rubble:d_init = (axis:read [rubble:fs] "df:data/init/d_init.txt")]
	
	(foreach [rubble:init_defaults] block key value {
		[rubble:d_init =
			(regex:replace [rubble:init_regexes [key]] [rubble:d_init] 
				(str:add "[" [key] ":" [rubble:init_defaults [key]] "]"))
		]
		
		(break true)
	})
	(axis:write [rubble:fs] "df:data/init/d_init.txt" [rubble:d_init])
}
