package worker

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/tajri15/go-pulse-monitoring/internal/db"
	"github.com/tajri15/go-pulse-monitoring/internal/ws"
)

// Checker sekarang juga memegang referensi ke Hub
type Checker struct {
	store *db.Store
	hub   *ws.Hub
}

// NewChecker diubah untuk menerima Hub
func NewChecker(store *db.Store, hub *ws.Hub) *Checker {
	return &Checker{
		store: store,
		hub:   hub,
	}
}

// Start tetap sama
func (c *Checker) Start() {
	log.Println("Starting health check worker...")
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		log.Println("Running health check cycle...")
		c.runChecks()
	}
}

// Tipe data untuk pesan update WebSocket
type WsUpdateMessage struct {
	SiteID         int64  `json:"site_id"`
	IsUp           bool   `json:"is_up"`
	ResponseTimeMs int    `json:"response_time_ms"`
	StatusCode     int    `json:"status_code"`
	CheckedAt      time.Time `json:"checked_at"`
}

func (c *Checker) runChecks() {
	ctx := context.Background()
	sites, err := c.store.GetAllSites(ctx)
	if err != nil {
		log.Printf("Error fetching sites: %v", err)
		return
	}

	if len(sites) == 0 {
		log.Println("No sites to check.")
		return
	}

	jobs := make(chan db.Site, len(sites))
	results := make(chan db.HealthCheck, len(sites)) // Hasilnya adalah HealthCheck

	numWorkers := 5
	for w := 1; w <= numWorkers; w++ {
		// worker sekarang menghasilkan db.HealthCheck
		go worker(w, jobs, results)
	}

	for _, site := range sites {
		jobs <- site
	}
	close(jobs)

	for a := 1; a <= len(sites); a++ {
		result := <-results
		
		// 1. Simpan hasil ke database
		savedCheck, err := c.store.CreateHealthCheck(ctx, db.CreateHealthCheckParams{
			SiteID:         result.SiteID,
			StatusCode:     result.StatusCode,
			ResponseTimeMs: result.ResponseTimeMs,
			IsUp:           result.IsUp,
		})
		if err != nil {
			log.Printf("Error saving health check result for site ID %d: %v", result.SiteID, err)
			continue
		}
		
		log.Printf("Successfully saved health check for site ID %d. Status UP: %t", result.SiteID, result.IsUp)

		// 2. Kirim pembaruan melalui WebSocket
		// Cari userID dari site yang sesuai
		var targetUserID int64
		for _, site := range sites {
			if site.ID == result.SiteID {
				targetUserID = site.UserID
				break
			}
		}

		if targetUserID != 0 {
			// Buat pesan JSON
			updateMsg := WsUpdateMessage{
				SiteID:         savedCheck.SiteID,
				IsUp:           savedCheck.IsUp,
				ResponseTimeMs: savedCheck.ResponseTimeMs,
				StatusCode:     savedCheck.StatusCode,
				CheckedAt:      savedCheck.CheckedAt,
			}
			jsonMsg, _ := json.Marshal(updateMsg)

			// Kirim ke Hub
			c.hub.Send(targetUserID, jsonMsg)
		}
	}
}

// worker diubah untuk mengembalikan struct db.HealthCheck
func worker(id int, jobs <-chan db.Site, results chan<- db.HealthCheck) {
	for site := range jobs {
		log.Printf("Worker %d started job for site %s", id, site.URL)
		
		startTime := time.Now()
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(site.URL)
		duration := time.Since(startTime).Milliseconds()

		result := db.HealthCheck{
			SiteID:         site.ID,
			ResponseTimeMs: int(duration),
		}

		if err != nil {
			log.Printf("Worker %d failed to check site %s: %v", id, site.URL, err)
			result.IsUp = false
			result.StatusCode = 0
		} else {
			result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}
		results <- result
	}
}