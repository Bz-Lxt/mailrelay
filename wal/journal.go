package wal

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type Journal struct {
	mu   sync.Mutex
	path string
	file *os.File
	seq  uint64
}

func Open(path string) (*Journal, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	j := &Journal{path: path, file: f}
	recs, err := ReplayFile(path)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if n := len(recs); n > 0 {
		j.seq = recs[n-1].Seq
	}
	return j, nil
}

func (j *Journal) Path() string { return j.path }

func (j *Journal) NextSeq() uint64 {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.seq++
	return j.seq
}

func (j *Journal) Append(rec Record) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	defer j.Seal()
	if rec.Seq == 0 {
		j.seq++
		rec.Seq = j.seq
	} else if rec.Seq > j.seq {
		j.seq = rec.Seq
	}
	frame, err := Encode(rec)
	if err != nil {
		return err
	}
	if _, err := j.file.Write(frame); err != nil {
		return err
	}
	return j.file.Sync()
}

func (j *Journal) Close() error {
	j.mu.Lock()
	defer j.mu.Unlock()
	if j.file == nil {
		return nil
	}
	err := j.file.Close()
	j.file = nil
	return err
}

func (j *Journal) TruncatePrefix(keepFrom uint64) error {
	j.mu.Lock()
	defer j.mu.Unlock()
	recs, err := ReplayFile(j.path)
	if err != nil {
		return err
	}
	keep := make([]Record, 0, len(recs))
	for _, rec := range recs {
		if rec.Seq >= keepFrom {
			keep = append(keep, rec)
		}
	}
	tmp := j.path + ".tmp"
	tf, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	for _, rec := range keep {
		frame, err := Encode(rec)
		if err != nil {
			_ = tf.Close()
			return err
		}
		if _, err := tf.Write(frame); err != nil {
			_ = tf.Close()
			return err
		}
	}
	if err := tf.Sync(); err != nil {
		_ = tf.Close()
		return err
	}
	if err := tf.Close(); err != nil {
		return err
	}
	if err := j.file.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp, j.path); err != nil {
		return err
	}
	f, err := os.OpenFile(j.path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Seek(0, 2); err != nil {
		return err
	}
	j.file = f
	return nil
}

func ReplayFile(path string) ([]Record, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Record
	i := 0
	for i+8 <= len(raw) {
		n := int(uint32(raw[i+4])<<24 | uint32(raw[i+5])<<16 | uint32(raw[i+6])<<8 | uint32(raw[i+7]))
		end := i + 8 + n
		if n < 0 || end > len(raw) {
			break
		}
		rec, err := Decode(raw[i:end])
		if err != nil {
			return nil, fmt.Errorf("replay offset %d: %w", i, err)
		}
		out = append(out, rec)
		i = end
	}
	return out, nil
}
