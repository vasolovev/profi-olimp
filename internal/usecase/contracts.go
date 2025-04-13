// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// Translation -.
	Translation interface {
		Translate(context.Context, entity.Translation) (entity.Translation, error)
		History(context.Context) ([]entity.Translation, error)
	}
	// Student -.
	Student interface {
		CreateStudent(ctx context.Context, student entity.Student) (entity.Student, error)
		GetStudents(ctx context.Context) ([]entity.Student, error)
		GetStudentByID(ctx context.Context, id int) (entity.Student, error)
		UpdateStudent(ctx context.Context, student entity.Student) error
		DeleteStudent(ctx context.Context, id int) error
		SearchStudents(ctx context.Context, query string) ([]entity.Student, error)
	}

	// Group -.
	Group interface {
		CreateGroup(ctx context.Context, group entity.Group) (entity.Group, error)
		GetGroups(ctx context.Context) ([]entity.Group, error)
		GetGroupByID(ctx context.Context, id int) (entity.Group, error)
		UpdateGroup(ctx context.Context, group entity.Group) error
		DeleteGroup(ctx context.Context, id int) error
		SearchGroups(ctx context.Context, query string) ([]entity.Group, error)
	}
)
