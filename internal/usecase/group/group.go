package group

import (
	"context"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/internal/repo"
)

// UseCase implements the group use case interface.
type UseCase struct {
	repo repo.GroupRepo
}

// New creates a new group use case.
func New(r repo.GroupRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

// CreateGroup creates a new group.
func (uc *UseCase) CreateGroup(ctx context.Context, group entity.Group) (entity.Group, error) {
	g, err := uc.repo.CreateGroup(ctx, group)
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupUseCase - CreateGroup - uc.repo.CreateGroup: %w", err)
	}

	return g, nil
}

// GetGroups retrieves all groups with their subgroups.
func (uc *UseCase) GetGroups(ctx context.Context) ([]entity.Group, error) {
	groups, err := uc.repo.GetGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("GroupUseCase - GetGroups - uc.repo.GetGroups: %w", err)
	}

	return groups, nil
}

// GetGroupByID retrieves a group by ID.
func (uc *UseCase) GetGroupByID(ctx context.Context, id int) (entity.Group, error) {
	group, err := uc.repo.GetGroupByID(ctx, id)
	if err != nil {
		return entity.Group{}, fmt.Errorf("GroupUseCase - GetGroupByID - uc.repo.GetGroupByID: %w", err)
	}

	return group, nil
}

// UpdateGroup updates an existing group.
func (uc *UseCase) UpdateGroup(ctx context.Context, group entity.Group) error {
	err := uc.repo.UpdateGroup(ctx, group)
	if err != nil {
		return fmt.Errorf("GroupUseCase - UpdateGroup - uc.repo.UpdateGroup: %w", err)
	}

	return nil
}

// DeleteGroup deletes a group by ID.
func (uc *UseCase) DeleteGroup(ctx context.Context, id int) error {
	// Check if the group has subgroups
	hasSubgroups, err := uc.repo.HasSubgroups(ctx, id)
	if err != nil {
		return fmt.Errorf("GroupUseCase - DeleteGroup - uc.repo.HasSubgroups: %w", err)
	}

	if hasSubgroups {
		return fmt.Errorf("cannot delete group with subgroups")
	}

	err = uc.repo.DeleteGroup(ctx, id)
	if err != nil {
		return fmt.Errorf("GroupUseCase - DeleteGroup - uc.repo.DeleteGroup: %w", err)
	}

	return nil
}

// SearchGroups searches for groups by name.
func (uc *UseCase) SearchGroups(ctx context.Context, query string) ([]entity.Group, error) {
	groups, err := uc.repo.SearchGroups(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GroupUseCase - SearchGroups - uc.repo.SearchGroups: %w", err)
	}

	return groups, nil
}
