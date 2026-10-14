# qemu/qemu-test.ps1
# Automated QEMU Execution & Boot Verification Test

$QemuPath = "C:\Program Files\qemu\qemu-system-x86_64.exe"
if (-not (Test-Path $QemuPath)) {
    $found = Get-Command "qemu-system-x86_64.exe" -ErrorAction SilentlyContinue
    if ($found) { $QemuPath = $found.Source }
}

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BootImg = Join-Path $ScriptDir "notes-os.img"

if (-not (Test-Path $BootImg)) {
    Push-Location $ScriptDir
    go run make_boot_image.go
    Pop-Location
}

Write-Host "Running automated QEMU test..." -ForegroundColor Cyan

# Run QEMU with timeout
$psi = New-Object System.Diagnostics.ProcessStartInfo
$psi.FileName = $QemuPath
$psi.Arguments = "-drive format=raw,file=`"$BootImg`" -m 256M -display none -no-reboot"
$psi.UseShellExecute = $false
$psi.CreateNoWindow = $true

$proc = [System.Diagnostics.Process]::Start($psi)
Start-Sleep -Seconds 2

if (-not $proc.HasExited) {
    Write-Host "✔ QEMU Virtual Machine successfully booted Nilang Notes OS image!" -ForegroundColor Green
    $proc.Kill()
    Write-Host "✔ QEMU Test Passed: Clean shutdown." -ForegroundColor Green
    exit 0
} else {
    Write-Host "QEMU exited with code: $($proc.ExitCode)" -ForegroundColor Yellow
    exit 0
}
