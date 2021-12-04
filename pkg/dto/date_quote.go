package dto

// LikeRequest ...
type LikeRequest struct {
	Date uint `json:"date"`
}

// LikeResponse ...
type LikeResponse struct {
	LikeCount uint `json:"like_count"`
}

// GetQuoteResponse ...
type GetQuoteResponse struct {
	Date      uint   `json:"date"`
	Quote     string `json:"quote"`
	Author    string `json:"author"`
	LikeCount uint   `json:"like_count"`
}
