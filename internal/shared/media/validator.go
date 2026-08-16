package media

import "fmt"

// ValidationRules - fayl validatsiyasi uchun qoidalar (hajm, format).
type ValidationRules struct {
	MaxSizeBytes     int64
	AllowedMimeTypes []string
}

// DefaultImageRules - rasm fayllar uchun standart qoidalar (5MB, jpeg/png/webp).
func DefaultImageRules() ValidationRules {
	return ValidationRules{
		MaxSizeBytes:     5 * 1024 * 1024, // 5MB
		AllowedMimeTypes: []string{"image/jpeg", "image/png", "image/webp"},
	}
}

func Validate(input UploadInput, rules ValidationRules) error {
	if len(input.Data) == 0 {
		return ErrEmptyFile
	}

	if int64(len(input.Data)) > rules.MaxSizeBytes {
		return fmt.Errorf("%w: max %d bytes", ErrFileTooLarge, rules.MaxSizeBytes)
	}

	allowed := false
	for _, t := range rules.AllowedMimeTypes {
		if t == input.ContentType {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("%w: %s", ErrUnsupportedType, input.ContentType)
	}

	return nil
}
