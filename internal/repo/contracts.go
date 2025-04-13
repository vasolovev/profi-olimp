// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/evrone/go-clean-template/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// TranslationRepo -.
	TranslationRepo interface {
		Store(context.Context, entity.Translation) error
		GetHistory(context.Context) ([]entity.Translation, error)
	}

	// TranslationWebAPI -.
	TranslationWebAPI interface {
		Translate(entity.Translation) (entity.Translation, error)
	}
)

type StudentRepo interface {
	CreateStudent(ctx context.Context, student entity.Student) (entity.Student, error)
	GetStudents(ctx context.Context) ([]entity.Student, error)
	GetStudentByID(ctx context.Context, id int) (entity.Student, error)
	UpdateStudent(ctx context.Context, student entity.Student) error
	DeleteStudent(ctx context.Context, id int) error
	SearchStudents(ctx context.Context, query string) ([]entity.Student, error)
}

// GroupRepo defines the group repository interface.
type GroupRepo interface {
	CreateGroup(ctx context.Context, group entity.Group) (entity.Group, error)
	GetGroups(ctx context.Context) ([]entity.Group, error)
	GetGroupByID(ctx context.Context, id int) (entity.Group, error)
	UpdateGroup(ctx context.Context, group entity.Group) error
	DeleteGroup(ctx context.Context, id int) error
	SearchGroups(ctx context.Context, query string) ([]entity.Group, error)
	HasSubgroups(ctx context.Context, id int) (bool, error)
	GetGroupWithSubgroups(ctx context.Context, id int) (entity.Group, error)
}
