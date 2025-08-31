package worker

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/tajri15/go-pulse-monitoring/internal/db"
)

// Checker adalah struct yang menampung dependensi worker
type Checker struct {
	store *db.Store
}

func NewChecker(store *db.Store) *Checker {
	return &Checker{store: store}
}

// Start memulai proses pemeriksaan berkala
func (c *Checker) Start() {
	log.Println("Starting health check worker...")
	// Ticker akan "berdetak" setiap satu menit
	ticker := time.NewTicker(1 * time.Minute)

	// Loop for-ever, menunggu detak dari ticker
	for range ticker.C {
		log.Println("Running health check cycle...")
		c.runChecks()
	}
}

func (c *Checker) runChecks() {
	ctx := context.Background()
	sites, err := c.store.GetAllSites(ctx)
	if err != nil {
		log.Printf("Error fetching sites: %v", err)
		return
	}

	// Implementasi Worker Pool
	// 1. Buat channel untuk pekerjaan (sites) dan hasil (health checks)
	jobs := make(chan db.Site, len(sites))
	results := make(chan db.CreateHealthCheckParams, len(sites))

	// 2. Jalankan beberapa worker goroutine (misalnya 5)
	numWorkers := 5
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// 3. Kirim semua pekerjaan ke channel jobs
	for _, site := range sites {
		jobs <- site
	}
	close(jobs) // Tutup channel jobs karena semua pekerjaan sudah dikirim

	// 4. Kumpulkan semua hasil dari channel results
	for a := 1; a <= len(sites); a++ {
		result := <-results
		_, err := c.store.CreateHealthCheck(ctx, result)
		if err != nil {
			log.Printf("Error saving health check result for site ID %d: %v", result.SiteID, err)
		} else {
			log.Printf("Successfully saved health check for site ID %d. Status UP: %t", result.SiteID, result.IsUp)
		}
	}
}

// worker adalah goroutine yang akan melakukan pekerjaan
func worker(id int, jobs <-chan db.Site, results chan<- db.CreateHealthCheckParams) {
	for site := range jobs {
		log.Printf("Worker %d started job for site %s", id, site.URL)
		
		startTime := time.Now()
		
		// Lakukan request HTTP dengan timeout 10 detik
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(site.URL)

		duration := time.Since(startTime).Milliseconds()

		result := db.CreateHealthCheckParams{
			SiteID:         site.ID,
			ResponseTimeMs: int(duration),
		}

		if err != nil {
			log.Printf("Worker %d failed to check site %s: %v", id, site.URL, err)
			result.IsUp = false
			result.StatusCode = 0 // Tidak ada status code jika gagal connect
		} else {
			result.IsUp = resp.StatusCode >= 200 && resp.StatusCode < 300 // Dianggap UP jika status 2xx
			result.StatusCode = resp.StatusCode
			resp.Body.Close()
		}

		results <- result
	}
}