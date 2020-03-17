package repository

// Covid19Subcription ...
type Covid19Subcription struct {
	IsSubscribe bool   `firestore:"is_subscribe" json:"is_subscribe"`
	Provider    string `firestore:"provider" json:"provider"`
	Token       string `firestore:"token" json:"token"`
}
