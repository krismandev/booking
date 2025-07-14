package utils

import "fmt"

func ConvertToSliceOfMapOfString(data []map[string]interface{}) []map[string]string {
	var output []map[string]string
	for _, row := range data {
		userMap := make(map[string]string)
		for key, value := range row {
			if str, ok := value.(string); ok {
				userMap[key] = str
			} else {
				userMap[key] = fmt.Sprintf("%v", value) // Konversi ke string jika bukan string
			}
		}
		output = append(output, userMap)
	}

	return output
}

func ConvertToMapOfString(data map[string]interface{}) map[string]string {
	// var output map[string]string
	output := make(map[string]string)
	for key, value := range data {
		if str, ok := value.(string); ok {
			output[key] = str
		} else {
			if value == nil {
				value = ""
			}
			output[key] = fmt.Sprintf("%v", value) // Konversi ke string jika bukan string
		}
	}

	return output
}
