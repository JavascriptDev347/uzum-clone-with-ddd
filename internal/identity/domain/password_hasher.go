package domain

type PasswordHasher interface {
	// Hash - ochiq parolni hash'ga aylantiradi (registratsiya paytida).
	Hash(plainPassword string) (string, error)

	// Compare - hash bilan ochiq parolni solishtiradi (login paytida).
	// Mos kelmasa xato qaytaradi (xato matni algoritmga bog'liq bo'lmasligi kerak).
	Compare(hashedPassword, plainPassword string) error
}
