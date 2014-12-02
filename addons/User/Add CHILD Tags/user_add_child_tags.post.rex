
var noCHILD = <map>
var curCreature = ""
var haschild = false

var creaturecount = 0
var creaturecounttotal = 0
var nochildcount = 0
var nochildcounttotal = 0
var fixedcount = 0
var fixedcounttotal = 0

# This loop makes a list of all creatures that do NOT have a CHILD tag...
(console:print "    Generating Table...\n")
(console:print "    (DOES_NOT_EXIST, NOT_LIVING, and EQUIPMENT_WAGON count as CHILD tags)\n")
(console:print "    (Vermin are also skipped.)\n")
(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		(df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				[creaturecount = (int:add [creaturecount] 1)]
				(if (int:eq (len [tag]) 1) {
					(if (bool:not (str:cmp [curCreature] "")){
						(if (bool:not [haschild]) {
							[noCHILD [curCreature] = true]
							[nochildcount = (int:add [nochildcount] 1)]
						})
					})
					[curCreature = [tag 0]]
					[haschild = false]
				}{
					(rubble:abort "Error: invalid param count to CREATURE raw tag in last file.")
				})
			})
			
			(if (str:cmp [tag id] "CHILD") {
				[haschild = true]
			})
			
			(if (str:cmp [tag id] "DOES_NOT_EXIST") {
				[haschild = true]
			})
			
			(if (str:cmp [tag id] "EQUIPMENT_WAGON") {
				[haschild = true]
			})
			
			(if (str:cmp [tag id] "NOT_LIVING") {
				[haschild = true]
			})
			
			(if (str:cmp [tag id] "VERMIN_GROUNDER") {
				[haschild = true]
			})
			
			(if (str:cmp [tag id] "VERMIN_SOIL") {
				[haschild = true]
			})
			
			(if (str:cmp [tag id] "VERMIN_SOIL_COLONY") {
				[haschild = true]
			})
			
			(break true)
		})
		
		# this should take care of the last creature in the file
		(if (bool:not (str:cmp [curCreature] "")){
			(if (bool:not [haschild]) {
				[noCHILD [curCreature] = true]
				[nochildcount = (int:add [nochildcount] 1)]
			})
		})
		(console:print "        Found Creatures  : " [creaturecount] "\n")
		(console:print "        Had no CHILD tag : " [nochildcount] "\n")
		[creaturecounttotal = (int:add [creaturecounttotal] [creaturecount])]
		[nochildcounttotal = (int:add [nochildcounttotal] [nochildcount])]
		[creaturecount = 0]
		[nochildcount = 0]
		[curCreature = ""]
		[haschild = false]
	})
	(break true)
})

(console:print "    Fixing raws...\n")

# and this loop fixes them.
(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				(if (int:eq (len [tag]) 1) {
					(if (exists [noCHILD] [tag 0]) {
						[tag append = "\n\t[CHILD:1]"]
						[fixedcount = (int:add [fixedcount] 1)]
					})
				}{
					(rubble:abort "Error: invalid param count to CREATURE raw tag in last file.")
				})
			})
			
			(break true)
		})]
		(console:print "        Fixed Creatures : " [fixedcount] "\n")
		[fixedcounttotal = (int:add [fixedcounttotal] [fixedcount])]
		[fixedcount = 0]
	})
	(break true)
})
(console:print "    Found Creatures total  : " [creaturecounttotal] "\n")
(console:print "    Had no CHILD tag total : " [nochildcounttotal] "\n")
(console:print "    Fixed total            : " [fixedcounttotal] "\n")
(console:print "    (The last two numbers should match.)\n")