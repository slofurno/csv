package csv

import (
	"bufio"
	"bytes"
	"io"
)

type UnquotedReader struct {
	Comma rune

	column       int
	r            *bufio.Reader
	lineBuffer   bytes.Buffer
	fieldIndexes []int
	lastRecord   []string
}

func NewUnquoted(r io.Reader) *UnquotedReader {
	return &UnquotedReader{
		r:     bufio.NewReader(r),
		Comma: ',',
	}
}

func (s *UnquotedReader) readRune() (rune, error) {
	r, _, err := s.r.ReadRune()
	s.column++
	return r, err

}

func (s *UnquotedReader) parseField() (bool, rune, error) {
	for {
		r, err := s.readRune()
		if err == io.EOF && s.column != 0 {
			return true, 0, err
		}

		if err != nil {
			return false, 0, err
		}

		switch r {
		case s.Comma:
			return true, s.Comma, nil
		case '\n':
			if s.column == 0 {
				return false, '\n', nil
			}
			return true, '\n', nil
		case '\r':
			continue
		default:
			s.lineBuffer.WriteRune(r)
		}
	}
}

func (s *UnquotedReader) Read() ([]string, error) {
	s.lineBuffer.Reset()
	s.fieldIndexes = s.fieldIndexes[:0]
	s.column = -1

	for {
		idx := s.lineBuffer.Len()

		haveField, delim, err := s.parseField()
		if haveField {
			s.fieldIndexes = append(s.fieldIndexes, idx)
		}

		if delim == '\n' || err == io.EOF {
			if len(s.fieldIndexes) == 0 {
				return nil, err
			}
			break
		}

		if err != nil {
			return nil, err
		}

	}

	fieldCount := len(s.fieldIndexes)
	line := s.lineBuffer.String()
	var fields []string

	if cap(s.lastRecord) >= fieldCount {
		fields = s.lastRecord[:fieldCount]
	} else {
		fields = make([]string, fieldCount)
	}

	for i, idx := range s.fieldIndexes {
		if i == fieldCount-1 {
			fields[i] = line[idx:]
		} else {
			fields[i] = line[idx:s.fieldIndexes[i+1]]
		}
	}

	return fields, nil
}
