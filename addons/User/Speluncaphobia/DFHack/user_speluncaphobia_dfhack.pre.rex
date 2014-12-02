
module rubble:speluncaphobia

var rubble:speluncaphobia:entities = (eval (rubble:configvar "USER_SPELUNCAPHOBIA_DFHACK_ENTITIES"))

(if (type [rubble:speluncaphobia:entities] index) {}{
	[rubble:speluncaphobia:entities = <array SAURIAN HUMAN>]
})

(rubble:template "SPELUNCAPHOBIA_HOOKS" {
	var out = ""
	(foreach [rubble:speluncaphobia:entities] block _ entity {
		[out = (str:add [out] "{REACTION_ADD_CLASS;ADDON_HOOK_" [entity] "}")]
		(break true)
	})
	(rubble:stageparse [out])
})

(rubble:dfhack_runcommand `modtools/item-trigger -itemType "ITEM_SPELUNCAPHOBIA_TORCH" -onEquip -command [ modtools/add-syndrome -syndrome "speluncaphobia immune" -resetPolicy DoNothing -target \\\\UNIT_ID ]`)

(rubble:dfhack_runcommand `modtools/item-trigger -itemType "ITEM_SPELUNCAPHOBIA_TORCH" -onUnequip -command [ modtools/add-syndrome -syndrome "speluncaphobia immune" -target \\\\UNIT_ID -eraseAll ]`)
