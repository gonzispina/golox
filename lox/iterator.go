package lox

import "unicode/utf8"

const linebreak = '\n'

func newIterator(s []rune) *Iterator {
	return &Iterator{
		source:  s,
		start:   0,
		current: 0,
		line:    1,
	}
}

type Iterator struct {
	source  []rune
	start   int
	current int
	line    int
	column  int
}

func (i *Iterator) startLexeme() {
	i.start = i.current
}

func (i *Iterator) currentLexeme() string {
	return string(i.source[i.start:i.current])
}

func (i *Iterator) isAtEnd() bool {
	return i.current == len(i.source)
}

func (i *Iterator) advance() rune {
	if i.isAtEnd() {
		return utf8.RuneError
	}

	r := i.source[i.current]
	i.current++

	if r == linebreak {
		i.column = 0
		i.line++
	}

	i.column++
	return r
}

func (i *Iterator) peek() rune {
	if i.isAtEnd() {
		return utf8.RuneError
	}
	return i.source[i.current]
}

func (i *Iterator) next() rune {
	if i.isAtEnd() {
		return utf8.RuneError
	}
	return i.source[i.current+1]
}

func (i *Iterator) match(s rune) bool {
	p := i.peek()

	if p == utf8.RuneError {
		return false
	}

	if p != s {
		return false
	}

	i.advance()
	return true
}
