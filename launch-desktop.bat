@echo off
title Nilang Notes - Desktop Launcher
color 0b
echo ============================================================
echo   NILANG NOTES APP - DESKTOP / WEB RUNNER
echo   Starting server at http://localhost:8085 ...
echo ============================================================
echo.

cd /d "%~dp0"
start http://localhost:8085
go run .
pause
