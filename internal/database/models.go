package database

// FeedResponse represents the complete feed endpoint response
type FeedResponse struct {
	Stats Stats      `json:"stats"`
	Feed  []FeedItem `json:"feed"`
}

// Stats contains request metadata
type Stats struct {
	URL  string  `json:"url"`
	Time float64 `json:"time"` // processing time in seconds
}

// FeedItem represents a single company in the feed
type FeedItem struct {
	CVR            string            `json:"cvr"`
	Name           string            `json:"name"`
	BusinessCode   int               `json:"business_code"`
	NumberOfOwners int               `json:"number_of_owners"`
	Increases      []CapitalIncrease `json:"increases"`
}

// CapitalIncrease represents a single capital increase event
type CapitalIncrease struct {
	Capital   float64 `json:"capital"`
	ValidFrom string  `json:"validFrom"` // YYYY-MM-DD format
}
