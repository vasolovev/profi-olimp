package student

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase implements the student use case interface.
type UseCase struct {
	repo repo.StudentRepo
}

// New creates a new student use case.
func New(r repo.StudentRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// CreateStudent creates a new student.
func (uc *UseCase) CreateStudent(ctx context.Context, student entity.Student) (entity.Student, error) {
	s, err := uc.repo.CreateStudent(ctx, student)
	if err != nil {
		return entity.Student{}, fmt.Errorf("StudentUseCase - CreateStudent - uc.repo.CreateStudent: %w", err)
	}

	return s, nil
}

// GetStudents retrieves all students.
func (uc *UseCase) GetStudents(ctx context.Context) ([]entity.Student, error) {
	students, err := uc.repo.GetStudents(ctx)
	if err != nil {
		return nil, fmt.Errorf("StudentUseCase - GetStudents - uc.repo.GetStudents: %w", err)
	}

	return students, nil
}

// GetStudentByID retrieves a student by ID.
func (uc *UseCase) GetStudentByID(ctx context.Context, id int) (entity.Student, error) {
	student, err := uc.repo.GetStudentByID(ctx, id)
	if err != nil {
		return entity.Student{}, fmt.Errorf("StudentUseCase - GetStudentByID - uc.repo.GetStudentByID: %w", err)
	}

	return student, nil
}

// UpdateStudent updates an existing student.
func (uc *UseCase) UpdateStudent(ctx context.Context, student entity.Student) error {
	err := uc.repo.UpdateStudent(ctx, student)
	if err != nil {
		return fmt.Errorf("StudentUseCase - UpdateStudent - uc.repo.UpdateStudent: %w", err)
	}

	return nil
}

// DeleteStudent deletes a student by ID.
func (uc *UseCase) DeleteStudent(ctx context.Context, id int) error {
	err := uc.repo.DeleteStudent(ctx, id)
	if err != nil {
		return fmt.Errorf("StudentUseCase - DeleteStudent - uc.repo.DeleteStudent: %w", err)
	}

	return nil
}

// SearchStudents searches for students by name or group name.
func (uc *UseCase) SearchStudents(ctx context.Context, query string) ([]entity.Student, error) {
	students, err := uc.repo.SearchStudents(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("StudentUseCase - SearchStudents - uc.repo.SearchStudents: %w", err)
	}

	return students, nil
}
