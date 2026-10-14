package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type CheckItem struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type Note struct {
	ID         string      `json:"id"`
	Title      string      `json:"title"`
	Content    string      `json:"content"`
	Category   string      `json:"category"`
	Tags       []string    `json:"tags"`
	Color      string      `json:"color"`
	Priority   string      `json:"priority"`
	IsPinned   bool        `json:"isPinned"`
	IsFavorite bool        `json:"isFavorite"`
	IsArchived bool        `json:"isArchived"`
	IsTrashed  bool        `json:"isTrashed"`
	Checklist  []CheckItem `json:"checklist"`
	CreatedAt  int64       `json:"createdAt"`
	UpdatedAt  int64       `json:"updatedAt"`
}

type Store struct {
	mu       sync.RWMutex
	filePath string
	notes    []Note
}

func NewStore(filePath string) (*Store, error) {
	s := &Store{
		filePath: filePath,
		notes:    []Note{},
	}

	if err := s.load(); err != nil {
		s.seedDefaults()
		if err := s.save(); err != nil {
			return nil, fmt.Errorf("failed to save initial store: %w", err)
		}
	}
	return s, nil
}

func (s *Store) seedDefaults() {
	now := time.Now().Unix()
	s.notes = []Note{
		{
			ID:         "note-1",
			Title:      "স্বাগতম নীলাং নোটবুক-এ! (Nilang Notes Guide)",
			Content:    "নীলাং ভাষায় তৈরি আধুনিক, দ্রুতগতির এবং সুন্দর নোট অ্যাপ্লিকেশন।\n\nপ্রধান সুবিধাসমূহ:\n- সম্পূর্ণ বাংলায় ও ইংরেজিতে নোট সংরক্ষণ\n- লাইভ সার্চ ও ক্যাটাগরি ফিল্টারিং\n- প্রতিটি নোটে কালার ব্যাজ ও প্রায়োরিটি\n- সম্পূর্ণ ক্রস-প্ল্যাটফর্ম (Windows, Linux, Onuron, QEMU)",
			Category:   "important",
			Tags:       []string{"নীলাং", "গাইড", "নোটবুক", "অনুরণ"},
			Color:      "indigo",
			Priority:   "high",
			IsPinned:   true,
			IsFavorite: true,
			IsArchived: false,
			IsTrashed:  false,
			Checklist: []CheckItem{
				{ID: "chk-1", Text: "ক্যাটাগরি ও ট্যাগ দিয়ে নোট সাজানো", Completed: true},
				{ID: "chk-2", Text: "ইন্টারেক্টিভ চেকলিস্ট ও টাস্ক তৈরি", Completed: true},
				{ID: "chk-3", Text: "QEMU এমুলেটরে ভার্চুয়াল ওএস টেস্ট করা", Completed: true},
				{ID: "chk-4", Text: "অটো-সেভ ও ডেটা রিকভারি যাচাই করা", Completed: false},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:         "note-2",
			Title:      "দৈনিক কাজের পরিকল্পনা ও ডেডলাইন",
			Content:    "আজকের সর্বোচ্চ অগ্রাধিকারের কাজগুলো সম্পন্ন করতে হবে। সন্ধ্যে ৬টার মধ্যে বিল্ড টেস্ট সমাপ্ত করা দরকার।",
			Category:   "todo",
			Tags:       []string{"কাজ", "প্ল্যানিং", "টাস্ক"},
			Color:      "emerald",
			Priority:   "high",
			IsPinned:   true,
			IsFavorite: false,
			IsArchived: false,
			IsTrashed:  false,
			Checklist: []CheckItem{
				{ID: "chk-21", Text: "নীলাং কোডবেজ রিফ্যাক্টরিং", Completed: true},
				{ID: "chk-22", Text: "আলাপ UI ডিক্লারেটিভ কম্পোনেন্ট তৈরি", Completed: true},
				{ID: "chk-23", Text: "QEMU x86_64 কার্নেল বুট টেস্ট", Completed: true},
				{ID: "chk-24", Text: "রিসার্চ পেপার ড্রাফট রিভিশন", Completed: false},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:         "note-3",
			Title:      "অনুরণ ওএস: মাইক্রোকার্নেল আর্কিটেকচার আইডিয়া",
			Content:    "অনুরণ অপারেটিং সিস্টেমে আল্ট্রা-লাইটওয়েট প্রসেস আইসোলেশন ও সফটবাস (SoftBus) আর্কিটেকচার বাস্তবায়ন করা।\n\nলক্ষ্য:\n১. দ্রুততম রেসপন্স টাইম (< 2ms)\n২. সুরক্ষিত মেমোরি স্যান্ডবক্সিং\n৩. নীলাং বাইটকোড ডিরেক্ট এক্সিকিউশন",
			Category:   "ideas",
			Tags:       []string{"আর্কিটেকচার", "ওএস", "রিসার্চ"},
			Color:      "cyan",
			Priority:   "medium",
			IsPinned:   false,
			IsFavorite: true,
			IsArchived: false,
			IsTrashed:  false,
			Checklist: []CheckItem{
				{ID: "chk-31", Text: "মাইক্রো-কার্নেল ড্রাইভার ডিজাইন", Completed: true},
				{ID: "chk-32", Text: "জিরো-কপি IPC মেমোরি চ্যানেল", Completed: false},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:         "note-4",
			Title:      "অনুপ্রেরণাদায়ক উক্তি ও জ্ঞান",
			Content:    "\"জ্ঞানের চেয়ে কল্পনাশক্তি অধিক গুরুত্বপূর্ণ, কারণ জ্ঞান সীমাবদ্ধ, কিন্তু কল্পনা জগতকে আলিঙ্গন করে।\" — অ্যালবার্ট আইনস্টাইন\n\nবইয়ের তালিকা:\n১. Clean Code (Robert C. Martin)\n২. Designing Data-Intensive Applications",
			Category:   "personal",
			Tags:       []string{"উক্তি", "বই", "শিক্ষা"},
			Color:      "violet",
			Priority:   "low",
			IsPinned:   false,
			IsFavorite: false,
			IsArchived: false,
			IsTrashed:  false,
			Checklist:  []CheckItem{},
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.notes)
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.notes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *Store) GetAll() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]Note, len(s.notes))
	copy(res, s.notes)
	return res
}

func (s *Store) Get(id string) (Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, n := range s.notes {
		if n.ID == id {
			return n, true
		}
	}
	return Note{}, false
}

func (s *Store) Create(note Note) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if note.ID == "" {
		note.ID = fmt.Sprintf("note-%d-%d", len(s.notes)+1, time.Now().Unix())
	}
	now := time.Now().Unix()
	note.CreatedAt = now
	note.UpdatedAt = now
	if note.Tags == nil {
		note.Tags = []string{}
	}
	if note.Checklist == nil {
		note.Checklist = []CheckItem{}
	}

	s.notes = append(s.notes, note)
	if err := s.save(); err != nil {
		return Note{}, err
	}
	return note, nil
}

func (s *Store) Update(id string, update Note) (Note, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.notes {
		if n.ID == id {
			update.ID = id
			update.CreatedAt = n.CreatedAt
			update.UpdatedAt = time.Now().Unix()
			if update.Tags == nil {
				update.Tags = []string{}
			}
			if update.Checklist == nil {
				update.Checklist = []CheckItem{}
			}
			s.notes[i] = update
			if err := s.save(); err != nil {
				return Note{}, err
			}
			return update, nil
		}
	}
	return Note{}, fmt.Errorf("note with id %s not found", id)
}

func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.notes {
		if n.ID == id {
			if n.IsTrashed {
				s.notes = append(s.notes[:i], s.notes[i+1:]...)
			} else {
				n.IsTrashed = true
				n.IsPinned = false
				n.UpdatedAt = time.Now().Unix()
				s.notes[i] = n
			}
			return s.save()
		}
	}
	return fmt.Errorf("note with id %s not found", id)
}

// ────────────────────────────────────────────────────────────────────
// POS Models & Store
// ────────────────────────────────────────────────────────────────────
type POSProduct struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	NameEn     string  `json:"nameEn"`
	SKU        string  `json:"sku"`
	Barcode    string  `json:"barcode"`
	Price      float64 `json:"price"`
	WAC        float64 `json:"wac"`
	Stock      int     `json:"stock"`
	Category   string  `json:"category"`
	Icon       string  `json:"icon"`
}

type POSCartItem struct {
	Product  POSProduct `json:"product"`
	Quantity int        `json:"quantity"`
}

type POSCheckoutRequest struct {
	Items         []POSCartItem `json:"items"`
	CustomerName  string        `json:"customerName"`
	CustomerPhone string        `json:"customerPhone"`
	DiscountPct   float64       `json:"discountPct"`
	PaymentMethod string        `json:"paymentMethod"` // Cash, Card, UPI/bKash, Due
}

type POSReceipt struct {
	InvoiceNo     string        `json:"invoiceNo"`
	Date          string        `json:"date"`
	Time          string        `json:"time"`
	CustomerName  string        `json:"customerName"`
	CustomerPhone string        `json:"customerPhone"`
	Items         []POSCartItem `json:"items"`
	Subtotal      float64       `json:"subtotal"`
	Discount      float64       `json:"discount"`
	VAT           float64       `json:"vat"`
	Total         float64       `json:"total"`
	PaymentMethod string        `json:"paymentMethod"`
	PrintedAt     int64         `json:"printedAt"`
	Cashier       string        `json:"cashier"`
}

func GetDefaultPOSCatalog() []POSProduct {
	return []POSProduct{
		{
			ID:       "p-01",
			Name:     "মিনিকেট চাল (৫০ কেজি)",
			NameEn:   "Miniket Rice (50kg)",
			SKU:      "RICE-MIN-50",
			Barcode:  "8901030012345",
			Price:    3200.0,
			WAC:      2850.0,
			Stock:    40,
			Category: "Groceries",
			Icon:     "🌾",
		},
		{
			ID:       "p-02",
			Name:     "তীর সয়াবিন তেল (৫ লিটার)",
			NameEn:   "Teer Soybean Oil (5L)",
			SKU:      "OIL-TEER-5L",
			Barcode:  "8901030098765",
			Price:    850.0,
			WAC:      780.0,
			Stock:    65,
			Category: "Oil",
			Icon:     "🛢️",
		},
		{
			ID:       "p-03",
			Name:     "ফ্রেশ আটা (২ কেজি)",
			NameEn:   "Fresh Atta (2kg)",
			SKU:      "ATTA-FRSH-2K",
			Barcode:  "8901030055555",
			Price:    125.0,
			WAC:      105.0,
			Stock:    120,
			Category: "Flour",
			Icon:     "🍞",
		},
		{
			ID:       "p-04",
			Name:     "মোল্লা লবণ (১ কেজি)",
			NameEn:   "Molla Salt (1kg)",
			SKU:      "SALT-MLLA-1K",
			Barcode:  "8901030077778",
			Price:    42.0,
			WAC:      31.0,
			Stock:    200,
			Category: "Salt",
			Icon:     "🧂",
		},
		{
			ID:       "p-05",
			Name:     "চিনি (১ কেজি)",
			NameEn:   "Sugar (1kg)",
			SKU:      "SUGAR-1K",
			Barcode:  "8901030066666",
			Price:    115.0,
			WAC:      98.0,
			Stock:    88,
			Category: "Sugar",
			Icon:     "🍬",
		},
		{
			ID:       "p-06",
			Name:     "নীলাং প্রিমিয়াম কফি ব্লেন্ড (২০০ গ্রাম)",
			NameEn:   "Nilang Coffee Blend (200g)",
			SKU:      "COFFEE-NIL-200",
			Barcode:  "8901030088888",
			Price:    350.0,
			WAC:      280.0,
			Stock:    50,
			Category: "Beverage",
			Icon:     "☕",
		},
	}
}

// ────────────────────────────────────────────────────────────────────
// Music Models
// ────────────────────────────────────────────────────────────────────
type MusicTrack struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration string `json:"duration"`
	Genre    string `json:"genre"`
	Cover    string `json:"cover"`
	AudioURL string `json:"audioUrl"`
}

func GetDefaultMusicTracks() []MusicTrack {
	return []MusicTrack{
		{
			ID:       "trk-1",
			Title:    "আমার সোনার বাংলা (Lofi Instrument)",
			Artist:   "নীলাং সিন্থেসাইজার টিম",
			Album:    "Onuron Ambient Sessions Vol. 1",
			Duration: "3:42",
			Genre:    "Lofi / Ambient",
			Cover:    "🎵",
			AudioURL: "",
		},
		{
			ID:       "trk-2",
			Title:    "গ্রাম ছাড়া ওই রাঙা মাটির পথ",
			Artist:   "অনুরণ সাউন্ড ল্যাব",
			Album:    "Bangla Heritage Modernized",
			Duration: "4:15",
			Genre:    "Folk Acoustic",
			Cover:    "🌾",
			AudioURL: "",
		},
		{
			ID:       "trk-3",
			Title:    "নীলাং সাইবারপাংক অ্যান্থেম",
			Artist:   "NilCore 8086",
			Album:    "Bangla 2040 Cyber OST",
			Duration: "2:58",
			Genre:    "Cyber Synthwave",
			Cover:    "⚡",
			AudioURL: "",
		},
		{
			ID:       "trk-4",
			Title:    "মেঘের কোলে রোদ হেসেছে",
			Artist:   "আলাপ অর্কেস্ট্রা",
			Album:    "Monsoon Serenade",
			Duration: "3:30",
			Genre:    "Classical Fusion",
			Cover:    "🌧️",
			AudioURL: "",
		},
	}
}

