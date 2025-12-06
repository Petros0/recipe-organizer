package handler

// Helper functions for safe type conversion

func getString(obj map[string]interface{}, key string) string {
	val, ok := obj[key]
	if !ok {
		return ""
	}
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func getStringPtr(obj map[string]interface{}, key string) *string {
	val, ok := obj[key]
	if !ok {
		return nil
	}
	if str, ok := val.(string); ok && str != "" {
		return &str
	}
	return nil
}

func getStringArray(obj map[string]interface{}, key string) []string {
	val, ok := obj[key]
	if !ok {
		return nil
	}

	var result []string
	switch v := val.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str != "" {
				result = append(result, str)
			}
		}
	case string:
		if v != "" {
			result = append(result, v)
		}
	}

	return result
}

