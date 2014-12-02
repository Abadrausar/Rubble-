
@echo off

rmdir /S /Q rubble
xcopy "C:/Docs/Projects/Go/src/rubble" rubble /e /i

cd dctech

rmdir /S /Q nca6
xcopy "C:/Docs/Projects/Go/src/dctech/nca6" nca6 /e /i

rmdir /S /Q ini
xcopy "C:/Docs/Projects/Go/src/dctech/ini" ini /e /i

rmdir /S /Q ncalex
xcopy "C:/Docs/Projects/Go/src/dctech/ncalex" ncalex /e /i

cd ..
