
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
	
	"TREE_ROOT_SLOPING"="127"
	"TREE_TRUNK_SLOPING"="127"
	"TREE_ROOT_SLOPING_DEAD"="127"
	"TREE_TRUNK_SLOPING_DEAD"="127"
	"TREE_ROOTS"="172"
	"TREE_ROOTS_DEAD"="172"
	"TREE_BRANCHES"="172"
	"TREE_BRANCHES_DEAD"="172"
	"TREE_SMOOTH_BRANCHES"="'#'"
	"TREE_SMOOTH_BRANCHES_DEAD"="'#'"
	"TREE_TRUNK_PILLAR"="'O'"
	"TREE_TRUNK_PILLAR_DEAD"="'O'"
	"TREE_CAP_PILLAR"="'O'"
	"TREE_CAP_PILLAR_DEAD"="'O'"
	"TREE_TRUNK_N"="205"
	"TREE_TRUNK_S"="205"
	"TREE_TRUNK_N_DEAD"="205"
	"TREE_TRUNK_S_DEAD"="205"
	"TREE_TRUNK_EW"="205"
	"TREE_TRUNK_EW_DEAD"="205"
	"TREE_CAP_WALL_N"="205"
	"TREE_CAP_WALL_S"="205"
	"TREE_CAP_WALL_N_DEAD"="205"
	"TREE_CAP_WALL_S_DEAD"="205"
	"TREE_TRUNK_E"="186"
	"TREE_TRUNK_W"="186"
	"TREE_TRUNK_E_DEAD"="186"
	"TREE_TRUNK_W_DEAD"="186"
	"TREE_TRUNK_NS"="186"
	"TREE_TRUNK_NS_DEAD"="186"
	"TREE_CAP_WALL_E"="186"
	"TREE_CAP_WALL_W"="186"
	"TREE_CAP_WALL_E_DEAD"="186"
	"TREE_CAP_WALL_W_DEAD"="186"
	"TREE_TRUNK_NW"="201"
	"TREE_CAP_WALL_NW"="201"
	"TREE_TRUNK_NW_DEAD"="201"
	"TREE_CAP_WALL_NW_DEAD"="201"
	"TREE_TRUNK_NE"="187"
	"TREE_CAP_WALL_NE"="187"
	"TREE_TRUNK_NE_DEAD"="187"
	"TREE_CAP_WALL_NE_DEAD"="187"
	"TREE_TRUNK_SW"="200"
	"TREE_CAP_WALL_SW"="200"
	"TREE_TRUNK_SW_DEAD"="200"
	"TREE_CAP_WALL_SW_DEAD"="200"
	"TREE_TRUNK_SE"="188"
	"TREE_CAP_WALL_SE"="188"
	"TREE_TRUNK_SE_DEAD"="188"
	"TREE_CAP_WALL_SE_DEAD"="188"
	"TREE_TRUNK_NSE"="204"
	"TREE_TRUNK_NSE_DEAD"="204"
	"TREE_TRUNK_NSW"="185"
	"TREE_TRUNK_NSW_DEAD"="185"
	"TREE_TRUNK_NEW"="202"
	"TREE_TRUNK_NEW_DEAD"="202"
	"TREE_TRUNK_SEW"="203"
	"TREE_TRUNK_SEW_DEAD"="203"
	"TREE_TRUNK_NSEW"="206"
	"TREE_TRUNK_NSEW_DEAD"="206"
	"TREE_TRUNK_BRANCH_N"="207"
	"TREE_TRUNK_BRANCH_N_DEAD"="207"
	"TREE_TRUNK_BRANCH_S"="209"
	"TREE_TRUNK_BRANCH_S_DEAD"="209"
	"TREE_TRUNK_BRANCH_E"="199"
	"TREE_TRUNK_BRANCH_E_DEAD"="199"
	"TREE_TRUNK_BRANCH_W"="182"
	"TREE_TRUNK_BRANCH_W_DEAD"="182"
	"TREE_BRANCH_NS"="179"
	"TREE_BRANCH_NS_DEAD"="179"
	"TREE_BRANCH_EW"="196"
	"TREE_BRANCH_EW_DEAD"="196"
	"TREE_BRANCH_NW"="217"
	"TREE_BRANCH_NW_DEAD"="217"
	"TREE_BRANCH_NE"="192"
	"TREE_BRANCH_NE_DEAD"="192"
	"TREE_BRANCH_SW"="191"
	"TREE_BRANCH_SW_DEAD"="191"
	"TREE_BRANCH_SE"="218"
	"TREE_BRANCH_SE_DEAD"="218"
	"TREE_BRANCH_NSE"="195"
	"TREE_BRANCH_NSE_DEAD"="195"
	"TREE_BRANCH_NSW"="180"
	"TREE_BRANCH_NSW_DEAD"="180"
	"TREE_BRANCH_NEW"="193"
	"TREE_BRANCH_NEW_DEAD"="193"
	"TREE_BRANCH_SEW"="194"
	"TREE_BRANCH_SEW_DEAD"="194"
	"TREE_BRANCH_NSEW"="197"
	"TREE_BRANCH_NSEW_DEAD"="197"
	"TREE_TWIGS"="';'"
	"TREE_TWIGS_DEAD"="';'"
	"TREE_CAP_RAMP"="30"
	"TREE_CAP_RAMP_DEAD"="30"
	"TREE_CAP_FLOOR1"="249"
	"TREE_CAP_FLOOR2"="249"
	"TREE_CAP_FLOOR1_DEAD"="249"
	"TREE_CAP_FLOOR2_DEAD"="249"
	"TREE_CAP_FLOOR3"="249"
	"TREE_CAP_FLOOR4"="249"
	"TREE_CAP_FLOOR3_DEAD"="249"
	"TREE_CAP_FLOOR4_DEAD"="249"
	"TREE_TRUNK_INTERIOR"="10"
	"TREE_TRUNK_INTERIOR_DEAD"="10"
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
	
	"TREE_ROOT_SLOPING"="\\[TREE_ROOT_SLOPING:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SLOPING"="\\[TREE_TRUNK_SLOPING:(?:[0-9]+|'.')\\]"
	"TREE_ROOT_SLOPING_DEAD"="\\[TREE_ROOT_SLOPING_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SLOPING_DEAD"="\\[TREE_TRUNK_SLOPING_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_ROOTS"="\\[TREE_ROOTS:(?:[0-9]+|'.')\\]"
	"TREE_ROOTS_DEAD"="\\[TREE_ROOTS_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCHES"="\\[TREE_BRANCHES:(?:[0-9]+|'.')\\]"
	"TREE_BRANCHES_DEAD"="\\[TREE_BRANCHES_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_SMOOTH_BRANCHES"="\\[TREE_SMOOTH_BRANCHES:(?:[0-9]+|'.')\\]"
	"TREE_SMOOTH_BRANCHES_DEAD"="\\[TREE_SMOOTH_BRANCHES_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_PILLAR"="\\[TREE_TRUNK_PILLAR:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_PILLAR_DEAD"="\\[TREE_TRUNK_PILLAR_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_PILLAR"="\\[TREE_CAP_PILLAR:(?:[0-9]+|'.')\\]"
	"TREE_CAP_PILLAR_DEAD"="\\[TREE_CAP_PILLAR_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_N"="\\[TREE_TRUNK_N:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_S"="\\[TREE_TRUNK_S:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_N_DEAD"="\\[TREE_TRUNK_N_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_S_DEAD"="\\[TREE_TRUNK_S_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_EW"="\\[TREE_TRUNK_EW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_EW_DEAD"="\\[TREE_TRUNK_EW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_N"="\\[TREE_CAP_WALL_N:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_S"="\\[TREE_CAP_WALL_S:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_N_DEAD"="\\[TREE_CAP_WALL_N_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_S_DEAD"="\\[TREE_CAP_WALL_S_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_E"="\\[TREE_TRUNK_E:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_W"="\\[TREE_TRUNK_W:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_E_DEAD"="\\[TREE_TRUNK_E_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_W_DEAD"="\\[TREE_TRUNK_W_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NS"="\\[TREE_TRUNK_NS:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NS_DEAD"="\\[TREE_TRUNK_NS_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_E"="\\[TREE_CAP_WALL_E:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_W"="\\[TREE_CAP_WALL_W:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_E_DEAD"="\\[TREE_CAP_WALL_E_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_W_DEAD"="\\[TREE_CAP_WALL_W_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NW"="\\[TREE_TRUNK_NW:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_NW"="\\[TREE_CAP_WALL_NW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NW_DEAD"="\\[TREE_TRUNK_NW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_NW_DEAD"="\\[TREE_CAP_WALL_NW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NE"="\\[TREE_TRUNK_NE:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_NE"="\\[TREE_CAP_WALL_NE:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NE_DEAD"="\\[TREE_TRUNK_NE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_NE_DEAD"="\\[TREE_CAP_WALL_NE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SW"="\\[TREE_TRUNK_SW:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_SW"="\\[TREE_CAP_WALL_SW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SW_DEAD"="\\[TREE_TRUNK_SW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_SW_DEAD"="\\[TREE_CAP_WALL_SW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SE"="\\[TREE_TRUNK_SE:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_SE"="\\[TREE_CAP_WALL_SE:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SE_DEAD"="\\[TREE_TRUNK_SE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_WALL_SE_DEAD"="\\[TREE_CAP_WALL_SE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NSE"="\\[TREE_TRUNK_NSE:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NSE_DEAD"="\\[TREE_TRUNK_NSE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NSW"="\\[TREE_TRUNK_NSW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NSW_DEAD"="\\[TREE_TRUNK_NSW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NEW"="\\[TREE_TRUNK_NEW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NEW_DEAD"="\\[TREE_TRUNK_NEW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SEW"="\\[TREE_TRUNK_SEW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_SEW_DEAD"="\\[TREE_TRUNK_SEW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NSEW"="\\[TREE_TRUNK_NSEW:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_NSEW_DEAD"="\\[TREE_TRUNK_NSEW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_N"="\\[TREE_TRUNK_BRANCH_N:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_N_DEAD"="\\[TREE_TRUNK_BRANCH_N_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_S"="\\[TREE_TRUNK_BRANCH_S:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_S_DEAD"="\\[TREE_TRUNK_BRANCH_S_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_E"="\\[TREE_TRUNK_BRANCH_E:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_E_DEAD"="\\[TREE_TRUNK_BRANCH_E_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_W"="\\[TREE_TRUNK_BRANCH_W:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_BRANCH_W_DEAD"="\\[TREE_TRUNK_BRANCH_W_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NS"="\\[TREE_BRANCH_NS:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NS_DEAD"="\\[TREE_BRANCH_NS_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_EW"="\\[TREE_BRANCH_EW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_EW_DEAD"="\\[TREE_BRANCH_EW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NW"="\\[TREE_BRANCH_NW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NW_DEAD"="\\[TREE_BRANCH_NW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NE"="\\[TREE_BRANCH_NE:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NE_DEAD"="\\[TREE_BRANCH_NE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_SW"="\\[TREE_BRANCH_SW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_SW_DEAD"="\\[TREE_BRANCH_SW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_SE"="\\[TREE_BRANCH_SE:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_SE_DEAD"="\\[TREE_BRANCH_SE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NSE"="\\[TREE_BRANCH_NSE:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NSE_DEAD"="\\[TREE_BRANCH_NSE_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NSW"="\\[TREE_BRANCH_NSW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NSW_DEAD"="\\[TREE_BRANCH_NSW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NEW"="\\[TREE_BRANCH_NEW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NEW_DEAD"="\\[TREE_BRANCH_NEW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_SEW"="\\[TREE_BRANCH_SEW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_SEW_DEAD"="\\[TREE_BRANCH_SEW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NSEW"="\\[TREE_BRANCH_NSEW:(?:[0-9]+|'.')\\]"
	"TREE_BRANCH_NSEW_DEAD"="\\[TREE_BRANCH_NSEW_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TWIGS"="\\[TREE_TWIGS:(?:[0-9]+|'.')\\]"
	"TREE_TWIGS_DEAD"="\\[TREE_TWIGS_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_RAMP"="\\[TREE_CAP_RAMP:(?:[0-9]+|'.')\\]"
	"TREE_CAP_RAMP_DEAD"="\\[TREE_CAP_RAMP_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR1"="\\[TREE_CAP_FLOOR1:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR2"="\\[TREE_CAP_FLOOR2:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR1_DEAD"="\\[TREE_CAP_FLOOR1_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR2_DEAD"="\\[TREE_CAP_FLOOR2_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR3"="\\[TREE_CAP_FLOOR3:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR4"="\\[TREE_CAP_FLOOR4:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR3_DEAD"="\\[TREE_CAP_FLOOR3_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_CAP_FLOOR4_DEAD"="\\[TREE_CAP_FLOOR4_DEAD:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_INTERIOR"="\\[TREE_TRUNK_INTERIOR:(?:[0-9]+|'.')\\]"
	"TREE_TRUNK_INTERIOR_DEAD"="\\[TREE_TRUNK_INTERIOR_DEAD:(?:[0-9]+|'.')\\]"
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
