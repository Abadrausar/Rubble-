package main

//import "fmt"
import "strings"
import "strconv"
import "sort"
import "dctech/nca5"

func SetupBuiltins() {
	NewNativeTemplate("!TEMPLATE", tempTemplate)
	NewNativeTemplate("!SCRIPT_TEMPLATE", tempScriptTemplate)
	
	NewNativeTemplate("!SCRIPT", tempScript)
	NewNativeTemplate("SCRIPT", tempScript)
	NewNativeTemplate("#SCRIPT", tempScript)
	
	NewNativeTemplate("ITEM", tempItem)
	NewNativeTemplate("ITEM_RARITY", tempItemRarity)
	NewNativeTemplate("#USES_ITEMS", tempUsesItems)

	NewNativeTemplate("BUILDING_WORKSHOP", tempBuildingWorkshop)
	NewNativeTemplate("BUILDING_FURNACE", tempBuildingFurnace)
	NewNativeTemplate("#USES_BUILDINGS", tempUsesBuildings)
	NewNativeTemplate("REACTION", tempReaction)
	NewNativeTemplate("#USES_REACTIONS", tempUsesReactions)
	NewNativeTemplate("#USES_TECH", tempUsesTech)
	
	NewNativeTemplate("#ADV_TIME", tempAdvTime)
	NewNativeTemplate("#FORT_TIME", tempFortTime)
	
	NewNativeTemplate("REGISTER_ORE", tempRegisterOre)
	NewNativeTemplate("#_REGISTERED_ORES", tempRegisteredOres)
	NewNativeTemplate("REGISTER_REACTION_CLASS", tempRegisterReactionClass)
	NewNativeTemplate("#_REGISTERED_REACTION_CLASSES", tempRegisteredReationClasses)
	NewNativeTemplate("REGISTER_REACTION_PRODUCT", tempRegisterReactionProduct)
	NewNativeTemplate("#_REGISTERED_REACTION_PRODUCTS", tempRegisteredReationProducts)
}

func tempTemplate(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to !TEMPLATE.")
	}
	
	name := params[0]
	text := params[len(params)-1]
	paramNames := params[1:len(params)-1]
	
	parsedParams := make([]*TemplateParam, 0, len(paramNames))
	
	for _, val := range paramNames {
		rtn := new(TemplateParam)
		if strings.Contains(val, "=") {
			parts := strings.SplitN(val, "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val
		parsedParams = append(parsedParams, rtn)
	}
	
	NewUserTemplate(name, text, parsedParams)
	
	return ""
}

func tempScriptTemplate(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to !SCRIPT_TEMPLATE.")
	}
	
	name := params[0]
	text := params[len(params)-1]
	paramNames := params[1:len(params)-1]
	
	parsedParams := make([]*TemplateParam, 0, len(paramNames))
	
	for _, val := range paramNames {
		rtn := new(TemplateParam)
		if strings.Contains(val, "=") {
			parts := strings.SplitN(val, "=", 2)
			rtn.Name = parts[0]
			rtn.Default = parts[1]
			parsedParams = append(parsedParams, rtn)
			continue
		}
		rtn.Name = val
		parsedParams = append(parsedParams, rtn)
	}
	
	NewScriptTemplate(name, text, parsedParams)
	
	return ""
}

func tempScript(params []string) string {
	if len(params) < 1 {
		panic("Wrong number of params to SCRIPT.")
	}
	
	GlobalNCAState.Code.Add(params[0])
	GlobalNCAState.Envs.Add(nca5.NewEnvironment())
	
	if len(params) > 1 {
		GlobalNCAState.AddParams(params[1:]...)
	}
	
	rtn, err := GlobalNCAState.Run()
	if err != nil {
		panic("Script Error: " + err.Error())
	}
	
	GlobalNCAState.Envs.Remove()
	
	if rtn == nil {
		return ""
	}
	return rtn.String()
}

var itemTypes = map[string]bool {
	"AMMO": false,
	"ARMOR": true,
	"DIGGER": false,
	"GLOVES": true,
	"HELM": true,
	"INSTRUMENT": false,
	"PANTS": true,
	"SHIELD": false,
	"SHOES": true,
	"SIEGEAMMO": false,
	"TOOL": false,
	"TOY": false,
	"TRAPCOMP": false,
	"WEAPON": false }

var itemRarities = map[string]int {
	"RARE": 1,
	"UNCOMMON": 2,
	"COMMON": 3,
	"FORCED": 4 }

type itemClassItem struct {
	Name string
	Type string
	Rarity string
}

var itemClasses = make(map[string][]*itemClassItem)
func tempItem(params []string) string {
	if len(params) < 3 {
		panic("Wrong number of params to ITEM.")
	}
	
	rtn := new(itemClassItem)
	rtn.Type = params[0]
	rtn.Name = params[1]
	if _, ok := itemTypes[rtn.Type]; !ok {
		panic("Invalid item type: " + rtn.Type)
	}
	
	
	
	classes := params[2:]
	if itemTypes[rtn.Type] {
		rtn.Rarity = params[2]
		if itemRarities[rtn.Rarity] != 0 {
			classes = params[3:]
			rtn.Rarity = "NULL"
		}
	}
	
	for _, class := range classes {
		if _, ok := itemClasses[class]; !ok {
			itemClasses[class] = make([]*itemClassItem, 0, 10)
		}
		itemClasses[class] = append(itemClasses[class], rtn)
	}
	
	return "[ITEM_" + rtn.Type + ":" + rtn.Name + "]"
}

func tempItemRarity(params []string) string {
	if len(params) != 3 {
		panic("Wrong number of params to ITEM_RARITY.")
	}
	class := params[1]
	if _, ok := itemClasses[class]; !ok {
		panic("Invalid class: " + class)
	}
	
	for _, item := range itemClasses[params[1]] {
		if item.Name == params[0] {
			item.Rarity = params[2]
			return ""
		}
	}
	panic("Invalid item: " + params[0])
}

func tempUsesItems(params []string) string {
	if len(params) < 1 {
		panic("Wrong number of params to #USES_ITEMS.")
	}
	
	permittedItems := make(map[string]*itemClassItem)
	nameList := make([]string, 0, 20)
	for _, class := range params {
		for _, item := range itemClasses[class] {
			if _, ok := permittedItems[item.Name]; ok {
				rtn := new(itemClassItem)
				rtn.Name = item.Name
				rtn.Type = item.Type
				
				if itemRarities[permittedItems[item.Name].Rarity] >= itemRarities[item.Rarity] {
					rtn.Rarity = permittedItems[item.Name].Rarity
				} else {
					rtn.Rarity = item.Rarity
				}
				permittedItems[item.Name] = rtn
				nameList = append(nameList, item.Name)
				continue
			}
			permittedItems[item.Name] = item
			nameList = append(nameList, item.Name)
		}
	}
	
	sort.Strings(nameList)
	out := ""
	for _, name := range nameList {
		item := permittedItems[name]
		if item.Rarity != "NULL" {
			out += "\n\t[" + item.Type + ":" + item.Name + ":" + item.Rarity + "]"
		} else {
			out += "\n\t[" + item.Type + ":" + item.Name + "]"
		}
	}
	return strings.TrimSpace(out)
}

var buildingData = make(map[string][]string)
func tempBuildingWorkshop(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to BUILDING_WORKSHOP.")
	}
	
	for _, name := range params[1:] {
		if _, ok := buildingData[name]; !ok {
			buildingData[name] = make([]string, 0, 5)
		}
		buildingData[name] = append(buildingData[name], params[0])
	}
	return "[BUILDING_WORKSHOP:" + params[0] + "]"
}

func tempBuildingFurnace(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to BUILDING_FURNACE.")
	}
	
	for _, name := range params[1:] {
		if _, ok := buildingData[name]; !ok {
			buildingData[name] = make([]string, 0, 5)
		}
		buildingData[name] = append(buildingData[name], params[0])
	}
	return "[BUILDING_FURNACE:" + params[0] + "]"
}

func tempUsesBuildings(params []string) string {
	if len(params) == 0 {
		panic("Wrong number of params to #USES_TECH_CLASSES.")
	}
	
	// An interesting (ab)use of maps...
	permittedBuildings := make(map[string]bool, 20)
	buildingNames := make([]string, 0, 20)
	for _, tag := range params {
		if _, ok := buildingData[tag]; ok {
			for _, name := range buildingData[tag] {
				if !permittedBuildings[name] {
					permittedBuildings[name] = true
					buildingNames = append(buildingNames, name)
				}
			}
		}
	}
	
	sort.Strings(buildingNames)
	out := ""
	for _, name := range buildingNames {
		out += "\n\t[PERMITTED_BUILDING:" + name + "]"
	}
	return strings.TrimSpace(out)
}

var reactionData = make(map[string][]string)
func tempReaction(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to REACTION.")
	}
	
	for _, class := range params[1:] {
		if _, ok := reactionData[class]; !ok {
			reactionData[class] = make([]string, 0, 5)
		}
		reactionData[class] = append(reactionData[class], params[0])
	}
	return "[REACTION:" + params[0] + "]"
}

func tempUsesReactions(params []string) string {
	if len(params) == 0 {
		panic("Wrong number of params to #USES_TECH_CLASSES.")
	}
	
	// An interesting (ab)use of maps...
	permittedReactions := make(map[string]bool, 20)
	reactionNames := make([]string, 0, 20)
	for _, tag := range params {
		if _, ok := reactionData[tag]; ok {
			for _, name := range reactionData[tag] {
				if !permittedReactions[name] {
					permittedReactions[name] = true
					reactionNames = append(reactionNames, name)
				}
			}
		}
	}
	
	sort.Strings(reactionNames)
	out := ""
	for _, name := range reactionNames {
		out += "\n\t[PERMITTED_REACTION:" + name + "]"
	}
	return strings.TrimSpace(out)
}

func tempUsesTech(params []string) string {
	if len(params) == 0 {
		panic("Wrong number of params to #USES_TECH.")
	}
	
	// An interesting (ab)use of maps...
	permittedReactions := make(map[string]bool, 20)
	reactionNames := make([]string, 0, 20)
	permittedBuildings := make(map[string]bool, 20)
	buildingNames := make([]string, 0, 20)
	for _, tag := range params {
		if _, ok := buildingData[tag]; ok {
			for _, name := range buildingData[tag] {
				if !permittedBuildings[name] {
					permittedBuildings[name] = true
					buildingNames = append(buildingNames, name)
				}
			}
		}
		if _, ok := reactionData[tag]; ok {
			for _, name := range reactionData[tag] {
				if !permittedReactions[name] {
					permittedReactions[name] = true
					reactionNames = append(reactionNames, name)
				}
			}
		}
	}
	
	sort.Strings(buildingNames)
	sort.Strings(reactionNames)
	out := ""
	for _, name := range buildingNames {
		out += "\n\t[PERMITTED_BUILDING:" + name + "]"
	}
	for _, name := range reactionNames {
		out += "\n\t[PERMITTED_REACTION:" + name + "]"
	}
	return strings.TrimSpace(out)
}

var advTime = map[string]int64 {
	"SECOND": 1,
	"SECONDS": 1,
	"MINUTE": 60,
	"MINUTES": 60,
	"HOUR": 3600,
	"HOURS": 3600,
	"DAY": 86400,
	"DAYS": 86400,
	"WEEK": 7 * 86400,
	"WEEKS": 7 * 86400,
	"MONTH": 28 * 86400,
	"MONTHS": 28 * 86400,
	"YEAR": 336 * 86400,
	"YEARS": 336 * 86400 }

func tempAdvTime(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to #ADV_TIME.")
	}
	
	amount, err := strconv.ParseInt(params[0], 0, 64)
	if err != nil {
		panic(err)
	}
	unit := params[1]
	
	if _, ok := advTime[unit]; !ok {
		panic("Invalid unit: " + unit + " to #ADV_TIME.")
	}
	
	amount = amount * advTime[unit]
	if amount < 0 {
		amount = 0
	}
	return strconv.FormatInt(amount, 10)
}

func tempFortTime(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to #ADV_TIME.")
	}
	
	amount, err := strconv.ParseInt(params[0], 0, 64)
	if err != nil {
		panic(err)
	}
	unit := params[1]
	
	if _, ok := advTime[unit]; !ok {
		panic("Invalid unit: " + unit + " to #ADV_TIME.")
	}
	
	amountf := float64(amount) * (float64(advTime[unit]) / 72)
	if amountf < 0 {
		amountf = 0
	}
	return strconv.FormatInt(int64(amountf), 10)
}

var oreData = make(map[string]map[string]int64)
func tempRegisterOre(params []string) string {
	if len(params) != 3 {
		panic("Wrong number of params to REGISTER_ORE.")
	}
	
	ore := params[0]
	metal := params[1]
	amount, err := strconv.ParseInt(params[2], 0, 64)
	if err != nil {
		panic(err)
	}
	
	if _, ok := oreData[ore]; !ok {
		oreData[ore] = make(map[string]int64)
	}
	
	if _, ok := oreData[ore][metal]; !ok {
		oreData[ore][metal] = amount
		return ""
	}
	if amount > oreData[ore][metal] {
		oreData[ore][metal] = amount
	}
	return ""
}

func tempRegisteredOres(params []string) string {
	if len(params) != 1 {
		panic("Wrong number of params to #_REGISTERED_ORES.")
	}
	
	out := ""
	ore := params[0]
	for metal := range oreData[ore] {
		amount := strconv.FormatInt(oreData[ore][metal], 10)
		out += "\n\t[METAL_ORE:" + metal + ":" + amount + "]"
	}
	return out
}

var reactionClassData = make(map[string][]string)
func tempRegisterReactionClass(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to REGISTER_REACTION_CLASS.")
	}
	
	inorganic := params[0]
	class := params[1]
	
	if _, ok := reactionClassData[inorganic]; !ok {
		reactionClassData[inorganic] = make([]string, 0, 5)
		reactionClassData[inorganic] = append(reactionClassData[inorganic], class)
		return ""
	}
	
	for _, val := range reactionClassData[inorganic] {
		if val == class {
			return ""
		}
	}
	reactionClassData[inorganic] = append(reactionClassData[inorganic], class)
	return ""
}

func tempRegisteredReationClasses(params []string) string {
	if len(params) != 1 {
		panic("Wrong number of params to #_REGISTERED_REACTION_CLASSES.")
	}
	
	out := ""
	for _, val := range reactionClassData[params[0]] {
		out += "\n\t[REACTION_CLASS:" + val + "]"
	}
	return out
}

type reactionProduct struct {
	Name string
	Mat string
}

var reactionProductData = make(map[string][]*reactionProduct)
// If for some odd reason more than one product is registered with the same id the first one "wins"
func tempRegisterReactionProduct(params []string) string {
	if len(params) != 3 {
		panic("Wrong number of params to REGISTER_REACTION_PRODUCT.")
	}
	
	inorganic := params[0]
	class := params[1]
	mat := StageParse(params[2])
	
	if _, ok := reactionProductData[inorganic]; !ok {
		reactionProductData[inorganic] = make([]*reactionProduct, 0, 5)
		prod := new(reactionProduct)
		prod.Name = class
		prod.Mat = mat
		reactionProductData[inorganic] = append(reactionProductData[inorganic], prod)
		return ""
	}
	
	for _, val := range reactionProductData[inorganic] {
		if val.Name == class {
			return ""
		}
	}
	prod := new(reactionProduct)
	prod.Name = class
	prod.Mat = mat
	reactionProductData[inorganic] = append(reactionProductData[inorganic], prod)
	return ""
}

func tempRegisteredReationProducts(params []string) string {
	if len(params) != 1 {
		panic("Wrong number of params to #_REGISTERED_REACTION_PRODUCTS.")
	}
	
	out := ""
	for _, val := range reactionProductData[params[0]] {
		out += "\n\t[MATERIAL_REACTION_PRODUCT:" + val.Name + ":" + val.Mat + "]"
	}
	return out
}
