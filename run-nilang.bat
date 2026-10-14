@echo off
title Nilang Notes - TUI Terminal Mode
color 0a
echo ============================================================
echo   NILANG NOTES APP - NILANG INTERPRETER / TUI MODE
echo   Running: nil run src/main.nil ...
echo ============================================================
echo.

cd /d "%~dp0"
"..\nilLang\bin\nil.exe" run src/main.nil
pause
