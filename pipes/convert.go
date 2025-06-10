package pipes

import "strconv"

// Int converts a string to an int. If conversion fails it returns 0.
func Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// String converts an int to a string.
func String(i int) string { return strconv.Itoa(i) }

// Ustring converts a uint to a string.
func Ustring(u uint) string { return String(int(u)) }

// Bstring converts a bool to its string representation ("true" / "false").
func Bstring(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// PtIntToString converts *int to string safely.
func PtIntToString(num *int) string {
	if num == nil {
		return ""
	}
	return strconv.Itoa(*num)
}

// Uint converts a numeric string to uint (0 on error).
func Uint(num string) uint { return uint(Int(num)) }

// UintToPtInt converts uint to *int.
func UintToPtInt(num uint) *int {
	i := int(num)
	return &i
}

// GetPointedUInt converts int to *uint.
func Puint(num int) *uint {
	u := uint(num)
	return &u
}

// GetPointedInt returns a pointer to the given int.
func Pint(num int) *int { return &num }
