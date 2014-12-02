
@echo off

rmdir /S /Q rubble
xcopy "C:/Docs/Projects/Go/src/rubble" rubble /e /i

cd dctech

rmdir /S /Q raptor
xcopy "C:/Docs/Projects/Go/src/dctech/raptor" raptor /e /i

rmdir /S /Q ini
xcopy "C:/Docs/Projects/Go/src/dctech/ini" ini /e /i

rmdir /S /Q dcfs
xcopy "C:/Docs/Projects/Go/src/dctech/dcfs" dcfs /e /i

rmdir /S /Q iconsole
xcopy "C:/Docs/Projects/Go/src/dctech/iconsole" iconsole /e /i

cd ..
