package utility

import "strings"

// NormalizePath 规范化URL路径，确保格式正确
func NormalizePath(parts ...string) string {
	return "/" + strings.TrimPrefix(strings.Join(
		append([]string{""},
			removeEmptyStrings(removeTrailingSlashes(parts))...),
		"/"), "/")
}

// removeTrailingSlashes 移除字符串末尾的斜杠
func removeTrailingSlashes(parts []string) []string {
	result := make([]string, len(parts))
	for i, part := range parts {
		result[i] = strings.TrimRight(part, "/")
	}
	return result
}

// removeEmptyStrings 移除空字符串
func removeEmptyStrings(parts []string) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
