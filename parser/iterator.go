package parser

import (
	"errors"
)

type Iterator[T any] struct {
	elems []T
	pos   int
}

func NewIterator[T any](arr []T) Iterator[T] {
	return Iterator[T]{
		elems: arr,
		pos: 0,
	}
}

func (i *Iterator[T]) Consume() {
	i.pos += 1;
}

func (i *Iterator[T]) Next() (*T, error) {
	if i.pos+1 >= len(i.elems) {
		return nil, errors.New("Iterator overflow")
	}

	i.pos += 1
	return &i.elems[i.pos], nil
}

func (i *Iterator[T]) Peek() (*T, error) {
	v, e := i.Next();

	i.pos -= 1;
	return v, e
}

func (i *Iterator[T]) PeekOr(or T) T {
	if v, e := i.Peek(); e != nil {
		return or
	} else {
		return *v
	}
}

func (i *Iterator[T]) Prev() (*T, error) {
	if i.pos-1 < 0 {
		return nil, errors.New("Iterator underflow")
	}

	i.pos -= 1
	return &i.elems[i.pos], nil
}

func (i *Iterator[T]) Curr() (*T, error) {
	if i.pos >= len(i.elems) {
		return nil, errors.New("Iterator overflow")
	}

	return &i.elems[i.pos], nil
}

func (i *Iterator[T]) Tell() int {
	return i.pos
}