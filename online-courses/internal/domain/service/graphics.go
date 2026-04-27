package service

import (
	"context"
	"online-courses/internal/domain/dto"
)

type GraphicsService interface {
	CountListenersOnProgram(ctx context.Context) ([]dto.CountListenersOnProgramDTO, error)
	PopularProgramType(ctx context.Context) ([]dto.PopularProgramTypeDTO, error)
	CountListenersOnProgramAccurate(ctx context.Context) ([]dto.CountListenersOnProgramDTO, error)
	WorthProgramAccurate(ctx context.Context) ([]dto.WorthProgramDTO, error)
	AgeDiff(ctx context.Context) ([]dto.AgeDiffDTO, error)
	WhoEnrolled(ctx context.Context) ([]dto.WhoEnrolledDTO, error)
	GroupMembers(ctx context.Context) ([]dto.GroupMembersDTO, error)
	DivisionMember(ctx context.Context) ([]dto.DivisionMemberDTO, error)
}
