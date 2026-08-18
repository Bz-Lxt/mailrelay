package catalog

var Boxes = []string{"inbound", "outbound", "deferred", "bounced", "sent"}

func KnownBox(name string) bool {
	for _, b := range Boxes {
		if b == name {
			return true
		}
	}
	return false
}
