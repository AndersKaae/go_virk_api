package database

// FeedResponse represents the complete feed endpoint response
type FeedResponse struct {
	Stats Stats      `json:"stats"`
	Feed  []FeedItem `json:"feed"`
}

// Stats contains request metadata
type Stats struct {
	URL         string  `json:"url"`
	Performance float64 `json:"performance"` // processing time in seconds
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

// CompanyInfoResponse represents the company_info endpoint response
type CompanyInfoResponse struct {
	Performance string  `json:"performance"`
	CVR         string  `json:"cvr"`
	Name        string  `json:"name"`
	Start       string  `json:"start"`
	Branchekode string  `json:"branchekode"`
	Website     *string `json:"website"`
	Address     string  `json:"adress"` // Note: typo in original API "adress"
}

// ManagementResponse represents the management endpoint response
type ManagementResponse struct {
	Performance float64          `json:"performance"`
	Result      ManagementResult `json:"result"`
}

// ManagementResult contains board and management lists
type ManagementResult struct {
	Board      []ManagementPerson `json:"board"`
	Management []ManagementPerson `json:"management"`
}

// ManagementPerson represents a person in management or board
type ManagementPerson struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

// OwnersResponse represents the owners endpoint response
type OwnersResponse struct {
	Performance float64 `json:"performance"`
	Result      []Owner `json:"result"`
}

// Owner represents a single owner
type Owner struct {
	Name  string `json:"name"`
	CVR   *int   `json:"cvr"`
	Value string `json:"value"`
}

// TopInvestorsResponse represents the top_investors endpoint response
type TopInvestorsResponse struct {
	Performance float64       `json:"performance"`
	Result      []TopInvestor `json:"result"`
}

// TopInvestor represents an investor with company count
type TopInvestor struct {
	Name      string `json:"name"`
	CVR       *int   `json:"cvr"`
	Companies int    `json:"companies"`
}

// SearchResponse represents the search endpoint response
type SearchResponse struct {
	Performance float64        `json:"performance"`
	Result      []SearchResult `json:"result"`
}

// SearchResult represents a single search result
type SearchResult struct {
	CVR  string `json:"cvr"`
	Name string `json:"name"`
}

// SitemapEntry represents a single sitemap entry
type SitemapEntry struct {
	Loc     string `json:"loc"`
	LastMod string `json:"lastmod"`
}

// IncreaseResponse represents the API response for capital increases
type IncreaseResponse struct {
	Stats     Stats             `json:"stats"`
	Increases []CapitalIncrease `json:"increases"`
}
