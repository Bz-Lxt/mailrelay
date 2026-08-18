package wal

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
)

const Magic uint32 = 0x50594152

type Record struct {
	Op    string          `json:"op"`
	Seq   uint64          `json:"seq"`
	Body  json.RawMessage `json:"body"`
	Stamp string          `json:"stamp"`
}

func Encode(rec Record) ([]byte, error) {
	payload, err := json.Marshal(rec)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 8+len(payload))
	binary.BigEndian.PutUint32(buf[0:4], Magic)
	binary.BigEndian.PutUint32(buf[4:8], uint32(len(payload)))
	copy(buf[8:], payload)
	return buf, nil
}

func Decode(frame []byte) (Record, error) {
	if len(frame) < 8 {
		return Record{}, fmt.Errorf("short frame")
	}
	if binary.BigEndian.Uint32(frame[0:4]) != Magic {
		return Record{}, fmt.Errorf("bad magic")
	}
	n := binary.BigEndian.Uint32(frame[4:8])
	if int(n) != len(frame)-8 {
		return Record{}, fmt.Errorf("length mismatch")
	}
	var rec Record
	if err := json.Unmarshal(frame[8:], &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}
