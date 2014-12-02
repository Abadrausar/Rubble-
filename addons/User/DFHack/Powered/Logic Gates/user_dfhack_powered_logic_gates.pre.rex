
(rubble:dfhack_loadscript "user_dfhack_powered_logic_gates.lua")

[rubble:dfhack_powered:wshops "LOGIC_GATES" = <struct
	proto=[rubble:dfhack_powered:wshop]
	
	name="Powered Logic Gate"
	
	outputs=<sarray
		<struct
			proto=[rubble:dfhack_powered:output]
			id="AND"
			name="AND"
		>
		<struct
			proto=[rubble:dfhack_powered:output]
			id="OR"
			name="OR"
		>
		<struct
			proto=[rubble:dfhack_powered:output]
			id="NOT"
			name="NOT"
		>
		<struct
			proto=[rubble:dfhack_powered:output]
			id="XOR"
			name="XOR"
		>
	>
	
	small=true
>]
