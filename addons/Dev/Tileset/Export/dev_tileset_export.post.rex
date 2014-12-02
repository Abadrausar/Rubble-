
# Write the tile info, the heavy lifting is done by the base templates.

var out = "\n"
(foreach [rubble:tileset] block _ line {
	[out = (str:add [out] [line] "\n")]
	(break true)
})
(axis:write [rubble:fs] "rubble:dev_export_tileset.tile_info.rbl" [out])
