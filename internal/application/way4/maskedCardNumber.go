// This validation is sued to ensure that all possible card numbers in Way4 DOC table data
// are correctly masked.

package way4validation

import (
	"errors"
	"regexp"
)

type Way4MaskedCardNumber struct {
	CardHoldKeys     []string
	MaskRegexPattern string
}

func NewWay4MaskedCardNumber(cardHoldKeys []string, maskRegexPattern string) (*Way4MaskedCardNumber, error) {
	if cardHoldKeys == nil || maskRegexPattern == "" {
		return nil, errors.New("missingRequirement: Way4MaskedCardNumber validation requires W4_MASKED_CHD_KEYS and W4_MASKED_REGEX_PATTERN environment variables")
	}

	return &Way4MaskedCardNumber{
		CardHoldKeys:     cardHoldKeys,
		MaskRegexPattern: maskRegexPattern,
	}, nil
}

func (w *Way4MaskedCardNumber) Validate(msgDecoded interface{}) (bool, error) {
	// Search for possible card numbers in every CardHoldKey
	// and verify if they are correctly masked according to given regex pattern.
	re := regexp.MustCompile(w.MaskRegexPattern)
	result := true

	for _, cardHoldKey := range w.CardHoldKeys {
		// Validate
		if before := msgDecoded.(map[string]interface{})["co.saltpay.acquiring.way4_curated_transaction"].(map[string]interface{})["before"]; before != nil {
			cardHoldData := before.(map[string]interface{})["co.saltpay.acquiring.way4_doc"].(map[string]interface{})[cardHoldKey]
			if cardHoldData != nil {
				element := cardHoldData.(map[string]interface{})["string"].(string)
				// If length match with card data, validate that its masked
				if len(element) >= 13 && len(element) <= 19 {
					result = result && re.MatchString(element)
				}
			}
		}

		if after := msgDecoded.(map[string]interface{})["co.saltpay.acquiring.way4_curated_transaction"].(map[string]interface{})["after"]; after != nil {
			cardHoldData := after.(map[string]interface{})["co.saltpay.acquiring.way4_doc_tokenized"].(map[string]interface{})[cardHoldKey]
			if cardHoldData != nil {
				element := cardHoldData.(map[string]interface{})["string"].(string)
				// If length match with card data, validate that its masked
				if len(element) >= 13 && len(element) <= 19 {
					result = result && re.MatchString(element)
				}
			}
		}
	}

	if result {
		return true, nil
	} else {
		return false, errors.New("failed validation")
	}
}
