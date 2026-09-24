package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"herbie-charge-api/internal/charge/dto"
	"herbie-charge-api/internal/config"
	"io"
	"net/http"
	"net/url"
	"time"
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
	result dto.ChargeResult,
) error {
	createdAtText := result.CreatedAt

	if parsedTime, err :=
		time.Parse(
			time.RFC3339,
			result.CreatedAt,
		); err == nil {

		createdAtText =
			parsedTime.Format(
				"02/01/2006 15:04",
			)
	}

	message := map[string]any{
		"type":    "flex",
		"altText": "เติมเครดิตสำเร็จ",
		"contents": map[string]any{
			"type": "bubble",

			// Header
			"header": map[string]any{
				"type":            "box",
				"layout":          "vertical",
				"backgroundColor": "#10B981",
				"paddingAll":      "20px",

				"contents": []any{
					map[string]any{
						"type":   "text",
						"text":   "เติมเครดิตสำเร็จ",
						"size":   "lg",
						"weight": "bold",
						"color":  "#FFFFFF",
					},

					map[string]any{
						"type":   "text",
						"text":   "รายการของคุณสำเร็จแล้ว",
						"size":   "sm",
						"color":  "#D1FAE5",
						"margin": "sm",
					},
				},
			},

			// Body
			"body": map[string]any{
				"type":       "box",
				"layout":     "vertical",
				"spacing":    "md",
				"paddingAll": "20px",

				"contents": []any{
					// Amount label
					map[string]any{
						"type":  "text",
						"text":  "จำนวนเงินที่เติม",
						"size":  "sm",
						"color": "#6B7280",
					},

					// Amount value
					map[string]any{
						"type": "text",
						"text": fmt.Sprintf(
							"%d.00 บาท",
							result.Amount,
						),
						"size":   "xxl",
						"weight": "bold",
						"color":  "#111827",
					},

					// Separator
					map[string]any{
						"type":   "separator",
						"margin": "lg",
						"color":  "#E5E7EB",
					},

					// Balance
					map[string]any{
						"type":   "box",
						"layout": "horizontal",
						"margin": "lg",

						"contents": []any{
							map[string]any{
								"type":  "text",
								"text":  "เครดิตคงเหลือ",
								"size":  "sm",
								"color": "#6B7280",
								"flex":  1,
							},

							map[string]any{
								"type": "text",
								"text": fmt.Sprintf(
									"%d.00 บาท",
									result.Balance,
								),
								"size":   "sm",
								"weight": "bold",
								"align":  "end",
								"color":  "#059669",
								"flex":   2,
							},
						},
					},

					// Transaction ID
					map[string]any{
						"type":   "box",
						"layout": "horizontal",

						"contents": []any{
							map[string]any{
								"type":  "text",
								"text":  "เลขที่รายการ",
								"size":  "sm",
								"color": "#6B7280",
								"flex":  1,
							},

							map[string]any{
								"type":  "text",
								"text":  result.TransactionID,
								"size":  "xs",
								"align": "end",
								"color": "#374151",
								"wrap":  true,
								"flex":  2,
							},
						},
					},

					// Created At
					map[string]any{
						"type":   "box",
						"layout": "horizontal",

						"contents": []any{
							map[string]any{
								"type":  "text",
								"text":  "วันที่ทำรายการ",
								"size":  "sm",
								"color": "#6B7280",
								"flex":  1,
							},

							map[string]any{
								"type":  "text",
								"text":  createdAtText,
								"size":  "sm",
								"align": "end",
								"color": "#374151",
								"flex":  2,
							},
						},
					},

					// Status
					map[string]any{
						"type":   "box",
						"layout": "horizontal",

						"contents": []any{
							map[string]any{
								"type":  "text",
								"text":  "สถานะ",
								"size":  "sm",
								"color": "#6B7280",
								"flex":  1,
							},

							map[string]any{
								"type":   "text",
								"text":   "สำเร็จ",
								"size":   "sm",
								"weight": "bold",
								"align":  "end",
								"color":  "#10B981",
								"flex":   2,
							},
						},
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

	body, err :=
		json.Marshal(payload)

	if err != nil {
		return fmt.Errorf(
			"marshal LINE message failed: %w",
			err,
		)
	}

	req, err :=
		http.NewRequest(
			http.MethodPost,
			"https://api.line.me/v2/bot/message/push",
			bytes.NewBuffer(body),
		)

	if err != nil {
		return fmt.Errorf(
			"create LINE request failed: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+s.config.LineChannelAccessToken,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err :=
		client.Do(req)

	if err != nil {
		return fmt.Errorf(
			"send LINE message failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		responseBody, readErr :=
			io.ReadAll(resp.Body)

		if readErr != nil {
			return fmt.Errorf(
				"LINE push failed with status %d",
				resp.StatusCode,
			)
		}

		return fmt.Errorf(
			"LINE push failed status=%d response=%s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}
