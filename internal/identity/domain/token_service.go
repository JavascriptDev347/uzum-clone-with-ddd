package domain

// TokenService - JWT (yoki boshqa token turi) yaratish va tekshirish uchun port.
type TokenService interface {
	// GenerateAccessToken - qisqa muddatli token (masalan 15 daqiqa), har bir
	// so'rovda foydalaniladi. Role ham ichiga qo'shiladi (RBAC uchun).
	GenerateAccessToken(userID string, role Role) (string, error)

	// GenerateRefreshToken - uzoq muddatli token (masalan 7 kun),
	// faqat yangi access token olish uchun ishlatiladi.
	GenerateRefreshToken(userID string) (string, error)

	// ValidateToken - tokenni tekshiradi va undan userID'ni chiqarib beradi.
	// Token yaroqsiz/muddati o'tgan bo'lsa xato qaytaradi.
	ValidateToken(token string) (userID string, role Role, err error)
}
