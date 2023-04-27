package utils

import "strings"

func FindKV(s string) ([3]string, bool) {
	if !strings.Contains(s, " ") {
		return [3]string{}, false
	}
	var cmd string
	var key string
	var value string
	result := [3]string{cmd, key, value}
	i := 0
	before := make(map[string]string, 1)
	for k, v := range s {

		if before["before"] == " " {
			before["before"] = string(v)
			if before["before"] != " " {
				i++
				if (result[0] == "忘掉" || result[0] == "查看") && i > 0 {
					result[i] += string(s[k:])
					break
				}
				if i > 1 {
					result[i] += string(s[k:])
					break
				}

			}
			continue
		}

		result[i] += before["before"]
		before["before"] = string(v)
	}
	return result, true
}
