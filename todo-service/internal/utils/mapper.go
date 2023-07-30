package utils

import "github.com/olad5/productive-pulse/todo-service/internal/domain"

func ToTodoDTO(todo domain.Todo) map[string]interface{} {
	return map[string]interface{}{
		"id":         todo.ID,
		"text":       todo.Text,
		"created_at": todo.CreatedAt,
		"updated_at": todo.UpdatedAt,
	}
}
