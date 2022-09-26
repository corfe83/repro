@echo off
if exist *.apk (
	del /Q *.apk
)
gomobile build -ldflags="-s -w" -target android -o repro.apk
zipalign 4 repro.apk repro_aligned.apk
call apksigner sign --ks key repro_aligned.apk
