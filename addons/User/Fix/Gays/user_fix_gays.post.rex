
# This config var is for experts only!
# Only set it to a map object literal.
# Example: -config="USER_FIX_GAYS=<map DWARF=true HUMAN=true>"
var gays = (eval (rubble:configvar "USER_FIX_GAYS"))
(if (type [gays] index) {}{
	[gays = <map 
		DWARF=true
		HUMAN=true
		# Elves, goblins, and other animals can keep their perverted habits...
	>]
})

var foundgay = false

(console:print "    Fixing Specified Gay Creatures...\n")
(foreach [rubble:raws] block name content {
	(if (str:cmp "creature_" (str:left [name] 9)) {
		(console:print "      " [name] "\n")
		
		[rubble:raws [name] = (df:raw:walk [content] block tag {
			(if (str:cmp [tag id] "CREATURE") {
				(if (int:eq (len [tag]) 1) {
					(if [gays [tag 0]]{
						[foundgay = true]
					}{
						[foundgay = false]
					})
				}{
					(rubble:abort "Error: invalid param count to CREATURE raw tag in last file.")
				})
			})
			
			(if (bool:and (str:cmp [tag id] "MALE") [foundgay]) {
				[tag replace = "[MALE]\n\t\t[ORIENTATION:MALE:100:0:0]\n\t\t[ORIENTATION:FEMALE:0:25:100]"]
			})
			
			(if (bool:and (str:cmp [tag id] "FEMALE") [foundgay]) {
				[tag replace = "[FEMALE]\n\t\t[ORIENTATION:MALE:0:25:100]\n\t\t[ORIENTATION:FEMALE:100:0:0]"]
			})
			
			(break true)
		})]
	})
	(break true)
})
