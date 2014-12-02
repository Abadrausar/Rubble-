package main

import "fmt"
import "strings"

func SetupBuiltins() {
	// This should be all the "fixed global templates" from Blast except !GLOBAL_TEMPLATE which is not applicable
	NewNativeTemplate("!TEMPLATE", tempTemplate)
	NewNativeTemplate("!T", tempTemplate)
	NewNativeTemplate("ONCE", tempOnce)
	NewNativeTemplate("STATIC", tempStatic)
	NewNativeTemplate("S", tempStatic)
	NewNativeTemplate("LOG", tempLog)
	NewNativeTemplate("L", tempLog)
	NewNativeTemplate("COMMENT", tempComment)
	NewNativeTemplate("C", tempComment)
	NewNativeTemplate("VOID", tempVoid)
	NewNativeTemplate("V", tempVoid)
	NewNativeTemplate("GETVAR", tempGetVar)
	NewNativeTemplate("SETVAR", tempSetVar)
	NewNativeTemplate("IFVAR", tempIfVar)
	
	// Now for other stuff...
	NewNativeTemplate("#ADVENTURE_TIER", tempAdventureTier)
	NewNativeTemplate("TECH_CLASS", tempTechClass)
	NewNativeTemplate("#USES_TECH_CLASSES", tempUsesTechClasses)
}

//=============================================
// Fixed global templates

func tempTemplate(params []string) string {
	if len(params) < 2 {
		panic("Wrong number of params to !TEMPLATE.")
	}
	
	name := params[0]
	paramNames := params[1:len(params)-1]
	text := strings.TrimSpace(params[len(params)-1])
	
	NewUserTemplate(name, text, paramNames)
	
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

func tempLog(params []string) string {
	for _, val := range params {
		fmt.Println(val)
	}
	return ""
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

var varData = make(map[string]string)
func tempGetVar(params []string) string {
	if len(params) != 1 {
		panic("Wrong number of params to GETVAR.")
	}
	
	return varData[params[0]]
}

func tempSetVar(params []string) string {
	if len(params) != 2 {
		panic("Wrong number of params to SETVAR.")
	}
	
	varData[params[0]] = params[1]
	return ""
}

func tempIfVar(params []string) string {
	if len(params) != 4 {
		panic("Wrong number of params to IFVAR.")
	}
	
	if varData[params[0]] == params[1] {
		return StageParse(strings.TrimSpace(params[2]))
	}
	return StageParse(strings.TrimSpace(params[3]))
}

//=============================================
// Other templates

var advTier = 0
func tempAdventureTier(params []string) string {
	rtn := fmt.Sprint("[ADVENTURE_TIER:", advTier, "]")
	advTier++
	return rtn
}

var reactionData = make(map[string][]string)
var buildingData = make(map[string][]string)
func tempTechClass(params []string) string {
	if len(params) != 3 {
		panic("Wrong number of params to TECH_CLASS.")
	}
	
	if params[0] == "REACTION" {
		if _, ok := reactionData[params[2]]; !ok {
			reactionData[params[2]] = make([]string, 0, 5)
		}
		reactionData[params[2]] = append(reactionData[params[2]], params[1])
		return ""
	} else if params[0] == "BUILDING" {
		if _, ok := buildingData[params[2]]; !ok {
			buildingData[params[2]] = make([]string, 0, 5)
		}
		buildingData[params[2]] = append(buildingData[params[2]], params[1])
		return ""
	}
	panic("Invalid Type:" + params[0])
}

func tempUsesTechClasses(params []string) string {
	if len(params) == 0 {
		panic("Wrong number of params to #USES_TECH_CLASSES.")
	}
	
	out := ""
	for _, tag := range params {
		if _, ok := buildingData[tag]; ok {
			for _, name := range buildingData[tag] {
				out += "\t[PERMITTED_BUILDING:" + name + "]\n"
			}
		}
		if _, ok := reactionData[tag]; ok {
			for _, name := range reactionData[tag] {
				out += "\t[PERMITTED_REACTION:" + name + "]\n"
			}
		}
	}
	return strings.TrimSpace(out) + "\n\t"
}
