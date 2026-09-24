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
	message := map[string]interface{}{
		"type":    "flex",
		"altText": "เติมเครดิตสำเร็จ",
		"contents": map[string]interface{}{
			"type": "bubble",

			"header": map[string]interface{}{
				"type":            "box",
				"layout":          "vertical",
				"backgroundColor": "#10B981",
				"paddingAll":      "20px",

				"contents": []interface{}{
					map[string]interface{}{
						"type":   "text",
						"text":   "เติมเครดิตสำเร็จ",
						"weight": "bold",
						"size":   "lg",
						"color":  "#FFFFFF",
					},
				},
			},

			"body": map[string]interface{}{
				"type":    "box",
				"layout":  "vertical",
				"spacing": "md",

				"contents": []interface{}{
					map[string]interface{}{
						"type":  "text",
						"text":  "จำนวนเงิน",
						"color": "#777777",
						"size":  "sm",
					},

					map[string]interface{}{
						"type":   "text",
						"text":   fmt.Sprintf("%d.00 บาท", amount),
						"weight": "bold",
						"size":   "xxl",
					},

					map[string]interface{}{
						"type": "separator",
					},

					map[string]interface{}{
						"type":   "text",
						"text":   "ทำรายการสำเร็จ",
						"color":  "#10B981",
						"weight": "bold",
					},
				},
			},
		},
	}

	payload := map[string]interface{}{
		"to": userID,

		"messages": []interface{}{
			message,
		},
	}

	body, err := json.Marshal(
		payload,
	)

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
		"Bearer "+
			s.config.LineChannelAccessToken,
	)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		body, _ := io.ReadAll(
			resp.Body,
		)

		return fmt.Errorf(
			"LINE push message failed: %s",
			string(body),
		)
	}

	return nil
}
