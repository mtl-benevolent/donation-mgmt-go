package apperrors

import "fmt"

func formatExtras(extras map[string]any) string {
	str := ""
	for key, value := range extras {
		str += fmt.Sprintf("%s=%v, ", key, value)
	}

	if str != "" {
		str = fmt.Sprintf("(%s)", extras)
	}

	return str
}

func formatID(id string) string {
	if id != "" {
		return fmt.Sprintf("with ID %s", id)
	}

	return ""
}
