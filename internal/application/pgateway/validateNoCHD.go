// This validation ensures that no CHD is present in Payment Gateway messages.
// User provided payment details (address, phone_number, PAN, CVC, etc) are present in 'instrument_params' and
// 'parameters' fields.
// Despite these fields being sanitized before the messages are published, this validation acts as a double precaution.
// The validator function receives a decoded Kafka message, iterates through it and checks that no sensitive parameter
// goes in 'instrument_params' and 'parameters' fields.

package pgateway

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type validateNoCHD struct {
	WhiteSpaceAndHyphensRegexPattern *regexp.Regexp
	DigitRegexPattern                *regexp.Regexp
}

func NewValidateNoCHD() (*validateNoCHD, error) {
	return &validateNoCHD{
		WhiteSpaceAndHyphensRegexPattern: regexp.MustCompile(`[-\s]`),
		DigitRegexPattern:                regexp.MustCompile(`\d+`),
	}, nil
}

func (w *validateNoCHD) Validate(msgDecoded interface{}) (bool, error) {
	// 1. Check each entry of 'instrument_param' and 'parameters' fields for PANs by:
	//    a) stripping all spaces and hyphens
	//    b) extracting all digits
	//    c) check against the Luhn algorithm
	//    * Except for the following fields which may yield false positives
	FieldsWhitelist := map[string]struct{}{
		"billing_address_postal_code":          {},
		"card_present_application_cryptogram":  {},
		"card_present_dedicated_file_name":     {},
		"card_present_form_factor_indicator":   {},
		"card_present_icc_dynamic_number":      {},
		"card_present_issuer_application_data": {},
		"card_present_pin_block_ksn":           {},
		"card_present_track1_ksn":              {},
		"card_present_track2_ksn":              {},
		"customer_id":                          {},
		"destination_account_id":               {},
		"dcc_data_id":                          {},
		"iban":                                 {},
		"identification_doc_id":                {},
		"instrument":                           {},
		"mobile":                               {},
		"order_id":                             {},
		"phone":                                {},
		"session_id":                           {},
		"shipping_address_postal_code":         {},
	}

	listInstrumentParams := getInstrumentParams(msgDecoded.(map[string]interface{}))
	for _, instrumentParams := range listInstrumentParams { // Iterate list of instrument params
		for instrumentParamName, instrumentParamValue := range instrumentParams { // Iterate keys and values of instrument params
			// Don't validate fields in whitelist (may yield false positives)
			if _, ok := FieldsWhitelist[instrumentParamName]; ok {
				continue
			}

			// Convert value to string if it isn't already a string
			instrumentParamValueStr := fmt.Sprint(instrumentParamValue)

			// Remove white spaces and hyphens from original value
			instrumentParamValueStr = w.WhiteSpaceAndHyphensRegexPattern.ReplaceAllString(instrumentParamValueStr, "")

			// Get all digit sequences
			digitSequences := w.DigitRegexPattern.FindAllString(instrumentParamValueStr, -1)

			// Iterate digit sequences
			for _, digitSequence := range digitSequences {
				// if digit sequence satisfies Luhn algorithm (may be a PAN), validation fails.
				if len(digitSequence) > 13 && isLuhnValid(digitSequence) {
					return false, fmt.Errorf("message contains possible pan in field '%s'", instrumentParamName)
				}
			}
		}
	}

	// 2. Check each key of 'instrument_params' to validate they don’t match any sensitive field (blacklisted)
	FieldsBlacklist := map[string]struct{}{
		"card_present_pin_block_encrypted": {},
		"card_present_track1":              {},
		"card_present_track1_encrypted":    {},
		"card_present_track2":              {},
		"card_present_track2_encrypted":    {},
		"cvc":                              {},
		"discretionary_data":               {},
		"number":                           {},
		"pin":                              {},
		"terminal_data_bdk":                {},
		"terminal_pin_bdk":                 {},
	}

	for _, instrumentParams := range listInstrumentParams { // Iterate list of instrument params
		for instrumentParamName := range instrumentParams { // Iterate keys and values of instrument params
			// If instrument param name is blacklisted (because it represents CHD), validation fails.
			if _, ok := FieldsBlacklist[instrumentParamName]; ok {
				return false, fmt.Errorf("message contains blacklisted field '%s'", instrumentParamName)
			}
		}
	}

	// No validation errors if code execution reached here
	return true, nil
}

// Given msg (a decoded Kafka message), this function returns a list of all instrument parameters present.
// This function works in a recursive manner, navigating through the tree-like structure of msg.
func getInstrumentParams(msg map[string]interface{}) []map[string]interface{} {
	var s []map[string]interface{}

	for k, v := range msg {
		// Base case: the key that is being iterated is 'instrument_params' or 'parameters'
		if (k == "instrument_params" || k == "parameters") && v != nil {
			// The value content is encoded in JSON and must be decoded
			decoded := v.(map[string]interface{})["co.saltpay.acquiring.gateway.acquiring_gateway_json_object"]
			jsonContent := decoded.(map[string]interface{})["body"].(string)
			var jsonDecoded map[string]interface{}
			if err := json.Unmarshal([]byte(jsonContent), &jsonDecoded); err == nil {
				s = append(s, jsonDecoded)
			}
		} else {
			// If value is a dictionary, iterate it
			if v, ok := v.(map[string]interface{}); ok { // Iterate over value (dictionary)
				r := getInstrumentParams(v)
				s = append(s, r...)
			}
			// If value is a list of dictionaries, iterate it
			if v, ok := v.([]map[string]interface{}); ok { // Iterate over value (list of dictionaries)
				for _, elem := range v {
					r := getInstrumentParams(elem)
					s = append(s, r...)
				}
			}
		}
	}

	return s
}

// Adapted from https://rosettacode.org/wiki/Luhn_test_of_credit_card_numbers#Go
func isLuhnValid(s string) bool {
	t := [...]int{0, 2, 4, 6, 8, 1, 3, 5, 7, 9}

	odd := len(s) & 1
	var sum int
	for i, c := range s {
		if i&1 == odd {
			sum += t[c-'0']
		} else {
			sum += int(c - '0')
		}
	}
	return sum%10 == 0
}
