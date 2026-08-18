package event

func FilterKind(evs []Event, kind string) []Event {
	out := make([]Event, 0, len(evs))
	for _, ev := range evs {
		if ev.Kind == kind {
			out = append(out, ev)
		}
	}
	return out
}
