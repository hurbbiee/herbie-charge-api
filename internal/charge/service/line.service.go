package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"herbie-charge-api/internal/config"
	"io"
	"net/http"
	"net/url"
)

type LineService struct {
	config config.Config
}

func NewLineService(
	cfg config.Config,
) *LineService {
	return &LineService{
		config: cfg,
	}
}

type LineVerifyResponse struct {
	Sub     string `json:"sub"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (s *LineService) VerifyIDToken(
	idToken string,
) (*LineVerifyResponse, error) {
	form := url.Values{}

	form.Set(
		"id_token",
		idToken,
	)

	form.Set(
		"client_id",
		s.config.LineLoginChannelID,
	)

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.line.me/oauth2/v2.1/verify",
		bytes.NewBufferString(
			form.Encode(),
		),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/x-www-form-urlencoded",
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(
		resp.Body,
	)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"LINE verify failed: %s",
			string(body),
		)
	}

	var result LineVerifyResponse

	if err := json.Unmarshal(
		body,
		&result,
	); err != nil {
		return nil, err
	}

	if result.Sub == "" {
		return nil, errors.New(
			"LINE user ID not found",
		)
	}

	return &result, nil
}

func (s *LineService) SendChargeSuccess(
	userID string,
	amount int,
) error {
	message := map[string]any{
		"type":    "flex",
		"altText": "เติมเครดิตสำเร็จ",
		"contents": map[string]any{
			"type": "bubble",

			"header": map[string]any{
				"type":            "box",
				"layout":          "vertical",
				"backgroundColor": "#10B981",
				"paddingAll":      "20px",

				"contents": []any{
					map[string]any{
						"type":   "text",
						"text":   "เติมเครดิตสำเร็จ",
						"weight": "bold",
						"size":   "lg",
						"color":  "#FFFFFF",
					},
				},
			},

			"body": map[string]any{
				"type":       "box",
				"layout":     "vertical",
				"spacing":    "lg",
				"paddingAll": "20px",

				"contents": []any{
					map[string]any{
						"type":   "text",
						"text":   "จำนวนเงินที่เติม",
						"size":   "sm",
						"color":  "#777777",
					},

					map[string]any{
						"type":   "text",
						"text":   fmt.Sprintf("%d.00 บาท", amount),
						"size":   "xxl",
						"weight": "bold",
						"color":  "#111111",
					},

					map[string]any{
						"type":   "separator",
						"margin": "lg",
					},

					map[string]any{
						"type":   "text",
						"text":   "สถานะ",
						"size":   "sm",
						"color":  "#777777",
						"margin": "lg",
					},

					map[string]any{
						"type":   "text",
						"text":   "ทำรายการสำเร็จ",
						"size":   "md",
						"weight": "bold",
						"color":  "#10B981",
					},
				},
			},
		},
	}

	payload := map[string]any{
		"to": userID,
		"messages": []any{
			message,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.line.me/v2/bot/message/push",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+s.config.LineChannelAccessToken,
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		responseBody, _ :=
			io.ReadAll(resp.Body)

		return fmt.Errorf(
			"LINE push message failed: %s",
			string(responseBody),
		)
	}

	return nil
}
