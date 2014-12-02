
var aquifercount = 0
var aquifercounttotal = 0

(foreach [rubble:raws] block name content {
	(if (str:cmp "inorganic_" (str:left [name] 10)) {
		(console:print "    " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "AQUIFER") {
				[tag disable = true]
				[aquifercount = (int:add [aquifercount] 1)]
			})
			
			(break true)
		})]
		(console:print "      Found Aquifers : " [aquifercount] "\n")
		[aquifercounttotal = (int:add [aquifercounttotal] [aquifercount])]
		[aquifercount = 0]
	})
	(break true)
})
(console:print "    Found Aquifers Total : " [aquifercounttotal] "\n")
