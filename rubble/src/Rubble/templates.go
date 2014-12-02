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

	NewNativeTemplate("#ADV_TIME", tempAdvTime)
	NewNativeTemplate("#FORT_TIME", tempFortTime)
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
