
(rubble:dfhack_loadscript "user_dfhack_powered_sawmill.lua")

[rubble:dfhack_powered:wshops "SAWMILL_POWERED" = <struct
	proto=[rubble:dfhack_powered:wshop]
	
	name="Powered Sawmill"
	
	extraitems=`
	[BUILD_ITEM:1:TRAPCOMP:ITEM_TRAPCOMP_LARGESERRATEDDISC:NONE:NONE]
		[CAN_USE_ARTIFACT]
`
>]
