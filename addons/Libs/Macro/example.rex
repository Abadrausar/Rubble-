
# Examples (in other words, my testing code :p)
var mak = (rubble:macro:new)

(rubble:macro:designate [mak]
	# The palette
	<sarray
		nil
		<sarray "DESIGNATE_DIG" "SELECT" "SELECT">
		<sarray "DESIGNATE_CHANNEL" "SELECT" "SELECT">
		<sarray "DESIGNATE_STAIR_UP" "SELECT" "SELECT">
	>
	
	# The data array, same basic format as returned by png:load
	# so designating based on an image is easy.
	<sarray
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 0 0 3 0 0 0 0 0 0 0 0>
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 0 0 0 0 3 0 0 0 0 0 0>
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 0 0 0 0 0 0 0 0 3 0 0>
		<sarray 0 0 0 0 0 0 0 0 0 0 3>
		<sarray 0 0 0 0 0 0 0 0 0 0 0>
		<sarray 2 1 2 1 2 1 2 1 2 1 2>
	>
)


#(rubble:macro:designate_rle [mak]
#	<sarray
#		nil
#		<sarray "DESIGNATE_DIG">
#	>
#	<sarray
#		<sarray 0 0 0 0 0 0 0 0 0 0 0>
#		<sarray 0 0 1 1 1 1 1 1 1 1 1>
#		<sarray 0 1 1 0 0 0 1 1 0 0 0>
#		<sarray 0 0 0 0 0 0 0 0 0 0 0>
#	>
#)

# For use with rexsh
#(fileio:write "Test.mak" (rubble:macro:export [mak] "Test"))

# For use with Rubble
(axis:write [rubble:fs] "df:data/init/macros/Test.mak" (rubble:macro:export [mak] "Test"))
