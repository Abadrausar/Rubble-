
# Fort time is value/72 and a growdur is value/7200, I think.
var rubble:timemap = <map
	"SECOND"=1
	"SECONDS"=1
	"MINUTE"=60
	"MINUTES"=60
	"HOUR"=3600
	"HOURS"=3600
	"DAY"=86400
	"DAYS"=86400
	"WEEK"=604800
	"WEEKS"=604800
	"MONTH"=2419200
	"MONTHS"=2419200
	"SEASON"=7257600
	"SEASONS"=7257600
	"YEAR"=29030400
	"YEARS"=29030400
>

(rubble:template "@ADV_TIME" block count unit {
	(if (exists [rubble:timemap] [unit]){
	}{
		(rubble:abort (str:add "Error: Attempt to use invalid unit: " [unit] " In #ADV_TIME."))
	})
	
	var out = (float:mul [rubble:timemap [unit]] [count])
	(if (int:lt [out] 1) {
		[out = 1]
	})
	(ret [out])
})

(rubble:template "@FORT_TIME" block count unit {
	(if (exists [rubble:timemap] [unit]){
	}{
		(rubble:abort (str:add "Error: Attempt to use invalid unit: " [unit] " In #FORT_TIME."))
	})
	
	# This should work...
	var out = (float:mul [rubble:timemap [unit]] [count])
	[out = (convert:int (float:div [out] 72))]
	(if (int:lt [out] 1) {
		[out = 1]
	})
	(ret [out])
})

(rubble:template "@GROWDUR" block count unit {
	(if (exists [rubble:timemap] [unit]){
	}{
		(rubble:abort (str:add "Error: Attempt to use invalid unit: " [unit] " In #GROWDUR."))
	})
	
	# This should work...
	var out = (float:mul [rubble:timemap [unit]] [count])
	[out = (convert:int (float:div [out] 7200))]
	(if (int:lt [out] 1) {
		[out = 1]
	})
	(ret (str:add "[GROWDUR:" [out] "]"))
})
