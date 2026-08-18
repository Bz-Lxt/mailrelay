package wal

import "encoding/json"

func DumpJSON(recs []Record) ([]byte, error) {
	return json.Marshal(recs)
}

func CountByOp(recs []Record) map[string]int {
	out := map[string]int{}
	for _, rec := range recs {
		out[rec.Op]++
	}
	return out
}
