package repository

// Covid19Subcription ...
type Covid19Subcription struct {
	IsSubscribe bool   `firestore:"is_subscribe" json:"is_subscribe"`
	Provider    string `firestore:"provider" json:"provider"`
	Token       string `firestore:"token" json:"token"`
}

// Covid19Data ...
type Covid19Data struct {
	Confirmed int `firestore:"confirmed" json:"confirmed"`
	Suspected int `firestore:"suspected" json:"suspected"`
	Recovered int `firestore:"recovered" json:"recovered"`
	Deceased  int `firestore:"deceased" json:"deceased"`
}
