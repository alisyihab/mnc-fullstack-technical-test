package shared

type RuneStack struct {
	items []rune
}

func NewRuneStack() *RuneStack {
	return &RuneStack{
		items: make([]rune, 0),
	}
}

func (s *RuneStack) Push(r rune) {
	s.items = append(s.items, r)
}

func (s *RuneStack) Pop() (rune, bool) {
	if len(s.items) == 0 {
		return 0, false
	}

	lastIndex := len(s.items) - 1
	value := s.items[lastIndex]

	s.items = s.items[:lastIndex]

	return value, true
}

func (s *RuneStack) IsEmpty() bool {
	return len(s.items) == 0
}
