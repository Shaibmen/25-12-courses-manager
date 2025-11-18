package repo

import (
	"auth-service/internal/domain/dto"
	"auth-service/internal/domain/entity"
	"context"

	"gorm.io/gorm"
)

type SessionRepo struct {
	db *gorm.DB
}

func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

func (s *SessionRepo) CreateSession(ctx context.Context, dto *dto.SessionDTO) (*dto.SessionDTO, error) {
	if err := gorm.G[entity.Session](s.db).Create(ctx, &entity.Session{
		ID:           dto.ID,
		Username:     dto.Username,
		RefreshToken: dto.RefreshToken,
		IsRevoked:    dto.IsRevoked,
		ExpiresAt:    dto.ExpiresAt,
	}); err != nil {
		return nil, err
	}

	return dto, nil

}

func (s *SessionRepo) ReadSession(ctx context.Context, id string) (*entity.Session, error) {
	sessions, err := gorm.G[entity.Session](s.db).Where("id_session = ?", id).First(ctx)

	if err != nil {
		return nil, err
	}

	return &sessions, nil
}

func (s *SessionRepo) RevokeSession(ctx context.Context, id string) error {

	if _, err := gorm.G[entity.Session](s.db).Where("id_session = ?", id).Update(ctx, "is_revoked", "1"); err != nil {
		return err
	}

	return nil
}

func (s *SessionRepo) DeleteSession(ctx context.Context, id string) error {

	if _, err := gorm.G[entity.Session](s.db).Where("id_session = ?", id).Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (s *SessionRepo) DeleteSessionByUsername(ctx context.Context, username string) error {

	if _, err := gorm.G[entity.Session](s.db).Where("user_name = ?", username).Delete(ctx); err != nil {
		return err
	}

	return nil
}
