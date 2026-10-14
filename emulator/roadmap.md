হ্যাঁ। রিপোজিটরিটা দেখে আমার মতে **আলাদা Python interpreter বানানো উচিত নয়**। আপনার আগের নকশায় Python-এ lexer/parser/evaluator নতুন করে লেখার প্রস্তাব ছিল পেস্ট করা টেক্সট, এবং সেই নকশাতেও মূল execution pipeline-কে Lexer → Parser → Evaluator হিসেবে ধরা হয়েছিল পেস্ট করা টেক্সট। কিন্তু এখনকার `nilLang`-এ এগুলো ইতিমধ্যে Go-তে আছে; আবার বানালে emulator আর আসল নীলাং runtime একসময় আলাদা আচরণ করতে শুরু করবে।

আমার প্রস্তাব:

```text
                    Nilang Emulator
                         │
        ┌────────────────┴────────────────┐
        │                                 │
 existing Nilang runtime             Device Shell
        │                                 │
 lexer/parser/typecheck                 Phone
 evaluator                              Tablet
        │                               Desktop
        ▼                                 │
 declarativeUIApp                         │
        │                                 │
   ┌────┴────┐                            │
   │         │                            │
 Render()  Dispatch() ◄── browser click ──┘
   │         │
   └────┬────┘
        ▼
   Existing Alap SSR
        ▼
      HTML
        ▼
   iframe device
```

অর্থাৎ নতুন emulator-এর কাজ হবে **runtime emulate করা নয়, runtime-কে ঘিরে development device simulate করা**।

## কী কী থাকবে

প্রথম সংস্করণে আমি এগুলো রাখছি:

- `nil emulate` কমান্ড
- Phone / Tablet / Desktop viewport
- স্থানীয় `127.0.0.1` server
- browser-এ emulator shell
- Alap declarative UI render
- Button → Nilang `on event(...)`
- event payload
- state update → rerender
- source edit → hot reload
- parse/type/evaluation error দেখানো
- Reload / Open App
- নতুন third-party Go dependency নেই
- Electron/Tauri/Chromium bundle নেই

এটা তাই খুব হালকা থাকবে।

### ব্যবহার হবে

```powershell
nil emulate
```

অথবা:

```powershell
nil emulate examples/ui-counter/src/main.nil
```

Phone:

```powershell
nil emulate examples/ui-counter/src/main.nil --device phone
```

Tablet:

```powershell
nil emulate examples/ui-counter/src/main.nil --device tablet
```

Desktop:

```powershell
nil emulate examples/ui-counter/src/main.nil --device desktop
```

নির্দিষ্ট port:

```powershell
nil emulate --port 8090
```

Browser নিজে খুলবে না:

```powershell
nil emulate --no-open
```

## সবচেয়ে গুরুত্বপূর্ণ বিষয়

আমি `cmd/nil/render.go`-তে ইতিমধ্যে থাকা `declarativeUIApp`, `loadDeclarativeUI()`, `Render()`, `Dispatch()` এবং `cmd/nil/web_cmd.go`-এর rendering path-ই ব্যবহার করেছি।

তাই নতুন emulator-এ এরকম দ্বিতীয় parser থাকবে না:

```text
.nil
 ↓
Python Lexer
 ↓
Python Parser
 ↓
Python VM
```

বরং:

```text
.nil
 ↓
Nilang-এর আসল Lexer/Parser
 ↓
Nilang-এর আসল Type Checker
 ↓
Nilang-এর আসল Evaluator
 ↓
Nilang-এর আসল Alap UI
 ↓
Emulator shell
```

এটাই দীর্ঘমেয়াদে সঠিক architecture।

## Emulator-এর মূল code

নতুন মূল ফাইল:

```text
cmd/nil/emulator.go
```

এতে মোটামুটি এই execution path:

```go
app, err := loadEmulatorUI(absEntry)

...

html, err := renderEmulatorApp(state.app)

...

if err := state.app.Dispatch(req.Event, req.Payload); err != nil {
    ...
}
```

আর browser-এর Button থেকে:

```javascript
const response = await fetch('/__nil/event', {
  method: 'POST',
  headers: {'Content-Type':'application/json'},
  body: JSON.stringify({
    event: eventName,
    payload
  })
});
```

তারপর আবার existing Nilang component-এর:

```text
on increment(payload) {
    ...
}
```

চলে যাবে।

## Hot reload

এখানে `fsnotify`-ও ব্যবহার করিনি। নতুন dependency না এনে 700ms পরপর `.nil` file-এর modification time দেখা হবে।

```text
source edit
    ↓
700ms poll
    ↓
file changed?
    ↓
loadDeclarativeUI()
    ↓
success → নতুন app
failure → পুরনো working app + error status
```

এটা deliberate decision। Emulator-এর development convenience-এর জন্য `fsnotify` টেনে dependency বাড়ানোর দরকার নেই।

## Device shell

Browser-এর ভেতর:

```text
┌──────────────────────────────────────────┐
│ Nilang Emulator       Phone ▼   Reload   │
├──────────────────────────────────────────┤
│                                          │
│          ┌──────────────────┐            │
│          │                  │            │
│          │   Nilang App     │            │
│          │                  │            │
│          │     Button       │            │
│          │                  │            │
│          │                  │            │
│          └──────────────────┘            │
│                                          │
└──────────────────────────────────────────┘
```

Phone/Tablet/Desktop শুধু CSS viewport simulation। Android kernel, CPU, sensor, GPU বা hardware emulation করছে না।

---

## এরপর কী করা উচিত

আমি emulator-কে ৩ স্তরে ভাগ করব।

### স্তর ১ — UI Emulator

এখনকার implementation:

```text
UI
State
Events
Rendering
Hot Reload
Viewport
```

এটাই আগে স্থির করা উচিত।

### স্তর ২ — Capability Emulator

তারপর:

```text
Filesystem
Network
Storage
Clock
Device Info
Permissions
Camera mock
Location mock
Notifications mock
```

যেমন:

```nilang
device.camera.capture()
```

emulator-এ real camera না খুলে:

```text
[ Mock Camera ]
```

দেবে।

একই API পরে Onuron OS-এ real implementation ব্যবহার করবে।

### স্তর ৩ — Onuron Runtime Emulator

সবশেষে:

```text
Nilang VM
   +
Alap Runtime
   +
Onuron System API
   +
Device Capability Layer
```

এটা হলে সত্যিকারের **Onuron application emulator** দাঁড়াবে।

এখনই CPU/hardware emulator বানানো একেবারেই দরকার নেই। সেটা অযথা বিশাল প্রকল্প হয়ে যাবে।

## আমার তৈরি করা সম্পূর্ণ code

আমি পুরো implementation package করে দিয়েছি। এর মধ্যে আছে:

```text
nil-emulator-patch/
├── cmd/
│   └── nil/
│       ├── emulator.go
│       └── main.go.patch
├── IMPLEMENTATION.md
└── README.md
```

nilang-lightweight-emulator.zip[**সম্পূর্ণ Nilang Lightweight Emulator Code ZIP**](sandbox:/mnt/data/nilang-lightweight-emulator.zip)

`main.go.patch`-এ শুধু existing `nil` CLI-তে:

```go
case "emulate", "emu":
    cmdEmulate()
```

যোগ করতে হবে এবং help text-এ command-টি যুক্ত করতে হবে।

তারপর:

```powershell
go build -o bin/nil.exe ./cmd/nil
```

এবং:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil
```

চালালেই browser-এ emulator খুলবে।

একটি সীমাবদ্ধতা পরিষ্কার করে বলি: এই পরিবেশে আপনার পুরো repository locally clone করে integration build চালাতে পারিনি, কারণ sandbox থেকে GitHub network access ব্যর্থ হয়েছে। নতুন `emulator.go`-কে `gofmt` দিয়ে যাচাই করেছি, কিন্তু repo-level `go build ./cmd/nil` এখানে সরাসরি চালানো সম্ভব হয়নি। তাই প্রথম local build-এ integration-level ছোটখাটো compile adjustment লাগলে সেটি বাস্তব repo-র বর্তমান state দেখে করতে হবে।

আমার মতে পরের ধাপে **`nil dev` এবং `nil emulate`-এর common UI runtime এক জায়গায় refactor করা** সবচেয়ে ভালো হবে—তাহলে একই event bridge, hot reload এবং renderer দুটো command-এই ব্যবহার হবে।

আপনার বর্তমান Nilang architecture-এ lexer/parser/typecheck/evaluator এবং Alap declarative UI ইতিমধ্যেই আছে; তাই আলাদা Python interpreter বানালে emulator ও আসল Nilang runtime ধীরে ধীরে আলাদা semantics দেখাতে পারে। আগের প্রস্তাবের Python-ভিত্তিক নকশাটিও এই আলাদা interpreter approach নিয়েছিল পেস্ট করা টেক্সট। এখনকার repo-র জন্য সেটি বাদ দেওয়াই ভালো।

# ১. Final architecture

```text
                    ┌─────────────────────────┐
                    │    nil emulate app.nil   │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │ Existing Nilang frontend│
                    │ lexer → parser → check  │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │ Existing evaluator      │
                    │ Tree Walking Runtime    │
                    └────────────┬────────────┘
                                 │
                                 ▼
                    ┌─────────────────────────┐
                    │ declarativeUIApp         │
                    │ Render / Dispatch        │
                    └────────────┬────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
                 Render                    Event
                    │                         │
                    ▼                         ▼
             Alap SSR HTML              /__nil/event
                    │                         │
                    └────────────┬────────────┘
                                 ▼
                       Browser Device Shell
                       ├── Phone
                       ├── Tablet
                       └── Desktop
```

আপনার repo-র `nil render` ইতিমধ্যেই declarative component খুঁজে `Render()` এবং event-এর ক্ষেত্রে `Dispatch()` ব্যবহার করে; emulator সেটাকেই ব্যবহার করবে। `nil dev`-এর server-side event ধারণাটাও এর সঙ্গে মিলে যায়।

---

# ২. নতুন ফাইল

এই ফাইলটি তৈরি করবেন:

```text
cmd/nil/emulator.go
```

পুরো code:

```go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/joysriramsarkar/nilLang/compiler/evaluator"
	"github.com/joysriramsarkar/nilLang/pkg/alap/ui"
	"github.com/joysriramsarkar/nilLang/pkg/config"
)

// emulatorState is intentionally small. It owns one local Nilang application
// instance, watches its source file for changes, and exposes only the HTTP
// endpoints required by the browser shell.
type emulatorState struct {
	mu sync.RWMutex

	app         *declarativeUIApp
	entryPath   string
	lastModTime time.Time
	lastError   string
	device      string
	startedAt   time.Time
	revision    atomic.Uint64
}

const (
	emulatorDefaultPort = 8090
	emulatorPollEvery   = 700 * time.Millisecond
	emulatorEventLimit  = 1 << 20 // 1 MiB JSON body cap
)

// cmdEmulate starts the lightweight browser-based device emulator.
//
// Examples:
//
//	nil emulate
//	nil emulate src/main.nil
//	nil emulate src/main.nil --device phone
//	nil emulate src/main.nil --device tablet
//	nil emulate src/main.nil --device desktop
//	nil emulate src/main.nil --port 8090
//	nil emulate src/main.nil --no-open
func cmdEmulate() {
	entryFile := ""
	port := emulatorDefaultPort
	device := "phone"
	noOpen := false

	for i := 2; i < len(os.Args); i++ {
		arg := os.Args[i]

		switch {
		case arg == "--no-open":
			noOpen = true

		case arg == "--port" && i+1 < len(os.Args):
			i++
			parsed, err := strconv.Atoi(os.Args[i])
			if err != nil || parsed < 1 || parsed > 65535 {
				fatalEmulator("invalid --port: %s", os.Args[i])
			}
			port = parsed

		case strings.HasPrefix(arg, "--port="):
			parsed, err := strconv.Atoi(strings.TrimPrefix(arg, "--port="))
			if err != nil || parsed < 1 || parsed > 65535 {
				fatalEmulator("invalid --port: %s", arg)
			}
			port = parsed

		case arg == "--device" && i+1 < len(os.Args):
			i++
			device = normalizeEmulatorDevice(os.Args[i])

		case strings.HasPrefix(arg, "--device="):
			device = normalizeEmulatorDevice(
				strings.TrimPrefix(arg, "--device="),
			)

		case strings.HasSuffix(arg, ".nil"):
			entryFile = arg

		case strings.HasPrefix(arg, "-"):
			fatalEmulator("unknown emulator option: %s", arg)

		case entryFile == "":
			entryFile = arg

		default:
			fatalEmulator("unexpected argument: %s", arg)
		}
	}

	if entryFile == "" {
		entryFile = resolveEmulatorEntry()
	}

	// Allow:
	//
	//   nil emulate my-app/
	//
	// when my-app contains nil.json.
	if entryFile != "" {
		if info, statErr := os.Stat(entryFile); statErr == nil && info.IsDir() {
			cfg, cfgErr := config.LoadConfig(entryFile)
			if cfgErr != nil {
				fatalEmulator(
					"%s is a directory but has no usable nil.json: %v",
					entryFile,
					cfgErr,
				)
			}

			entryFile = cfg.GetEntryPath(entryFile)
		}
	}

	if entryFile == "" {
		fatalEmulator(
			"no .nil entry file found; pass one explicitly, e.g. `nil emulate src/main.nil`",
		)
	}

	absEntry, err := filepath.Abs(entryFile)
	if err != nil {
		fatalEmulator("cannot resolve entry path: %v", err)
	}

	info, err := os.Stat(absEntry)
	if err != nil {
		fatalEmulator("entry file not found: %s", absEntry)
	}
	if info.IsDir() {
		fatalEmulator("entry path is a directory: %s", absEntry)
	}

	app, err := loadEmulatorUI(absEntry)
	if err != nil {
		fatalEmulator(
			"cannot load %s: %v",
			absEntry,
			err,
		)
	}

	state := &emulatorState{
		app:         app,
		entryPath:   absEntry,
		lastModTime: info.ModTime(),
		device:      device,
		startedAt:   time.Now(),
	}

	state.revision.Store(1)

	mux := http.NewServeMux()
	registerEmulatorRoutes(mux, state)

	server := &http.Server{
		Addr:              fmt.Sprintf("127.0.0.1:%d", port),
		Handler:           emulatorSecurityMiddleware(mux),
		ReadHeaderTimeout: 2 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go emulatorWatchSource(state)

	url := fmt.Sprintf(
		"http://127.0.0.1:%d/?device=%s",
		port,
		state.device,
	)

	fmt.Println("⚡ Nilang Lightweight Emulator")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("  App:     %s\n", absEntry)
	fmt.Printf("  Device:  %s\n", state.device)
	fmt.Printf("  URL:     %s\n", url)
	fmt.Println("  Engine:  existing Nilang frontend + evaluator + Alap renderer")
	fmt.Println("  Reload:  source polling every 700ms")
	fmt.Println("  Network: localhost only")
	fmt.Println("  Stop:    Ctrl+C")
	fmt.Println()

	if !noOpen {
		go func() {
			time.Sleep(300 * time.Millisecond)
			_ = openEmulatorBrowser(url)
		}()
	}

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		fatalEmulator(
			"emulator server failed: %v",
			err,
		)
	}
}

func fatalEmulator(format string, args ...interface{}) {
	fmt.Fprintf(
		os.Stderr,
		"❌ Emulator: "+format+"\n",
		args...,
	)
	os.Exit(1)
}

func normalizeEmulatorDevice(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "phone", "mobile":
		return "phone"

	case "tablet":
		return "tablet"

	case "desktop", "pc":
		return "desktop"

	default:
		fatalEmulator(
			"unsupported device %q; use phone, tablet, or desktop",
			value,
		)
		return "phone"
	}
}

// Resolve an entry automatically:
//
//   src/main.nil
//   main.nil
//
// If neither exists, look for exactly one child directory containing
// nil.json and assume its src/main.nil is the entry.
func resolveEmulatorEntry() string {
	if _, err := os.Stat("src/main.nil"); err == nil {
		return "src/main.nil"
	}

	if _, err := os.Stat("main.nil"); err == nil {
		return "main.nil"
	}

	entries := []string{}

	if list, err := os.ReadDir("."); err == nil {
		for _, item := range list {
			if !item.IsDir() {
				continue
			}

			candidate := filepath.Join(
				item.Name(),
				"nil.json",
			)

			if _, err := os.Stat(candidate); err == nil {
				entries = append(
					entries,
					filepath.Join(
						item.Name(),
						"src",
						"main.nil",
					),
				)
			}
		}
	}

	if len(entries) == 1 {
		return entries[0]
	}

	return ""
}

// Load a .nil file using the existing Nilang declarative UI loader.
func loadEmulatorUI(entryPath string) (*declarativeUIApp, error) {
	source, err := os.ReadFile(entryPath)
	if err != nil {
		return nil, err
	}

	return loadEmulatorUIFromSource(
		entryPath,
		string(source),
	)
}

// Match `nil run` import resolution.
//
// Example:
//
//	src/main.nil
//	   |
//	   └── import "./components/foo.nil"
//
// resolves relative to src/.
func loadEmulatorUIFromSource(
	entryPath string,
	source string,
) (*declarativeUIApp, error) {

	evaluator.PushScriptDir(
		filepath.Dir(entryPath),
	)

	defer evaluator.PopScriptDir()

	return loadDeclarativeUI(source)
}

// Hot reload watcher.
//
// Deliberately uses standard-library polling instead of fsnotify so the
// emulator does not acquire another dependency merely for development UX.
func emulatorWatchSource(state *emulatorState) {
	ticker := time.NewTicker(
		emulatorPollEvery,
	)

	defer ticker.Stop()

	for range ticker.C {
		info, err := os.Stat(state.entryPath)

		if err != nil {
			state.mu.Lock()
			state.lastError = fmt.Sprintf(
				"source unavailable: %v",
				err,
			)
			state.mu.Unlock()

			continue
		}

		state.mu.RLock()

		changed := info.ModTime().After(
			state.lastModTime,
		)

		state.mu.RUnlock()

		if !changed {
			continue
		}

		// Parse/evaluate outside the state lock so a bad edit does not block
		// an already-running UI.
		newApp, err := loadEmulatorUI(
			state.entryPath,
		)

		state.mu.Lock()

		if err != nil {
			// Keep the last good application alive.
			state.lastError = err.Error()
		} else {
			state.app = newApp
			state.lastError = ""
			state.lastModTime = info.ModTime()

			state.revision.Add(1)
		}

		state.mu.Unlock()
	}
}

func registerEmulatorRoutes(
	mux *http.ServeMux,
	state *emulatorState,
) {

	// Browser shell.
	mux.HandleFunc("/", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		device := normalizeDeviceFromQuery(
			r.URL.Query().Get("device"),
		)

		writeHTML(
			w,
			emulatorShellHTML(state, device),
		)
	})

	// Actual Nilang application screen.
	mux.HandleFunc("/app", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		state.mu.RLock()

		app := state.app
		lastErr := state.lastError
		revision := state.revision.Load()

		state.mu.RUnlock()

		if app == nil {
			writePlain(
				w,
				http.StatusServiceUnavailable,
				"NilLang app is not loaded",
			)
			return
		}

		html, err := renderEmulatorApp(
			app,
			revision,
			lastErr,
		)

		if err != nil {
			writePlain(
				w,
				http.StatusInternalServerError,
				err.Error(),
			)
			return
		}

		writeHTML(w, html)
	})

	// Hot-reload/status endpoint.
	mux.HandleFunc("/__nil/status", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodGet {
			methodNotAllowed(w)
			return
		}

		state.mu.RLock()

		errText := state.lastError
		entry := state.entryPath
		device := state.device
		started := state.startedAt

		state.mu.RUnlock()

		payload := map[string]interface{}{
			"ok":        errText == "",
			"revision":  state.revision.Load(),
			"error":     errText,
			"entry":     entry,
			"device":    device,
			"uptime_ms": time.Since(started).Milliseconds(),
		}

		writeJSON(
			w,
			http.StatusOK,
			payload,
		)
	})

	// Nilang UI event endpoint.
	mux.HandleFunc("/__nil/event", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}

		var request struct {
			Event   string      `json:"event"`
			Payload interface{} `json:"payload,omitempty"`
		}

		limited := http.MaxBytesReader(
			w,
			r.Body,
			emulatorEventLimit,
		)

		defer limited.Close()

		decoder := json.NewDecoder(limited)

		if err := decoder.Decode(&request); err != nil {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]interface{}{
					"ok": false,
					"error": "invalid event JSON: " +
						err.Error(),
				},
			)
			return
		}

		if strings.TrimSpace(request.Event) == "" {
			writeJSON(
				w,
				http.StatusBadRequest,
				map[string]interface{}{
					"ok":    false,
					"error": "event is required",
				},
			)
			return
		}

		state.mu.Lock()
		defer state.mu.Unlock()

		if state.app == nil {
			writeJSON(
				w,
				http.StatusServiceUnavailable,
				map[string]interface{}{
					"ok":    false,
					"error": "app is not loaded",
				},
			)
			return
		}

		if request.Payload != nil {
			if err := state.app.Dispatch(
				request.Event,
				request.Payload,
			); err != nil {
				writeJSON(
					w,
					http.StatusBadRequest,
					map[string]interface{}{
						"ok":    false,
						"error": err.Error(),
					},
				)
				return
			}
		} else {
			if err := state.app.Dispatch(
				request.Event,
			); err != nil {
				writeJSON(
					w,
					http.StatusBadRequest,
					map[string]interface{}{
						"ok":    false,
						"error": err.Error(),
					},
				)
				return
			}
		}

		state.lastError = ""

		writeJSON(
			w,
			http.StatusOK,
			map[string]interface{}{
				"ok":       true,
				"revision": state.revision.Load(),
			},
		)
	})
}

func normalizeDeviceFromQuery(value string) string {
	if value == "" {
		return "phone"
	}

	switch strings.ToLower(value) {
	case "tablet":
		return "tablet"

	case "desktop":
		return "desktop"

	default:
		return "phone"
	}
}

func renderEmulatorApp(
	app *declarativeUIApp,
	revision uint64,
	lastError string,
) (string, error) {

	page, err := app.Render()
	if err != nil {
		return "", err
	}

	state := map[string]interface{}(nil)

	if app.State() != nil {
		state = hashToGoMap(
			app.State(),
		)
	}

	// IMPORTANT:
	// Use the existing Alap SSR renderer.
	html := page.RenderSSR(
		ui.OnuronTheme(),
		state,
	)

	// Add the emulator's event bridge after Alap has generated its HTML.
	html = injectBeforeBody(
		html,
		emulatorBridgeScript(
			revision,
			lastError,
		),
	)

	return html, nil
}

func emulatorBridgeScript(
	revision uint64,
	lastError string,
) string {

	errorJSON, _ := json.Marshal(lastError)

	return fmt.Sprintf(`<script>
(() => {
  const REVISION = %d;
  const LAST_ERROR = %s;

  const root =
    document.querySelector('[data-alap-root]')
    || document.body;

  function payloadFrom(target) {
    const raw =
      target &&
      target.getAttribute('data-alap-payload');

    if (!raw) {
      return null;
    }

    try {
      return JSON.parse(raw);
    } catch (_) {
      return raw;
    }
  }

  async function dispatch(
    eventName,
    payload
  ) {
    try {
      const response =
        await fetch('/__nil/event', {
          method: 'POST',
          headers: {
            'Content-Type':
              'application/json'
          },
          body: JSON.stringify({
            event: eventName,
            payload: payload
          })
        });

      const data =
        await response.json();

      if (!response.ok || !data.ok) {
        window.parent.postMessage(
          {
            type: 'nil-emulator-error',
            error:
              data.error ||
              'event failed'
          },
          '*'
        );

        return;
      }

      window.parent.postMessage(
        {
          type: 'nil-emulator-event',
          revision:
            data.revision
        },
        '*'
      );

      // Re-render from the same server-side Nilang state.
      window.location.reload();

    } catch (err) {
      window.parent.postMessage(
        {
          type:
            'nil-emulator-error',
          error:
            String(err)
        },
        '*'
      );
    }
  }

  root.addEventListener(
    'click',
    (event) => {
      let node = event.target;

      while (
        node &&
        node !== root
      ) {
        const name =
          node.getAttribute &&
          node.getAttribute(
            'data-alap-click'
          );

        if (name) {
          event.preventDefault();

          if (node.disabled) {
            return;
          }

          dispatch(
            name,
            payloadFrom(node)
          );

          return;
        }

        node = node.parentNode;
      }
    }
  );

  window.parent.postMessage(
    {
      type:
        'nil-emulator-ready',
      revision: REVISION,
      error: LAST_ERROR
    },
    '*'
  );

})();
</script>`,
		revision,
		string(errorJSON),
	)
}

func injectBeforeBody(
	html string,
	fragment string,
) string {

	idx := strings.LastIndex(
		strings.ToLower(html),
		"</body>",
	)

	if idx < 0 {
		return html + fragment
	}

	return html[:idx] +
		fragment +
		html[idx:]
}

func emulatorShellHTML(
	state *emulatorState,
	device string,
) string {

	state.mu.RLock()

	entry := state.entryPath
	revision := state.revision.Load()
	errorText := state.lastError

	state.mu.RUnlock()

	entryJS, _ := json.Marshal(entry)
	errorJS, _ := json.Marshal(errorText)

	deviceClass := "device-phone"

	switch device {
	case "tablet":
		deviceClass = "device-tablet"

	case "desktop":
		deviceClass = "device-desktop"
	}

	return fmt.Sprintf(`<!doctype html>
<html>
<head>
<meta charset="utf-8">
<meta
  name="viewport"
  content="width=device-width,initial-scale=1"
>
<title>Nilang Emulator</title>

<style>
* {
  box-sizing: border-box;
}

html,
body {
  height: 100%%;
}

body {
  margin: 0;
  background: #0b0d12;
  color: #e5e7eb;
  font: 14px/1.45
    system-ui,
    -apple-system,
    "Segoe UI",
    sans-serif;
  overflow: hidden;
}

.toolbar {
  height: 52px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 14px;
  border-bottom:
    1px solid #202530;
  background: #11141b;
}

.brand {
  font-weight: 750;
  letter-spacing: .2px;
  margin-right: auto;
}

.brand span {
  color: #4dd0e1;
}

.control {
  background: #191d26;
  color: #e5e7eb;
  border:
    1px solid #303644;
  border-radius: 7px;
  padding: 7px 9px;
}

.btn {
  cursor: pointer;
}

.status {
  min-width: 95px;
  text-align: center;
  font-size: 12px;
  border-radius: 99px;
  padding: 5px 10px;
  background: #17231a;
  color: #71e294;
}

.status.error {
  background: #2d1717;
  color: #ff8b8b;
}

.main {
  height: calc(100%% - 52px);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
  padding: 24px;
}

.stage {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%%;
  height: 100%%;
}

.device {
  position: relative;
  background: #050608;

  box-shadow:
    0 18px 60px
    rgba(0,0,0,.45),
    0 0 0 1px #30343d;

  border: 7px solid #1b1f27;
  overflow: hidden;

  transition:
    width .18s,
    height .18s,
    border-radius .18s;
}

.device-phone {
  width: 410px;
  height: 872px;
  border-radius: 36px;
}

.device-tablet {
  width: 860px;
  height: 1180px;
  max-height:
    calc(100vh - 100px);
  border-radius: 28px;
}

.device-desktop {
  width: min(
    1240px,
    92vw
  );

  height: min(
    760px,
    calc(100vh - 92px)
  );

  border-radius: 14px;
  border-width: 5px;
}

.camera {
  position: absolute;
  z-index: 3;
  top: 0;
  left: 50%%;
  transform:
    translateX(-50%%);

  width: 84px;
  height: 16px;

  background: #06070a;
  border-radius:
    0 0 12px 12px;
}

.screen {
  position: absolute;
  inset: 0;

  border: 0;
  width: 100%%;
  height: 100%%;

  background: #0a0a0f;
}

.errorbar {
  position: fixed;
  left: 50%%;
  bottom: 14px;

  transform:
    translateX(-50%%);

  max-width:
    min(800px, 90vw);

  background: #351919;
  color: #ffb1b1;

  border:
    1px solid #793333;

  padding: 9px 13px;

  border-radius: 8px;

  font-size: 12px;

  display: none;

  z-index: 9;

  white-space: pre-wrap;
}

@media (max-width: 700px) {

  .main {
    padding: 10px;
  }

  .device-phone {
    transform: scale(.84);
    transform-origin:
      center;
  }

  .toolbar {
    gap: 6px;
  }

  .brand {
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
</head>

<body>

<div class="toolbar">

  <div class="brand">
    নীলাং <span>Emulator</span>
  </div>

  <select
    id="device"
    class="control"
  >
    <option
      value="phone"
      %s
    >
      Phone
    </option>

    <option
      value="tablet"
      %s
    >
      Tablet
    </option>

    <option
      value="desktop"
      %s
    >
      Desktop
    </option>
  </select>

  <button
    id="reload"
    class="control btn"
  >
    ↻ Reload
  </button>

  <div
    id="status"
    class="status"
  >
    RUNNING
  </div>

</div>

<div class="main">

  <div class="stage">

    <div class="device %s">

      <div class="camera"></div>

      <iframe
        id="screen"
        class="screen"
        src="/app?rev=%d"
        title="Nilang application"
      ></iframe>

    </div>

  </div>

</div>

<div
  id="errorbar"
  class="errorbar"
></div>

<script>

const ENTRY = %s;

let lastRevision = %d;
let initialError = %s;

const frame =
  document.getElementById(
    'screen'
  );

const status =
  document.getElementById(
    'status'
  );

const errorbar =
  document.getElementById(
    'errorbar'
  );

function setStatus(
  ok,
  msg
) {

  status.textContent =
    ok ? 'RUNNING' : 'ERROR';

  status.className =
    'status' +
    (ok ? '' : ' error');

  if (msg) {

    errorbar.textContent =
      msg;

    errorbar.style.display =
      'block';

  } else {

    errorbar.style.display =
      'none';
  }
}

function reloadApp() {

  frame.src =
    '/app?rev=' +
    Date.now();
}

document
  .getElementById('reload')
  .addEventListener(
    'click',
    reloadApp
  );

document
  .getElementById('device')
  .addEventListener(
    'change',
    e => {

      const q =
        new URLSearchParams(
          location.search
        );

      q.set(
        'device',
        e.target.value
      );

      location.search =
        q.toString();
    }
  );

window.addEventListener(
  'message',
  e => {

    if (
      !e.data ||
      typeof e.data !==
        'object'
    ) {
      return;
    }

    if (
      e.data.type ===
      'nil-emulator-ready'
    ) {

      if (e.data.error) {

        setStatus(
          false,
          e.data.error
        );

      } else {

        setStatus(
          true,
          ''
        );
      }
    }

    if (
      e.data.type ===
      'nil-emulator-event'
    ) {

      lastRevision =
        e.data.revision ||
        lastRevision;

      setStatus(
        true,
        ''
      );
    }

    if (
      e.data.type ===
      'nil-emulator-error'
    ) {

      setStatus(
        false,
        e.data.error ||
          'event failed'
      );
    }

  }
);

async function poll() {

  try {

    const r =
      await fetch(
        '/__nil/status',
        {
          cache:
            'no-store'
        }
      );

    const s =
      await r.json();

    if (
      s.revision !==
      lastRevision
    ) {

      lastRevision =
        s.revision;

      reloadApp();
    }

    setStatus(
      s.ok,
      s.error || ''
    );

  } catch (err) {

    setStatus(
      false,
      String(err)
    );
  }
}

if (initialError) {

  setStatus(
    false,
    initialError
  );
}

setInterval(
  poll,
  700
);

poll();

</script>

</body>
</html>`,
		selected(device == "phone"),
		selected(device == "tablet"),
		selected(device == "desktop"),
		deviceClass,
		revision,
		string(entryJS),
		revision,
		string(errorJS),
	)
}

func selected(value bool) string {
	if value {
		return "selected"
	}

	return ""
}

// Bind only to localhost.
//
// This is a development emulator; it should never accidentally become
// reachable from another device on the LAN.
func emulatorSecurityMiddleware(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {

			if !strings.HasPrefix(
				r.RemoteAddr,
				"127.0.0.1:",
			) &&
				!strings.HasPrefix(
					r.RemoteAddr,
					"[::1]:",
				) &&
				!strings.HasPrefix(
					r.RemoteAddr,
					"::1:",
				) {

				writePlain(
					w,
					http.StatusForbidden,
					"Nilang emulator is localhost-only",
				)

				return
			}

			w.Header().Set(
				"Cache-Control",
				"no-store",
			)

			w.Header().Set(
				"X-Content-Type-Options",
				"nosniff",
			)

			w.Header().Set(
				"X-Frame-Options",
				"SAMEORIGIN",
			)

			next.ServeHTTP(w, r)
		},
	)
}

func methodNotAllowed(
	w http.ResponseWriter,
) {
	writePlain(
		w,
		http.StatusMethodNotAllowed,
		"method not allowed",
	)
}

func writeHTML(
	w http.ResponseWriter,
	value string,
) {
	w.Header().Set(
		"Content-Type",
		"text/html; charset=utf-8",
	)

	w.WriteHeader(
		http.StatusOK,
	)

	_, _ = w.Write(
		[]byte(value),
	)
}

func writePlain(
	w http.ResponseWriter,
	status int,
	value string,
) {
	w.Header().Set(
		"Content-Type",
		"text/plain; charset=utf-8",
	)

	w.WriteHeader(status)

	_, _ = w.Write(
		[]byte(value),
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value interface{},
) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(
		value,
	)
}

func openEmulatorBrowser(
	url string,
) error {

	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "windows":
		cmd = exec.Command(
			"rundll32",
			"url.dll,FileProtocolHandler",
			url,
		)

	case "darwin":
		cmd = exec.Command(
			"open",
			url,
		)

	default:
		cmd = exec.Command(
			"xdg-open",
			url,
		)
	}

	return cmd.Start()
}
```

---

# ৩. `cmd/nil/main.go`

আপনার বর্তমান `main.go` পুরোটা পাল্টাতে হবে না। শুধু switch-এর মধ্যে এই অংশ যোগ করুন।

বর্তমানে আছে:

```go
case "dev":
    cmdDev()
case "routes":
    cmdRoutes()
```

এটাকে করুন:

```go
case "dev":
    cmdDev()

case "emulate", "emu":
    cmdEmulate()

case "routes":
    cmdRoutes()
```

অর্থাৎ `emu`-ও alias থাকবে।

তারপর `printUsage()`-এ এই line যোগ করুন:

```go
fmt.Println("  emulate [file.nil] [--device phone|tablet|desktop] [--port 8090] [--no-open]  লাইটওয়েট ডিভাইস emulator চালান")
```

আর examples-এর শেষে:

```go
fmt.Println("  nil emulate examples/ui-counter/src/main.nil --device phone")
```

## সম্পূর্ণ পরিবর্তিত `main.go`

আপনি সরাসরি existing `main.go`-এর command dispatch অংশে উপরের পরিবর্তন করলেই যথেষ্ট। অন্য সব command untouched থাকবে।

---

# ৪. আপনার `ui-counter` দিয়ে কী হবে

আপনার বর্তমান example:

```nil
component Counter {
    state count: i32 = 0;

    render {
        return {
            "type": "Page",
            "title": "Alap UI Counter",
            "content": [
                {
                    "type": "Card",
                    "title": "Current Count",
                    "body": "\(count)"
                },
                {
                    "type": "Button",
                    "id": "increment-count",
                    "label": "Increment",
                    "event": "increment",
                    "payload": {"step": 1}
                }
            ],
            "footer": "NilLang declarative component"
        };
    }

    on increment(payload) {
        let step = 1;

        if (payload != null) {
            step = payload["step"];
        }

        count = count + step;

        emit("changed", count);
    }
}
```

Emulator-এ:

```text
┌──────────────────────────┐
│                          │
│     Alap UI Counter      │
│                          │
│  ┌────────────────────┐  │
│  │ Current Count      │  │
│  │        0           │  │
│  └────────────────────┘  │
│                          │
│     [  Increment  ]      │
│                          │
└──────────────────────────┘
```

`Increment` চাপলে:

```text
Button click
    ↓
data-alap-click="increment"
    ↓
POST /__nil/event
    ↓
app.Dispatch("increment", payload)
    ↓
count = count + step
    ↓
app.Render()
    ↓
Alap SSR
    ↓
iframe reload
    ↓
Count = 1
```

আবার চাপলে:

```text
1 → 2 → 3 → 4 → ...
```

---

# ৫. Hot reload

ধরুন আপনি লিখলেন:

```nil
"title": "আমার প্রথম নীলাং অ্যাপ"
```

save করলেন।

Emulator:

```text
700ms
   ↓
source modified?
   ↓
yes
   ↓
loadEmulatorUI()
   ↓
lexer
parser
typecheck
evaluator
   ↓
success
   ↓
revision++
   ↓
browser reload
```

আর syntax error করলে পুরনো working application বন্ধ হবে না।

যেমন ভুল করে:

```nil
"title": "আমার প্রথম নীলাং অ্যাপ
```

তাহলে emulator-এর নিচে:

```text
ERROR
[Parser] ...
```

দেখাবে, কিন্তু আগের working UI-টাও রেখে দেবে।

এটা development-এর জন্য বেশ গুরুত্বপূর্ণ।

---

# ৬. Emulator device profile

## Phone

```text
410 × 872
```

## Tablet

```text
860 × 1180
```

## Desktop

```text
1240 × 760
```

এগুলো hardware emulation নয়। এগুলো visual viewport।

অর্থাৎ এটা:

```text
UI responsive কি না?
button কোথায় বসছে?
overflow হচ্ছে?
text ভেঙে যাচ্ছে?
layout ঠিক আছে?
```

এসব পরীক্ষা করার জন্য।

---

# ৭. চালানোর command

Build:

```powershell
go build -o bin/nil.exe ./cmd/nil
```

তারপর:

```powershell
.\bin\nil.exe emulate
```

Project root-এ `src/main.nil` থাকলে নিজে সেটা নেবে।

সরাসরি:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil
```

Phone:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil --device phone
```

Tablet:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil --device tablet
```

Desktop:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil --device desktop
```

Port:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil --port 8095
```

Browser auto-open বন্ধ:

```powershell
.\bin\nil.exe emulate examples\ui-counter\src\main.nil --no-open
```

তারপর নিজে browser-এ:

```text
http://127.0.0.1:8095/
```

---

# ৮. কেন এটা lightweight

এখানে এগুলো নেই:

```text
Electron              ❌
Tauri                 ❌
Chromium bundled      ❌
Android Emulator      ❌
QEMU                  ❌
VirtualBox             ❌
WebView dependency     ❌
Selenium               ❌
Playwright             ❌
fsnotify dependency    ❌
নতুন Go runtime        ❌
নতুন Nilang parser     ❌
নতুন Nilang VM         ❌
```

শুধু:

```text
Go standard library
+
existing Nilang runtime
+
existing Alap UI
+
user-এর browser
```

এটাই মূল সুবিধা।

---

# ৯. কিন্তু একটা জিনিস পরিষ্কার

এটাকে এখনই **Android Emulator** বলব না।

এটা:

> **Nilang Lightweight Application Emulator**

অর্থাৎ:

```text
App logic
✓
Nilang execution
✓
Alap UI
✓
state
✓
events
✓
responsive viewport
✓
hot reload

CPU emulation
✗

Android kernel
✗

real sensor
✗

real filesystem sandbox
✗

camera
✗

GPS
✗

Bluetooth
✗

SoftBus hardware
✗
```

এটাই প্রথম version-এর সঠিক সীমা।

---

# ১০. ভবিষ্যতে এটাকে আরও শক্তিশালী করা

এরপর `emulator.go`-র মধ্যে সব কিছু ঢোকাব না। একটা capability layer বানাব:

```text
pkg/emulator/
│
├── environment.go
├── filesystem.go
├── network.go
├── storage.go
├── device.go
├── permissions.go
└── emulator.go
```

তখন:

```text
Nilang Application
        │
        ▼
   Alap Runtime
        │
        ▼
 ┌────────────────────┐
 │ Emulator API Layer │
 ├────────────────────┤
 │ Filesystem Mock    │
 │ Network Mock       │
 │ Storage Mock       │
 │ Device Info        │
 │ Permission Mock    │
 │ Camera Mock        │
 │ Location Mock      │
 │ Clock              │
 └────────────────────┘
```

যেমন Nilang app:

```nilang
let camera = device.camera();
```

emulator-এ সেটা:

```text
[ Camera Mock ]
```

দেবে।

আর Onuron OS-এর আসল runtime-এ:

```text
device.camera()
      ↓
real camera driver
```

চলবে।

এভাবে **একই Nilang API** emulator ও real OS—দুটোতেই থাকবে।

---

# ১১. আরও ভালো ভবিষ্যৎ architecture

শেষ পর্যন্ত আমি structure-টা এভাবে করব:

```text
                     Nilang
                       │
              ┌────────┴────────┐
              │                 │
        Developer Mode      Production
              │                 │
        Tree Evaluator      Bytecode VM
              │                 │
              └────────┬────────┘
                       │
                  Alap Runtime
                       │
             ┌─────────┴─────────┐
             │                   │
        Emulator Target      Onuron Target
             │                   │
      Browser Device UI      Native UI
             │                   │
      Mock Capabilities     Real Capabilities
```

এতে emulator কোনো আলাদা ভাষা হবে না; **Nilang-এর আরেকটি execution target** হয়ে যাবে।

এটাই আমার মতে আপনার Onuron ecosystem-এর জন্য সবচেয়ে পরিষ্কার পথ।

---

## ১২. প্রস্তুত ZIP

আমি code-গুলো আলাদা ফাইলে রেখেও দিয়েছি:

nilang-lightweight-emulator-full.zip[**Nilang Lightweight Emulator — Full Source ZIP**](sandbox:/mnt/data/nilang-lightweight-emulator-full.zip)

ভেতরে:

```text
nilang-emulator-full/
├── README.md
└── cmd/
    └── nil/
        ├── emulator.go
        └── main.go
```

`emulator.go`-তে প্রায় ৬৩০ লাইন আছে, আর `main.go`-তে emulator command-সহ সম্পূর্ণ CLI file রাখা হয়েছে।

আমি এখানে `gofmt`-level validation করেছি; পুরো `nilLang` repository clone করে `go test ./cmd/nil` চালানো সম্ভব হয়নি, কারণ এই runtime থেকে GitHub repository clone করার network access ছিল না। তাই এটাকে **repo-compatible implementation** হিসেবে দিচ্ছি, কিন্তু আপনার Windows checkout-এ প্রথম `go build`-এ integration compile error থাকলে সেটাই পরবর্তী সংশোধনের ভিত্তি হবে।

সবচেয়ে গুরুত্বপূর্ণভাবে, **আলাদা Python emulator এখন আর বানানোর দরকার নেই**। আপনার existing Nilang evaluator + Alap runtime-কে ঘিরে এই browser-device emulator-ই lightweight এবং ভবিষ্যৎ-proof পথ।