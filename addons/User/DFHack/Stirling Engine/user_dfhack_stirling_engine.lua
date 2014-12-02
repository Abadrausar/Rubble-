
local buildings = require 'plugins.building-hacks'

buildings.registerBuilding{
	name="STIRLING_ENGINE",
	produce=6,
	gears={{x=0,y=0}},
	animate={
		isMechanical=true,
		frames={
			{{x=0,y=0,42,7,0,0}},
			{{x=0,y=0,15,7,0,0}}
		}
	}
}