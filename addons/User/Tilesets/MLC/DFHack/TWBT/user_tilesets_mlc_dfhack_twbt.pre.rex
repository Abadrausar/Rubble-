
(rubble:template "TWBT_TRACK_OVERRIDES" block tilemap {
	var shapes = <sarray
		FloorTrackN
		FloorTrackS
		FloorTrackE
		FloorTrackW
		FloorTrackNS
		FloorTrackNE
		FloorTrackNW
		FloorTrackSE
		FloorTrackSW
		FloorTrackEW
		FloorTrackNSE
		FloorTrackNSW
		FloorTrackNEW
		FloorTrackSEW
		FloorTrackNSEW
		RampTrackN
		RampTrackS
		RampTrackE
		RampTrackW
		RampTrackNS
		RampTrackNE
		RampTrackNW
		RampTrackSE
		RampTrackSW
		RampTrackEW
		RampTrackNSE
		RampTrackNSW
		RampTrackNEW
		RampTrackSEW
		RampTrackNSEW
	>
	var mats = <sarray
		Constructed
		Feature
		Frozen
		Lava
		Mineral
		Stone
	>
	var shape_to_dtile = <smap
		FloorTrackN = "208"
		FloorTrackS = "210"
		FloorTrackE = "198"
		FloorTrackW = "181"
		FloorTrackNS = "186"
		FloorTrackNE = "200"
		FloorTrackNW = "188"
		FloorTrackSE = "201"
		FloorTrackSW = "187"
		FloorTrackEW = "205"
		FloorTrackNSE = "204"
		FloorTrackNSW = "185"
		FloorTrackNEW = "202"
		FloorTrackSEW = "203"
		FloorTrackNSEW = "206"
		RampTrackN = "30"
		RampTrackS = "30"
		RampTrackE = "30"
		RampTrackW = "30"
		RampTrackNS = "30"
		RampTrackNE = "30"
		RampTrackNW = "30"
		RampTrackSE = "30"
		RampTrackSW = "30"
		RampTrackEW = "30"
		RampTrackNSE = "30"
		RampTrackNSW = "30"
		RampTrackNEW = "30"
		RampTrackSEW = "30"
		RampTrackNSEW = "30"
	>
	var shape_to_otile = (eval [tilemap])
	
	var out = ""
	(foreach [shapes] block _ shape {
		(foreach [mats] block _ mat {
			[out = (str:add [out] "[OVERRIDE:" [shape_to_dtile [shape]] ":T:" [mat] [shape] ":ex1:" [shape_to_otile [shape]] "]\n")]
		})
		[out = (str:add [out] "\n")]
	})
	(str:trimspace [out])
})
