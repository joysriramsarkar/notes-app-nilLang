# 📝 নীলাং নোটবুক (Nilang Notes App)

সম্পূর্ণ **নীলাং (NilLang)** ভাষা ও **আলাপ (Alap)** ফ্রেমওয়ার্কে তৈরি একটি আধুনিক, দ্রুতগতির ও সুন্দর নোটবুক অ্যাপ্লিকেশন। এটি উইন্ডোজ, লিনাক্স, অনুরণ ওএস এবং **QEMU ভার্চুয়াল মেশিন এমুলেটর**-এ সম্পূর্ণ কার্যকর।

---

## ✨ প্রধান বৈশিষ্ট্যসমূহ (Features)

1. **📑 সমৃদ্ধ নোট ম্যানেজমেন্ট (Rich Note Management)**:
   - নতুন নোট তৈরি, এডিট, পিন (Pin), ফেভারিট (Favorite) এবং ট্র্যাশে ডিলিট/পুনরুদ্ধার।
   - ৬টি সুন্দর কালার ব্যাজ (Indigo, Emerald, Amber, Rose, Violet, Cyan) এবং প্রায়োরিটি লেভেল।
2. **✅ ইন্টারেক্টিভ চেকলিস্ট ও টু-ডু (Interactive Checklist)**:
   - নোটে লাইভ চেকলিস্ট তৈরি ও চেকমার্ক টগল।
   - রিয়েল-টাইম অগ্রগতি বার (Progress Bar) ও শতাংশ গণনা।
3. **📁 স্মার্ট ক্যাটাগরি ও ফোল্ডার সিস্টেম**:
   - `all` (সকল নোট), `important` (জরুরি), `work` (কাজ), `personal` (ব্যক্তিগত), `ideas` (আইডিয়া), `study` (শিক্ষা), `todo` (টাস্ক), `archive` (আর্কাইভ), `trash` (ট্র্যাশ)।
4. **🔍 তাৎক্ষণিক সার্চ ও ট্যাগ ক্লাউড (Instant Search & Tags)**:
   - টাইটেল, বডি ও ট্যাগ দিয়ে রিয়েল-টাইম ফিল্টারিং (কীবোর্ড শর্টকাট `Ctrl + K`)।
5. **👁️ মার্কডাউন লাইভ প্রিভিউ ও স্প্লিট মোড**:
   - মার্কডাউন হেডিংস, বোল্ড/ইটালিক, কোড ব্লক, লিস্ট এবং এক্সসার্পট লাইভ ভিউ।
   - শব্দ সংখ্যা (Word count), অক্ষর সংখ্যা ও পড়ার আনুমানিক সময় গণনা।
6. **📊 অ্যানালিটিক্স ড্যাশবোর্ড**:
   - সক্রিয় নোট, পিন্ড নোট, সর্বমোট শব্দ সংখ্যা এবং সম্পন্ন টাস্কের পরিসংখ্যান।
7. **🎨 আধুনিক থিম সাপোর্ট**:
   - Deep Night (Dark), Cyberpunk (Cyan), Emerald Forest (Green), Bangla Heritage (Red-Green), Pure Clean (Light)।
8. **🖥️ QEMU এমুলেটর ইন্টিগ্রেশন**:
   - উইন্ডোজ ল্যাপটপে `qemu-system-x86_64` দিয়ে অনুরণ ভার্চুয়াল কার্নেল ও নীলাং নোটবুক ওএস বুট টেস্ট।

---

## 📂 প্রজেক্টের গঠন (Project Structure)

```text
notes-app-nilLang/
├── src/
│   ├── main.nil              # নীলাং টার্মিনাল TUI ও মাস্টার এন্ট্রি পয়েন্ট
│   ├── notes_app.nil         # Alap ডিক্লারেটিভ UI কম্পোনেন্ট (@Component)
│   ├── domain/
│   │   ├── types.nil         # টাইপ, ক্যাটাগরি ও কালার ডেফিনিশন
│   │   ├── notes.nil         # NotesEngine (CRUD, সার্চ, ফিল্টারিং, সর্টিং)
│   │   ├── checklist.nil     # চেকলিস্ট ও প্রোগ্রেস ইঞ্জিন
│   │   ├── markdown.nil      # মার্কডাউন অ্যানালাইসিস ও টেক্সট মেট্রিক্স
│   │   ├── theme.nil         # থিম টোকেন ও ANSI কালার প্যালেট
│   │   └── storage.nil       # ফাইল ও JSON পারসিস্টেন্স
│   └── ui/
│       ├── Header.nil        # হেডার ও সার্চ বার কম্পোনেন্ট
│       ├── Sidebar.nil       # ফোল্ডার ও সাইডবার নেভিগেশন
│       ├── NoteList.nil      # নোট কার্ড গ্রিড/লিস্ট ভিউ
│       ├── NoteEditor.nil    # বিস্তারিত এডিটর ও মেটাডাটা ভিউ
│       └── StatsModal.nil    # পরিসংখ্যান ড্যাশবোর্ড মোডাল
├── qemu/
│   ├── qemu-run.ps1          # QEMU ভার্চুয়াল মেশিন লঞ্চার (PowerShell)
│   ├── run-qemu.bat          # ডাবল-ক্লিক উইন্ডোজ লঞ্চার
│   ├── qemu-test.ps1         # স্বয়ংক্রিয় বুট ভেরিফিকেশন টেস্ট
│   ├── make_boot_image.go    # x86 বুটেবল ওএস ইমেজ জেনারেটর
│   └── notes-os.img          # বুটেবল ভার্চুয়াল ওএস ইমেজ (1.44MB)
├── public/
│   ├── index.html            # আধুনিক গ্লাস-মরফিজম ওয়েব UI
│   ├── styles.css            # আল্ট্রা-মডার্ন ডার্ক ডিজাইন সিস্টেম
│   └── app.js                # রিঅ্যাক্টিভ ক্লায়েন্ট লজিক ও API সিঙ্ক
├── tests/
│   └── notes_test.go         # Go টেস্ট সুইট (CRUD, QEMU, Engine)
├── main.go                   # লোকাল ডেস্কটপ ও API সার্ভার
├── db.go                     # ডেটাবেস ও ফাইল সিঙ্ক ইঞ্জিন
├── ui_renderer.go            # ASCII ও টার্মিনাল রেন্ডারার
├── nil.json                  # নীলাং প্যাকেজ কনফিগারেশন
├── launch-desktop.bat        # ডেস্কটপ ওয়েব অ্যাপ লঞ্চার
└── run-nilang.bat            # টার্মিনাল নীলাং মোড লঞ্চার
```

---

## 🚀 কীভাবে চালাবেন (How to Run)

### ১. QEMU এমুলেটরে ভার্চুয়াল ওএস মোডে টেস্ট করা (Windows QEMU)

পাওয়ারশেলে চালান:
```powershell
cd "c:\Users\joysr\Documents\programming language\notes-app-nilLang\qemu"
.\qemu-run.ps1
```
অথবা সরাসরি ডাবল-ক্লিক করুন: `qemu\run-qemu.bat`

স্বয়ংক্রিয় হেডলেস ভেরিফিকেশন টেস্ট:
```powershell
.\qemu-test.ps1
```

### ২. নীলাং টার্মিনাল TUI মোডে চালানো

```powershell
cd "c:\Users\joysr\Documents\programming language\notes-app-nilLang"
..\nilLang\bin\nil.exe run src/main.nil
```
অথবা ডাবল-ক্লিক করুন: `run-nilang.bat`

**টার্মিনাল কীবোর্ড শর্টকাট:**
- `1-9` : নোট নির্বাচন
- `n` : নতুন নোট তৈরি
- `e` : নির্বাচিত নোট এডিট
- `d` : নোট ট্র্যাশে পাঠানো / ডিলিট
- `p` : পিন টগল
- `f` : ফেভারিট টগল
- `t` : চেকলিস্ট টাস্ক টিক/আনটিক
- `+` : নতুন টাস্ক যোগ
- `s` : সার্চ
- `c` : ক্যাটাগরি ফিল্টার
- `v` : পরিসংখ্যান ড্যাশবোর্ড
- `x` : ব্যাকআপ এক্সপোর্ট
- `q` : বের হওয়া

### ৩. 📱 নীল ওএস (NilOS) রিয়ালিস্টিক টাচ এমুলেটরে চালানো

মোবাইল ও ট্যাবলেট ওএস-এর মতো সম্পূর্ণ টাচস্ক্রিন ডিভাইসে নোটবুক চালাতে:
- সরাসরি ডাবল-ক্লিক করুন: [launch-device-emulator.bat](file:///c:/Users/joysr/Documents/programming%20language/notes-app-nilLang/launch-device-emulator.bat)
- অথবা ব্রাউজারে খুলুন: **http://localhost:8085/device.html**

**টাচ ডিভাইস ফিচারসমূহ:**
- 📱 টাইটানিয়াম বেজেল ও ডায়নামিক আইল্যান্ড সহ রিয়ালিস্টিক স্মার্টফোন বডি।
- 👆 প্রতিটি নোট কার্ড ও চেকলিস্টে টাচ/ক্লিক করে তাৎক্ষণিক ওপেন ও এডিটিং।
- ⌨️ অন-স্ক্রিন বাংলা ভার্চুয়াল টাচ কীবোর্ড (সাজেশন ও স্পেস/এন্টার সহ)।
- ➕ গোল ফ্লোটিং অ্যাকশন বাটন (FAB) দিয়ে এক ট্যাপে নতুন নোট।
- 🔒 সোয়াইপ-টু-আনলক এবং বটম জেসচার হোম বার।
- 🎛️ ফিজিক্যাল পাওয়ার ও ভলিউম বাটন সিমুলেশন।

### ৪. আধুনিক ডেস্কটপ ও ওয়েব GUI মোডে চালানো

```powershell
cd "c:\Users\joysr\Documents\programming language\notes-app-nilLang"
go run .
```
অথবা ডাবল-ক্লিক করুন: [launch-desktop.bat](file:///c:/Users/joysr/Documents/programming%20language/notes-app-nilLang/launch-desktop.bat)  
এরপর ব্রাউজারে খুলুন: **http://localhost:8085**

---

## 🧪 টেস্ট সুইট চালানো

```powershell
cd "c:\Users\joysr\Documents\programming language\notes-app-nilLang"
go test -v ./tests
```

---

## 📜 লাইসেন্স

MIT License © 2026 Joy Sriram Sarkar

