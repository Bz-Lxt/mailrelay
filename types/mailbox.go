package types

type Mailbox struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

const (
	BoxInbound  = "inbound"
	BoxOutbound = "outbound"
	BoxDeferred = "deferred"
	BoxBounced  = "bounced"
	BoxSent     = "sent"
)
