package character

import (
	"strings"
)

/**
 * 字符串拼接
 */
func Joint(str ...string) string {
	var buf strings.Builder
	for i := 0; i < len(str); i++ {
		buf.WriteString(str[i])
	}
	return buf.String()
}
