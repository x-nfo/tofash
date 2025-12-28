package conv

import "strconv"

func StringToInt64(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	return strconv.ParseInt(s, 10, 64)
}

func StringToUint(s string) (uint, error) {
	if s == "" {
		return 0, nil
	}
	val, err := strconv.ParseUint(s, 10, 64)
	return uint(val), err
}
