package Set

type Set struct {
	buf map[string]struct{}
}

func New() Set {
	return Set{
		buf: make(map[string]struct{}),
	}
}
func (s *Set) Has(key string) bool {
	_, exist := s.buf[key]
	return exist
}
func (s *Set) Add(key string) *Set {
	s.buf[key] = struct {}{}
	return s
}
func (s *Set) Adds(keys []string) *Set {
	for _, key := range keys {
		s.Add(key)
	}
	return s
}
func (s *Set) Delete(key string) *Set {
	delete(s.buf, key)
	return s
}
func (s *Set) Clear() *Set {
	for key, _ := range s.buf {
		delete(s.buf, key)
	}
	return s
}
func (s *Set) Keys() []string {
	result := make([]string, 0)
	for key, _ := range s.buf {
		result = append(result, key)
	}
	return result
}