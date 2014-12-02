package main

import "fmt"
import "strings"
import "strconv"
import "sort"
import "regexp"

func SetupBuiltins() {
	NewNativeTemplate("!TEMPLATE", tempTemplate)
	NewNativeTemplate("ONCE", tempOnce)
	NewNativeTemplate("STATIC", tempStatic)
	NewNativeTemplate("COMMENT", tempComment)
	NewNativeTemplate("C", tempComment)
	NewNativeTemplate("VOID", tempVoid)
	
	NewNativeTemplate("SET", tempSet)
	NewNativeTemplate("IF", tempIf)
	
	// Not very useful yet, needs custom nca commands
	NewNativeTemplate("!SCRIPT", tempScript)
	NewNativeTemplate("SCRIPT", tempScript)
	NewNativeTemplate("#SCRIPT", tempScript)
	
	NewNativeTemplate("#ADVENTURE_TIER", tempAdventureTier)
	
	// Replace?
	NewNativeTemplate("ITEM_CLASS", tempItemClass)
	NewNativeTemplate("#USES_ITEM_CLASSES", tempUsesItemClasses)
	NewNativeTemplate("TECH_CLASS", tempTechClass)
	NewNativeTemplate("#USES_TECH_CLASSES", tempUsesTechClasses)
	
	NewNativeTemplate("#ADV_TIME", tempAdvTime)
	NewNativeTemplate("#FORT_TIME", tempFortTime)
	
	NewNativeTemplate("SHARED_INORGANIC", tempSharedInorganic)
	NewNativeTemplate("SHARED_MATERIAL_TEMPLATE", tempSharedInorganic)
	
	NewNativeTemplate("REGISTER_ORE", tempRegisterOre)
	NewNativeTemplate("#_REGISTERED_ORES", tempRegisteredOres)
	NewNativeTemplate("REGISTER_REACTION_CLASS", tempRegisterReactionClass)
	NewNativeTemplate("#_REGISTERED_REACTION_CLASSES", tempRegisteredReationClasses)
	NewNativeTemplate("REGISTER_REACTION_PRODUCT", tempRegisterReactionProduct)
	NewNativeTemplate("#_REGISTERED_REACTION_PRODUCTS", tempRegisteredReationProducts)
	
	// Replace
	NewNativeTemplate("SHARED_ITEM", tempSharedItem)
}

func tempTemplate(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to !TEMPLATE.")
	}
	
	name := params[0]
	text := strings.TrimSpace(params[len(params)-1])
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

var onceData = make(map[string]bool)
func tempOnce(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to ONCE.")
	}
	
	if onceData[params[0]] {
		return ""
	}
	onceData[params[0]] = true
	return StageParse(strings.TrimSpace(params[1]))
}

var staticData = make(map[string]string)
func tempStatic(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to STATIC.")
	}
	
	if staticData[params[0]] != "" {
		return staticData[params[0]]
	}
	rtn := StageParse(strings.TrimSpace(params[1]))
	staticData[params[0]] = rtn
	return rtn
}

func tempComment(params []string) string {
	return ""
}

func tempVoid(params []string) string {
	for _, val := range params {
		StageParse(val)
	}
	return ""
}

var varNameValidateRegEx = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
func tempSet(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to SET.")
	}
	
	if !varNameValidateRegEx.MatchString(params[0]) {
		panic("Variable name supplied to SET is invalid.")
	}
	
	VariableData[params[0]] = params[1]
	return ""
}

func tempIf(params []string) string {
	if len(params) != 3 && len(params) != 4 {
		panic("Wrong number of params to IF.")
	}
	
	if strings.ToLower(params[0]) == strings.ToLower(params[1]) {
		return StageParse(params[2])
	}
	if len(params) == 4 {
		return StageParse(params[3])
	}
	return ""
}

func tempScript(params []string) string {
	if len(params) < 1 {
		panic("Wrong number of params to SCRIPT.")
	}
	
	GlobalNCAState.Code.Add(params[0])
	
	if len(params) > 1 {
		GlobalNCAState.AddParams(params[1:]...)
	}
	
	rtn, err := GlobalNCAState.Run()
	if err != nil {
		panic("Script Error:" + err.Error())
	}
	
	GlobalNCAState.DeleteMap("params") // Something of a hack to keep from trashing the env
	return rtn.String()
}

var advTier = 0
func tempAdventureTier(params []string) string {
	rtn := fmt.Sprint("[ADVENTURE_TIER:", advTier, "]")
	advTier++
	return rtn
}

var itemTypes = map[string]bool {
	"AMMO": true,
	"ARMOR": true,
	"DIGGER": true,
	"GLOVES": true,
	"HELM": true,
	"INSTRUMENT": true,
	"PANTS": true,
	"SHIELD": true,
	"SHOES": true,
	"SIEGEAMMO": true,
	"TOOL": true,
	"TOY": true,
	"TRAPCOMP": true,
	"WEAPON": true }

var itemRarities = map[string]int {
	"RARE": 1,
	"UNCOMMON": 2,
	"COMMON": 3,
	"FORCED": 4 }

type itemClassItem struct {
	Name string
	Type string
	Rarity string
	Tags []string
}

var itemClasses = make(map[string][]*itemClassItem)
func tempItemClass(params []string) string {
	if len(params) < 4 {
		panic("Wrong number of params to ITEM_CLASS.")
	}
	
	rtn := new(itemClassItem)
	rtn.Type = params[0]
	rtn.Name = params[1]
	if !itemTypes[rtn.Type] {
		panic("Invalid item type: " + rtn.Type)
	}
	rtn.Rarity = params[2]
	if itemRarities[rtn.Rarity] == 0 {
		panic("Invalid item rarity: " + rtn.Rarity)
	}
	class := params[3]
	rtn.Tags = make([]string, 0, 5)
	
	for i := range params[4:] {
		rtn.Tags = append(rtn.Tags, params[i+4])
	}
	
	if _, ok := itemClasses[class]; !ok {
		itemClasses[class] = make([]*itemClassItem, 0, 10)
	}
	itemClasses[class] = append(itemClasses[class], rtn)
	
	return ""
}

func tempUsesItemClasses(params []string) string {
	if len(params) < 1 {
		panic("Wrong number of params to #USES_ITEM_CLASSES.")
	}
	
	classes := make([]string, 0, 5)
	tags := make([]string, 0, 5)
	for _, val := range params {
		if strings.HasPrefix(val, "#") {
			tags = append(tags, val)
			continue
		}
		classes = append(classes, val)
	}
	
	permittedItems := make(map[string]*itemClassItem)
	nameList := make([]string, 0, 20)
	for _, class := range classes {
ItemLoop: for _, item := range itemClasses[class] {
			// if the all of the item tags are present
			for _, itemtag := range item.Tags {
				foundTag := false
				for _, tag := range tags {
					if itemtag == tag {
						foundTag = true
						break
					}
				}
				if !foundTag {
					continue ItemLoop
				}
			}
			
			// add item to list
			if _, ok := permittedItems[item.Name]; ok {
				
				rtn := new(itemClassItem)
				rtn.Name = item.Name
				rtn.Type = item.Type
				
				if itemRarities[permittedItems[item.Name].Rarity] >= itemRarities[item.Rarity] {
					rtn.Rarity = permittedItems[item.Name].Rarity
				} else {
					rtn.Rarity = item.Rarity
				}
				rtn.Tags = item.Tags
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
		out += "\n\t[" + item.Type + ":" + item.Name + ":" + item.Rarity + "]"
	}
	return strings.TrimSpace(out)
}

var reactionData = make(map[string][]string)
var buildingData = make(map[string][]string)
func tempTechClass(params []string) string {
	if len(params) < 3 {
		panic("Wrong number of params to TECH_CLASS.")
	}
	
	if params[0] == "REACTION" {
		for i := range params[2:] {
			if _, ok := reactionData[params[i+2]]; !ok {
				reactionData[params[i+2]] = make([]string, 0, 5)
			}
			reactionData[params[i+2]] = append(reactionData[params[i+2]], params[1])
		}
		return ""
	} else if params[0] == "BUILDING" {
		for i := range params[2:] {
			if _, ok := buildingData[params[i+2]]; !ok {
				buildingData[params[i+2]] = make([]string, 0, 5)
			}
			buildingData[params[i+2]] = append(buildingData[params[i+2]], params[1])
		}
		return ""
	}
	panic("Invalid Type:" + params[0])
}

func tempUsesTechClasses(params []string) string {
	if len(params) == 0 {
		panic("Wrong number of params to #USES_TECH_CLASSES.")
	}
	
	// An interesting (ab)use of maps...
	permittedReactions := make(map[string]bool, 20)
	reactionNames := make([]string, 0, 20)
	permittedBuildings := make(map[string]bool, 20)
	buildingNames := make([]string, 0, 20)
	for _, tag := range params {
		if _, ok := buildingData[tag]; ok {
			for _, name := range buildingData[tag] {
				if !permittedReactions[name] {
					permittedReactions[name] = true
					reactionNames = append(reactionNames, name)
				}
			}
		}
		if _, ok := reactionData[tag]; ok {
			for _, name := range reactionData[tag] {
				if !permittedBuildings[name] {
					permittedBuildings[name] = true
					buildingNames = append(buildingNames, name)
				}
			}
		}
	}
	
	sort.Strings(reactionNames)
	sort.Strings(buildingNames)
	out := ""
	for _, name := range reactionNames {
		out += "\n\t[PERMITTED_BUILDING:" + name + "]"
	}
	for _, name := range buildingNames {
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

var sharedInorganicData = make(map[string]bool)
// This does not have all the error checking that Blast has, watch out!
func tempSharedInorganic(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to SHARED_INORGANIC.")
	}
	
	if sharedInorganicData[params[0]] {
		StageParse(params[1])
		return ""
	}
	
	sharedInorganicData[params[0]] = true
	rtn := "{#_REGISTERED_ORES;" + params[0] + "}\n"
	rtn += "{#_REGISTERED_REACTION_CLASSES;" + params[0] + "}\n"
	rtn += "{#_REGISTERED_REACTION_PRODUCTS;" + params[0] + "}\n"
	return StageParse(strings.TrimSpace(params[1])) + rtn
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
	mat := params[2]
	
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

var sharedItemData = make(map[string]bool)
// This does not have all the error checking that Blast has, watch out!
func tempSharedItem(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to SHARED_ITEM.")
	}
	
	if sharedItemData[params[0]] {
		StageParse(params[1])
		return ""
	}
	
	sharedItemData[params[0]] = true
	return StageParse(params[1])
}
