package main

//import "fmt"
import "strings"
import "strconv"
import "dctech/nca6"

func SetupBuiltins() {
	NewNativeTemplate("!TEMPLATE", tempTemplate)
	NewNativeTemplate("!SCRIPT_TEMPLATE", tempScriptTemplate)
	
	NewNativeTemplate("!SCRIPT", tempScript)
	NewNativeTemplate("SCRIPT", tempScript)
	NewNativeTemplate("#SCRIPT", tempScript)

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
	GlobalNCAState.Envs.Add(nca6.NewEnvironment())
	
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
