
# This is a template generic patch file.

# 
# To install:
#	Copy this file into your addons folder and name it "patch_a_5.2.load.rex" (the "load.rex" extension is important!).
# 

# The Rubble version this patch is for.
var version = "5.2"

# Format should be "<sequence letter> (<short description>)"
var name = "A (Example patch)"

# Patch structure prototype, do not modify!
var patch = <map
	# The name of the effected addon
	addon=""
	
	# The file to modify in the addon
	file=""
	
	# The actions, only set one!
	# diff should be a simple difference patch
	# replace should be the entire new file text (encoded with the Rubble file encoder)
	diff=nil
	replace=nil
>

# An array of patch structures, modify as needed.
var patches = <array
	# This example updates the 40.6 base to work with 40.8
	# Obviously this is way out of date, but it still makes a good example.
	<struct proto=[patch]
		addon="Base"
		file="plant_standard.txt"
		diff="@@ -1237,14 +1237,16 @@
 	[GROWTH:FLOWERS]
 		[GROWTH_NAME:acacia flower:STP]
 		[GROWTH_ITEM:PLANT_GROWTH:NONE:LOCAL_PLANT_MAT:FLOWER]
 		[GROWTH_DENSITY:1000]
+		[GROWTH_HOST_TILE:BRANCHES_AND_TWIGS]
 		[GROWTH_TIMING:60000:119999]
 		[GROWTH_PRINT:5:5:7:0:1:60000:119999:2]
 	[GROWTH:POD]
 		[GROWTH_NAME:acacia seed pod:STP]
 		[GROWTH_ITEM:PLANT_GROWTH:NONE:LOCAL_PLANT_MAT:STRUCTURAL]
 		[GROWTH_DENSITY:1000]
+		[GROWTH_HOST_TILE:BRANCHES_AND_TWIGS]
 		[GROWTH_TIMING:120000:200000]
 		[GROWTH_PRINT:'%':'%':2:0:1:120000:200000:3]
 		[GROWTH_HAS_SEED]
 }"
	>
	
	# You can add as many patch objects as needed.
>

# Code below this point is common to all patches, you should not need to modify it.

(if (str:cmp [rubble:version] [version]) {}{
	(rubble:abort `This patch cannot be applied to this Rubble version, remove the patch file and try again.`)
})

(console:print "    Installing Patch " [name] "\n    for Rubble v" [version] ":\n")
(foreach [patches] block _ patch {
	(console:print '      "' [patch addon] '": ' [patch file] '\n')
	
	var addon = (rubble:fetchaddon [patch addon])
	(if (isnil [addon]) {
		(rubble:abort `This patch refers to a non-existent addon, remove the patch file and try again.`)
	})
	
	(if (isnil [patch diff]) {}{
		var content = (genii:bytes_string [addon Files [patch file] Content])
		[content = (rubble:patch [content] [patch diff])]
		[addon Files [patch file] Content = (genii:string_bytes [content])]
	})
	
	(if (isnil [patch replace]) {}{
		[addon Files [patch file] Content = (genii:string_bytes (rubble:decompress [patch replace]))]
	})
})
