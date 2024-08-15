package util

// StringP converts string to *string
func StringP(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// IntP converts int to *int
func IntP(i int) *int {
	return &i
}

// Float32P converts float32 to *float32
func Float32P(f float32) *float32 {
	return &f
}

// BoolP converts bool to *bool
func BoolP(b bool) *bool {
	return &b
}

// Float64P converts float64 to *float64
func Float64P(f float64) *float64 {
	return &f
}

// Int64P converts int64 to *int64
func Int64P(i int64) *int64 {
	return &i
}

// Uint64P converts uint64 to *uint64
func Uint64P(i uint64) *uint64 {
	return &i
}

// Uint32P converts uint32 to *uint32
func Uint32P(i uint32) *uint32 {
	return &i
}

// Uint16P converts uint16 to *uint16
func Uint16P(i uint16) *uint16 {
	return &i
}

// Uint8P converts uint8 to *uint8
func Uint8P(i uint8) *uint8 {
	return &i
}
