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
		"[charge] request amount=%d",
		req.Amount,
	)

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
		log.Printf(
			"[charge] LINE verify failed: %v",
			err,
		)

		return err
	}

	log.Println(
		"[charge] LINE user verified",
	)

	if err :=
		s.lineService.SendChargeSuccess(
			profile.Sub,
			req.Amount,
		); err != nil {

		log.Printf(
			"[charge] LINE push failed: %v",
			err,
		)

		return err
	}

	log.Printf(
		"[charge] success amount=%d",
		req.Amount,
	)

	return nil
}

