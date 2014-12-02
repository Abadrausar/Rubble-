
@echo off

rmdir /S /Q rubble
xcopy "C:/Docs/Projects/Go/src/rubble" rubble /e /i

cd dctech

rmdir /S /Q nca7
xcopy "C:/Docs/Projects/Go/src/dctech/nca7" nca7 /e /i

rmdir /S /Q ini
xcopy "C:/Docs/Projects/Go/src/dctech/ini" ini /e /i

cd ..
