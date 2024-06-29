package api

import (
	"context"
	"sync"
)

type SpreadsheetsMock struct {
	Rows map[string][][]string
	lock sync.Mutex
}

func (s *SpreadsheetsMock) AppendRow(ctx context.Context, sheetName string, row []string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.Rows == nil {
		s.Rows = make(map[string][][]string)
	}

	s.Rows[sheetName] = append(s.Rows[sheetName], row)

	return nil
}
