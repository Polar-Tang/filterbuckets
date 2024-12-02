package main

// IMPORTS
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"pdf_greyhat_go/api"
	"pdf_greyhat_go/download"
)

const outputFile = "results.json"

func main() {
	// Initialize session and keywords
	sessionCookie := "__stripe_mid=af2965ba-4f1d-4bf0-a073-1aa7b3987d61ccd651; _gid=GA1.2.1907481890.1732978699; REMEMBERME=QXBwXEVudGl0eVxVc2VyOk1HUmhlV055WlhjPToxNzMzNjk5MTYwOjgzNmUyZjljZDIxYjViMzZlNDk0MTNiYjE5MzRjZWQyNzUwZTRjNzc2ZmEwZmVhNjkzNTNiMTk2NDg0NDJmZGM%3D; SFSESSID=traa1qmokdpvv9k89d7p8eh30i; _gat_gtag_UA_121795267_1=1; _ga_QGK3VF4QHK=GS1.1.1733108461.58.1.1733108463.0.0.0; _ga=GA1.1.1930062620.1731774259"
	keywords := []string{"Business"}
	extensions := []string{"pdf"}
	pdfKeywords := []string{"mercado libre", "mercado pago", "amazon", "3cx", "actbluetech", "agoric", "airlock", "airtable", "akamai", "alchemyplatform", "algolia", "aliexpress", "alsco", "aptible", "arkose", "artsy", "automattic", "avast", "barracuda", "basecamp", "bigbasket", "bitmex", "blogger", "boozt", "brave", "bugcrowd", "bugify", "bumble", "buzzfeed", "careerzooms", "chaturbate", "chime", "chromium", "clario", "cloudbees", "cloudflare", "clubhouse", "cobalt", "creditkarma", "crowdstrike", "crypto", "cryptocat", "csg-public", "databricks", "datastax", "deribit", "deriv", "deutschetelekom", "directv", "discourse", "doppler", "eero", "elastic", "elisa", "epicgames", "ethereum", "etoro", "etsy", "eufy", "eureka jp", "exodus", "expediagroup", "fetlife", "figma", "filezilla", "fireblocks", "fireblocks", "flexispyltd", "flickr", "floqast", "flutteruki", "freshworks", "fsecure", "gearbest", "ghostscript", "github", "gitlab", "glassdoor", "gocardless", "goodrx", "google", "grammarly", "highrise", "hiro", "hootsuite", "hostelworld", "hostfact", "hyatt", "hybridsaas", "hyperledger", "hypr", "irccloud", "kayak", "kfc", "khealth", "kiwicom", "koho", "kolesagroup", "krisp", "kyupcloud", "lastpass", "leather wallet", "localizejs", "lockkey", "logitech", "lyft", "lyst", "magic", "makemytrip", "malwarebytes", "mapbox", "marriott", "matomo", "meetic", "megaprivacy", "meraki", "mergify", "metamask", "miro", "mozilla", "mux", "oculus", "okcupid", "okg", "olacabs", "olark", "omise", "ovchipkaart", "ovh", "parseit", "paypal", "paysafecard", "phabricator", "pillarproject", "pixiv", "plaid", "playtika", "pleo", "plotly", "porsche", "portswigger", "qualcomm", "quora", "razorpay", "recorded future", "reddit", "redox", "releasewire", "remitly", "resourcespace", "ridewithvia", "rijksoverheid", "riot", "ripple", "riskio", "roblox", "roche", "rockstargames", "rootstocklabs", "schubergphilis", "scopely", "sec1337", "semrush", "shapeshift", "sheer", "shipt", "shopify", "sidefx", "simplybuilt", "six group", "slack", "smtp2go", "sofi", "sorare", "sourcegraph", "spotify", "starbucks", "statuspageio", "stripe", "superbet", "swiggy", "telegram", "temu", "tide", "tiktok", "torproject", "tripadvisor", "trivago", "truecaller", "tweakers", "uber", "udemy", "valve", "veridu", "vimeo", "visa", "wamba", "wazirx", "wealthsimple", "websecurify", "wellsfargo", "whmcs", "wisdomtree", "withpaper", "wordpress", "xiaomi", "xverse", "xvideos", "yahoo", "yammer", "yelp", "yoti", "youtube", "yuga", "labs", "zabbix", "zapier", "zomato", "zoom-private", "zte", "TeamViewer", "VFS", "Kiwai", "ATG", "DANA", "Deezer", "DRACOON", "Telenor", "Dailymotion", "Coindcx", "DataDome", "SNCF", "OVH", "Infomaniak", "ownCloud", "Alasco", "Withings", "Doctolib", "Qwant", "Okto", "Coindcx", "Toom", "Baumarkt", "Swapcard", "Spacelift", "ezviz", "Dassault", "BookBeat", "Comuto", "Legapass", "Sovereign", "Sogexia", "Ooredoo", "Cryptobox", "Ercom", "PowerDNS", "Open-Xchange", "Moneybox", "Maya", "MCN", "CycloneDX", "DataDome", "Dovecot", "Tencent", "MatrixALM", "Salt", "Thüringer", "Aufbaubank", "DINUM", "Sequoia", "ZECIBLE"}

	// Query files THROUGH THE API
	files, err := api.QueryFiles(sessionCookie, keywords, extensions)
	// generic error handling
	if err != nil {
		log.Fatalf("Failed to query files: %v", err)
	}

	// Create a semaphore for concurrent downloads
	var wg sync.WaitGroup

	// Initialize results
	results := make([]map[string]interface{}, 0)
	// RESULTS:
	/*
			{"Filename": "file1.pdf", "URL": "http://example.com/file1", "Matches": 10},
		    {"Filename": "file2.pdf", "URL": "http://example.com/file2", "Matches": 5},
	*/

	mutex := &sync.Mutex{}

	// Set the concurrency limit
	concurrencyLimit := 6
	// use semaphore var to set a maximum number of concurrent goroutines
	semaphore := make(chan struct{}, concurrencyLimit)

	// Creates a timer that triggers every 60 seconds.
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	// iterate just on the values of the files
	// go routine
	// HANDLE THE GO ROUTINE CONCURRENCY

	// ensure to goes periodicly saving it's making to avoid big lost if the process is interrupted
	go func() {
		// A channel that emits a signal every time the ticker fires.
		for range ticker.C {
			// save the file periodically
			mutex.Lock()
			saveResults(results)
			fmt.Printf("Result added: %+v\n", results) // Add this line for debugging
			// MUTEX write the file but priventing race conditions
			mutex.Unlock()
		}
	}()

	for _, fileInfo := range files {
		// increment the wait counter
		wg.Add(1)
		go func(file api.FileInfo) {
			// DECREMENT the wait routine when it's done
			defer wg.Done()
			// send an empty struct into the sempahore channel
			semaphore <- struct{}{} // Acquire a semaphore slot
			// semaphoro green!
			defer func() { <-semaphore }() // Release slot after processing

			// fmt.Printf("Found file: %s (URL: %s, Size: %d bytes)\n", file.Filename, file.URL, file.Size)

			result := download.ProcessFile(file, pdfKeywords) // redefine result
			// redefine the results with the function proces file
			if result != nil {
				// append the result (no overwrite)
				mutex.Lock()
				results = append(results, result)

				// MUTEX write the file but priventing race conditions
				mutex.Unlock()
			}
		}(fileInfo)
	}
	// The file info is a struct
	/* type FileInfo struct {
		URL      string
		Filename string
		Size     int
	} */
	// The results are saved as JSON in results.json, after the whole fucking process ends:
	saveResults(results)
	wg.Wait() // Wait for all goroutines to complete

}

func saveResults(results []map[string]interface{}) {
	fmt.Printf("Saving %d results...\n", len(results)) // Debug log

	file, err := os.Create(outputFile) // Create (or overwrite) results.json
	if err != nil {
		log.Printf("Failed to create output file: %v", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Add indentation for readability

	if err := encoder.Encode(results); err != nil {
		log.Printf("Failed to write results: %v", err)
	} else {
		fmt.Printf("Results saved to %s\n", outputFile)
	}
}
