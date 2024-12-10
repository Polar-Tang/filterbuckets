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

func main() {
	// Initialize session and keywords
	sessionCookie := "01931a3ff4929fa0e8d8c93ba9dac24c"
	keywords := []string{"private", "restricted", "classified", "internal-use", "for_official_use_only", "docs", "document_file", "file", "note", "writeup", "write-up", "file_record", "draft", "text_file", "paper", "form", "worksheet", "sensitive", "secret", "non_public", "classified", "proprietary", "privileged", "summary", "log", "statement", "record", "analysis", "assessment", "review", "findings", "private_file", "classified_doc", "restricted_document", "internal_note", "internal_record", "tech_docs", "technical_manual", "spec_doc", "engineering_doc", "system_guide", "implementation_guide", "guide", "instruction", "user_guide", "reference_manual", "handbook", "how_to", "procedure_guide", "audit_log", "compliance_audit", "inspection_report", "review_summary", "audit_summary", "inspection", "assessment", "evaluation", "check", "compliance_check", "review_process", "learning_material", "training_manual", "educational_guide", "training_resource", "instruction_guide", "requirements", "spec_doc", "system_specs", "tech_requirements", "implementation_details", "changelog", "version_log", "update_notes", "patch_notes", "deployment_notes", "revision_log", "notice", "internal_message", "bulletin", "announcement", "reminder", "brief", "conference", "call_notes", "discussion", "session", "standup", "minutes_of_meeting", "plan", "roadmap", "approach", "tactics", "blueprint", "agenda", "promotion_plan", "campaign_strategy", "marketing_strategy", "advertising_plan", "business_promotion", "case_report", "project_analysis", "success_story", "research_report", "example_study", "regulatory_report", "policy_check", "compliance_summary", "standards_report", "audit_summary", "q1_report", "q2_report", "q3_report", "q4_report", "financial_summary", "quarterly_review", "performance_report"}
	extensions := map[string][]string{
		"json": {"(?i)(password|passwd|pwd|pass|secret)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(api[_-]?key|apikey|api[_-]?secret|secret[_-]?key|key|access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(token|access[_-]?token|auth[_-]?token|jwt|bearer)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(db[_-]?password|database[_-]?password|db[_-]?user|database[_-]?user|db[_-]?host|database[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"-----BEGIN(.*)PRIVATE KEY-----[\\s\\S]*?-----END(.*)PRIVATE KEY-----",
			"(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key|aws[_-]?secret[_-]?access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9+/=]+['\"\\s]",
			"(?i)(firebase[_-]?api[_-]?key|firebase[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(google[_-]?api[_-]?key|google[_-]?client[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(smtp[_-]?password|smtp[_-]?user|smtp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(ftp[_-]?password|ftp[_-]?user|ftp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(oauth[_-]?token|oauth[_-]?secret|oauth2)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(eyJ[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+)",
			"(?i)(azure[_-]?key|azure[_-]?secret|azure[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(stripe[_-]?key|stripe[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(paypal[_-]?key|paypal[_-]?secret|paypal[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"https://hooks\\.slack\\.com/services/[A-Za-z0-9/-]+",
			"(?i)(github[_-]?token|github[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(gitlab[_-]?token|gitlab[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(webhook[_-]?url|webhook[_-]?token|webhook[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-/]+['\"\\s]",
			"s3://[A-Za-z0-9._\\-/]+",
			"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			"(?i)(https?://)[^\\s]+:[^\\s]+@[^\\s]+",
			"(?i)(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519) [A-Za-z0-9+/=]+",
			"Adidas",
			"Algemeen Dagblad",
			"Allegro",
			"Amazon",
			"Apple",
			"Axel-Springer",
			"Azena",
			"Bank-of-America",
			"BMW",
			"Bpost",
			"Buhler",
			"CM.com",
			"Canada-Post",
			"Capital.com",
			"Citigroup",
			"Cloudways-by-DigitalOcean",
			"Coca-Cola",
			"Cross-Border-Fines",
			"Cyber-Security-Coalition",
			"DPG-Media",
			"De-Lijn",
			"De-Morgen",
			"De-Volkskrant",
			"Delen-Private Bank",
			"Dell",
			"Digitaal-Vlaanderen",
			"DigitalOcean",
			"Disney",
			"Discovery",
			"Donorbox",
			"E-Gor",
			"EURid",
			"Fing",
			"HRS Group",
			"Henkel",
			"Here Technologies",
			"Het Laatste Nieuws",
			"Het Parool",
			"Humo",
			"Kinepolis Group",
			"Lansweeper",
			"Libelle",
			"Mobile Vikings",
			"Moralis",
			"Nestlé",
			"Nexuzhealth",
			"Nexuzhealth Web PACS",
			"OVO",
			"PDQ bug bounty program",
			"PeopleCert",
			"Personio",
			"Port of Antwerp-Bruges",
			"Purolator",
			"RGF BE",
			"RIPE NCC",
			"Randstad",
			"Red Bull",
			"Revolut",
			"SimScale",
			"Sixt",
			"Social Deal",
			"Soundtrack Your Brand",
			"Sqills",
			"Stravito",
			"Suivo bug bounty",
			"Sustainable",
			"Telenet",
			"Tempo-Team",
			"Tomorrowland",
			"Torfs",
			"Trouw",
			"TrueLayer",
			"Twago",
			"Tweakers",
			"UZ-Leuven",
			"Ubisoft",
			"Universalg",
			"Unilever",
			"Visa",
			"Vlerick-Business-School",
			"VMT GO",
			"Voi-Scooters",
			"Volkswagen",
			"VRT",
			"WP-Engine",
			"Yacht",
			"Yahoo",
			"e-tracker"},
		"xml": {"Adidas",
			"Algemeen Dagblad",
			"Allegro",
			"Amazon",
			"Apple",
			"Axel-Springer",
			"Azena",
			"Bank-of-America",
			"BMW",
			"Bpost",
			"Buhler",
			"CM.com",
			"Canada-Post",
			"Capital.com",
			"Citigroup",
			"Cloudways-by-DigitalOcean",
			"Coca-Cola",
			"Cross-Border-Fines",
			"Cyber-Security-Coalition",
			"DPG-Media",
			"De-Lijn",
			"De-Morgen",
			"De-Volkskrant",
			"Delen-Private Bank",
			"Dell",
			"Digitaal-Vlaanderen",
			"DigitalOcean",
			"Disney",
			"Discovery",
			"UZ Leuven",
			"Ubisoft",
			"VRT",
			"VTM GO",
			"Venly",
			"Vlerick Business School",
			"Voi Scooters",
			"WP Engine",
			"Yacht",
			"Yahoo",
			"e-tracker",
			"token",
			"eHealth Hub VZN KUL"},
		"pdf": {
			"Adidas",
			"Algemeen Dagblad",
			"Allegro",
			"Amazon",
			"Apple",
			"Axel-Springer",
			"Azena",
			"Bank-of-America",
			"BMW",
			"Bpost",
			"Buhler",
			"CM.com",
			"Canada-Post",
			"Capital.com",
			"Citigroup",
			"Cloudways-by-DigitalOcean",
			"Coca-Cola",
			"Cross-Border-Fines",
			"Cyber-Security-Coalition",
			"DPG-Media",
			"De-Lijn",
			"De-Morgen",
			"De-Volkskrant",
			"Delen-Private Bank",
			"Dell",
			"Digitaal-Vlaanderen",
			"DigitalOcean",
			"Disney",
			"Discovery",
			"Donorbox",
			"E-Gor",
			"EURid",
			"Fing",
			"Ford",
			"Henkel",
			"Here-Technologies",
			"Het-Laatste Nieuws",
			"Het-Parool",
			"Honda",
			"HSBC",
			"HRS-Group",
			"Humo",
			"Hyundai",
			"IBM",
			"Intel",
			"Kinepolis-Group",
			"Lansweeper",
			"Libelle",
			"Mastercard",
			"McDonald",
			"Mercedes-Benz",
			"Meta",
			"Mobile-Vikings",
			"Moralis",
			"Nestle",
			"Nextuzhealth",
			"Nexuzhealth-Web-PACS",
			"Nike",
			"OVO",
			"PDQ",
			"PepsiCo",
			"PeopleCert",
			"Personio",
			"Porsche",
			"Port-of-Antwerp-Bruges",
			"Purolator",
			"RGF-BE",
			"RIPE-NCC",
			"Randstad",
			"Red-Bull",
			"Revolut",
			"Samsung",
			"Shell",
			"Siemens",
			"SimScale",
			"Sixt",
			"Social-Deal",
			"Sony",
			"Soundtrack-Your-Brand",
			"Spotify",
			"Sqills",
			"Stravito",
			"Suivo",
			"Sustainable",
			"Telenet",
			"Tempo-Team",
			"Tesla",
			"Tomorrowland",
			"Torfs",
			"Toyota",
			"Trouw",
			"TrueLayer",
			"Twago",
			"Tweakers",
			"UZ-Leuven",
			"Ubisoft",
			"Universalg",
			"Unilever",
			"Visa",
			"Vlerick-Business-School",
			"VMT GO",
			"Voi-Scooters",
			"Volkswagen",
			"VRT",
			"WP-Engine",
			"Yacht",
			"Yahoo",
			"e-tracker"},
		"php": {"(?i)(password|passwd|pwd|pass|secret)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(api[_-]?key|apikey|api[_-]?secret|secret[_-]?key|key|access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(token|access[_-]?token|auth[_-]?token|jwt|bearer)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(db[_-]?password|database[_-]?password|db[_-]?user|database[_-]?user|db[_-]?host|database[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"-----BEGIN(.*)PRIVATE KEY-----[\\s\\S]*?-----END(.*)PRIVATE KEY-----",
			"(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key|aws[_-]?secret[_-]?access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9+/=]+['\"\\s]",
			"(?i)(firebase[_-]?api[_-]?key|firebase[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(google[_-]?api[_-]?key|google[_-]?client[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(smtp[_-]?password|smtp[_-]?user|smtp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(ftp[_-]?password|ftp[_-]?user|ftp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(oauth[_-]?token|oauth[_-]?secret|oauth2)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(eyJ[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+)",
			"(?i)(azure[_-]?key|azure[_-]?secret|azure[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(stripe[_-]?key|stripe[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(paypal[_-]?key|paypal[_-]?secret|paypal[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"https://hooks\\.slack\\.com/services/[A-Za-z0-9/-]+",
			"(?i)(github[_-]?token|github[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(gitlab[_-]?token|gitlab[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(webhook[_-]?url|webhook[_-]?token|webhook[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-/]+['\"\\s]",
			"s3://[A-Za-z0-9._\\-/]+",
			"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			"(?i)(https?://)[^\\s]+:[^\\s]+@[^\\s]+",
			"(?i)(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519) [A-Za-z0-9+/=]+",
		},
		"js": {"(?i)(password|passwd|pwd|pass|secret)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(api[_-]?key|apikey|api[_-]?secret|secret[_-]?key|key|access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(token|access[_-]?token|auth[_-]?token|jwt|bearer)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(db[_-]?password|database[_-]?password|db[_-]?user|database[_-]?user|db[_-]?host|database[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"-----BEGIN(.*)PRIVATE KEY-----[\\s\\S]*?-----END(.*)PRIVATE KEY-----",
			"(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key|aws[_-]?secret[_-]?access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9+/=]+['\"\\s]",
			"(?i)(firebase[_-]?api[_-]?key|firebase[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(google[_-]?api[_-]?key|google[_-]?client[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(smtp[_-]?password|smtp[_-]?user|smtp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(ftp[_-]?password|ftp[_-]?user|ftp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(oauth[_-]?token|oauth[_-]?secret|oauth2)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(eyJ[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+)",
			"(?i)(azure[_-]?key|azure[_-]?secret|azure[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(stripe[_-]?key|stripe[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(paypal[_-]?key|paypal[_-]?secret|paypal[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"https://hooks\\.slack\\.com/services/[A-Za-z0-9/-]+",
			"(?i)(github[_-]?token|github[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(gitlab[_-]?token|gitlab[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(webhook[_-]?url|webhook[_-]?token|webhook[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-/]+['\"\\s]",
			"s3://[A-Za-z0-9._\\-/]+",
			"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			"(?i)(https?://)[^\\s]+:[^\\s]+@[^\\s]+",
			"(?i)(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519) [A-Za-z0-9+/=]+",
		},
		"py": {"(?i)(password|passwd|pwd|pass|secret)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(api[_-]?key|apikey|api[_-]?secret|secret[_-]?key|key|access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(token|access[_-]?token|auth[_-]?token|jwt|bearer)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(db[_-]?password|database[_-]?password|db[_-]?user|database[_-]?user|db[_-]?host|database[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"-----BEGIN(.*)PRIVATE KEY-----[\\s\\S]*?-----END(.*)PRIVATE KEY-----",
			"(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key|aws[_-]?secret[_-]?access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9+/=]+['\"\\s]",
			"(?i)(firebase[_-]?api[_-]?key|firebase[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(google[_-]?api[_-]?key|google[_-]?client[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(smtp[_-]?password|smtp[_-]?user|smtp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(ftp[_-]?password|ftp[_-]?user|ftp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(oauth[_-]?token|oauth[_-]?secret|oauth2)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(eyJ[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+)",
			"(?i)(azure[_-]?key|azure[_-]?secret|azure[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(stripe[_-]?key|stripe[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(paypal[_-]?key|paypal[_-]?secret|paypal[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"https://hooks\\.slack\\.com/services/[A-Za-z0-9/-]+",
			"(?i)(github[_-]?token|github[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(gitlab[_-]?token|gitlab[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(webhook[_-]?url|webhook[_-]?token|webhook[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-/]+['\"\\s]",
			"s3://[A-Za-z0-9._\\-/]+",
			"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			"(?i)(https?://)[^\\s]+:[^\\s]+@[^\\s]+",
			"(?i)(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519) [A-Za-z0-9+/=]+",
		},
		"java": {"(?i)(password|passwd|pwd|pass|secret)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(api[_-]?key|apikey|api[_-]?secret|secret[_-]?key|key|access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(token|access[_-]?token|auth[_-]?token|jwt|bearer)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(db[_-]?password|database[_-]?password|db[_-]?user|database[_-]?user|db[_-]?host|database[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"-----BEGIN(.*)PRIVATE KEY-----[\\s\\S]*?-----END(.*)PRIVATE KEY-----",
			"(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key|aws[_-]?secret[_-]?access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9+/=]+['\"\\s]",
			"(?i)(firebase[_-]?api[_-]?key|firebase[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(google[_-]?api[_-]?key|google[_-]?client[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(smtp[_-]?password|smtp[_-]?user|smtp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(ftp[_-]?password|ftp[_-]?user|ftp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(oauth[_-]?token|oauth[_-]?secret|oauth2)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(eyJ[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+)",
			"(?i)(azure[_-]?key|azure[_-]?secret|azure[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(stripe[_-]?key|stripe[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(paypal[_-]?key|paypal[_-]?secret|paypal[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"https://hooks\\.slack\\.com/services/[A-Za-z0-9/-]+",
			"(?i)(github[_-]?token|github[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(gitlab[_-]?token|gitlab[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(webhook[_-]?url|webhook[_-]?token|webhook[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-/]+['\"\\s]",
			"s3://[A-Za-z0-9._\\-/]+",
			"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			"(?i)(https?://)[^\\s]+:[^\\s]+@[^\\s]+",
			"(?i)(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519) [A-Za-z0-9+/=]+",
		},
		"go": {"(?i)(password|passwd|pwd|pass|secret)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(api[_-]?key|apikey|api[_-]?secret|secret[_-]?key|key|access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(token|access[_-]?token|auth[_-]?token|jwt|bearer)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(db[_-]?password|database[_-]?password|db[_-]?user|database[_-]?user|db[_-]?host|database[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"-----BEGIN(.*)PRIVATE KEY-----[\\s\\S]*?-----END(.*)PRIVATE KEY-----",
			"(?i)(aws[_-]?access[_-]?key|aws[_-]?secret[_-]?key|aws[_-]?secret[_-]?access[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9+/=]+['\"\\s]",
			"(?i)(firebase[_-]?api[_-]?key|firebase[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(google[_-]?api[_-]?key|google[_-]?client[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(smtp[_-]?password|smtp[_-]?user|smtp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(ftp[_-]?password|ftp[_-]?user|ftp[_-]?host)[^\\n]*[=:][^\\n]*['\"\\s][^'\"\\s]+['\"\\s]",
			"(?i)(oauth[_-]?token|oauth[_-]?secret|oauth2)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(eyJ[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+\\.[a-zA-Z0-9-_]+)",
			"(?i)(azure[_-]?key|azure[_-]?secret|azure[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(stripe[_-]?key|stripe[_-]?secret)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(paypal[_-]?key|paypal[_-]?secret|paypal[_-]?token)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"https://hooks\\.slack\\.com/services/[A-Za-z0-9/-]+",
			"(?i)(github[_-]?token|github[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(gitlab[_-]?token|gitlab[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-]+['\"\\s]",
			"(?i)(webhook[_-]?url|webhook[_-]?token|webhook[_-]?key)[^\\n]*[=:][^\\n]*['\"\\s][A-Za-z0-9._\\-/]+['\"\\s]",
			"s3://[A-Za-z0-9._\\-/]+",
			"[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}",
			"(?i)(https?://)[^\\s]+:[^\\s]+@[^\\s]+",
			"(?i)(ssh-rsa|ssh-dss|ecdsa-sha2-nistp256|ssh-ed25519) [A-Za-z0-9+/=]+",
		},
	}
	for _, keyword := range keywords {
		outputFile := fmt.Sprintf("results-%s.json", keyword)
		fmt.Printf("Searching for files with keyword: %s\n", keyword)
		var files []api.FileInfo
		var err error
		maxRetries := 3
		for retries := 0; retries < maxRetries; retries++ {
			files, err = api.QueryFiles(sessionCookie, []string{keyword}, extensions)
			if err == nil {
				break // Exit the retry loop if successful
			}
			log.Printf("Retry %d/%d for keyword '%s' failed: %v", retries+1, maxRetries, keyword, err)
			time.Sleep(2 * time.Second)
		}
		if err != nil {
			log.Printf("All retries failed for keyword '%s'\n", keyword)
			continue
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

		// ensure to goes periodicly saving it's making to avoid big lost if the process is interrupted
		go func() {
			// A channel that emits a signal every time the ticker fires.
			for range ticker.C {
				// save the file periodically
				mutex.Lock()
				err := saveResults(results, outputFile)
				if err != nil {
					log.Printf("Error saving periodic results for keyword '%s': %v", keyword, err)
				}
				// fmt.Printf("Result added: %+v\n", results) // Add this line for debugging
				// MUTEX write the file but priventing race conditions
				mutex.Unlock()
			}
		}()

		for _, fileInfo := range files {

			fmt.Println("Processing file:", fileInfo.Filename)
			if fileInfo.Size > 50*1024*1024 { // Skip files larger than 50 MB
				fmt.Printf("Skipping large file: %s\n", fileInfo.Filename)
				continue
			}

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

				result := download.ProcessFile(file, extensions) // redefine result
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
		wg.Wait() // Wait for all goroutines to complete
		mutex.Lock()
		err = saveResults(results, outputFile)
		if err != nil {
			log.Printf("Error saving final results for keyword '%s': %v", keyword, err)
		}
		mutex.Unlock()
	}
}

func saveResults(results []map[string]interface{}, outputFile string) error {
	fmt.Printf("Saving %d results...\n", len(results)) // Debug log

	file, err := os.Create(outputFile) // Create (or overwrite) results.json
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", outputFile, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Add indentation for readability

	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("failed to write JSON to file '%s': %w", outputFile, err)
	}
	fmt.Printf("Results saved to %s\n", outputFile)
	return nil
}
