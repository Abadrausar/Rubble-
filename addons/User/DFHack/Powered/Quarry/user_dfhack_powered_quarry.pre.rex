
(rubble:dfhack_loadscript "user_dfhack_powered_quarry.lua")

[rubble:dfhack_powered:wshops "QUARRY_POWERED" = <struct
	proto=[rubble:dfhack_powered:wshop]
	
	name="Powered Quarry"
	
	outputs=<sarray
		<struct
			proto=[rubble:dfhack_powered:output]
			id="CLAY"
			name="clay"
		>
		<struct
			proto=[rubble:dfhack_powered:output]
			id="SAND"
			name="sand"
		>
	>
>]
