// Package main implements the standalone NilOS Studio Device Emulator launcher.
// Usage: go run ./emulator
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	port := 8085
	targetURL := fmt.Sprintf("http://localhost:%d/device.html", port)

	// Determine root directory where public/ is located
	rootDir := "."
	if _, err := os.Stat("public"); os.IsNotExist(err) {
		// If running from ./emulator subfolder
		rootDir = ".."
	}

	fmt.Println("============================================================")
	fmt.Println("  📱 নীল ওএস (NilOS) ভার্চুয়াল ডিভাইস ও অ্যাপ্লিকেশন এমুলেটর")
	fmt.Println("  ⚡ NilOS 2.0 Studio Emulator — Onuron Virtual Environment")
	fmt.Println("============================================================")
	fmt.Printf("  🌐 লোকাল পোর্ট: %d\n", port)
	fmt.Printf("  📱 এমুলেটর ইউআরএল: %s\n", targetURL)
	fmt.Println("  🛠️  ডিভাইস প্রোফাইল: Phone, Tablet, POS, Desktop")
	fmt.Println("============================================================")
	fmt.Println("  [Ctrl+C চাপলে সার্ভার বন্ধ হবে]")
	fmt.Println("")

	// Launch Standalone Browser Window
	go launchAppWindow(targetURL)

	// Delegate to root HTTP server or static server
	publicDir := filepath.Join(rootDir, "public")
	fs := http.FileServer(http.Dir(publicDir))
	http.Handle("/", fs)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Printf("সার্ভার ইতোমধ্যে ব্যাকগ্রাউন্ডে চলছে অথবা পোর্ট দখলকৃত: %v", err)
	}
}

func launchAppWindow(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Try MS Edge app window mode, then Chrome, then default browser
		cmd = exec.Command("cmd", "/c", "start", "msedge", fmt.Sprintf("--app=%s", url), "--window-size=1200,880")
		if err := cmd.Start(); err != nil {
			exec.Command("cmd", "/c", "start", "chrome", fmt.Sprintf("--app=%s", url), "--window-size=1200,880").Start()
		}
	case "darwin":
		cmd = exec.Command("open", url)
		cmd.Start()
	default:
		cmd = exec.Command("xdg-open", url)
		cmd.Start()
	}
}
