package service

import "context"

type DB interface {
	PingContext(ctx context.Context) error
}

type TenantValidator interface {
	TenantIDValidation(tenantID string) error
}

type PosService struct {
	db        DB
	validator TenantValidator
}

func NewPosService(db DB, validator TenantValidator) *PosService {
	return &PosService{
		db:        db,
		validator: validator,
	}
}

func (s *PosService) GetHealth(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *PosService) GetHealthByTenantID(ctx context.Context, tenantID string) error {
	if err := s.validator.TenantIDValidation(tenantID); err != nil {
		return err
	}
	return s.db.PingContext(ctx)
}