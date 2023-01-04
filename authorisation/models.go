package authorisation

type Claims struct {
	Id   string   `json:"sub"`
	Aud  []string `json:"aud"`
	Name string   `json:"preferred_username"`
}
