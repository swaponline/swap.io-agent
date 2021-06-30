package Set

type Set struct {
	buf map[string]struct{}
}

func New() Set {
	return Set{
		buf: make(map[string]struct{}),
	}
}

func (s *Set) Add(key string) *Set {
	s.buf[key] = struct {}{}
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
	result := make([]string, len(s.buf))
	for key, _ := range s.buf {
		result = append(result, key)
	}
	return result
}