
module rubble:macro

# rubble:macro:new is a convenience to help ensure your macro maps have the correct structure.
command rubble:macro:new {
	(ret <map x=0 y=0 z=0 actions=<array>>)
}

# rubble:macro:export returns a string containing the macro.
command rubble:macro:export mak name {
	var out = (str:add [name] "\n")
	
	(foreach [mak actions] block _ action {
		[out = (str:add [out] "\t\t" [action] "\n\tEnd of group\n")]
	})
	
	(ret (str:add [out] "End of macro\n"))
}

# rubble:macro:install writes a macro file to the correct directory.
command rubble:macro:install mak name {
	(axis:write [rubble:fs] (str:add "df:data/init/macros/" [name] ".mak") (rubble:macro:export [mak] [name]))
}

# rubble:macro:setcursor allows you to set the absolute position of the cursor (relative to the origin, eg 0,0,0).
# Cursor movement ordered by this command is fairly optimized, and makes use of all types of cursor movements,
# including diagonal and "fast" movements where appropriate.
# TODO: Better path finder, diagonals are not used everywhere they should be.
command rubble:macro:setcursor mak x y z {
	# Some local "commands" that we will need in a bit.
	var abs = block i {(if (int:lt [i] 0){(int:sub 0 [i])}{[i]})}
	var move = block cur dest curO destO dir1 dir2 actions {
		(if (int:eq [cur] [dest]) {}{
			(loop {
				(if (int:eq ([abs] (int:sub [dest] [cur])) ([abs] (int:sub [destO] [curO]))) {
					(ret [cur])
				})
				
				(if (int:lt [dest] [cur]) {
					(if (expr "a-b >= c" [cur] [dest] 10) {
						[actions append = (str:add "CURSOR_" [dir1] "_FAST")]
						[cur = (int:sub [cur] 10)]
					}{
						[actions append = (str:add "CURSOR_" [dir1])]
						[cur = (int:sub [cur] 1)]
					})
					(break true)
				}{
					(break false)
				})
			})
			(loop {
				(if (int:eq ([abs] (int:sub [dest] [cur])) ([abs] (int:sub [destO] [curO]))) {
					(ret [cur])
				})
				
				(if (int:gt [dest] [cur]) {
					(if (expr "a-b >= c" [dest] [cur] 10) {
						[actions append = (str:add "CURSOR_" [dir2] "_FAST")]
						[cur = (int:add [cur] 10)]
					}{
						[actions append = (str:add "CURSOR_" [dir2])]
						[cur = (int:add [cur] 1)]
					})
					(break true)
				}{
					(break false)
				})
			})
		})
		(ret [cur])
	}
	
	# First get to the correct Z-level, It really doesn't matter when this is done
	# so let's get it out of the way.
	(if (int:eq [mak z] [z]) {}{
		(loop {
			(if (int:lt [z] [mak z]) {
				[mak actions append = "CURSOR_DOWN_Z"]
				[mak z = (int:sub [mak z] 1)]
				(break true)
			}{
				(break false)
			})
		})
		(loop {
			(if (int:gt [z] [mak z]) {
				[mak actions append = "CURSOR_UP_Z"]
				[mak z = (int:add [mak z] 1)]
				(break true)
			}{
				(break false)
			})
		})
	})
	
	# Now we move along whichever axis has the largest total movement
	# until both axises have the same number of tiles to go
	# (or we have moved as far as needed on this axis)
	(if (int:gt ([abs] (int:sub [y] [mak y])) ([abs] (int:sub [x] [mak x]))) {
		[mak y = ([move] [mak y] [y] [mak x] [x] "UP" "DOWN" [mak actions])]
	}{
		[mak x = ([move] [mak x] [x] [mak y] [y] "LEFT" "RIGHT" [mak actions])]
	})
	
	# As fast movements may have caused us to overshoot the point where diagonal movements will
	# put us on our target we may need to correct a little bit
	(if (int:gt ([abs] (int:sub [y] [mak y])) ([abs] (int:sub [x] [mak x]))) {
		[mak y = ([move] [mak y] [y] [mak x] [x] "UP" "DOWN" [mak actions])]
	}{
		[mak x = ([move] [mak x] [x] [mak y] [y] "LEFT" "RIGHT" [mak actions])]
	})
	
	# Then we move diagonally to our destination (if we are not already there from the previous step)
	(if (expr "a != b && c != d" [mak x] [x] [mak y] [y]) {
		# Up left
		(loop {
			(if (expr "a > b && c > d" [mak x] [x] [mak y] [y]) {
				(if (expr "a-b >= e && c-d >= e" [mak x] [x] 10 [mak y] [y] 10) {
					[mak actions append = "CURSOR_UPLEFT_FAST"]
					[mak x = (int:sub [mak x] 10)]
					[mak y = (int:sub [mak y] 10)]
				}{
					[mak actions append = "CURSOR_UPLEFT"]
					[mak x = (int:sub [mak x] 1)]
					[mak y = (int:sub [mak y] 1)]
				})
				(break true)
			}{
				(break false)
			})
		})
		
		# Down left
		(loop {
			(if (expr "a > b && c < d" [mak x] [x] [mak y] [y]) {
				(if (expr "a-b >= e && c-d >= e" [mak x] [x] 10 [y] [mak y] 10) {
					[mak actions append = "CURSOR_DOWNLEFT_FAST"]
					[mak x = (int:sub [mak x] 10)]
					[mak y = (int:add [mak y] 10)]
				}{
					[mak actions append = "CURSOR_DOWNLEFT"]
					[mak x = (int:sub [mak x] 1)]
					[mak y = (int:add [mak y] 1)]
				})
				(break true)
			}{
				(break false)
			})
		})
		
		# Up right
		(loop {
			(if (expr "a < b && c > d" [mak x] [x] [mak y] [y]) {
				(if (expr "a-b >= e && c-d >= e" [x] [mak x] 10 [mak y] [y] 10) {
					[mak actions append = "CURSOR_UPRIGHT_FAST"]
					[mak x = (int:add [mak x] 10)]
					[mak y = (int:sub [mak y] 10)]
				}{
					[mak actions append = "CURSOR_UPRIGHT"]
					[mak x = (int:add [mak x] 1)]
					[mak y = (int:sub [mak y] 1)]
				})
				(break true)
			}{
				(break false)
			})
		})
		
		# Down right
		(loop {
			(if (expr "a < b && c < d" [mak x] [x] [mak y] [y]) {
				(if (expr "a-b >= e && c-d >= e" [x] [mak x] 10 [y] [mak y] 10) {
					[mak actions append = "CURSOR_DOWNRIGHT_FAST"]
					[mak x = (int:add [mak x] 10)]
					[mak y = (int:add [mak y] 10)]
				}{
					[mak actions append = "CURSOR_DOWNRIGHT"]
					[mak x = (int:add [mak x] 1)]
					[mak y = (int:add [mak y] 1)]
				})
				(break true)
			}{
				(break false)
			})
		})
	})
}

# rubble:macro:movecursor allows you to move the cursor relative to it's current location
command rubble:macro:movecursor mak x y z {
	(rubble:macro:setcursor [mak] (int:add [x] [mak x]) (int:add [y] [mak y]) (int:add [z] [mak z]))
}

# rubble:macro:designate takes a palette of actions and a 2D data array and generates a macro,
# see the example file.
command rubble:macro:designate mak palette data {
	var h = (int:sub (len [data]) 1)
	var w = (int:sub (len [data 0]) 1)
	
	(for 0 [h] 1 block y {
		(for 0 [w] 1 block x {
			var actions = [palette [data [y] [x]]]
			
			(if (isnil [actions]) {
				# The current value has no action in the palette, do nothing.
			}{
				(rubble:macro:setcursor [mak] [x] [y] [mak z])
				
				(foreach [actions] block _ action {
					[mak actions append = [action]]
					(break true)
				})
			})
			
			(break true)
		})
		(break true)
	})
}

# rubble:macro:designate_rle is just like the normal rubble:macro:designate except it assumes that
# any action list in the palette that only has one action should be consolidated with
# adjacent actions of the same type (only works with designate type actions, like dig and channel).
command rubble:macro:designate_rle mak palette data {
	var h = (int:sub (len [data]) 1)
	var w = (int:sub (len [data 0]) 1)
	
	var lastaction
	var aStart = 0
	(for 0 [h] 1 block y {
		(for 0 [w] 1 block x {
			var actIndex = [data [y] [x]]
			
			(if (expr "a != b && !c" [actIndex] [lastaction] (isnil [lastaction])) {
				(rubble:macro:setcursor [mak] [aStart] [y] [mak z])
				[mak actions append = [palette [lastaction] 0]]
				[mak actions append = "SELECT"]
				(rubble:macro:setcursor [mak] (int:sub [x] 1) [y] [mak z])
				[mak actions append = "SELECT"]
				[lastaction = nil]
				[aStart = 0]
			})
			
			(if (isnil [palette [actIndex]]) {
				# The current value has no action in the palette, do nothing.
			}{
				(if (int:eq (len [palette [actIndex]]) 1) {
					(if (int:eq [lastaction] [actIndex]) {}{
						[lastaction = [actIndex]]
						[aStart = [x]]
					})
				}{
					(rubble:macro:setcursor [mak] [x] [y] [mak z])
					
					(foreach [palette [actIndex]] block _ action {
						[mak actions append = [action]]
						(break true)
					})
				})
			})
			
			(break true)
		})
		(if (isnil [lastaction]) {}{
			(rubble:macro:setcursor [mak] [aStart] [y] [mak z])
			[mak actions append = [palette [lastaction] 0]]
			[mak actions append = "SELECT"]
			(rubble:macro:setcursor [mak] [w] [y] [mak z])
			[mak actions append = "SELECT"]
			[lastaction = nil]
			[aStart = 0]
		})
		(break true)
	})
}
