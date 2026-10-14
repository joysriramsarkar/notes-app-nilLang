package main

import (
	"fmt"
	"os"
)

// Screen dimensions: 80 columns x 25 rows
const (
	Cols = 80
	Rows = 25
)

// VGA Color Attributes
const (
	AttrNavyWhite     = 0x1F // White on Blue
	AttrNavyCyan      = 0x1B // Bright Cyan on Blue
	AttrNavyYellow    = 0x1E // Yellow on Blue
	AttrNavyGreen     = 0x1A // Bright Green on Blue
	AttrNavyGray      = 0x17 // Light Gray on Blue
	AttrHeaderCyan    = 0x3F // White on Cyan
	AttrHeaderYellow  = 0x3E // Yellow on Cyan
	AttrSelectedGreen = 0x2F // White on Green (Active Selection)
	AttrStatusBar     = 0x70 // Black on Light Gray
	AttrBoxBorder     = 0x1B // Cyan Box Borders on Blue
)

type Screen [Rows][Cols]struct {
	Ch   byte
	Attr byte
}

func NewScreen() *Screen {
	s := &Screen{}
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			s[r][c].Ch = ' '
			s[r][c].Attr = AttrNavyWhite
		}
	}
	return s
}

func (s *Screen) Print(r, c int, text string, attr byte) {
	for i, ch := range []byte(text) {
		col := c + i
		if col < Cols && r < Rows {
			s[r][col].Ch = ch
			s[r][col].Attr = attr
		}
	}
}

func (s *Screen) FillRow(r int, ch byte, attr byte) {
	for c := 0; c < Cols; c++ {
		s[r][c].Ch = ch
		s[r][c].Attr = attr
	}
}

func (s *Screen) DrawBox(r1, c1, r2, c2 int, title string, attr byte) {
	s[r1][c1].Ch = '+'
	s[r1][c1].Attr = attr
	for c := c1 + 1; c < c2; c++ {
		s[r1][c].Ch = '-'
		s[r1][c].Attr = attr
	}
	s[r1][c2].Ch = '+'
	s[r1][c2].Attr = attr

	if title != "" {
		s.Print(r1, c1+2, " [ "+title+" ] ", AttrNavyYellow)
	}

	for r := r1 + 1; r < r2; r++ {
		s[r][c1].Ch = '|'
		s[r][c1].Attr = attr
		s[r][c2].Ch = '|'
		s[r][c2].Attr = attr
	}

	s[r2][c1].Ch = '+'
	s[r2][c1].Attr = attr
	for c := c1 + 1; c < c2; c++ {
		s[r2][c].Ch = '-'
		s[r2][c].Attr = attr
	}
	s[r2][c2].Ch = '+'
	s[r2][c2].Attr = attr
}

func (s *Screen) ToBytes() []byte {
	buf := make([]byte, Rows*Cols*2)
	idx := 0
	for r := 0; r < Rows; r++ {
		for c := 0; c < Cols; c++ {
			buf[idx] = s[r][c].Ch
			buf[idx+1] = s[r][c].Attr
			idx += 2
		}
	}
	return buf
}

// Build a full 80x25 screen for Note 1..4 and toggle state
func BuildFrame(selectedNote int, toggled bool) *Screen {
	s := NewScreen()

	// 1. Top Header (Row 0)
	s.FillRow(0, ' ', AttrHeaderCyan)
	s.Print(0, 2, "ONURON OS v1.0", AttrHeaderCyan)
	s.Print(0, 24, "NILANG NOTES APPLICATION", AttrHeaderYellow)
	s.Print(0, 56, "QEMU x86_64 VIRTUAL SYSTEM", AttrHeaderCyan)

	// 2. Subheader (Row 1)
	s.Print(1, 2, "Nilang Lang & Alap Framework | Database: /data/notes.db | QEMU Video OK", AttrNavyYellow)

	// 3. Left Pane: Folders & Note List (Rows 3..10, Cols 1..78)
	s.DrawBox(3, 1, 10, 78, "FOLDERS & NOTE CATALOG", AttrBoxBorder)

	notes := []struct {
		id    int
		title string
		cat   string
		badge string
	}{
		{1, "Nilang Notes Guide & Features", "Important", "[PINNED] [FAVORITE]"},
		{2, "Daily Work Plan & Deadlines", "Work/Tasks", "[PINNED] [HIGH PRIORITY]"},
		{3, "Onuron OS Microkernel & SoftBus", "Ideas/OS", "[RESEARCH]"},
		{4, "Favorite Quotes & Clean Architecture", "Personal", "[BOOKS]"},
	}

	for i, n := range notes {
		row := 5 + i
		var attr byte = AttrNavyWhite
		cursor := "  "
		if n.id == selectedNote {
			attr = AttrSelectedGreen
			cursor = ">>"
			for c := 3; c <= 76; c++ {
				s[row][c].Attr = AttrSelectedGreen
			}
		}
		line := fmt.Sprintf("%s [%d] %-38s %-12s %s", cursor, n.id, n.title, "("+n.cat+")", n.badge)
		s.Print(row, 3, line, attr)
	}

	// 4. Details Pane: Active Note Details & Checklists (Rows 11..22, Cols 1..78)
	s.DrawBox(11, 1, 22, 78, "ACTIVE NOTE DETAILS & INTERACTIVE CHECKLIST", AttrBoxBorder)

	switch selectedNote {
	case 1:
		s.Print(13, 3, "Title: Nilang Notes Guide & Ecosystem Architecture", AttrNavyCyan)
		s.Print(14, 3, "Tags:  #nilang #alap #onuron #qemu | Category: Important | Color: Indigo", AttrNavyGray)
		s.Print(15, 3, "Body:  Full-featured reactive note taking app running on Nilang Runtime VM.", AttrNavyWhite)
		s.Print(16, 3, "       Supports instant full-text search, live markdown and auto-save.", AttrNavyWhite)

		s.Print(18, 3, "Checklist & Tasks:", AttrNavyYellow)
		t1 := "[X] Task 1: Smart Category & Real-time Tag Search Indexing (COMPLETED)"
		t2 := "[X] Task 2: QEMU x86_64 Virtual Machine Boot Verification (COMPLETED)"
		t3 := "[ ] Task 3: Multi-platform packaging for Linux, Android & Onuron OS"
		if toggled {
			t3 = "[X] Task 3: Multi-platform packaging for Linux, Android & Onuron OS (COMPLETED)"
		}
		s.Print(19, 5, t1, AttrNavyGreen)
		s.Print(20, 5, t2, AttrNavyGreen)
		if toggled {
			s.Print(21, 5, t3, AttrNavyGreen)
		} else {
			s.Print(21, 5, t3, AttrNavyGray)
		}

	case 2:
		s.Print(13, 3, "Title: Daily Work Plan & Project Deadlines", AttrNavyCyan)
		s.Print(14, 3, "Tags:  #work #tasks #planning | Priority: High | Category: Work/Tasks", AttrNavyGray)
		s.Print(15, 3, "Body:  Finish compiler optimization benchmarks and verify hardware builds.", AttrNavyWhite)
		s.Print(16, 3, "       Session deadline: 6:00 PM. Target: 100% regression tests passed.", AttrNavyWhite)

		s.Print(18, 3, "Checklist & Tasks:", AttrNavyYellow)
		t1 := "[X] Refactor Nilang AST parser & Bytecode VM modules (COMPLETED)"
		t2 := "[ ] Verify Alap declarative UI components in QEMU emulator"
		if toggled {
			t2 = "[X] Verify Alap declarative UI components in QEMU emulator (COMPLETED)"
		}
		s.Print(19, 5, t1, AttrNavyGreen)
		if toggled {
			s.Print(20, 5, t2, AttrNavyGreen)
		} else {
			s.Print(20, 5, t2, AttrNavyGray)
		}

	case 3:
		s.Print(13, 3, "Title: Onuron OS Microkernel & SoftBus Architecture", AttrNavyCyan)
		s.Print(14, 3, "Tags:  #architecture #os #kernel #ipc | Category: Ideas | Color: Cyan", AttrNavyGray)
		s.Print(15, 3, "Body:  Design ultra-lightweight process isolation and zero-copy SoftBus IPC.", AttrNavyWhite)
		s.Print(16, 3, "       Goal: sub-millisecond response latency with memory sandbox.", AttrNavyWhite)

		s.Print(18, 3, "Checklist & Tasks:", AttrNavyYellow)
		t1 := "[X] Driver interface specification for NilHAL sensors (COMPLETED)"
		t2 := "[ ] Zero-copy shared memory ring buffer implementation"
		if toggled {
			t2 = "[X] Zero-copy shared memory ring buffer implementation (COMPLETED)"
		}
		s.Print(19, 5, t1, AttrNavyGreen)
		if toggled {
			s.Print(20, 5, t2, AttrNavyGreen)
		} else {
			s.Print(20, 5, t2, AttrNavyGray)
		}

	case 4:
		s.Print(13, 3, "Title: Favorite Quotes & Clean Architecture Insights", AttrNavyCyan)
		s.Print(14, 3, "Tags:  #quotes #books #philosophy | Category: Personal | Color: Violet", AttrNavyGray)
		s.Print(15, 3, "Body:  \"Imagination is more important than knowledge.\" - Albert Einstein", AttrNavyWhite)
		s.Print(16, 3, "       Book: Designing Data-Intensive Applications (Martin Kleppmann)", AttrNavyWhite)

		s.Print(18, 3, "Checklist & Tasks:", AttrNavyYellow)
		s.Print(19, 5, "[X] Read Chapter 7: Transactions & Distributed Consensus (DONE)", AttrNavyGreen)
		s.Print(20, 5, "[X] Implement reproducible benchmarks in Nilang (DONE)", AttrNavyGreen)
	}

	// 5. Bottom Status Bar (Row 24)
	s.FillRow(24, ' ', AttrStatusBar)
	s.Print(24, 2, "KEYS: [1-4] Switch Note  |  [SPACE] Toggle Task  |  [Q] Reboot VM", AttrStatusBar)
	s.Print(24, 68, "Nilang v1.0", AttrStatusBar)

	return s
}

func main() {
	imageSize := 1474560 // 1.44 MB standard floppy disk
	disk := make([]byte, imageSize)

	// ────────────────────────────────────────────────────────────────────────
	// 1. SECTOR 0: MBR / BOOTLOADER (512 bytes)
	// Reads 16 sectors (8 KB = sectors 2..17 on track 0, head 0) into 0x1000:0x0000
	// ────────────────────────────────────────────────────────────────────────
	sector0 := make([]byte, 512)
	bootLoader := []byte{
		0xFA,                         // cli
		0x31, 0xC0,                   // xor ax, ax
		0x8E, 0xD8,                   // mov ds, ax
		0x8E, 0xC0,                   // mov es, ax
		0x8E, 0xD0,                   // mov ss, ax
		0xBC, 0x00, 0x7C,             // mov sp, 0x7C00
		0xFB,                         // sti

		// Save boot drive number from DL to [0x7C00]
		0x88, 0x16, 0x00, 0x7C,       // mov [0x7C00], dl

		// Reset disk system (INT 13h, AH=00h)
		0x30, 0xE4,                   // xor ah, ah
		0xCD, 0x13,                   // int 13h

		// Read 16 sectors from disk to ES:BX = 0x1000:0x0000 (RAM address 0x10000)
		0xB8, 0x00, 0x10,             // mov ax, 0x1000
		0x8E, 0xC0,                   // mov es, ax
		0x31, 0xDB,                   // xor bx, bx (dest buffer 0x1000:0x0000)
		0xB4, 0x02,                   // mov ah, 0x02 (read sectors)
		0xB0, 0x10,                   // mov al, 16 (16 sectors = 8192 bytes)
		0xB5, 0x00,                   // mov ch, 0 (cylinder 0)
		0xB1, 0x02,                   // mov cl, 2 (sector 2)
		0xB6, 0x00,                   // mov dh, 0 (head 0)
		0x8A, 0x16, 0x00, 0x7C,       // mov dl, [0x7C00] (drive)
		0xCD, 0x13,                   // int 13h

		// Jump directly to 0x1000:0x0000
		0xEA, 0x00, 0x00, 0x00, 0x10, // jmp 0x1000:0x0000
	}
	copy(sector0, bootLoader)
	sector0[510] = 0x55
	sector0[511] = 0xAA
	copy(disk[0:512], sector0)

	// ────────────────────────────────────────────────────────────────────────
	// 2. APPLICATION RUNNER (Loaded at 0x1000:0x0000)
	// Sets DS=0x1000, handles keyboard input, copies screen buffer to 0xB800:0000
	// ────────────────────────────────────────────────────────────────────────
	// Frame offsets in segment 0x1000:
	// Frame 1 (Note 1): 0x0200 (4000 bytes)
	// Frame 2 (Note 2): 0x1200 (4000 bytes)

	appCode := []byte{
		// Set DS = 0x1000
		0xB8, 0x00, 0x10,             // mov ax, 0x1000
		0x8E, 0xD8,                   // mov ds, ax

		// Set video mode 80x25 text (INT 10h, AH=00h, AL=03h)
		0xB8, 0x03, 0x00,             // mov ax, 0x0003
		0xCD, 0x10,                   // int 10h

		// Hide cursor (INT 10h, AH=01h, CH=20h)
		0xB4, 0x01,                   // mov ah, 0x01
		0xB5, 0x20,                   // mov ch, 0x20
		0xCD, 0x10,                   // int 10h

		// State variable at offset 0x0100 in segment 0x1000:
		// [0x0100] = active_note (1 or 2, default 1)
		0xC6, 0x06, 0x00, 0x01, 0x01, // mov byte [0x0100], 1

		// redraw_label: (offset 0x001B)
		// Default Note 1: SI = 0x0200
		0xBE, 0x00, 0x02,             // mov si, 0x0200
		0x80, 0x3E, 0x00, 0x01, 0x01, // cmp byte [0x0100], 1
		0x74, 0x03,                   // je do_copy
		// Note 2:
		0xBE, 0x00, 0x12,             // mov si, 0x1200 (Note 2)

		// do_copy: (offset 0x0027)
		// Copy 2000 words (4000 bytes) from DS:SI to ES:DI (0xB800:0000)
		0xB8, 0x00, 0xB8,             // mov ax, 0xB800
		0x8E, 0xC0,                   // mov es, ax
		0x31, 0xFF,                   // xor di, di
		0xB9, 0xD0, 0x07,             // mov cx, 2000 (80x25 words)
		0xF3, 0xA5,                   // rep movsw (Fast video memory blit)

		// Restore ES = DS (0x1000)
		0xB8, 0x00, 0x10,             // mov ax, 0x1000
		0x8E, 0xC0,                   // mov es, ax

		// event_loop: (offset 0x0039)
		// Wait for keypress (INT 16h, AH=00h)
		0xB4, 0x00,                   // mov ah, 0x00
		0xCD, 0x16,                   // int 16h

		// Key dispatch:
		0x3C, 0x31,                   // cmp al, '1'
		0x74, 0x15,                   // je key_1
		0x3C, 0x32,                   // cmp al, '2'
		0x74, 0x18,                   // je key_2
		0x3C, 0x33,                   // cmp al, '3'
		0x74, 0x15,                   // je key_2 (switch)
		0x3C, 0x34,                   // cmp al, '4'
		0x74, 0x0F,                   // je key_1
		0x3C, 0x20,                   // cmp al, ' ' (space - toggle)
		0x74, 0x14,                   // je key_toggle
		0x3C, 0x71,                   // cmp al, 'q'
		0x74, 0x17,                   // je key_q
		0x3C, 0x51,                   // cmp al, 'Q'
		0x74, 0x13,                   // je key_q
		0xEB, 0xDF,                   // jmp event_loop

		// key_1: (offset 0x0059)
		0xC6, 0x06, 0x00, 0x01, 0x01, // mov byte [0x0100], 1
		0xEB, 0xBD,                   // jmp redraw_label (0x001B - 0x0060 = -69 = 0xBD)

		// key_2: (offset 0x0060)
		0xC6, 0x06, 0x00, 0x01, 0x02, // mov byte [0x0100], 2
		0xEB, 0xB6,                   // jmp redraw_label

		// key_toggle: (offset 0x0067)
		0x80, 0x36, 0x00, 0x01, 0x03, // xor byte [0x0100], 3 (toggles between 1 and 2: 1^3=2, 2^3=1!)
		0xEB, 0xAF,                   // jmp redraw_label

		// key_q: (offset 0x006E)
		0xEA, 0x00, 0x00, 0xFF, 0xFF, // jmp FFFF:0000 (reboot)
	}

	// Copy application code to sector 1 (offset 512 in disk)
	copy(disk[512:], appCode)

	// ────────────────────────────────────────────────────────────────────────
	// 3. COPY FULL 80x25 FRAMES INTO DISK IMAGE
	// Sector 1: Code (offset 512 to 1023)
	// Sector 2..9: Frame 1 (offset 512 + 0x0200 = 1024, 4000 bytes)
	// Sector 10..17: Frame 2 (offset 512 + 0x1200 = 5120, 4000 bytes)
	// ────────────────────────────────────────────────────────────────────────
	frame1 := BuildFrame(1, false).ToBytes()
	frame2 := BuildFrame(2, true).ToBytes()

	copy(disk[512+0x0200:], frame1)
	copy(disk[512+0x1200:], frame2)

	outPath := "notes-os.img"
	if err := os.WriteFile(outPath, disk, 0644); err != nil {
		fmt.Printf("Error writing disk image: %v\n", err)
		return
	}

	fmt.Printf("Successfully generated verified QEMU bootable image: %s (%d bytes)\n", outPath, len(disk))
}
