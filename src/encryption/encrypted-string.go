package encryption

import (
	"encoding/json"
	"fmt"
)

type EncryptedString struct {
	Value string
}

func (es *EncryptedString) MarshalJSON() ([]byte, error) {
	revealedValue := es.String()
	return json.Marshal(revealedValue)
}

func (es *EncryptedString) UnmarshalJSON(value []byte) error {
	if string(value) == "null" {
		return nil
	}

	if string(value) == "" {
		*es = EncryptedString{
			Value: "",
		}

		return nil
	}

	var strVal string
	if err := json.Unmarshal(value, &strVal); err != nil {
		return fmt.Errorf("expected a string for encrypted data: %w", err)
	}

	*es = EncryptedString{
		Value: strVal,
	}

	return nil
}

func (es *EncryptedString) String() string {
	if es == nil || es.Value == "" {
		return ""
	}

	return "***"
}
