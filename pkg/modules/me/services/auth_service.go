package services

import (
	dto "github.com/michaelchandrag/botfood-go/pkg/modules/me/dto"
	entities "github.com/michaelchandrag/botfood-go/pkg/modules/me/entities"
	user_repository "github.com/michaelchandrag/botfood-go/pkg/modules/me/repositories/user"
	user_branch_repository "github.com/michaelchandrag/botfood-go/pkg/modules/me/repositories/user_branch"
	user_brand_repository "github.com/michaelchandrag/botfood-go/pkg/modules/me/repositories/user_brand"
)

func (s *service) FormatAuthFromMiddleware(payload dto.MeAuthRequestPayload) (response dto.MeAuthResponse) {
	middlewareAuthBrand := payload.AuthBrand
	authBrand := entities.Brand{
		ID:   middlewareAuthBrand.ID,
		Slug: middlewareAuthBrand.Slug,
		Name: middlewareAuthBrand.Name,
	}
	response.Auth.Brand = authBrand
	if middlewareAuthBrand.UserBrandThumbnail != nil {
		authUserBrandThumbnail := middlewareAuthBrand.UserBrandThumbnail
		userBrandRepository := user_brand_repository.NewRepository(s.db)
		userBrandFilter := user_brand_repository.Filter{
			UserID:   &authUserBrandThumbnail.UserID,
			IsActive: true,
		}
		userBrands, err := userBrandRepository.FindAll(userBrandFilter)
		if err != nil {
			response.Errors.AddHTTPError(500, err)
			return response
		}
		var userBrandThumbnail entities.UserBrand
		for _, userBrand := range userBrands {
			if userBrand.IsThumbnail == 1 {
				userBrandThumbnail = userBrand
			}
		}
		response.Auth.UserBrands = &userBrands
		response.Auth.UserBrandThumbnail = &userBrandThumbnail

		userBranchRepository := user_branch_repository.NewRepository(s.db)
		userBranchFilter := user_branch_repository.Filter{
			UserID:   &authUserBrandThumbnail.UserID,
			IsActive: true,
			BrandID:  &userBrandThumbnail.BrandID,
		}
		userBranchs, err := userBranchRepository.FindAll(userBranchFilter)
		if err != nil {
			response.Errors.AddHTTPError(500, err)
			return response
		}
		var branchIDs []int
		for _, userBranch := range userBranchs {
			branchIDs = append(branchIDs, int(userBranch.BranchID))
		}
		response.Auth.BranchIDs = branchIDs
		response.Auth.UserBranchs = &userBranchs
		response.Auth.IsMaster = true

		userRepository := user_repository.NewRepository(s.db)
		userFilter := user_repository.Filter{
			ID: &userBrandThumbnail.UserID,
		}
		user, err := userRepository.FindOne(userFilter)
		if err != nil {
			response.Errors.AddHTTPError(500, err)
			return response
		}
		response.Auth.User = &user

	} else {
		response.Auth.IsMaster = false
	}
	return response
}
