
# name -> dimensions
var buildings = <map>
var curBuilding = ""

(console:print "    Generating Table...\n")
(foreach [rubble:raws] block name contents {
	(if (str:cmp "building_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		(df:raw:walk [contents] block tag {
			(if (bool:or (str:cmp [tag id] "BUILDING_FURNACE") (str:cmp [tag id] "BUILDING_WORKSHOP")) {
				(if (int:eq (len [tag]) 1) {
					[curBuilding = [tag 0]]
				}{
					(rubble:abort "Error: invalid param count to BUILDING_XXXX raw tag in last file.")
				})
			})
			
			(if (str:cmp [tag id] "DIM") {
				(if (int:eq (len [tag]) 2) {
					[buildings [curBuilding] = (str:add " " [tag 0] "x" [tag 1])]
				}{
					(rubble:abort "Error: invalid param count to DIM raw tag in last file.")
				})
			})
			
			(break true)
		})
		
		[curBuilding = ""]
	})
	(break true)
})

(console:print "    Fixing raws...\n")

(foreach [rubble:raws] block name contents {
	(if (str:cmp "building_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [contents] block tag {
			(if (bool:or (str:cmp [tag id] "BUILDING_FURNACE") (str:cmp [tag id] "BUILDING_WORKSHOP")) {
				(if (int:eq (len [tag]) 1) {
					[curBuilding = [tag 0]]
				}{
					(rubble:abort "Error: invalid param count to BUILDING_XXXX raw tag in last file.")
				})
			})
			
			(if (str:cmp [tag id] "NAME") {
				(if (int:eq (len [tag]) 1) {
					(if (str:cmp [curBuilding] "") {
						# nop
					}{
						(if (str:cmp [buildings [curBuilding]] nil) {
							(break)
						})
						
						[tag replace = (str:add "[NAME:" [tag 0] [buildings [curBuilding]] "]")]
					})
				}{
					(rubble:abort "Error: invalid param count to NAME raw tag in last file.")
				})
			})
			
			(break true)
		})]
		[curBuilding = ""]
	})
	(break true)
})
