package authorisation

type Claims struct {
	Id  string   `json:"nameid"`
	Aud []string `json:"-"`
}
