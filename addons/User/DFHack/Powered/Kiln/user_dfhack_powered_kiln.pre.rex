
(rubble:dfhack_loadscript "user_dfhack_powered_kiln.lua")

[rubble:dfhack_powered:wshops "KILN_POWERED" = <struct
	proto=[rubble:dfhack_powered:wshop]
	
	name="Powered Kiln"
	
	outputs=<sarray
		<struct
			proto=[rubble:dfhack_powered:output]
			id="CERAMIC"
			name="ceramic boulders"
		>
		<struct
			proto=[rubble:dfhack_powered:output]
			id="PEARLASH"
			name="pearlash"
		>
	>
>]
