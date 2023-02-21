package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/storage"
	"github.com/ImpressionableRaccoon/aircirculatorServer/internal/utils"
)

func (s *Service) GetUserCompanies(ctx context.Context, user storage.User) (companies []storage.Company, err error) {
	return s.st.GetUserCompanies(ctx, user, user.IsAdmin)
}

func (s *Service) AddCompany(ctx context.Context, user storage.User, name string, offset utils.Offset) (
	company storage.Company, err error) {
	return s.st.AddCompany(ctx, user, name, offset)
}

func (s *Service) GetUserCompany(ctx context.Context, user storage.User, companyID uuid.UUID) (
	company storage.Company, err error) {
	return s.st.GetUserCompany(ctx, user, companyID)
}

func (s *Service) GetCompanyDevices(ctx context.Context, user storage.User, companyID uuid.UUID) (
	devices []storage.Device, err error) {
	return s.st.GetCompanyDevices(ctx, user, companyID)
}
