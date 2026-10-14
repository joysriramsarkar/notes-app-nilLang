# qemu/qemu-run.ps1
# নীলাং নোট অ্যাপ — QEMU এমুলেটর লঞ্চার (Windows PowerShell)

param (
    [string]$Mode = "gui",      # 'gui' (গ্রাফিক্যাল উইন্ডো) অথবা 'serial' (টার্মিনাল মোড)
    [int]$Memory = 512,         # মেমোরি সাইজ (MB)
    [string]$QemuPath = "C:\Program Files\qemu\qemu-system-x86_64.exe"
)

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "  🖥️  NILANG NOTES APP — QEMU EMULATOR LAUNCHER" -ForegroundColor Yellow
Write-Host "  অনুরণ ওএস / ভার্চুয়াল x86_64 প্ল্যাটফর্ম এমুলেশন" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan
Write-Host ""

# ১. QEMU পাথ যাচাই
if (-not (Test-Path $QemuPath)) {
    $found = Get-Command "qemu-system-x86_64.exe" -ErrorAction SilentlyContinue
    if ($found) {
        $QemuPath = $found.Source
    } else {
        Write-Host "❌ ত্রুটি: QEMU ইনস্টলেশন পাওয়া যায়নি।" -ForegroundColor Red
        Write-Host "দয়া করে নিশ্চিত করুন QEMU ইনস্টল করা আছে ($QemuPath)" -ForegroundColor Yellow
        exit 1
    }
}

Write-Host "✔ QEMU পাওয়া গেছে: $QemuPath" -ForegroundColor Green

# ২. বুট ইমেজ যাচাই ও রি-বিল্ড
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BootImg = Join-Path $ScriptDir "notes-os.img"

if (-not (Test-Path $BootImg)) {
    Write-Host "⚙️ বুট ইমেজ তৈরি করা হচ্ছে..." -ForegroundColor Yellow
    Push-Location $ScriptDir
    go run make_boot_image.go
    Pop-Location
}

if (-not (Test-Path $BootImg)) {
    Write-Host "❌ বুট ইমেজ তৈরি ব্যর্থ হয়েছে।" -ForegroundColor Red
    exit 1
}

Write-Host "✔ বুট ইমেজ প্রস্তুত: $BootImg" -ForegroundColor Green
Write-Host "🚀 QEMU ভার্চুয়াল মেশিন চালু হচ্ছে (RAM: ${Memory}MB, Mode: $Mode)..." -ForegroundColor Cyan
Write-Host "💡 টিপস: QEMU গ্রাফিক্যাল উইন্ডো থেকে মাউস আনলক করতে Ctrl+Alt চাপুন।" -ForegroundColor DarkYellow
Write-Host ""

if ($Mode -eq "serial") {
    & $QemuPath -drive "format=raw,file=$BootImg" -m "${Memory}M" -nographic -serial mon:stdio -boot a
} else {
    & $QemuPath -name "Nilang Notes - Onuron Virtual OS" -drive "format=raw,file=$BootImg" -m "${Memory}M" -vga std -boot a
}
