
(rubble:dfhack_loadscript "user_dfhack_powered_block_cutter.lua")

[rubble:dfhack_powered:wshops "BLOCK_CUTTER_POWERED" = <struct
	proto=[rubble:dfhack_powered:wshop]
	
	name="Powered Block Cutter"
	outputs=<sarray
		<struct
			proto=[rubble:dfhack_powered:output]
			id="BLACK" name="black">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="BROWN" name="brown">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="LIGHT_BLUE" name="light blue">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="LIGHT_CYAN" name="light cyan">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="LIGHT_GREEN" name="light green">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="LIGHT_GREY" name="light grey">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="LIGHT_RED" name="light red">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="LIGHT_VIOLET" name="light violet">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="DARK_BLUE" name="dark blue">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="DARK_CYAN" name="dark cyan">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="DARK_GREEN" name="dark green">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="DARK_GREY" name="dark grey">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="DARK_RED" name="dark red">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="DARK_VIOLET" name="dark violet">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="WHITE" name="white">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="YELLOW" name="yellow">
		<struct
			proto=[rubble:dfhack_powered:output]
			id="INPUT_MAT" name="input material">
	>
	
	extraitems=`
	[BUILD_ITEM:1:TRAPCOMP:ITEM_TRAPCOMP_LARGESERRATEDDISC:NONE:NONE]
		[CAN_USE_ARTIFACT]
`
>]
