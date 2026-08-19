package wal

// Seal flushes the pending writes and releases the on-disk handle so the
// journal can be replayed safely from the file alone. It is the counterpart
// to Open for callers that want an explicit finalization point.
func (j *Journal) Seal() error {
	if err := j.file.Sync(); err != nil {
		return err
	}
	return j.file.Close()
}
