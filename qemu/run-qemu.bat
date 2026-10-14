@echo off
title Nilang Notes - QEMU Emulator Launcher
color 0b
echo ============================================================
echo   NILANG NOTES APP - QEMU EMULATOR LAUNCHER (WINDOWS)
echo   Onuron OS / Virtual Platform Verification
echo ============================================================
echo.

set QEMU_BIN="C:\Program Files\qemu\qemu-system-x86_64.exe"
if not exist %QEMU_BIN% (
    echo [ERROR] QEMU was not found at %QEMU_BIN%
    echo Please ensure QEMU is installed.
    pause
    exit /b 1
)

if not exist "%~dp0notes-os.img" (
    echo [INFO] Generating QEMU bootable disk image...
    cd /d "%~dp0"
    go run make_boot_image.go
)

echo [OK] Starting QEMU Virtual Machine...
echo [INFO] Press Ctrl+Alt to toggle mouse focus in QEMU window.
echo.

%QEMU_BIN% -name "Nilang Notes - Onuron Virtual OS" -drive "format=raw,file=%~dp0notes-os.img" -m 512M -vga std -boot a

pause
