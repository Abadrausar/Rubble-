
(rubble:dfhack_loadscript "user_dfhack_powered_wood_furnace.lua")

[rubble:dfhack_powered:wshops "WOOD_FURNACE_POWERED" = <struct
	proto=[rubble:dfhack_powered:wshop]
	
	name="Powered Wood Furnace"
	outputs=<sarray
		<struct
			proto=[rubble:dfhack_powered:output]
			id="CHARCOAL"
			name="charcoal"
		>
		<struct
			proto=[rubble:dfhack_powered:output]
			id="ASH"
			name="ash"
		>
	>
>]
