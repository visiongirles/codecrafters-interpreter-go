package main

type Stack struct {
	items []rune
}

// Push добавляет элемент на вершину стека
func (s *Stack) Push(item rune) {
	s.items = append(s.items, item)
}

// Pop удаляет и возвращает элемент с вершины стека
func (s *Stack) Pop() rune {
	if len(s.items) == 0 {
		return 0 // возвращаем ноль, если стек пуст
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

// IsEmpty проверяет, пуст ли стек
func (s *Stack) IsEmpty() bool {
	return len(s.items) == 0
}

// Top возвращает элемент с вершины стека без удаления
func (s *Stack) Top() rune {
	if len(s.items) == 0 {
		return 0 // возвращаем ноль, если стек пуст
	}
	return s.items[len(s.items)-1]
}
func (s *Stack) isBalanced(expression string) bool {
	matchingParentheses := map[rune]rune{')': '(', '}': '{', ']': '['}

	for _, char := range expression {
		if char == '(' || char == '{' || char == '[' {
			s.Push(char)
		} else if char == ')' || char == '}' || char == ']' {
			if s.IsEmpty() || s.Pop() != matchingParentheses[char] {
				return false
			}
		}
	}
	return s.IsEmpty()
}
