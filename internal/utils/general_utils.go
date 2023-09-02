package utils

import "encoding/json"

func IsEmpty(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return v == 0
	case int64:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case float32:
		return v == 0.0
	case float64:
		return v == 0.0
	case string:
		return v == ""
	case []int:
		return len(v) == 0
	case map[string]int:
		return len(v) == 0
	case bool:
		return !v
	case nil:
		return true
	default:
		return false
	}
}

// Utility to convert struct data types to maps
func StructToMap(obj interface{}) (newMap map[string]interface{}, err error) {
	data, err := json.Marshal(obj) // Convert to a json string

	if err != nil {
		return
	}

	err = json.Unmarshal(data, &newMap) // Convert to a map
	return
}
