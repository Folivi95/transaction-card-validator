//go:build unit

package pgateway

import (
	"errors"
	"testing"

	"github.com/matryer/is"
)

func TestPGatewayValidation(t *testing.T) {
	is := is.New(t)
	validator, _ := NewValidateNoCHD()
	t.Run("Test Validate() ok with valid messages", func(t *testing.T) {
		messages := []map[string]interface{}{
			// Message 1: instrument_params without CHD
			map[string]interface{}{
				"id":            "transaction_id",
				"charge_id":     "charge_id",
				"instrument_id": "instrument_id",
				"payment_id":    "payment_id",
				"status":        "status",
				"transaction_info": map[string]interface{}{
					"destination_business_id": "12345",
					"amount":                  1.2345,
				},
				"amount":     1.56789,
				"currency":   "EUR",
				"return_url": "https://google.com",
				"next_step": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{"json": 12345}`,
					},
				},
				"payment_method": "paysafecard",
				"instrument_params": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{
							"address": "Porto 4100",
							"phone_number": "987654321",
							"wallet_type": "virtual"
						}`,
					},
				},
				"created_at": 123456789,
				"updated_at": 123456789,
			},
			// Message 2: null instrument_params
			map[string]interface{}{
				"instrument_params": nil,
			},
			// Message 3: PAN-like string value but associated with a whitelisted field
			map[string]interface{}{
				"instrument_params": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{
							"billing_address_postal_code": "4012-8888-8888-1881"
						}`,
					},
				},
			},
			// Message 4: indented parameters without CHD
			map[string]interface{}{
				"message": map[string]interface{}{
					"parameters": map[string]interface{}{
						"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
							"body": `{
								"address": "Porto 4100",
								"phone_number": "987654321",
								"wallet_type": "virtual"
							}`,
						},
					},
				},
			},
		}
		for _, message := range messages {
			validationOk, err := validator.Validate(message)
			is.Equal(validationOk, true)
			is.Equal(err, nil)
		}
	})
	t.Run("Test Validate() fails with invalid messages", func(t *testing.T) {
		messages := []map[string]interface{}{
			// Message 1
			// instrument_params.body.cellphone contains a Pan
			map[string]interface{}{
				"id":            "transaction_id",
				"charge_id":     "charge_id",
				"instrument_id": "instrument_id",
				"payment_id":    "payment_id",
				"status":        "status",
				"transaction_info": map[string]interface{}{
					"destination_business_id": "12345",
					"amount":                  1.2345,
				},
				"amount":     1.56789,
				"currency":   "EUR",
				"return_url": "https://google.com",
				"next_step": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{"json": 12345}`,
					},
				},
				"payment_method": "paysafecard",
				"instrument_params": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{
							"address": "Porto 4100",
							"phone_number": "987654321",
							"wallet_type": "virtual",
							"cellphone": "4012-8888-8888-1881"
						}`,
					},
				},
				"created_at": 123456789,
				"updated_at": 123456789,
			},
			// Message 2
			// instrument_params.body contains blacklisted parameter (number)
			map[string]interface{}{
				"id":            "transaction_id",
				"charge_id":     "charge_id",
				"instrument_id": "instrument_id",
				"payment_id":    "payment_id",
				"status":        "status",
				"transaction_info": map[string]interface{}{
					"destination_business_id": "12345",
					"amount":                  1.2345,
				},
				"amount":     1.56789,
				"currency":   "EUR",
				"return_url": "https://google.com",
				"next_step": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{"json": 12345}`,
					},
				},
				"payment_method": "paysafecard",
				"instrument_params": map[string]interface{}{
					"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
						"body": `{
							"number": "field_value",
							"phone_number": "987654321"
						}`,
					},
				},
				"created_at": 123456789,
				"updated_at": 123456789,
			},
			// Message 3
			// An indented instrument_params.body contains a blacklisted parameter (number)
			map[string]interface{}{
				"some_field": map[string]interface{}{
					"other_field": map[string]interface{}{
						"instrument_params": map[string]interface{}{
							"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
								"body": `{
									"number": "field_value"
								}`,
							},
						},
					},
				},
			},
			// Message 4
			// An indented instrument_params.body inside a slice contains a blacklisted parameter (number)
			map[string]interface{}{
				"some_field": []map[string]interface{}{
					map[string]interface{}{
						"other_field": map[string]interface{}{
							"instrument_params": map[string]interface{}{
								"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
									"body": `{
										"number": "field_value"
									}`,
								},
							},
						},
					},
				},
			},
			// Message 5
			// indented parameters.body.random_field contains a Pan
			map[string]interface{}{
				"indented_field": map[string]interface{}{
					"instrument_params": map[string]interface{}{
						"co.saltpay.acquiring.gateway.acquiring_gateway_json_object": map[string]interface{}{
							"body": `{
								"random_field": "4012-8888-8888-1881"
							}`,
						},
					},
				},
			},
		}

		expectedErrors := []string{
			"message contains possible pan in field 'cellphone'",
			"message contains blacklisted field 'number'",
			"message contains blacklisted field 'number'",
			"message contains blacklisted field 'number'",
			"message contains possible pan in field 'random_field'",
		}

		for idx, message := range messages {
			validationOk, err := validator.Validate(message)
			is.Equal(validationOk, false)
			is.Equal(err, errors.New(expectedErrors[idx]))
		}
	})
}

func TestLuhnValidation(t *testing.T) {
	is := is.New(t)
	t.Run("Test isLuhnValid() func", func(t *testing.T) {
		// Test with valid Luhn numbers
		validLuhnNumbers := [...]string{
			"376100000000004",
			"6761000000000006",
			"5123450000000008",
			"4111111111111111",
			"4761739001012222",
		}
		for _, number := range validLuhnNumbers {
			is.True(isLuhnValid(number))
		}
		// Test with invalid Luhn numbers
		invalidLuhnNumbers := [...]string{
			"376100000000005",
			"6761000000000007",
			"5123450000000009",
			"4111111111111112",
			"4761739001012223",
		}
		for _, number := range invalidLuhnNumbers {
			is.Equal(isLuhnValid(number), false)
		}
	})
}
