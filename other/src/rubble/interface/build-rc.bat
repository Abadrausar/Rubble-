
@call mingwpath
@echo Building res.rc (interfaces)...
@cd cli
@windres -o res_windows.syso ./res.rc

@cd ../gui
@windres -o res_windows.syso ./res.rc

@cd ../web
@windres -o res_windows.syso ./res.rc
