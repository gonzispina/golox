package lox

func NewScanner(source []byte) *Scanner {
	return &Scanner{source: source}
}

type Scanner struct {
	source []byte
}

func (s *Scanner) ScanTokens() ([]string, error) {
	if s.source == nil {
		return []string{}, nil
	}

	var res []string
	for _, t := range s.source {
		res = append(res, string(t))
	}

	return res, nil
}
