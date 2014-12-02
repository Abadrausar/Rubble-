
Better Dorfs - Like Masterwork but smaller :P

Better Dorfs is my stab at overhauling dwarf fortress without getting too far from vanilla.

==============================================
Overview:
==============================================
Better Dorfs adds a bunch of industries as well as a new (possibly playable) race, the saurians.
Saurians have some special items/materials/plants that are unavailable to dwarves while dwarves have some special industries that saurians can't use. Most features can be toggled on or off in the configuration file prior to generating the raws.

==============================================
Install:
==============================================
	First install Blast as detailed here <Blast url here>
	Back up the "base" folder from Blast, then delete "base/objects/entity_default.txt".
	Now copy all files/folders from the Better Dorfs download into the base folder, overwriting as needed.
	You should new have a copy of this readme in the base folder along side "base.pm". If these two files are not in the base folder then you copied something to the wrong place, start over.
	Now you may edit "base/config.txt" to match your preferences, almost everything in Better Dorfs has an option to toggle it on or off here.
	Run Blast! This will generate the raws for Better Dorfs using the options you specified. When Blast is finished (with no errors) you may run Dwarf Fortress and generate a world. 
	
	If Blast runs into any errors during generation please report them in the thread dedicated to this package NOT in the Blast thread (unless you are sure the problem is Blast and not Better Dorfs).

==============================================
Plants:
==============================================
	If Harder farming is on:
		All underground plants take two seasons to grow and can grow all year round.
		All above ground plants take 1 season to grow and can grow every season except winter.
	
	Wire Weed - Saurians only
		Wire weed is a strong fibrous vine used by saurians to make clothing and armor. To make armor sheets of wire weed cloth are glued together with leather to reinforce it. The resulting leather can be used to make armor at a leather worker. 
		Wire weed takes about 3 seasons to grow and may be grown in winter.
	
	Stick Reed - Saurians only
		Stick reed has a sticky sap that can be refined into a glue at the laminator or it may be bundled into wicker at the basketry.
	
	Young Trees - many, many kinds
		Farmable wood. They must be processed before use at the sawmill. Surface trees take 1 season to grow (but they can't grow in winter), underground trees take 2 or 4 seasons and may grow all year round.	
	
==============================================
Major Materials:
==============================================
	Glue (3 kinds for dwarves, 4 for saurians)
		Glue is made at a laminators shop from hide, hoof, bone or stick reed (if playing saurians). Glue is used for laminating various materials together.
	
	Laminated Log
		Made at a laminators shop from glue and planks, laminated logs are more valuable than normal wood.
	
	Wicker
		Wicker is a light material made from woven plants at the basketry. Good for furniture if no clay is available.
	
	Wire Weed Fiber Cloth - Saurians only
		Wire weed fiber cloth is good for protective clothing and laminating with leather to make wire weed reinforced leather.
		
	Wire Weed Reinforced Leather - Saurians only
		The queen of armor materials (the king is water steel) wire weed reinforced leather is equal to iron in strength but much lighter.
		
	Water Steel
		Water steel (a name for damascus steel) is the best weapon material other than candy because it takes an edge twice as sharp as normal steel. Water steel may only be made by/bought from the saurians.
	
==============================================
Buildings:
==============================================
	Kiln
		There are many more things to be made from clay here. Glazing is done via powders made from various minerals or ash at the quern or millstone.
	
	Batch Kiln - 8 blocks
		The Batch Kiln has exactly the same reactions as the kiln except they all take 4X the materials and produce 5X the output.
		The collect clay order must still be issued from the normal kiln and you can not glaze items at the batch kiln.
		
	Block Kiln - 4 blocks
		Here you can produce two kinds of blocks in 16 different colors.
		Normal blocks produce 10 blocks from 5 clay.
		Fine blocks produce 5 more valuable blocks from 5 clay.
		Note that if you just want normal (uncolored) blocks, the batch kiln is a little more efficient.
	
	Sawmill - 4 blocks, 2 mechanisms, 1 serrated disk
		You can cut wood into planks or scrap and process young trees into wood or scrap here. Scrap can be burned into charcoal or ashes here as well. Remember to bring a serrated disk!

	Laminator - 4 blocks
		This workshop is where everything involving glue takes place. Here you can make laminated wood from planks. Note that jobs here use the strand extractor skill.
	
	Basketry - 1 generic building mat
		Here you can bundle some plants into bales and use bales to make various items, mostly furniture.
		Note that wicker is treated like stone for stockpiling purposes, and bales are blocks.
		The bundleable plants are: longland grass, swamp grass, and stick reed.
	
	Bonfire - nothing
		Start two kinds of fires, a short wood fueled one or a long wood/charcoal fueled one.

==============================================
Special thanks to the makers of:
==============================================
	Masterwork
	Kobold Camp
	Dwarven Lamination
	Expanded Glazes
	Broken Arrow
	
	These mods provided many ideas that I then reworked into the Better Dorfs you have today.

==============================================
BUGS:
==============================================
	Underground trees do not show on the list of plants at a (underground) farm plot even if you have seeds, I have no idea why this is.
	
	Please report any (other) abnormal behavior.

==============================================
Change Log:
==============================================
v1.0 for DF 34.11
	Initial release