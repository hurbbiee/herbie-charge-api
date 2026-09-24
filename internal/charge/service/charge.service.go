package service

import (
	"errors"
	"herbie-charge-api/internal/charge/dto"
)

type ChargeService struct {
	lineService *LineService
}

func NewChargeService(
	lineService *LineService,
) *ChargeService {
	return &ChargeService{
		lineService: lineService,
	}
}

func (s *ChargeService) Charge(
	req dto.ChargeRequest,
) error {
	if req.IDToken == "" {
		return errors.New(
			"idToken is required",
		)
	}

	allowedAmounts := map[int]bool{
		100: true,
		200: true,
		300: true,
		400: true,
		500: true,
		600: true,
	}

	if !allowedAmounts[req.Amount] {
		return errors.New(
			"invalid amount",
		)
	}

	profile, err :=
		s.lineService.VerifyIDToken(
			req.IDToken,
		)

	if err != nil {
		return err
	}

	/*
		ตรงนี้ภายหลัง:

		1. Create transaction
		2. Update credit
		3. Save PostgreSQL
	*/

	err =
		s.lineService.SendChargeSuccess(
			profile.Sub,
			req.Amount,
		)

	if err != nil {
		return err
	}

	return nil
}
