package main

type StackASTNode struct {
	items []ASTNode
}

// Push добавляет элемент на вершину стека
func (s *StackASTNode) Push(item ASTNode) {
	s.items = append(s.items, item)
}

// Pop удаляет и возвращает элемент с вершины стека
func (s *StackASTNode) Pop() ASTNode {
	if len(s.items) == 0 {
		return nil // возвращаем ноль, если стек пуст
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

// IsEmpty проверяет, пуст ли стек
func (s *StackASTNode) IsEmpty() bool {
	return len(s.items) == 0
}

// Top возвращает элемент с вершины стека без удаления
func (s *StackASTNode) Top() ASTNode {
	if len(s.items) == 0 {
		return nil // возвращаем ноль, если стек пуст
	}
	return s.items[len(s.items)-1]
}

// func (s *Stack) isBalanced(expression string) bool {
// 	matchingParentheses := map[rune]rune{')': '(', '}': '{', ']': '['}

// 	for _, char := range expression {
// 		if char == '(' || char == '{' || char == '[' {
// 			s.Push(char)
// 		} else if char == ')' || char == '}' || char == ']' {
// 			if s.IsEmpty() || s.Pop() != matchingParentheses[char] {
// 				return false
// 			}
// 		}
// 	}
// 	return s.IsEmpty()
// }
