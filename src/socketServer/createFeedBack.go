package socketServer

import (
	"encoding/json"
	"fmt"
)

func feedBackOk(data interface{}) string {
	if str, ok := data.(string); ok {
		if len(str) == 0 {
			return `{"status": 1}`
		} else {
			return fmt.Sprintf(`{"status": 1, "data": "%v"}`, str)
		}
	}
	// todo: next types

	if dataStd, err := json.Marshal(data); err == nil {
		return fmt.Sprintf(
			`{"status": 1, "data": %v}`,
			string(dataStd),
		)
	} else {
		return `{"status": 1}`
	}
}
func feedBackErr(message string) string {
	return fmt.Sprintf(`{"status": 1, "message": "%v"}`, message)
}
