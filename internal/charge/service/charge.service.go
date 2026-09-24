package service

import (
	"errors"
	"fmt"
	"log"
	"time"

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
) (*dto.ChargeResult, error) {
	log.Printf(
		"[charge] request amount=%d",
		req.Amount,
	)

	if req.IDToken == "" {
		return nil, errors.New(
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
		return nil, errors.New(
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

		return nil, err
	}

	log.Println(
		"[charge] LINE user verified",
	)

	/*
		MOCK DATABASE

		ตอนนี้สมมติว่าผู้ใช้มีเครดิตเดิม 100 บาท

		ของจริงภายหลังจะเป็น:

		currentBalance :=
		    walletRepository.GetBalance(profile.Sub)

		newBalance :=
		    currentBalance + req.Amount

		walletRepository.Update(...)
	*/

	currentBalance := 100

	newBalance :=
		currentBalance + req.Amount

	now := time.Now().UTC()

	result := dto.ChargeResult{
		TransactionID: fmt.Sprintf(
			"TXN-%d",
			now.UnixMilli(),
		),
		Amount:    req.Amount,
		Balance:   newBalance,
		Status:    "success",
		CreatedAt: now.Format(time.RFC3339),
	}

	if err :=
		s.lineService.SendChargeSuccess(
			profile.Sub,
			result,
		); err != nil {

		log.Printf(
			"[charge] LINE push failed: %v",
			err,
		)

		return nil, err
	}

	log.Printf(
		"[charge] success transaction=%s amount=%d balance=%d",
		result.TransactionID,
		result.Amount,
		result.Balance,
	)

	return &result, nil
}