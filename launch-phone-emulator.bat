@echo off
title NilPhone 16 Pro Max - NilOS Studio Emulator
color 0b

echo ============================================================
echo   NILPHONE 16 PRO MAX - NILOS EMULATOR (NATIVE WINDOW)
echo   Initializing NilOS Runtime and Launching Dedicated Window...
echo ============================================================
echo.

cd /d "%~dp0"

rem 1. Check if backend server is already running on port 8085
netstat -ano | findstr :8085 >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [info] Starting NilOS Emulator backend server in background...
    if exist "nilos-emulator" (
        cd /d "%~dp0nilos-emulator"
        start /b go run . >nul 2>&1
        cd /d "%~dp0"
    ) else (
        start /b go run . >nul 2>&1
    )
    ping 127.0.0.1 -n 3 >nul
) else (
    echo [ok] NilOS server is already running on port 8085.
)

rem 2. Launch Dedicated Standalone Phone Window (Frameless App Window Mode)
echo [launch] Opening dedicated NilPhone Emulator window...
start "" msedge --app="http://127.0.0.1:8085/phone.html" --window-size=520,950 --window-position=550,40 --no-first-run --no-default-browser-check || start "" chrome --app="http://127.0.0.1:8085/phone.html" --window-size=520,950 --window-position=550,40

echo.
echo ============================================================
echo   [SUCCESS] NilPhone Emulator is active on your screen!
echo ============================================================
ping 127.0.0.1 -n 2 >nul
exit
