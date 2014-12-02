
# Adds the BUILD_KEY template

# returns nil for invalid keys
var rubble:bkey_valid = <map>
# returns the key index in the order array
var rubble:bkey_index = <map>
# the key order, used for picking alternate keys
var rubble:bkey_order = <array>

# Three cheers for foreach!
<sarray "A" "B" "C" "D" "E" "F" "G" "H" "I" "J" "K" "L" "M" "N" "O" "P" "Q" "R" "S" "T" "U" "V" "W" "X" "Y" "Z"> 
(foreach (nop) block _ char {
	(foreach <sarray "" "SHIFT_" "CTRL_" "ALT_"> block _ prefix {
		var full_key = (str:add [prefix] [char])
		[rubble:bkey_valid [full_key] = true]
		[rubble:bkey_index [full_key] = (len [rubble:bkey_order])]
		[rubble:bkey_order append = [full_key]]
		(break true)
	})
	(break true)
})

var rubble:bkey_used_w = <map
	"E" = true
	"Q" = true
	"SHIFT_M" = true
	"O" = true
	"K" = true
	"B" = true
	"C" = true
	"F" = true
	"V" = true
	"J" = true
	"M" = true
	"U" = true
	"N" = true
	"R" = true
	"S" = true
	"T" = true
	"L" = true
	"W" = true
	"Z" = true
	"H" = true
	"Y" = true
	"D" = true
>

var rubble:bkey_used_f = <map
	"W" = true
	"S" = true
	"G" = true
	"K" = true
	"L" = true
	"A" = true
	"N" = true
>

(rubble:template "@BUILD_KEY" block key furnace=false {
	# Remember, all script vars are "by reference".
	# Sometimes I'm amazed this stuff works...
	var used_list
	(if [furnace] {
		[used_list = [rubble:bkey_used_f]]
	}{
		[used_list = [rubble:bkey_used_w]]
	})
	
	(if [rubble:bkey_valid [key]] {}{
		(rubble:abort (str:add "Error: Invalid key: " [key] " passed to @BUILD_KEY"))
	})
	
	# if key is in use, pick next key that is not being used
	(if [used_list [key]] {
		var wrap = false
		var next_index = (int:add [rubble:bkey_index [key]] 1)
		(if (int:eq [next_index] (len [rubble:bkey_order])) {
			[wrap = true]
			[next_index = 0]
		})
		var next_key = [rubble:bkey_order [next_index]]
		
		(loop {
			(if [used_list [next_key]] {
				[next_index = (int:add [rubble:bkey_index [next_key]] 1)]
				(if (int:eq [next_index] (len [rubble:bkey_order])) {
					(if [wrap] {
						(rubble:abort "Error: @BUILD_KEY Cannot find valid key, all keys used?!!")
					})
					[wrap = true]
					[next_index = 0]
				})
				[next_key = [rubble:bkey_order [next_index]]]
				(break true)
			}{
				[key = [next_key]]
				(break false)
			})
		})
	})
	
	[used_list [key] = true]
	(str:add "[BUILD_KEY:CUSTOM_" [key] "]")
})
