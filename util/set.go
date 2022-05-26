package util

type Set struct {
	data map[string]bool
}

func (s *Set) Contains(e string) bool {
	_, contains := s.data[e]
	return contains
}

func (s *Set) Add(e string) {
	s.data[e] = true
}

func (s *Set) Delete(e string) {
	delete(s.data, e)
}

func (s *Set) Length() int {
	return len(s.data)
}

func SetFromArray(A []string) Set {
	s := Set{data: make(map[string]bool)}
	for _, a := range A {
		s.Add(a)
	}
	return s
}
