
var isfish = false
var fishcount = 0
var totalfishcount = 0
var popschanged = 0
var totalpopschanged = 0

(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "    " [name] "\n")
		
		[totalfishcount = (int:add [totalfishcount] [fishcount])]
		[totalpopschanged = (int:add [totalpopschanged] [popschanged])]
		[fishcount = 0]
		[popschanged = 0]
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				[isfish = false]
			})
			
			# It may be better to use VERMIN_FISH, but with vanilla raws this works just fine.
			(if (str:cmp [tag id] "FISHITEM") {
				[isfish = true]
				[fishcount = (int:add [fishcount] 1)]
			})
			
			# If the POPULATION_NUMBER tag comes before the FISHITEM tag this script will fail to up 
			# the population of that fish. Good thing that never happens in any vanilla creatures!
			# Such a case is detectable by the detected fish count being larger than the changed population count.
			(if [isfish] {
				(if (str:cmp [tag "id"] "POPULATION_NUMBER") {
					(if (int:eq (len [tag]) 2) {
						[tag replace = "[POPULATION_NUMBER:30000:30000]"]
						[popschanged = (int:add [popschanged] 1)]
					}{
						(rubble:abort "Error: invalid param count to POPULATION_NUMBER raw tag in last file.")
					})
				})
			})
			
			(break true)
		})]
		(console:print "      Found:   " [fishcount] " Fish\n")
		(console:print "      Changed: " [popschanged] " Populations\n")
	})
	(break true)
})

(console:print "    Found:   " [totalfishcount] " Fish Total\n")
(console:print "    Changed: " [totalpopschanged] " Populations Total\n")
(console:print "    (The above two numbers should match.)\n")
