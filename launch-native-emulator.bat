@echo off
title NilPhone 16 Pro Max — 100%% Pure Native Go Emulator
color 0a

echo ============================================================
echo   📱 NILPHONE 16 PRO MAX — 100%% PURE NATIVE GO EMULATOR
echo   Architecture: Native Windows Win32 GDI (Zero HTML/Browser)
echo   Launching Standalone Native Desktop Window...
echo ============================================================
echo.

cd /d "%~dp0native-phone-emulator"

if not exist "nilphone.exe" (
    echo [build] Compiling native Go binary...
    go build -o nilphone.exe .
)

start "" "%~dp0native-phone-emulator\nilphone.exe"
echo [OK] Pure Native Go Emulator is running!
