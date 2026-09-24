package service

import (
	"errors"
	"herbie-charge-api/internal/charge/dto"
	"log"
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
	log.Printf(
		"charge requested: amount=%d",
		req.Amount,
	)

	if req.IDToken == "" {
		return errors.New(
			"idToken is required",
		)
	}

	profile, err :=
		s.lineService.VerifyIDToken(
			req.IDToken,
		)

	if err != nil {
		log.Printf(
			"LINE verify failed: %v",
			err,
		)

		return err
	}

	log.Println(
		"LINE user verified",
	)

	if err :=
		s.lineService.SendChargeSuccess(
			profile.Sub,
			req.Amount,
		); err != nil {

		log.Printf(
			"LINE push failed: %v",
			err,
		)

		return err
	}

	log.Printf(
		"charge success: amount=%d",
		req.Amount,
	)

	return nil
}
