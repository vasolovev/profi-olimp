package persistent

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
)

// StudentRepo implements the student repository interface
type StudentRepo struct {
	*postgres.Postgres
}

// GroupRepo implements the group repository interface
type GroupRepo struct {
	*postgres.Postgres
}

// NewStudentRepo creates a new student repository
func NewStudentRepo(pg *postgres.Postgres) *StudentRepo {
	return &StudentRepo{pg}
}

// NewGroupRepo creates a new group repository
func NewGroupRepo(pg *postgres.Postgres) *GroupRepo {
	return &GroupRepo{pg}
}

// CreateStudent creates a new student
func (r *StudentRepo) CreateStudent(ctx context.Context, student entity.Student) (entity.Student, error) {
	sql, args, err := r.Builder.
		Insert("students").
		Columns("name", "email", "group_id").
		Values(student.Name, student.Email, student.GroupID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.Student{}, fmt.Errorf("StudentRepo - CreateStudent - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&student.ID)
	if err != nil {
		return entity.Student{}, fmt.Errorf("StudentRepo - CreateStudent - r.Pool.QueryRow: %w", err)
	}

	return student, nil
}

// GetStudents retrieves all students
func (r *StudentRepo) GetStudents(ctx context.Context) ([]entity.Student, error) {
	sql, _, err := r.Builder.
		Select("id", "name", "group_id").
		From("students").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StudentRepo - GetStudents - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("StudentRepo - GetStudents - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var students []entity.Student
	for rows.Next() {
		var s entity.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.GroupID); err != nil {
			return nil, fmt.Errorf("StudentRepo - GetStudents - rows.Scan: %w", err)
		}
		students = append(students, s)
	}

	return students, nil
}

// GetStudentByID retrieves a student by ID
func (r *StudentRepo) GetStudentByID(ctx context.Context, id int) (entity.Student, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "group_id").
		From("students").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.Student{}, fmt.Errorf("StudentRepo - GetStudentByID - r.Builder: %w", err)
	}

	var student entity.Student
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&student.ID, &student.Name, &student.GroupID)
	if err != nil {
		return entity.Student{}, fmt.Errorf("StudentRepo - GetStudentByID - r.Pool.QueryRow: %w", err)
	}

	return student, nil
}

// UpdateStudent updates an existing student
func (r *StudentRepo) UpdateStudent(ctx context.Context, student entity.Student) error {
	sql, args, err := r.Builder.
		Update("students").
		Set("name", student.Name).
		Set("group_id", student.GroupID).
		Where("id = ?", student.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("StudentRepo - UpdateStudent - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("StudentRepo - UpdateStudent - r.Pool.Exec: %w", err)
	}

	return nil
}

// DeleteStudent deletes a student by ID
func (r *StudentRepo) DeleteStudent(ctx context.Context, id int) error {
	sql, args, err := r.Builder.
		Delete("students").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("StudentRepo - DeleteStudent - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("StudentRepo - DeleteStudent - r.Pool.Exec: %w", err)
	}

	return nil
}

// SearchStudents searches for students by name or group name
func (r *StudentRepo) SearchStudents(ctx context.Context, query string) ([]entity.Student, error) {
	sql, args, err := r.Builder.
		Select("s.id", "s.name", "s.group_id").
		From("students s").
		Join("groups g ON s.group_id = g.id").
		Where("LOWER(s.name) LIKE LOWER(?) OR LOWER(g.name) LIKE LOWER(?)", "%"+query+"%", "%"+query+"%").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("StudentRepo - SearchStudents - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("StudentRepo - SearchStudents - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var students []entity.Student
	for rows.Next() {
		var s entity.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.GroupID); err != nil {
			return nil, fmt.Errorf("StudentRepo - SearchStudents - rows.Scan: %w", err)
		}
		students = append(students, s)
	}

	return students, nil
}

// CreateGroup creates a new group
func (r *GroupRepo) CreateGroup(ctx context.Context, group entity.Group) (entity.Group, error) {
	sql, args, err := r.Builder.
		Insert("groups").
		Columns("name", "parent_id").
		Values(group.Name, group.ParentID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupRepo - CreateGroup - r.Builder: %w", err)
	}

	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&group.ID)
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupRepo - CreateGroup - r.Pool.QueryRow: %w", err)
	}

	return group, nil
}

// GetGroups retrieves all groups with their subgroups
func (r *GroupRepo) GetGroups(ctx context.Context) ([]entity.Group, error) {
	// First get all root groups (those without parent)
	sql, _, err := r.Builder.
		Select("id", "name", "parent_id").
		From("groups").
		Where("parent_id IS NULL").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GroupRepo - GetGroups - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("GroupRepo - GetGroups - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var rootGroups []entity.Group
	for rows.Next() {
		var g entity.Group
		var parentID *int
		if err := rows.Scan(&g.ID, &g.Name, &parentID); err != nil {
			return nil, fmt.Errorf("GroupRepo - GetGroups - rows.Scan: %w", err)
		}
		g.ParentID = parentID

		// For each root group, get its subgroups recursively
		fullGroup, err := r.GetGroupWithSubgroups(ctx, g.ID)
		if err != nil {
			return nil, err
		}
		rootGroups = append(rootGroups, fullGroup)
	}

	return rootGroups, nil
}

// GetGroupByID retrieves a group by ID
func (r *GroupRepo) GetGroupByID(ctx context.Context, id int) (entity.Group, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "parent_id").
		From("groups").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupRepo - GetGroupByID - r.Builder: %w", err)
	}

	var group entity.Group
	var parentID *int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&group.ID, &group.Name, &parentID)
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupRepo - GetGroupByID - r.Pool.QueryRow: %w", err)
	}
	group.ParentID = parentID

	return group, nil
}

// GetGroupWithSubgroups retrieves a group with all its subgroups recursively
func (r *GroupRepo) GetGroupWithSubgroups(ctx context.Context, id int) (entity.Group, error) {
	// First get the group itself
	group, err := r.GetGroupByID(ctx, id)
	if err != nil {
		return entity.Group{}, err
	}

	// Then get all direct subgroups
	sql, args, err := r.Builder.
		Select("id", "name", "parent_id").
		From("groups").
		Where("parent_id = ?", id).
		ToSql()
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupRepo - GetGroupWithSubgroups - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupRepo - GetGroupWithSubgroups - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	// For each subgroup, recursively get its subgroups
	for rows.Next() {
		var subGroup entity.Group
		var parentID *int
		if err := rows.Scan(&subGroup.ID, &subGroup.Name, &parentID); err != nil {
			return entity.Group{}, fmt.Errorf("GroupRepo - GetGroupWithSubgroups - rows.Scan: %w", err)
		}
		subGroup.ParentID = parentID

		// Recursively get subgroups
		fullSubGroup, err := r.GetGroupWithSubgroups(ctx, subGroup.ID)
		if err != nil {
			return entity.Group{}, err
		}

		group.SubGroups = append(group.SubGroups, fullSubGroup)
	}

	return group, nil
}

// UpdateGroup updates an existing group
func (r *GroupRepo) UpdateGroup(ctx context.Context, group entity.Group) error {
	sql, args, err := r.Builder.
		Update("groups").
		Set("name", group.Name).
		Set("parent_id", group.ParentID).
		Where("id = ?", group.ID).
		ToSql()
	if err != nil {
		return fmt.Errorf("GroupRepo - UpdateGroup - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("GroupRepo - UpdateGroup - r.Pool.Exec: %w", err)
	}

	return nil
}

// DeleteGroup deletes a group by ID
func (r *GroupRepo) DeleteGroup(ctx context.Context, id int) error {
	sql, args, err := r.Builder.
		Delete("groups").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return fmt.Errorf("GroupRepo - DeleteGroup - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("GroupRepo - DeleteGroup - r.Pool.Exec: %w", err)
	}

	return nil
}

// SearchGroups searches for groups by name
func (r *GroupRepo) SearchGroups(ctx context.Context, query string) ([]entity.Group, error) {
	sql, args, err := r.Builder.
		Select("id", "name", "parent_id").
		From("groups").
		Where("LOWER(name) LIKE LOWER(?)", "%"+query+"%").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("GroupRepo - SearchGroups - r.Builder: %w", err)
	}

	rows, err := r.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GroupRepo - SearchGroups - r.Pool.Query: %w", err)
	}
	defer rows.Close()

	var groups []entity.Group
	for rows.Next() {
		var g entity.Group
		var parentID *int
		if err := rows.Scan(&g.ID, &g.Name, &parentID); err != nil {
			return nil, fmt.Errorf("GroupRepo - SearchGroups - rows.Scan: %w", err)
		}
		g.ParentID = parentID
		groups = append(groups, g)
	}

	return groups, nil
}

// HasSubgroups checks if a group has any subgroups
func (r *GroupRepo) HasSubgroups(ctx context.Context, id int) (bool, error) {
	sql, args, err := r.Builder.
		Select("COUNT(*)").
		From("groups").
		Where("parent_id = ?", id).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("GroupRepo - HasSubgroups - r.Builder: %w", err)
	}

	var count int
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("GroupRepo - HasSubgroups - r.Pool.QueryRow: %w", err)
	}

	return count > 0, nil
}
