# Uzum Clone API — Frontend uchun qo'llanma

Bu hujjat backendda hozircha tayyor bo'lgan barcha endpointlarni tasvirlaydi. DDD arxitekturasida qurilgan (auth — `identity`, kategoriya/mahsulot — `catalog` bounded context).

> Swagger UI: `http://localhost:8080/swagger/index.html` (backend ishga tushirilgandan keyin shu yerdan interaktiv sinab ko'rish mumkin)

---

## 1. Umumiy ma'lumotlar

| | |
|---|---|
| **Base URL** | `http://localhost:8080/api/v1` |
| **Content-Type (JSON so'rovlar)** | `application/json` |
| **Content-Type (fayl yuklash)** | `multipart/form-data` |
| **Auth** | JWT Bearer token (`Authorization: Bearer <access_token>`) |
| **CORS** | Hozircha dev rejimda barcha originlarga ochiq (`*`) |

### Javob formati (Envelope)

**Barcha** endpointlar bir xil "envelope" ko'rinishida javob qaytaradi:

```json
{
  "data": { },
  "error": "",
  "message": ""
}
```

- Muvaffaqiyatli bo'lsa → `data` to'ldiriladi, `error` bo'lmaydi.
- Xatolik bo'lsa → `error` maydonida matn bo'ladi, `data` bo'lmaydi (HTTP status kodiga qarang, pastda jadval bor).

> Eslatma: `data`, `error`, `message` maydonlari `omitempty` — ya'ni bo'sh bo'lsa JSON'da umuman ko'rinmaydi.

### Xatolik status kodlari

| Status | Ma'nosi |
|---|---|
| `400` | Noto'g'ri so'rov / validatsiya xatosi (masalan bo'sh nom, manfiy narx, email formati noto'g'ri) |
| `401` | Token yo'q / noto'g'ri / muddati o'tgan, yoki login/parol xato |
| `403` | Token bor, lekin huquq yetarli emas (admin bo'lmagan user admin endpointga murojaat qilsa) |
| `404` | Resurs topilmadi (masalan mavjud bo'lmagan category ID) |
| `409` | Conflict — email band, yoki bir xil nomli category allaqachon mavjud |
| `500` | Server ichki xatosi |

---

## 2. Autentifikatsiya (Auth)

Barcha auth endpointlar prefiksi: **`/api/v1/auth`**

Yangi ro'yxatdan o'tgan foydalanuvchi avtomatik `customer` rolida yaratiladi. `admin` rolini faqat backend/DB orqali qo'lda berish mumkin (frontendda buni tanlash imkoniyati yo'q).

### 2.1 Ro'yxatdan o'tish

```
POST /api/v1/auth/register
```

**Body:**
```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

- Parol uchun minimal uzunlik yoki murakkablik tekshiruvi backendda **yo'q** — istalgan bo'sh bo'lmagan string qabul qilinadi (frontendda o'zingiz validatsiya qo'shishni tavsiya qilamiz).
- Email formati regex bilan tekshiriladi va kichik harflarga normalize qilinadi.

**Muvaffaqiyatli javob — `201 Created`:**
```json
{
  "data": {
    "user_id": "uuid",
    "email": "user@example.com"
  }
}
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | email formati noto'g'ri yoki JSON body noto'g'ri |
| 409 | bu email bilan foydalanuvchi allaqachon mavjud |

---

### 2.2 Tizimga kirish

```
POST /api/v1/auth/login
```

**Body:**
```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

**Muvaffaqiyatli javob — `200 OK`:**
```json
{
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi..."
  }
}
```

- `access_token` — TTL: `.env` dagi `JWT_ACCESS_TTL` (default 15m)
- `refresh_token` — TTL: `.env` dagi `JWT_REFRESH_TTL` (default 168h / 7 kun)

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | body noto'g'ri |
| 401 | email yoki parol noto'g'ri |

---

### 2.3 Tokenni yangilash

```
POST /api/v1/auth/refresh
```

`access_token` muddati tugaganda, foydalanuvchini qayta login qildirmasdan yangi juftlik olish uchun ishlatiladi.

**Body:**
```json
{
  "refresh_token": "eyJhbGciOi..."
}
```

**Muvaffaqiyatli javob — `200 OK`:** (login bilan bir xil formatda, yangi `access_token` + `refresh_token`)
```json
{
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi..."
  }
}
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | body noto'g'ri |
| 401 | refresh token yaroqsiz yoki muddati o'tgan → foydalanuvchini qayta login sahifasiga yo'naltiring |

---

### 2.4 Joriy foydalanuvchini olish

```
GET /api/v1/auth/me
```

🔒 **Autentifikatsiya talab qilinadi** — `Authorization: Bearer <access_token>`

**Muvaffaqiyatli javob — `200 OK`:**
```json
{
  "data": {
    "user_id": "uuid",
    "email": "user@example.com",
    "role": "customer"
  }
}
```

- `role` — `"customer"` yoki `"admin"`. Frontendda admin panelni ko'rsatish/yashirish uchun shu maydondan foydalaning.

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 401 | token yo'q, noto'g'ri yoki muddati o'tgan |

---

## 3. Kategoriyalar (Categories)

Prefiks: **`/api/v1/categories`** (bu endpointlar `/api/v1/auth` ostida emas, to'g'ridan-to'g'ri `/api/v1` ostida)

### Category obyekti (response shakli)

```json
{
  "id": "uuid",
  "name": "Elektronika",
  "parent_id": null,
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z"
}
```

> **Muhim:** `parent_id` maydoni DTO'da bor, lekin hozirgi backend logikasida hech qachon to'ldirilmaydi — har doim `null`/mavjud emas bo'lib qaytadi. Ya'ni **kategoriyalar ierarxiyasi (parent/child daraxti) hozircha backendda ishlamaydi**, barcha kategoriyalar "flat" (tekis) ro'yxat sifatida keladi. Rasm URL'i (`image_url`) response'da yo'qligiga ham e'tibor bering — bu quyida alohida ko'rsatilgan.

Yuqoridagi namunada ko'rinmasa-da, category yaratishda rasm yuklanadi, lekin hozirgi `CreateCategoryResponse` javobida **`image_url` maydoni yo'q** (backendda saqlanadi, lekin frontendga hali qaytarilmayapti). Kategoriya rasmini ko'rsatish kerak bo'lsa, hozircha buni backend jamoasiga bildirish kerak.

---

### 3.1 Kategoriyalar ro'yxatini olish (public)

```
GET /api/v1/categories?search=<matn>
```

- Auth talab qilinmaydi.
- `search` — ixtiyoriy query parametr, nom bo'yicha qidirish uchun.
- Faqat **o'chirilmagan** (`deleted_at IS NULL`) kategoriyalarni qaytaradi.

**Javob — `200 OK`:**
```json
{
  "data": [
    { "id": "uuid", "name": "Elektronika", "image_url": "url", "created_at": "...", "updated_at": "..." }
  ]
}
```

---

### 3.2 Bitta kategoriyani olish

```
GET /api/v1/categories/{id}
```

- Auth talab qilinmaydi.

**Javob — `200 OK`:**
```json
{ "data": { "id": "uuid", "name": "Elektronika", "created_at": "...", "updated_at": "..." } }
```

---

### 3.3 Kategoriyalarni olish — admin (o'chirilganlar bilan birga)

```
GET /api/v1/categories/admin?search=<matn>
```

🔒 **Faqat admin** (`Authorization: Bearer <access_token>`, foydalanuvchi roli `admin` bo'lishi shart)

- Oddiy `/categories`dan farqi: soft-delete qilingan (o'chirilgan) kategoriyalarni ham qaytaradi. Admin panelda "o'chirilgan kategoriyalar" bo'limi uchun ishlating.

**Xatoliklar:** `401` (token yo'q), `403` (admin emas)

---

### 3.4 Yangi kategoriya yaratish

```
POST /api/v1/categories
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Turi | Majburiymi | Izoh |
|---|---|---|---|
| `name` | string | ✅ ha | Kategoriya nomi |
| `image` | file | ✅ ha | Kategoriya rasmi — **majburiy**, bo'lmasa 400 xato qaytadi |

**Rasm cheklovlari (categories va products uchun bir xil):**
- Maksimal hajm: **3 MB**
- Ruxsat etilgan formatlar: `image/jpeg`, `image/png`, `image/webp`
- Bulardan tashqarisi (masalan gif, boshqa hajm) → `400` xato

**Muvaffaqiyatli javob — `201 Created`:**
```json
{
  "data": {
    "id": "uuid",
    "name": "Elektronika",
    "created_at": "2026-08-18T10:00:00Z",
    "updated_at": "2026-08-18T10:00:00Z"
  }
}
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `name` bo'sh, `image` yuborilmagan/noto'g'ri format/3MB dan katta |
| 401 | token yo'q |
| 403 | admin emas |
| 409 | shu nomdagi kategoriya allaqachon mavjud |

---

### 3.5 Kategoriyani yangilash

```
PUT /api/v1/categories/{id}
```

🔒 **Faqat admin**

**Content-Type:** `application/json`

**Body:**
```json
{
  "name": "Yangi nom"
}
```

> Diqqat: bu endpoint **JSON body** qabul qiladi (rasm yangilash uchun `multipart/form-data` emas!). Hozircha faqat `name` maydonini yangilash mumkin — rasmni yangilash uchun alohida endpoint yo'q.

**Javob — `200 OK`:**
```json
{ "data": null }
```

**Xatoliklar:** `401`, `403`, `500` (backendda 400/404 mapping hozircha to'liq ulanmagan — bo'sh nom yoki topilmagan ID kelsa ham xatolik matni `error` maydonida keladi, lekin status kodi 500 bo'lishi mumkin — frontendda `error` matnini ham ko'rsating).

---

### 3.6 Kategoriyani o'chirish

```
DELETE /api/v1/categories/{id}
```

🔒 **Faqat admin**

- Bu **soft delete** — ya'ni ma'lumot bazadan butunlay o'chmaydi, faqat `deleted_at` belgilanadi. Shu sababli o'chirilgan kategoriya oddiy `GET /categories` ro'yxatida chiqmaydi, lekin `GET /categories/admin` orqali ko'rish mumkin.

**Javob — `200 OK`:**
```json
{ "data": "Kategoriya o'chirildi" }
```

**Xatoliklar:** `401`, `403`, `500`

---

## 4. Mahsulotlar (Products)

Prefiks: **`/api/v1/products`**

> ⚠️ **Muhim eslatma frontendchi uchun:** Hozircha backendda **faqat mahsulot yaratish (`POST`)** endpointi tayyor. Mahsulotlar ro'yxatini olish, bitta mahsulotni olish, yangilash va o'chirish endpointlari **hali yozilmagan**. Ya'ni katalog/vitrina sahifasi (mahsulotlarni ko'rsatish) uchun hozircha API yo'q — bu backend tomonda navbatda turgan ish. Iltimos buni backend jamoasi bilan kelishib oling.

### 4.1 Yangi mahsulot yaratish

```
POST /api/v1/products
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Turi | Majburiymi | Izoh |
|---|---|---|---|
| `name` | string | ✅ ha | Mahsulot nomi |
| `amount` | integer | ✅ ha | Narx — **tiyin/kopeykada** (masalan 150000 so'm bo'lsa `15000000` deb yuboriladi, chunki `amount` eng kichik pul birligida saqlanadi — bu haqda backend jamoasidan tasdiqlab oling, chunki hozircha aniq koeffitsient kod ichida izohlangan emas) |
| `currency` | string | ✅ ha | Valyuta kodi, masalan `"UZS"` |
| `categoryId` | string (uuid) | ✅ ha | Mahsulot tegishli kategoriya ID'si |
| `image` | file | ❌ yo'q | Mahsulot rasmi — **ixtiyoriy** (categorydan farqli o'laroq, bu yerda rasm shart emas) |

**Rasm cheklovi:** yuklansa — maks. 3MB, formatlar: `jpeg`/`png`/`webp` (categorydagi bilan bir xil).

**Muvaffaqiyatli javob — `201 Created`:**
```json
{
  "data": {
    "id": "uuid",
    "name": "iPhone 15",
    "category_id": "uuid",
    "amount": 15000000,
    "image_url": "https://res.cloudinary.com/.../product-images/....jpg",
    "currency": "UZS",
    "created_at": "2026-08-18T10:00:00Z"
  }
}
```

- Rasm yuklanmasa `image_url` bo'sh string bo'lib qaytadi.

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `name`/`categoryId` bo'sh, `amount` manfiy yoki raqam emas, `currency` bo'sh, rasm formati/hajmi noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 500 | server xatosi |

---

## 5. Fayl yuklash haqida umumiy qoidalar

Categories va Products ikkalasida ham rasm quyidagi qoidalarga bo'ysunadi:

- **Maksimal hajm:** 3 MB (`multipart` form umumiy hajm chegarasi ham 10MB, lekin rasm faylining o'zi 3MB dan oshmasligi kerak)
- **Ruxsat etilgan formatlar:** `image/jpeg`, `image/png`, `image/webp`
- Rasmlar **Cloudinary**'ga yuklanadi, qaytadigan `image_url` — to'liq CDN havolasi (frontendda to'g'ridan-to'g'ri `<img src>` sifatida ishlatavering).

---

## 6. Rollar (Roles)

| Rol | Qanday beriladi | Nima qila oladi |
|---|---|---|
| `customer` | Har bir yangi `register` shu rolda yaratiladi (default) | Faqat public GET endpointlar (`/categories`, `/categories/{id}`, `/auth/me`) |
| `admin` | Faqat DB orqali qo'lda beriladi, frontendda tanlash yo'q | Category/Product yaratish-o'chirish-yangilash, `/categories/admin` |

Frontendda: login qilingandan keyin `GET /auth/me` chaqirib, javobdagi `role` maydoniga qarab admin panelni ko'rsatish/yashirishni belgilang.

---

## 7. Tezkor cheat-sheet

| Endpoint | Method | Auth | Rol |
|---|---|---|---|
| `/api/v1/auth/register` | POST | ❌ | — |
| `/api/v1/auth/login` | POST | ❌ | — |
| `/api/v1/auth/refresh` | POST | ❌ | — |
| `/api/v1/auth/me` | GET | ✅ | har qanday |
| `/api/v1/categories` | GET | ❌ | — |
| `/api/v1/categories/{id}` | GET | ❌ | — |
| `/api/v1/categories` | POST | ✅ | admin |
| `/api/v1/categories/{id}` | PUT | ✅ | admin |
| `/api/v1/categories/{id}` | DELETE | ✅ | admin |
| `/api/v1/categories/admin` | GET | ✅ | admin |
| `/api/v1/products` | POST | ✅ | admin |

---

## 8. Hali tayyor bo'lmagan (backendda yo'q) narsalar

Frontend ishini rejalashtirishda hisobga oling:

- ❌ `GET /products` — mahsulotlar ro'yxati
- ❌ `GET /products/{id}` — bitta mahsulot
- ❌ `PUT /products/{id}` — mahsulot yangilash
- ❌ `DELETE /products/{id}` — mahsulot o'chirish
- ❌ Kategoriya ierarxiyasi (parent/child daraxti) — `parent_id` maydoni bor, lekin ishlamaydi
- ❌ Category rasm URL'i response'da yo'q
- ❌ Savat, buyurtma (order) — `internal/ordering` papkasi mavjud, lekin ichida hali HTTP endpoint yo'q
- ❌ Parolni tiklash / o'zgartirish, logout endpointi
