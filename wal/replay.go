package wal

func FilterOp(recs []Record, op string) []Record {
	out := make([]Record, 0, len(recs))
	for _, rec := range recs {
		if rec.Op == op {
			out = append(out, rec)
		}
	}
	return out
}

func LastSeq(recs []Record) uint64 {
	var n uint64
	for _, rec := range recs {
		if rec.Seq > n {
			n = rec.Seq
		}
	}
	return n
}
