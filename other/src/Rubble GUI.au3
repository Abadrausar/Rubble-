#AutoIt3Wrapper_Icon=..\..\Rubble.ico
#AutoIt3Wrapper_UseX64=n
#AutoIt3Wrapper_Res_Description=Rubble GUI Launcher
#AutoIt3Wrapper_Res_Fileversion=2.0
#AutoIt3Wrapper_Res_LegalCopyright=Copyright 2013-2014 Milo Christiansen
#AutoIt3Wrapper_Res_Language=1033
#AutoIt3Wrapper_Res_Comment=Rubble Version (at last GUI build): pre4

#include <GUIConstantsEx.au3>
#include <TreeViewConstants.au3>
#include <WindowsConstants.au3>

#include <File.au3>
#include <Array.au3>

Opt("TrayMenuMode",1)
Opt("TrayIconHide",1)

Dim $AddonList[1]
$AddonList[0] = 0
Dim $AddonItemList[1]
$AddonItemList[0] = 0

Dim $RegenAddonItemList[1]
$RegenAddonItemList[0] = 0

Dim $Minimal = IniRead(@ScriptDir & "\gui.ini", "general", "minimal", "false")
If $Minimal == "true" Then
	$Minimal = True
Else
	$Minimal = False
EndIf

; Create Main Window
Global $MainWindow = GUICreate("Rubble GUI v2", 430, 430, -1, -1, BitOR($GUI_SS_DEFAULT_GUI,$WS_MAXIMIZEBOX,$WS_SIZEBOX,$WS_THICKFRAME,$WS_TABSTOP))

Global $TabMain = GUICtrlCreateTab(5, 5, 420, 420)
GUICtrlSetResizing(-1, $GUI_DOCKBORDERS)

; Addons Tab
Global $TabAddons = GUICtrlCreateTabItem("Addons")
GUICtrlSetState(-1,$GUI_SHOW)

Global $Addons = GUICtrlCreateTreeView(14, 35, 400, 327, BitOR($GUI_SS_DEFAULT_TREEVIEW,$TVS_CHECKBOXES,$WS_HSCROLL,$WS_VSCROLL,$WS_BORDER))
GUICtrlSetResizing(-1, $GUI_DOCKBORDERS)

Global $RunRubble = GUICtrlCreateButton("Generate Raws With Selected Addons", 14, 375, 400, 41)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKBOTTOM+$GUI_DOCKHEIGHT)

; Prep Tab
If Not $Minimal Then
	Global $TabPrep = GUICtrlCreateTabItem("Prep")

	Global $PrepRegions = GUICtrlCreateList("", 14, 35, 400, 327)
	GUICtrlSetResizing(-1, $GUI_DOCKBORDERS)

	Global $RunPrep = GUICtrlCreateButton("Run Prep for Selected Region", 14, 375, 400, 41)
	GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKBOTTOM+$GUI_DOCKHEIGHT)

	; Regen Tab
	Global $TabRegen = GUICtrlCreateTabItem("Regen")

	GUICtrlCreateLabel("Regen Region", 14, 35, 400, 17)
	GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKTOP+$GUI_DOCKHEIGHT)
	Global $RegenRegions = GUICtrlCreateList("", 14, 50, 400, 97)
	GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKTOP+$GUI_DOCKHEIGHT)

	GUICtrlCreateLabel("Regen Addons", 14, 155, 400, 17)
	GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKTOP)
	Global $RegenAddons = GUICtrlCreateTreeView(14, 170, 400, 200, BitOR($GUI_SS_DEFAULT_TREEVIEW,$TVS_CHECKBOXES,$WS_HSCROLL,$WS_VSCROLL,$WS_BORDER))
	GUICtrlSetResizing(-1, $GUI_DOCKBORDERS)

	Global $RunRegen = GUICtrlCreateButton("Regen Selected Region With Selected Addons", 14, 375, 400, 41)
	GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKBOTTOM+$GUI_DOCKHEIGHT)
EndIf

; Misc. Tab
Global $TabOther = GUICtrlCreateTabItem("Other")

GUICtrlCreateLabel("Extra Command Line Options:", 14, 35, 400, 17)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKTOP+$GUI_DOCKHEIGHT)

Global $OtherOptions = GUICtrlCreateInput("", 14, 55, 400, 21)
GUICtrlSetResizing(-1, $GUI_DOCKLEFT+$GUI_DOCKRIGHT+$GUI_DOCKTOP+$GUI_DOCKHEIGHT)

GUICtrlCreateTabItem("")

; Load settings
Dim $RubbleBin = IniRead(@ScriptDir & "\gui.ini", "general", "rubble_bin", "rubble.exe")
Dim $RubbleWd = IniRead(@ScriptDir & "\gui.ini", "general", "rubble_cwd", ".")

Dim $DFDir = ".."
$DFDir = IniRead($RubbleWd & "/rubble.ini", "rubble", "dfdir", $DFDir)
StringReplace($DFDir, "rubble:", "./")

Dim $OutDir = "df:raw"
$OutDir = IniRead($RubbleWd & "/rubble.ini", "rubble", "outputdir", $OutDir)
StringReplace($OutDir, "rubble:", "./")
StringReplace($OutDir, "df:", $DFDir & "/")

Dim $AddonDir = ".\addons"
$AddonDir = IniRead($RubbleWd & "/rubble.ini", "rubble", "addonsdir", $AddonDir)
StringReplace($AddonDir, "rubble:", "./")
StringReplace($AddonDir, "df:", $DFDir & "/")
StringReplace($AddonDir, "out:", $OutDir & "/")

Dim $AddonListOverride = IniRead(@ScriptDir & "\gui.ini", "general", "addonlist", "")

; Install Mode
If $CmdLineRaw <> "" Then
	Run(@ComSpec & ' /k ""./' & $RubbleBin & '" -install=' & $CmdLineRaw & '"', $RubbleWd)
	Exit
EndIf

; Handle writing settings
Func OnExit()
	If FileExists($AddonListOverride) Then
		For $X = 1 To $AddonItemList[0]
			If BitAnd(GUICtrlRead($AddonItemList[$X]),$GUI_CHECKED) Then
				IniWrite($AddonListOverride, "addons", $AddonList[$X], "true")
			Else
				IniWrite($AddonListOverride, "addons", $AddonList[$X], "false")
			EndIf
		Next
	EndIf

	IniWrite(@ScriptDir & "\gui.ini", "general", "otheroptions", GUICtrlRead($OtherOptions))
	IniWrite(@ScriptDir & "\gui.ini", "general", "rubble_bin", $RubbleBin)
	IniWrite(@ScriptDir & "\gui.ini", "general", "rubble_cwd", $RubbleWd)
	IniWrite(@ScriptDir & "\gui.ini", "general", "addonsdir", $AddonDir)
	IniWrite(@ScriptDir & "\gui.ini", "general", "dfdir", $DFDir)
	IniWrite(@ScriptDir & "\gui.ini", "general", "addonlist", $AddonListOverride)
EndFunc
OnAutoItExitRegister("OnExit")

; Update Addonlist
If $AddonListOverride == "" Then
	RunWait($RubbleWd & "/" & $RubbleBin & ' -addonlist', $RubbleWd , @SW_HIDE)

	If Not FileExists($AddonDir & "\addonlist.ini") Then
		MsgBox(16, "Loading Error", "Rubble failed to generate/update the addonlist.ini, please fix the gui.ini file and try again.")
		Exit 1
	EndIf
Else
	If Not FileExists($AddonListOverride) Then
		MsgBox(16, "Loading Error", "Rubble GUI failed to find the addon list override file, please fix the gui.ini file and try again.")
		Exit 1
	EndIf
EndIf

; Load Addonlist
If Not FileExists($AddonListOverride) Then
	Dim $section = IniReadSection($AddonDir & "\addonlist.ini", "addons")
Else
	Dim $section = IniReadSection($AddonListOverride, "addons")
EndIf
If Not  @error Then
	For $X = 1 To $section[0][0]
		_ArrayAdd($AddonList, $section[$X][0])
		_ArrayAdd($AddonItemList, GUICtrlCreateTreeViewItem($section[$X][0], $Addons))
		$AddonList[0] += 1
		$AddonItemList[0] +=1

		If $section[$X][1] == "true" Then
			GUICtrlSetState($AddonItemList[$AddonItemList[0]], $GUI_CHECKED)
		EndIf
	Next
Else
	MsgBox(16, "Loading Error", "Rubble GUI failed to read the addon list (this is probably bad).")
EndIf

; Load other command line options
GUICtrlSetData($OtherOptions, IniRead(@ScriptDir & "\gui.ini", "general", "otheroptions", ""))

If Not $Minimal Then
	Dim $Regions = _FileListToArray($DFDir & "/data/save", "*", 2)
	If @error Then
		MsgBox(16, "Error", "Rubble GUI failed to find your saves directory, please fix the gui.ini file and try again.")
		Exit 1
	EndIf

	; Populate region lists
	GUICtrlSetData($PrepRegions, "raw")

	For $X = 1 To $Regions[0]
		; "current" is not a region
		If $Regions[$X] == "current" Then
			ContinueLoop
		EndIf

		GUICtrlSetData($PrepRegions, $Regions[$X])
		GUICtrlSetData($RegenRegions, $Regions[$X])
	Next

	GUICtrlSetData($PrepRegions, "raw")
EndIf

GUISetState(@SW_SHOW)

While 1
	$msg = GUIGetMsg()

	If $Minimal Then
		Switch $msg
			Case $GUI_EVENT_CLOSE
				Exit

			Case $RunRubble
				Local $addons = ""
				For $X = 1 To $AddonItemList[0]
					If BitAnd(GUICtrlRead($AddonItemList[$X]),$GUI_CHECKED) Then
						$addons &= ";" & $AddonList[$X]
					EndIf
				Next
				$addons = StringTrimLeft($addons,1)
				Run(@ComSpec & ' /k ""./' & $RubbleBin & '" -addons="' & $addons & '" ' & GUICtrlRead($OtherOptions) & '"', $RubbleWd)
				Exit
		EndSwitch
		ContinueLoop
	EndIf

	Switch $msg
		Case $GUI_EVENT_CLOSE
			Exit

	Case $RunRubble
			Local $addons = ""
			For $X = 1 To $AddonItemList[0]
				If BitAnd(GUICtrlRead($AddonItemList[$X]),$GUI_CHECKED) Then
					$addons &= ";" & $AddonList[$X]
				EndIf
			Next
			$addons = StringTrimLeft($addons,1)
			Run(@ComSpec & ' /k ""./' & $RubbleBin & '" -addons="' & $addons & '" ' & GUICtrlRead($OtherOptions) & '"', $RubbleWd)
			Exit

	Case $RunPrep
			Run(@ComSpec & ' /k ""./' & $RubbleBin & '" -prep="' & GUICtrlRead($PrepRegions) & '"', $RubbleWd)

	Case $RegenRegions
			If FileExists($DFDir & '/data/save/' & GUICtrlRead($RegenRegions) & "/raw/addonlist.ini") Then

				Dim $section = IniReadSection($DFDir & '/data/save/' & GUICtrlRead($RegenRegions) & "/raw/addonlist.ini", "addons")
				If Not  @error Then
					For $X = 1 To $AddonList[0]
						_ArrayAdd($RegenAddonItemList, GUICtrlCreateTreeViewItem($AddonList[$X], $RegenAddons))
					Next
					$RegenAddonItemList[0] = $AddonList[0]

					For $X = 1 To $section[0][0]
						If $section[$X][1] == "true" Then
							; Ouch...
							Local $found = False
							For $Y = 1 To $AddonList[0]
								If $AddonList[$Y] == $section[$X][0] Then
									GUICtrlSetState($RegenAddonItemList[$Y], $GUI_CHECKED)
									$found = True
									ExitLoop
								EndIf
							Next
							If Not $found Then
								MsgBox(16, "Warning!", "You do not have all the addons needed to regenerate this world! Depending on the nature of the missing addons this may cause the world to be corrupted! It is HIGHLY recommended that you make a backup before proceeding." & @CRLF & "Missing Addon: " & $section[$X][0])
							EndIf
						EndIf
					Next
				Else
					MsgBox(16, "Loading Error", "Rubble GUI failed to read the addon list (this is probably bad).")
				EndIf
			EndIf

	Case $RunRegen
		If GUICtrlRead($RegenRegions) <> "" Then
			Local $addons = ""
			For $X = 1 To $RegenAddonItemList[0]
				If BitAnd(GUICtrlRead($RegenAddonItemList[$X]),$GUI_CHECKED) Then
					$addons &= ";" & $AddonList[$X]
				EndIf
			Next
			$addons = StringTrimLeft($addons,1)

			local $rtn = MsgBox(52, "Warning!", "Regenerating a world only works correctly in some cases, it is recommended that this be used ONLY for changing tilesets and the like! As world corruption is a possibility, it is HIGHLY recommended that you make a backup before proceeding." & @CRLF & "Do you wish to continue?")
			If $rtn <> 7 Then
				Local $save = 'df:data/save/' & GUICtrlRead($RegenRegions)
				Run(@ComSpec & ' /k ""./' & $RubbleBin & '" -zapaddons -zapconfig -addons="' & $addons & '" -config="' & $save & '/raw/genconfig.ini" -outputdir="' & $save & '/raw"', $RubbleWd)
			EndIf
		EndIf
	EndSwitch
WEnd
