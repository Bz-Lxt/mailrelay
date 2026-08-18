package types

type Envelope struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Status  string `json:"status"`
	Box     string `json:"box"`
	Tries   int    `json:"tries"`
	DeferTo string `json:"defer_to,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

func CloneEnvelope(e Envelope) Envelope { return e }
