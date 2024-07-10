package entities

type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value string `json:"value"`
}
