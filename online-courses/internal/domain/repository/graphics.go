package repository

import (
	"context"
	"online-courses/internal/domain/entity"
)

type GraphicsRepository interface {
	CountListenersOnProgram(ctx context.Context) ([]entity.CountListenersOnProgramStruct, error)
	PopularProgramType(ctx context.Context) ([]entity.PopularProgramTypeStruct, error)
	CountListenersOnProgramAccurate(ctx context.Context) ([]entity.CountListenersOnProgramStruct, error)
	WorthProgramAccurate(ctx context.Context) ([]entity.WorthProgramAccurateStruct, error)
	AgeDiff(ctx context.Context) ([]entity.AgeDiffStruct, error)
	WhoEnrolled(ctx context.Context) ([]entity.WhoEnrolledStruct, error)
	GroupMembers(ctx context.Context) ([]entity.GroupMembersStruct, error)
	DivisionMember(ctx context.Context) ([]entity.DivisionMemberStruct, error)
}
