
@echo off
setlocal

rmdir /s /q addons
rmdir /s /q other
del /q rubble.exe

:: DO NOT DELETE THE GUI BINARY! It cannot be rebuilt from the command line!

SET "GO_SRC=C:/Docs/Projects/Go/src"
SET "RBL_STDLIB=%GO_SRC%/rubble/std_lib"

call gopath.bat -quiet

:: ==============================================================
echo Fetching Rubble Documentation...

xcopy "%RBL_STDLIB%/docs" other /e /i /q

:: ==============================================================
echo Building Rubble...

call gobuild rubble windows 386

mkdir "other/linux_386"
cd "other/linux_386"
call gobuild rubble linux 386

mkdir "../darwin_386"
cd "../darwin_386"
call gobuild rubble darwin 386

cd "../.."

:: ==============================================================
echo Fetching Source Code...
cd other

xcopy "%RBL_STDLIB%/src" src /e /i /q
cd src

xcopy "%GO_SRC%/rubble" rubble /e /i /q
cd rubble
rmdir /s /q std_lib
cd ..

mkdir dctech
cd dctech

xcopy "%GO_SRC%/dctech/rex" rex /e /i /q

xcopy "%GO_SRC%/dctech/dfrex/dfraw" "dfrex/dfraw" /e /i /q

xcopy "%GO_SRC%/dctech/axis" axis /e /i /q

xcopy "%GO_SRC%/dctech/iconsole" iconsole /e /i /q

cd "../../.."

:: ==============================================================
echo Generating Rex Documentation...
cd other

xcopy "%RBL_STDLIB%/rex_docs" "Rex Docs" /i /q
xcopy "%GO_SRC%/dctech/rex/docs" "Rex Docs" /i /q

cd "Rex Docs"

godoc dctech/rex/commands/base Command_+ >> base.txt

godoc dctech/rex/commands/boolean Command_+ >> boolean.txt

godoc dctech/rex/commands/console Command_+ >> console.txt

godoc dctech/rex/commands/convert Command_+ >> convert.txt

godoc dctech/rex/commands/debug Command_+ >> debug.txt

godoc dctech/rex/commands/expr Command_+ >> expr.txt

godoc dctech/rex/commands/float Command_+ >> float.txt

godoc dctech/rex/commands/integer Command_+ >> integer.txt

godoc dctech/rex/commands/regex Command_+ >> regex.txt

godoc dctech/rex/commands/sort Command_+ >> sort.txt

godoc dctech/rex/commands/str Command_+ >> str.txt

godoc dctech/rex/genii "Command_+" >> genii.txt

godoc dctech/axis/axisrex "Command_+" >> axis.txt

godoc dctech/dfrex/dfraw "Command_+" >> df_raw.txt

godoc rubble/guts "Command_+" >> custom.txt

cd "../.."

:: ==============================================================
echo Fetching Rubble Standard Addons...

xcopy "%RBL_STDLIB%/addons" addons /e /i /q

pause
