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

Kategoriya nomi endi (mahsulot bilan bir xil qoidada) **3 tilda** (`uz`, `eng`, `ru`) saqlanadi. Ommaviy (public) endpointlar `?lang=` query parametriga qarab **bitta tildagi** javob qaytaradi; admin endpointi (`/categories/admin`) esa **har doim barcha 3 tilni** to'liq qaytaradi. Til mexanizmi mahsulotdagi bilan bir xil — [4.0](#40-til-lang-qanday-ishlaydi)ga qarang: `?lang=uz|eng|ru`, berilmasa yoki noto'g'ri bo'lsa `uz` (default), so'ralgan tilda qiymat bo'sh bo'lsa `uz`ga fallback.

### Category obyekti — ommaviy (public) javob shakli

```json
{
  "id": "uuid",
  "name": "Elektronika",
  "image_url": "https://your-bucket.s3.your-region.amazonaws.com/category-images/....jpg",
  "image_public_id": "category-images/xxxxxxx",
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z"
}
```

### Category obyekti — admin javob shakli (`GET /categories/admin`)

```json
{
  "id": "uuid",
  "name_uz": "Elektronika",
  "name_eng": "Electronics",
  "name_ru": "Электроника",
  "image_url": "https://your-bucket.s3.your-region.amazonaws.com/category-images/....jpg",
  "image_public_id": "category-images/xxxxxxx",
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z",
  "deleted_at": null
}
```

> **Muhim:** Kategoriyada **`parent_id` maydoni umuman yo'q** — hozircha kategoriyalar ierarxiyasi (parent/child daraxti) backendda mavjud emas, barcha kategoriyalar "flat" (tekis) ro'yxat sifatida keladi. `image_url` va `image_public_id` har doim javobda bor (AWS S3'ga yuklangan rasm). Nomlar yaratishda **barcha 3 til majburiy** (bo'sh bo'lsa `400`).

---

### 3.1 Kategoriyalar ro'yxatini olish (public)

```
GET /api/v1/categories?search=<matn>&lang=<uz|eng|ru>
```

- Auth talab qilinmaydi.
- `search` — ixtiyoriy query parametr, nom bo'yicha qidirish uchun (`name_uz`/`name_eng`/`name_ru` — uchala til ustuni bo'yicha birga qidiradi).
- `lang` — ixtiyoriy, default `uz`.
- Faqat **o'chirilmagan** (`deleted_at IS NULL`) kategoriyalarni qaytaradi.

**Javob — `200 OK`:**
```json
{
  "data": [
    { "id": "uuid", "name": "Elektronika", "image_url": "url", "image_public_id": "...", "created_at": "...", "updated_at": "..." }
  ]
}
```

---

### 3.2 Bitta kategoriyani olish

```
GET /api/v1/categories/{id}?lang=<uz|eng|ru>
```

- Auth talab qilinmaydi.
- `lang` — ixtiyoriy, default `uz`.

**Javob — `200 OK`:**
```json
{ "data": { "id": "uuid", "name": "Elektronika", "image_url": "...", "image_public_id": "...", "created_at": "...", "updated_at": "..." } }
```

---

### 3.3 Kategoriyalarni olish — admin (o'chirilganlar bilan birga)

```
GET /api/v1/categories/admin?search=<matn>
```

🔒 **Faqat admin** (`Authorization: Bearer <access_token>`, foydalanuvchi roli `admin` bo'lishi shart)

- Oddiy `/categories`dan farqi: soft-delete qilingan (o'chirilgan) kategoriyalarni ham qaytaradi va **har doim barcha 3 tilni to'liq** qaytaradi (`lang` qabul qilinmaydi — yuqoridagi admin shaklga qarang). Admin panelda ro'yxat va tahrirlash formasini shu javobdan to'ldiring.

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
| `name_uz` | string | ✅ ha | Kategoriya nomi (o'zbekcha) |
| `name_eng` | string | ✅ ha | Kategoriya nomi (inglizcha) |
| `name_ru` | string | ✅ ha | Kategoriya nomi (ruscha) |
| `image` | file | ✅ ha | Kategoriya rasmi — **majburiy**, bo'lmasa 400 xato qaytadi |

**Rasm cheklovlari (categories va products uchun bir xil):**
- Maksimal hajm: **3 MB**
- Ruxsat etilgan formatlar: `image/jpeg`, `image/png`, `image/webp`
- Bulardan tashqarisi (masalan gif, boshqa hajm) → `400` xato

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <Category obyekti — public shakl, lang=uz> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `name_uz`/`name_eng`/`name_ru`dan biri bo'sh, `image` yuborilmagan/noto'g'ri format/3MB dan katta |
| 401 | token yo'q |
| 403 | admin emas |
| 500 | server xatosi |

> Eslatma: bir xil nomli kategoriya yaratishga hozircha **cheklov yo'q** (nom unikal bo'lishi shart emas) — `409` bu endpointda qaytmaydi.

---

### 3.5 Kategoriyani yangilash

```
PUT /api/v1/categories/{id}
```

🔒 **Faqat admin**

**Content-Type:** `application/json`

**Body** (barcha maydonlar ixtiyoriy — faqat yubormoqchi bo'lgan tillarni jo'nating, qolganlari o'zgarmaydi):
```json
{
  "name_uz": "Yangi nom",
  "name_eng": "New name",
  "name_ru": "Новое название"
}
```

> Diqqat: bu endpoint **JSON body** qabul qiladi (rasm yangilash uchun `multipart/form-data` emas!). Rasmni yangilash uchun pastdagi **3.6 Kategoriya rasmini yangilash** endpointidan foydalaning. Nomni yangilashda `name_uz`/`name_eng`/`name_ru`dan **kamida bittasini** yuborsangiz, backend o'zgarmagan tillarni joriy qiymati bilan birga qayta tekshiradi (uchalasi ham bo'sh bo'lmasligi kerak) — shuning uchun xavfsizroq usul barcha 3 tilni birga yuborishdir.

**Javob — `200 OK`:**
```json
{ "data": null }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | nom(lar) bo'sh string sifatida yuborilgan |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | kategoriya topilmadi |
| 500 | server xatosi |

---

### 3.6 Kategoriya rasmini yangilash

```
PUT /api/v1/categories/{id}/image
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Tur | Majburiymi | Izoh |
|---|---|---|---|
| `image` | file | ✅ ha | Yangi kategoriya rasmi |

> Eski rasm S3'dan avtomatik o'chiriladi, faqat yangi rasm muvaffaqiyatli saqlangandan keyin.

**Javob — `200 OK`:** yangilangan `CategoryOutput` (3.1-bo'limdagi shakl bilan bir xil)

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | rasm yuborilmagan, hajmi katta yoki formati noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | kategoriya topilmadi |
| 500 | server xatosi |

---

### 3.7 Kategoriyani o'chirish

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

Mahsulot nomi va tavsifi **3 tilda** (`uz`, `eng`, `ru`) saqlanadi. Ommaviy (public) endpointlar `?lang=` query parametriga qarab **bitta tildagi** javob qaytaradi; admin endpointi (`/products/admin`, tahrirlash formasi uchun) esa **har doim barcha 3 tilni** to'liq qaytaradi.

### 4.0 Til (`lang`) qanday ishlaydi

Quyidagi public GET endpointlarning barchasi `?lang=uz|eng|ru` query parametrini qabul qiladi:
`GET /products`, `GET /products/{id}`, `GET /products/slug/{slug}`, `GET /categories/{id}/products`.

- Berilmasa yoki noto'g'ri qiymat yuborilsa → `uz` ishlatiladi (default).
- Agar so'ralgan tilda `name`/`description`/`tag` bo'sh bo'lsa → backend avtomatik `uz` qiymatiga fallback qiladi (frontendda bo'sh matn ko'rinmasligi uchun).
- Misol: `GET /api/v1/products?lang=ru&page=1`

> `GET /products/admin` (admin ro'yxati) `lang` qabul qilmaydi — u har doim 3 ta tilni birga qaytaradi, chunki admin panel tahrirlash formasi barcha tillarni bir vaqtda ko'rsatishi/yangilashi kerak.

---

### Product obyekti — ommaviy (public) javob shakli

`lang`ga mos ravishda **bitta tildagi** `name`/`description`/`tag` bilan qaytadi:

```json
{
  "id": "uuid",
  "name": "51 ta qizil atirgul",
  "description": "Premium Ekvador atirgullaridan yig'ilgan buket",
  "tag": "bestseller",
  "images": [
    "https://your-bucket.s3.your-region.amazonaws.com/product-images/....jpg",
    "https://your-bucket.s3.your-region.amazonaws.com/product-images/....jpg"
  ],
  "category_id": "uuid",
  "price_amount": 150000,
  "price_currency": "UZS",
  "discount_amount": 120000,
  "final_price_amount": 120000,
  "slug": "51-ta-qizil-atirgul",
  "is_available": true,
  "rating": 4.5,
  "stock": 12,
  "sold_count": 34,
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z"
}
```

### Product obyekti — admin javob shakli (`GET /products/admin`)

Barcha 3 til birga, tahrirlash formasini to'ldirish uchun:

```json
{
  "id": "uuid",
  "name_uz": "51 ta qizil atirgul",
  "name_eng": "51 red roses",
  "name_ru": "51 красная роза",
  "description_uz": "Premium Ekvador atirgullaridan yig'ilgan buket",
  "description_eng": "A bouquet of premium Ecuadorian roses",
  "description_ru": "Букет из премиальных эквадорских роз",
  "tag_uz": "bestseller",
  "tag_eng": "Bestseller",
  "tag_ru": "хит продаж",
  "images": ["https://your-bucket.s3.your-region.amazonaws.com/product-images/....jpg"],
  "category_id": "uuid",
  "price_amount": 150000,
  "price_currency": "UZS",
  "discount_amount": 120000,
  "final_price_amount": 120000,
  "slug": "51-ta-qizil-atirgul",
  "is_available": true,
  "rating": 4.5,
  "stock": 12,
  "sold_count": 34,
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z",
  "deleted_at": null
}
```

Muhim izohlar:
- **`name` / `description`** — 3 tilda saqlanadi (`name_uz/eng/ru`, `description_uz/eng/ru`). Yaratishda **nom uchun barcha 3 til majburiy** (bo'sh bo'lsa `400`), tavsif ixtiyoriy (bo'sh string bo'lishi mumkin).
- **`tag`** — ixtiyoriy belgi/badge, masalan `"bestseller"`, 3 tilda (`tag_uz/eng/ru`). Berilmasa `null`/`omitempty`.
- Video (YouTube/Instagram) va gul-do'koniga xos maydonlar (`flower_types`, `color`, `stem_count`, `packaging_type`, `freshness_lifespan`, `care_instructions`, `allow_custom_card`, `compatible_addons`, `occasions`) **butunlay olib tashlangan** — endi mavjud emas.
- **`images`** — AWS S3'dagi rasm URL'lari ro'yxati, **eng ko'pi bilan 5 ta**. Bo'sh bo'lishi ham mumkin (`[]`).
- **`price_amount` / `discount_amount` / `final_price_amount`** — endi **tiyinda emas, to'g'ridan-to'g'ri so'mda** (butun son) saqlanadi va qaytadi — backend hech qanday konversiya qilmaydi, frontend qanday yuborsa, shundayligicha saqlanadi va qaytadi. `discount_amount` bo'lmasa, javobda bu maydon umuman ko'rinmaydi (`omitempty`) va `final_price_amount` = `price_amount` bilan teng bo'ladi. `discount_amount` bo'lsa, u har doim `price_amount`dan kichik bo'ladi va **narxni ko'rsatishda `final_price_amount`dan foydalaning**.
- **`slug`** — URL uchun (masalan `/product/51-ta-qizil-atirgul`). Yaratishda yubormasangiz, backend `name_uz`dan avtomatik hosil qiladi. Har bir mahsulotda **unikal** bo'lishi shart — band bo'lgan slug yuborilsa `409` qaytadi.
- **`rating`** — 1 dan 5 gacha, default `1`. Hozircha foydalanuvchi sharhlaridan avtomatik hisoblanmaydi — admin qo'lda kiritadi (izoh/sharh tizimi hali yo'q).
- **`stock`** / **`sold_count`** — ombordagi son va sotilganlar soni. `is_available` bilan bir xil narsa emas: `stock=0` bo'lsa ham `is_available` alohida `true`/`false` bo'lishi mumkin — frontendda ikkalasini alohida hisobga oling ("tugadi" belgisi uchun `is_available`ni ishlating).
- `category_id` — **backend create/update paytida bu ID chindan mavjud kategoriyaga tegishli ekanligini tekshiradi** — mavjud bo'lmagan `category_id` yuborilsa xato qaytadi.

---

### 4.1 Mahsulotlar ro'yxatini olish (public)

```
GET /api/v1/products?search=<matn>&category_id=<uuid>&lang=<uz|eng|ru>&page=<son>&page_size=<son>
```

- Auth talab qilinmaydi.
- `search` — ixtiyoriy, nom bo'yicha qidirish (`ILIKE`, uchala til ustuni — `name_uz`/`name_eng`/`name_ru` — bo'yicha birga qidiradi).
- `category_id` — ixtiyoriy, faqat shu kategoriyaga tegishli mahsulotlarni qaytaradi. Ikkalasini birga ham berish mumkin.
- `lang` — ixtiyoriy, [4.0](#40-til-lang-qanday-ishlaydi)ga qarang, default `uz`.
- `page` — ixtiyoriy, sahifa raqami, **1 dan boshlanadi**. Berilmasa yoki `1`dan kichik bo'lsa `1` deb olinadi.
- `page_size` — ixtiyoriy, sahifadagi elementlar soni. Berilmasa `20`. Maksimal `100` — undan katta qiymat yuborilsa `100`ga qisqartiriladi.
- Faqat **o'chirilmagan** (`deleted_at IS NULL`) mahsulotlarni qaytaradi. `is_available=false` bo'lgan mahsulotlar ham shu ro'yxatda keladi (yashirilmaydi) — "tugagan" holatini frontendda `is_available` orqali ko'rsating.

**Javob — `200 OK`:**
```json
{
  "data": {
    "items": [ <Product obyekti — public shakl>, ... ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total_items": 42,
      "total_pages": 3
    }
  }
}
```

> Diqqat: bu **ro'yxat qaytaradigan barcha `/products` endpointlarida bir xil shakl** (pastdagi 4.2 va 4.5 ham shu). `data` to'g'ridan-to'g'ri massiv emas, balki `items` + `pagination` bo'lgan obyekt. Mahsulotlar ro'yxatini chizishda `data.items`ni, "keyingi sahifa" tugmasi uchun `data.pagination.total_pages`ni ishlating.

---

### 4.2 Kategoriya bo'yicha mahsulotlarni olish (public)

```
GET /api/v1/categories/{id}/products?search=<matn>&lang=<uz|eng|ru>&page=<son>&page_size=<son>
```

- Auth talab qilinmaydi.
- Xuddi `GET /products?category_id={id}` bilan bir xil natija — kategoriya sahifasida (masalan "Atirgullar" kategoriyasi) qulay bo'lishi uchun alohida yo'l sifatida ham ochilgan.
- `lang` / `page` / `page_size` — 4.1 bilan bir xil qoidalar.

**Javob — `200 OK`:** 4.1dagi bilan bir xil `{ "data": { "items": [...], "pagination": {...} } }` shakli.

---

### 4.3 Bitta mahsulotni olish (ID bo'yicha)

```
GET /api/v1/products/{id}?lang=<uz|eng|ru>
```

- Auth talab qilinmaydi.
- `lang` — ixtiyoriy, default `uz`.

**Javob — `200 OK`:** `{ "data": <Product obyekti — public shakl> }`

**Xatoliklar:** `404` — mahsulot topilmadi (yoki o'chirilgan)

---

### 4.4 Bitta mahsulotni olish (slug bo'yicha)

```
GET /api/v1/products/slug/{slug}?lang=<uz|eng|ru>
```

- Auth talab qilinmaydi. Mahsulot sahifasi (`/product/51-ta-qizil-atirgul`) uchun SEO-friendly URL'da shundan foydalaning.
- `lang` — ixtiyoriy, default `uz`.

**Javob — `200 OK`:** `{ "data": <Product obyekti — public shakl> }`

**Xatoliklar:** `404` — mahsulot topilmadi

---

### 4.5 Mahsulotlarni olish — admin (o'chirilganlar bilan birga)

```
GET /api/v1/products/admin?search=<matn>&category_id=<uuid>&page=<son>&page_size=<son>
```

🔒 **Faqat admin**

- Oddiy `/products`dan farqi: soft-delete qilingan mahsulotlarni ham qaytaradi va **har doim barcha 3 tilni to'liq** qaytaradi (`lang` qabul qilinmaydi — [4.0](#40-til-lang-qanday-ishlaydi)ga qarang). Har bir elementda qo'shimcha `deleted_at` maydoni bo'ladi (o'chirilmagan bo'lsa `null`).
- `page` / `page_size` — 4.1 bilan bir xil qoidalar. Javob shakli ham bir xil: `{ "data": { "items": [...], "pagination": {...} } }`, lekin har bir element yuqoridagi **admin shaklida** (`name_uz/eng/ru`, `description_uz/eng/ru`, `tag_uz/eng/ru`, `deleted_at`).
- Admin panelda mahsulotni tahrirlash formasini shu ro'yxatdagi qatordan to'ldiring — alohida "bitta mahsulotni admin ko'rinishida olish" endpointi yo'q.

**Xatoliklar:** `401` (token yo'q), `403` (admin emas)

---

### 4.6 Yangi mahsulot yaratish

```
POST /api/v1/products
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Turi | Majburiymi | Izoh |
|---|---|---|---|
| `name_uz` | string | ✅ ha | Mahsulot nomi (o'zbekcha) |
| `name_eng` | string | ✅ ha | Mahsulot nomi (inglizcha) |
| `name_ru` | string | ✅ ha | Mahsulot nomi (ruscha) |
| `description_uz` | string | ❌ yo'q | Tavsif (o'zbekcha) |
| `description_eng` | string | ❌ yo'q | Tavsif (inglizcha) |
| `description_ru` | string | ❌ yo'q | Tavsif (ruscha) |
| `category_id` | string (uuid) | ✅ ha | Mavjud kategoriya ID'si — backend tekshiradi |
| `amount` | number | ✅ ha | Narx, **so'mda** (kasr son, masalan `19999.99`, tiyinga aylantirilmaydi) |
| `currency` | string | ✅ ha | Valyuta kodi, masalan `"UZS"` |
| `discount_amount` | number | ❌ yo'q | Chegirma narxi, so'mda. Berilsa `amount`dan kichik va bir xil valyutada bo'lishi shart |
| `slug` | string | ❌ yo'q | Bo'sh qoldirilsa `name_uz`dan avtomatik hosil qilinadi |
| `is_available` | `"true"`/`"false"` | ❌ yo'q | Berilmasa `true` deb olinadi |
| `rating` | number | ❌ yo'q | 1–5, berilmasa `1` |
| `stock` | integer | ❌ yo'q | Berilmasa `0` |
| `tag_uz` | string | ❌ yo'q | Belgi/badge, masalan `"bestseller"` (o'zbekcha, ixtiyoriy) |
| `tag_eng` | string | ❌ yo'q | Belgi/badge (inglizcha, ixtiyoriy) |
| `tag_ru` | string | ❌ yo'q | Belgi/badge (ruscha, ixtiyoriy) |
| `images` | file (bir nechta) | ❌ yo'q | Bir nechta faylni **bir xil `images` maydon nomi bilan** yuboring; eng ko'pi bilan 5 ta |

> Ko'p faylni bitta form-data maydonida yuborish: `FormData.append('images', file1); FormData.append('images', file2); ...` (brauzer/`fetch`da bir nechta marta shu nomda `append` qiling — array belgisi `images[]` shart emas, backend `images` nomidagi barcha fayllarni oladi).

**Rasm cheklovi:** har bir rasm — maks. 3MB, formatlar: `jpeg`/`png`/`webp` (categorydagi bilan bir xil).

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <Product obyekti — public shakl, lang=uz> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `name_uz`/`name_eng`/`name_ru`/`category_id` bo'sh yoki noto'g'ri, `amount`/`discount_amount` noto'g'ri yoki chegirma asosiy narxdan katta, `rating` diapazondan tashqari, 5 tadan ortiq rasm, rasm formati/hajmi noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 409 | shu `slug` allaqachon band |
| 500 | server xatosi |

---

### 4.7 Mahsulotni yangilash

```
PUT /api/v1/products/{id}
```

🔒 **Faqat admin**

**Content-Type:** `application/json`

**Body** (barcha maydonlar ixtiyoriy — faqat yubormoqchi bo'lgan maydonlarni jo'nating, qolganlari o'zgarmaydi):
```json
{
  "name_uz": "Yangi nom",
  "name_eng": "New name",
  "name_ru": "Новое название",
  "description_uz": "Yangi tavsif",
  "description_eng": "New description",
  "description_ru": "Новое описание",
  "category_id": "boshqa-uuid",
  "amount": 160000,
  "currency": "UZS",
  "discount_amount": 130000,
  "clear_discount": false,
  "slug": "yangi-slug",
  "is_available": true,
  "rating": 4.8,
  "stock": 20,
  "sold_count": 40,
  "tag_uz": "bestseller",
  "tag_eng": "Bestseller",
  "tag_ru": "хит продаж",
  "clear_tag_uz": false,
  "clear_tag_eng": false,
  "clear_tag_ru": false
}
```

> Diqqat: bu endpoint **JSON body** qabul qiladi (rasm yangilash uchun `multipart/form-data` emas). Rasmlarni yangilash uchun pastdagi **4.8 Mahsulotga rasm(lar) qo'shish** va **4.9 Mahsulotning bitta rasmini almashtirish** endpointlaridan foydalaning. `category_id` yuborilsa, backend uni ham mavjudligiga tekshiradi. Nomni yangilashda `name_uz`/`name_eng`/`name_ru`dan **kamida bittasini** yuborsangiz, backend o'zgarmagan tillarni joriy qiymati bilan birga qayta tekshiradi — shuning uchun agar 3 tildan birortasini yangilamoqchi bo'lsangiz ham, xavfsizroq usul barcha 3 tilni birga yuborishdir.
>
> Maxsus bayroqlar: `clear_discount: true` — chegirmani butunlay o'chiradi (`discount_amount`ni yubormasdan); `clear_tag_uz`/`clear_tag_eng`/`clear_tag_ru: true` — mos tildagi tag'ni `null` qiladi. Bular berilmasa, mos maydon yuborilgan taqdirdagina o'zgaradi.

**Javob — `200 OK`:** `{ "data": null }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | validatsiya xatosi (bo'sh nom, noto'g'ri narx/chegirma va h.k.) |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | mahsulot topilmadi |
| 409 | yangi `slug` allaqachon band |
| 500 | server xatosi |

---

### 4.8 Mahsulotga rasm(lar) qo'shish

```
POST /api/v1/products/{id}/images
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Tur | Majburiymi | Izoh |
|---|---|---|---|
| `images` | file (bir nechta) | ✅ ha | Qo'shiladigan yangi rasmlar |

> Bu endpoint mavjud rasmlarni **o'chirmaydi** — yangi rasmlarni ularning ustiga qo'shadi. Jami rasmlar soni (eskilar + yangilar) 5 tadan oshsa, `400` qaytadi.

**Javob — `200 OK`:** yangilangan `ProductOutput` (4.1-bo'limdagi shakl bilan bir xil, `images` massivi endi to'liq — eski + yangi rasmlar bilan)

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | rasm yuborilmagan, jami 5 tadan ortiq bo'lib qoladi, hajmi katta yoki formati noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | mahsulot topilmadi |
| 500 | server xatosi |

---

### 4.9 Mahsulotning bitta rasmini almashtirish

```
PUT /api/v1/products/{id}/images/{index}
```

🔒 **Faqat admin**

`{index}` — almashtiriladigan rasmning tartib raqami, **0 dan boshlanadi**, `GET /products/{id}` javobidagi `images` massividagi shu rasmning o'rniga mos keladi (masalan, ro'yxatdagi 2-rasmni almashtirish uchun `index=1`).

**Content-Type:** `multipart/form-data`

| Maydon | Tur | Majburiymi | Izoh |
|---|---|---|---|
| `image` | file | ✅ ha | Yangi rasm (shu index'dagi eski rasm o'rniga) |

> Faqat shu bitta rasm almashtiriladi, qolgan rasmlar o'zgarmaydi. Eski rasm muvaffaqiyatli saqlangandan keyin S3'dan avtomatik o'chiriladi.

**Javob — `200 OK`:** yangilangan `ProductOutput` (4.1-bo'limdagi shakl bilan bir xil)

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | rasm yuborilmagan, `index` noto'g'ri/mavjud bo'lmagan, hajmi katta yoki formati noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | mahsulot topilmadi |
| 500 | server xatosi |

---

### 4.10 Mahsulotning bitta rasmini o'chirish

```
DELETE /api/v1/products/{id}/images/{index}
```

🔒 **Faqat admin**

`{index}` — o'chiriladigan rasmning tartib raqami, **0 dan boshlanadi**, `GET /products/{id}` javobidagi `images` massividagi shu rasmning o'rniga mos keladi.

> Faqat shu bitta rasm o'chiriladi, qolgan rasmlar o'zgarmaydi. Mahsulotda **kamida bitta rasm** qolishi shart — agar bu oxirgi (yagona) rasm bo'lsa, `400` xatosi qaytadi (avval yangi rasm qo'shing yoki almashtiring). Rasm bazadan muvaffaqiyatli o'chirilgandan keyin S3'dan ham o'chiriladi.

**Javob — `200 OK`:** yangilangan `ProductOutput` (4.1-bo'limdagi shakl bilan bir xil)

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `index` noto'g'ri/mavjud bo'lmagan, yoki mahsulotdagi yagona rasmni o'chirishga urinish |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | mahsulot topilmadi |
| 500 | server xatosi |

---

### 4.11 Mahsulotni o'chirish

```
DELETE /api/v1/products/{id}
```

🔒 **Faqat admin**

- Bu **soft delete** — `deleted_at` belgilanadi, yozuv bazadan o'chmaydi. O'chirilgan mahsulot oddiy `GET /products` ro'yxatida chiqmaydi, lekin `GET /products/admin` orqali ko'rish mumkin.

**Javob — `200 OK`:**
```json
{ "data": "Mahsulot o'chirildi" }
```

**Xatoliklar:** `401`, `403`, `500`

---

## 5. Eventlar (Events)

Bosh sahifadagi banner/aksiya bloklari (masalan "Bugungi taklif — Sevimlilar uchun gullar") uchun. Prefiks: **`/api/v1/events`**

`eyebrow`, `title`, `subtitle`, `cta` — to'rttasi ham endi **3 tilda** (`uz`, `eng`, `ru`) saqlanadi (product/category bilan bir xil qoida). Ommaviy endpointlar `?lang=` query parametriga qarab **bitta tildagi** javob qaytaradi; admin endpointi (`/events/admin`) esa **har doim barcha 3 tilni** to'liq qaytaradi. Til mexanizmi — [4.0](#40-til-lang-qanday-ishlaydi)ga qarang: `?lang=uz|eng|ru`, berilmasa yoki noto'g'ri bo'lsa `uz` (default), so'ralgan tilda qiymat bo'sh bo'lsa `uz`ga fallback.

### Event obyekti — ommaviy (public) javob shakli

```json
{
  "id": "uuid",
  "eyebrow": "Bugungi taklif",
  "title": "Sevimlilar uchun gullar",
  "subtitle": "Bugun buyurtma bering, bugun yetkazamiz",
  "cta": "Mahsulotlarni ko'rish",
  "image": "https://your-bucket.s3.your-region.amazonaws.com/event-images/....jpg",
  "category_id": "uuid",
  "is_root": true,
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z"
}
```

### Event obyekti — admin javob shakli (`GET /events/admin`)

```json
{
  "id": "uuid",
  "eyebrow_uz": "Bugungi taklif",
  "eyebrow_eng": "Today's offer",
  "eyebrow_ru": "Предложение дня",
  "title_uz": "Sevimlilar uchun gullar",
  "title_eng": "Flowers for your loved ones",
  "title_ru": "Цветы для любимых",
  "subtitle_uz": "Bugun buyurtma bering, bugun yetkazamiz",
  "subtitle_eng": "Order today, delivered today",
  "subtitle_ru": "Закажите сегодня — доставим сегодня",
  "cta_uz": "Mahsulotlarni ko'rish",
  "cta_eng": "View products",
  "cta_ru": "Смотреть товары",
  "image": "https://your-bucket.s3.your-region.amazonaws.com/event-images/....jpg",
  "category_id": "uuid",
  "is_root": true,
  "created_at": "2026-08-18T10:00:00Z",
  "updated_at": "2026-08-18T10:00:00Z",
  "deleted_at": null
}
```

- `image` — AWS S3'ga yuklangan rasmning to'liq URL'i (local `/images/...` fayl yo'li emas — bu maydonga to'g'ridan-to'g'ri backenddan qaytgan URL keladi, frontendda `<img src>` sifatida shuni ishlating).
- `category_id` — event qaysi kategoriyaga tegishli ekanligi. **Create va update paytida backend bu ID chindan mavjud kategoriyaga tegishli ekanligini tekshiradi** — mavjud bo'lmagan/noto'g'ri `category_id` yuborilsa `400` xato qaytadi.
- `is_root` — `true` bo'lsa, `GET /events` va `GET /events/admin` ro'yxatlarida **birinchi bo'lib** chiqadi (backendda `ORDER BY is_root DESC, created_at DESC`). Bir nechta event `is_root: true` bo'lishi mumkin — ular orasida eng yangisi birinchi keladi.
- `title` — yaratishda **barcha 3 til majburiy** (bo'sh bo'lsa `400`). `eyebrow`/`subtitle`/`cta` — ixtiyoriy (bo'sh string bo'lishi mumkin).

---

### 5.1 Eventlar ro'yxatini olish (public)

```
GET /api/v1/events?lang=<uz|eng|ru>
```

- Auth talab qilinmaydi.
- `lang` — ixtiyoriy, default `uz`.
- Faqat **o'chirilmagan** (`deleted_at IS NULL`) eventlarni qaytaradi — bosh sahifada shundan foydalaning.

**Javob — `200 OK`:**
```json
{
  "data": [
    {
      "id": "uuid",
      "eyebrow": "Bugungi taklif",
      "title": "Sevimlilar uchun gullar",
      "subtitle": "Bugun buyurtma bering, bugun yetkazamiz",
      "cta": "Mahsulotlarni ko'rish",
      "image": "https://your-bucket.s3.your-region.amazonaws.com/gul2.jpg",
      "category_id": "uuid",
      "is_root": true,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

---

### 5.2 Bitta eventni olish

```
GET /api/v1/events/{id}?lang=<uz|eng|ru>
```

- Auth talab qilinmaydi.
- `lang` — ixtiyoriy, default `uz`.

**Javob — `200 OK`:**
```json
{ "data": { "id": "uuid", "eyebrow": "...", "title": "...", "...": "..." } }
```

**Xatoliklar:** `404` — event topilmadi (yoki o'chirilgan)

---

### 5.3 Eventlarni olish — admin (o'chirilganlar bilan birga)

```
GET /api/v1/events/admin
```

🔒 **Faqat admin**

- Oddiy `/events`dan farqi: soft-delete qilingan eventlarni ham qaytaradi va **har doim barcha 3 tilni to'liq** qaytaradi (`lang` qabul qilinmaydi — yuqoridagi admin shaklga qarang). Har bir elementda qo'shimcha `deleted_at` maydoni bo'ladi (o'chirilmagan bo'lsa `null`, o'chirilgan bo'lsa sana).

**Xatoliklar:** `401` (token yo'q), `403` (admin emas)

---

### 5.4 Yangi event yaratish

```
POST /api/v1/events
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Turi | Majburiymi | Izoh |
|---|---|---|---|
| `eyebrow_uz` | string | ❌ yo'q | Kichik ustki matn (o'zbekcha) |
| `eyebrow_eng` | string | ❌ yo'q | Kichik ustki matn (inglizcha) |
| `eyebrow_ru` | string | ❌ yo'q | Kichik ustki matn (ruscha) |
| `title_uz` | string | ✅ ha | Sarlavha (o'zbekcha) |
| `title_eng` | string | ✅ ha | Sarlavha (inglizcha) |
| `title_ru` | string | ✅ ha | Sarlavha (ruscha) |
| `subtitle_uz` | string | ❌ yo'q | Sarlavha ostidagi matn (o'zbekcha) |
| `subtitle_eng` | string | ❌ yo'q | Sarlavha ostidagi matn (inglizcha) |
| `subtitle_ru` | string | ❌ yo'q | Sarlavha ostidagi matn (ruscha) |
| `cta_uz` | string | ❌ yo'q | Tugma matni, masalan "Mahsulotlarni ko'rish" (o'zbekcha) |
| `cta_eng` | string | ❌ yo'q | Tugma matni (inglizcha) |
| `cta_ru` | string | ❌ yo'q | Tugma matni (ruscha) |
| `category_id` | string (uuid) | ✅ ha | Mavjud kategoriya ID'si — backend tekshiradi |
| `is_root` | `"true"` / `"false"` | ❌ yo'q | Berilmasa `false` deb olinadi |
| `image` | file | ✅ ha | Event rasmi — majburiy |

**Rasm cheklovi:** maks. 3MB, formatlar `jpeg`/`png`/`webp` (categories/products bilan bir xil — [12-bo'lim](#12-fayl-yuklash-haqida-umumiy-qoidalar)ga qarang).

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <Event obyekti — public shakl, lang=uz> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `title_uz`/`title_eng`/`title_ru`dan biri bo'sh, `category_id` bo'sh/mavjud bo'lmagan kategoriyaga ishora qilyapti, `image` yuborilmagan yoki formati/hajmi noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 500 | server xatosi |

---

### 5.5 Eventni yangilash

```
PUT /api/v1/events/{id}
```

🔒 **Faqat admin**

**Content-Type:** `application/json`

**Body** (barcha maydonlar ixtiyoriy — faqat yubormoqchi bo'lgan tillarni jo'nating, qolganlari o'zgarmaydi):
```json
{
  "eyebrow_uz": "Yangi eyebrow",
  "eyebrow_eng": "New eyebrow",
  "eyebrow_ru": "Новый eyebrow",
  "title_uz": "Yangi sarlavha",
  "title_eng": "New title",
  "title_ru": "Новый заголовок",
  "subtitle_uz": "Yangi subtitle",
  "subtitle_eng": "New subtitle",
  "subtitle_ru": "Новый подзаголовок",
  "cta_uz": "Yangi tugma matni",
  "cta_eng": "New button text",
  "cta_ru": "Новый текст кнопки",
  "category_id": "boshqa-uuid",
  "is_root": false
}
```

> Diqqat: bu endpoint **JSON body** qabul qiladi. Rasmni bu orqali yangilab bo'lmaydi — pastdagi **5.6 Event rasmini yangilash** endpointidan foydalaning. `category_id` yuborilsa, backend uni ham mavjudligiga tekshiradi. Nomni (`title`) yangilashda `title_uz`/`title_eng`/`title_ru`dan **kamida bittasini** yuborsangiz, backend o'zgarmagan tillarni joriy qiymati bilan birga qayta tekshiradi (uchalasi ham bo'sh bo'lmasligi kerak) — shuning uchun xavfsizroq usul barcha 3 tilni birga yuborishdir.

**Javob — `200 OK`:**
```json
{ "data": null }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `title_uz`/`title_eng`/`title_ru` bo'sh string sifatida yuborilgan, yoki `category_id` mavjud bo'lmagan kategoriyaga ishora qilyapti |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | event topilmadi |
| 500 | server xatosi |

---

### 5.6 Event rasmini yangilash

```
PUT /api/v1/events/{id}/image
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Tur | Majburiymi | Izoh |
|---|---|---|---|
| `image` | file | ✅ ha | Yangi event rasmi |

> Eski rasm S3'dan avtomatik o'chiriladi, faqat yangi rasm muvaffaqiyatli saqlangandan keyin.

**Javob — `200 OK`:** yangilangan `EventOutput` (5.1-bo'limdagi shakl bilan bir xil)

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | rasm yuborilmagan, hajmi katta yoki formati noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | event topilmadi |
| 500 | server xatosi |

---

### 5.7 Eventni o'chirish

```
DELETE /api/v1/events/{id}
```

🔒 **Faqat admin**

- Bu **soft delete** — `deleted_at` belgilanadi, yozuv bazadan o'chmaydi. O'chirilgan event oddiy `GET /events` ro'yxatida chiqmaydi, lekin `GET /events/admin` orqali (o'chirilgan holatda, `deleted_at` sana bilan) ko'rinadi.

**Javob — `200 OK`:**
```json
{ "data": "Event o'chirildi" }
```

**Xatoliklar:** `401`, `403`, `404` (event topilmadi), `500`

---

## 6. Sevimlilar (Wishlist)

Foydalanuvchining "sevimlilar" ro'yxati — mahsulotni yurakcha bosib saqlab qo'yish. Prefiks: **`/api/v1/wishlist`**

🔒 **Uchala endpoint ham autentifikatsiya talab qiladi** (`Authorization: Bearer <access_token>`) — `customer` yoki `admin`, farqi yo'q, faqat token to'g'ri va muddati o'tmagan bo'lishi kerak. Har bir foydalanuvchining **faqat bitta** wishlist'i bo'ladi (token ichidagi `user_id` bo'yicha avtomatik topiladi/yaratiladi — frontend wishlist ID yubormaydi, faqat `product_id`).

> Diqqat: bu bo'lim `/api/v1/wishlist` ostida — yuqoridagi Categories/Products/Events kabi to'g'ridan-to'g'ri `/api/v1` ostida emas, alohida `/api/v1/wishlist` prefiksida joylashgan.

### Wishlist item obyekti

```json
{
  "product_id": "uuid",
  "added_at": "2026-09-08T10:00:00Z"
}
```

---

### 6.1 Sevimlilar ro'yxatini olish

```
GET /api/v1/wishlist
```

🔒 Autentifikatsiya talab qilinadi.

- Foydalanuvchining hali wishlist'i yaratilmagan bo'lsa ham xato qaytmaydi — bo'sh `items` massivi bilan `200 OK` qaytadi.

**Javob — `200 OK`:**
```json
{
  "data": {
    "user_id": "uuid",
    "items": [
      { "product_id": "uuid", "added_at": "2026-09-08T10:00:00Z" }
    ]
  }
}
```

**Xatoliklar:** `401` — token yo'q/noto'g'ri

---

### 6.2 Mahsulotni sevimlilarga qo'shish

```
POST /api/v1/wishlist/items/{product_id}
```

🔒 Autentifikatsiya talab qilinadi.

- `{product_id}` — qo'shiladigan mahsulot ID'si (URL path'da).
- Foydalanuvchining wishlist'i hali mavjud bo'lmasa, avtomatik yaratiladi.
- Backend hozircha `product_id`ning haqiqatan mavjud mahsulotga tegishli ekanligini **tekshirmaydi** — frontend faqat haqiqiy mahsulot ID'sini yuborishi kerak (masalan mahsulot sahifasidagi `id` maydonidan).

**Javob — `200 OK`:**
```json
{ "data": "Muvaffaqiyatli wishlistga qo'shildi" }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `product_id` bo'sh |
| 401 | token yo'q/noto'g'ri |
| 409 | mahsulot allaqachon wishlist'da bor |

---

### 6.3 Mahsulotni sevimlilardan o'chirish

```
DELETE /api/v1/wishlist/items/{product_id}
```

🔒 Autentifikatsiya talab qilinadi.

**Javob — `200 OK`:**
```json
{ "data": "Muvaffaqiyatli wishlistdan o'chirildi" }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 401 | token yo'q/noto'g'ri |
| 404 | foydalanuvchining wishlist'i mavjud emas, yoki bu mahsulot wishlist'da yo'q |

---

## 7. Savat (Cart)

Foydalanuvchining shaxsiy savati. Prefiks: **`/api/v1/cart`**

🔒 **Barcha endpoint ham autentifikatsiya talab qiladi** (`customer` yoki `admin`, farqi yo'q). Har bir foydalanuvchining **faqat bitta** savati bo'ladi (token ichidagi `user_id` bo'yicha avtomatik topiladi/lazy-yaratiladi — frontend cart ID yubormaydi).

> Muhim: savatdagi narxlar **snapshot emas** — `GET /cart` chaqirilganda narx har doim Catalog'dan **jonli (live)** o'qiladi. Ya'ni admin mahsulot narxini o'zgartirsa, foydalanuvchi savatini ochganda darhol yangi narxni ko'radi. Narx checkout paytidagina (buyurtmaga aylanganda) "surat" qilib saqlanadi — [8-bo'lim](#8-buyurtmalar-ordering)ga qarang.

### Cart item obyekti

```json
{
  "product_id": "uuid",
  "product_name": "51 ta qizil atirgul",
  "unit_price": 150000,
  "discount_price": 120000,
  "currency": "UZS",
  "quantity": 2,
  "subtotal": 240000,
  "available": true
}
```

- **`discount_price`** — bo'lsa, `subtotal` shundan hisoblanadi (`discount_price * quantity`); bo'lmasa (`omitempty`, javobda umuman ko'rinmaydi) `subtotal` = `unit_price * quantity`.
- **`available`** — `false` bo'lsa, mahsulot o'chirilgan yoki topilmadi degani; bunday holatda `product_name`/`unit_price`/`currency`/`subtotal` bo'sh keladi, faqat `product_id` va `quantity` bor. Frontendda bunday item'larni "bu mahsulot endi mavjud emas, olib tashlang" tarzida alohida ko'rsating.

---

### 7.1 Savatni olish

```
GET /api/v1/cart
```

🔒 Autentifikatsiya talab qilinadi.

- Foydalanuvchining hali savati yaratilmagan bo'lsa ham xato qaytmaydi — bo'sh `items` bilan `200 OK` qaytadi (wishlist bilan bir xil pattern).

**Javob — `200 OK`:**
```json
{
  "data": {
    "items": [ { "product_id": "uuid", "product_name": "...", "unit_price": 150000, "quantity": 2, "subtotal": 300000, "currency": "UZS", "available": true } ],
    "total_items": 2,
    "total_price": 300000
  }
}
```

- `total_items` — barcha item'lar `quantity`sining yig'indisi (savat badge'i uchun).
- `total_price` — barcha item'lar `subtotal`ining yig'indisi.

**Xatoliklar:** `401` — token yo'q/noto'g'ri

---

### 7.2 Savatga mahsulot qo'shish

```
POST /api/v1/cart/items
```

🔒 Autentifikatsiya talab qilinadi.

**Body:**
```json
{ "product_id": "uuid", "quantity": 2 }
```

> **Diqqat — upsert xulqi:** agar bu `product_id` savatda allaqachon bo'lsa, yuborilgan `quantity` mavjud miqdorning **ustiga qo'shiladi** (increment), ustidan yozilmaydi. Masalan savatda `quantity: 2` bo'lsa va yana `{"quantity": 3}` yuborilsa, natija `quantity: 5` bo'ladi. Miqdorni **aynan shu songa belgilash** kerak bo'lsa, o'rniga [7.3](#73-savatdagi-mahsulot-miqdorini-yangilash)dagi `PUT` endpointidan foydalaning.
>
> Stock tekshiruvi ham kumulyativ: savatdagi joriy miqdor + yangi yuborilgan miqdor mahsulotning `stock`idan oshsa, `409` qaytadi.

**Javob — `200 OK`:**
```json
{ "data": "Mahsulot savatga qo'shildi" }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `quantity` 0 yoki manfiy, `product_id` bo'sh |
| 401 | token yo'q/noto'g'ri |
| 404 | `product_id` bo'yicha mahsulot topilmadi (yoki o'chirilgan) |
| 409 | (savatdagi joriy + yangi) miqdor mahsulot `stock`idan oshib ketadi |

---

### 7.3 Savatdagi mahsulot miqdorini yangilash

```
PUT /api/v1/cart/items/{product_id}
```

🔒 Autentifikatsiya talab qilinadi.

**Body:**
```json
{ "quantity": 5 }
```

> 7.2'dan farqi: bu yerda `quantity` **absolyut qiymat** sifatida belgilanadi (eskisining ustiga qo'shilmaydi, to'g'ridan-to'g'ri almashtiriladi).

**Javob — `200 OK`:**
```json
{ "data": "Savatdagi mahsulot miqdori yangilandi" }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `quantity` 0 yoki manfiy |
| 401 | token yo'q/noto'g'ri |
| 404 | bu `product_id` savatda yo'q |

---

### 7.4 Savatdan mahsulotni o'chirish

```
DELETE /api/v1/cart/items/{product_id}
```

🔒 Autentifikatsiya talab qilinadi.

**Javob — `200 OK`:**
```json
{ "data": "Mahsulot savatdan o'chirildi" }
```

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 401 | token yo'q/noto'g'ri |
| 404 | bu `product_id` savatda yo'q |

---

## 8. Buyurtmalar (Ordering)

Checkout (savatni buyurtmaga aylantirish), buyurtmalar tarixi va admin buyurtma boshqaruvi. **Uchta alohida prefiks** ostida: `/api/v1/checkout`, `/api/v1/orders`, `/api/v1/admin/orders` (Catalog bare `/api/v1`ni egallagani sabab, ordering context'i bitta router o'rniga uchta alohida router'ga bo'lingan — bu faqat backend ichki tuzilishi, frontend uchun ahamiyatsiz, shunchaki har uch prefiks ham mavjud ekanini bilib qo'ying).

### Order obyekti

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "address": "Toshkent sh., Chilonzor tumani, ...",
  "phone": "+998901234567",
  "note": "Domofon kodi 1234",
  "items": [
    { "product_id": "uuid", "product_name": "51 ta qizil atirgul", "unit_price": 150000, "currency": "UZS", "quantity": 2 }
  ],
  "payment_status": "unpaid",
  "delivery_status": "preparing",
  "total_amount": 300000,
  "total_currency": "UZS",
  "created_at": "2026-09-01T10:00:00Z"
}
```

- **`user_id`** — admin qo'lda yaratgan (cart'siz) buyurtmalarda `null`/`omitempty` bo'lishi mumkin (offline mijoz uchun user hisobi yo'q).
- **`items`** — checkout/admin-order yaratish paytidagi mahsulot nomi va narxi **"surat" (snapshot) qilib saqlanadi** — keyinchalik mahsulot narxi yoki nomi o'zgarsa ham, eski buyurtmadagi qiymatlar o'zgarmaydi (Cart'dagi live-narx bilan bu yerning asosiy farqi shu).
- **`payment_status`** — `"unpaid"` yoki `"paid"`.
- **`delivery_status`** — `"preparing"` → `"handed_to_courier"` → `"delivered"`, yoki alohida terminal holat `"cancelled"`. Holat faqat **oldinga** qarab o'zgarishi mumkin (masalan `delivered`dan `preparing`ga qaytarib bo'lmaydi) — bekor qilish alohida endpoint orqali ([8.8](#88-buyurtmani-bekor-qilish)) va faqat `preparing` bosqichida mumkin.
- **`note`** — ixtiyoriy, bo'sh bo'lishi mumkin (`omitempty`).
- **`total_amount`/`total_currency`** — barcha `items`ning `unit_price * quantity` yig'indisi.

---

### 8.1 Checkout — savatni buyurtmaga aylantirish

```
POST /api/v1/checkout
```

🔒 Autentifikatsiya talab qilinadi.

- Joriy foydalanuvchining **savatidagi** mahsulotlarni buyurtmaga aylantiradi: Catalog'dan joriy nom/narxni "surat" qiladi, zaxirani (`stock`) atomik ravishda kamaytiradi va savatni bo'shatadi.
- Yangi buyurtma har doim `payment_status: "unpaid"`, `delivery_status: "preparing"` bilan boshlanadi.

**Body:**
```json
{
  "address": "Toshkent sh., Chilonzor tumani, ...",
  "phone": "+998901234567",
  "note": "Domofon kodi 1234"
}
```

- `phone` — qat'iy formatda tekshiriladi: `+998` bilan boshlanib, keyin 9 ta raqam (`+998901234567`). Boshqa format `400` qaytaradi.
- `address` — bo'sh bo'lishi mumkin emas.
- `note` — ixtiyoriy.

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <Order obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | manzil bo'sh, telefon formati noto'g'ri, yoki savat bo'sh |
| 401 | token yo'q |
| 404 | savatdagi mahsulotlardan biri topilmadi (o'chirilgan) |
| 409 | savatdagi mahsulotlardan biri uchun yetarli `stock` yo'q |
| 500 | server xatosi |

---

### 8.2 Mening buyurtmalarim

```
GET /api/v1/orders
```

🔒 Autentifikatsiya talab qilinadi.

- Joriy foydalanuvchining barcha buyurtmalari ro'yxati ("Mening buyurtmalarim" sahifasi uchun).

**Javob — `200 OK`:** `{ "data": [ <Order obyekti>, ... ] }`

**Xatoliklar:** `401`, `500`

---

### 8.3 Bitta buyurtma tafsilotlari

```
GET /api/v1/orders/{id}
```

🔒 Autentifikatsiya talab qilinadi.

- Faqat buyurtma **egasi** yoki **admin** ko'ra oladi — boshqa foydalanuvchining buyurtmasini ochishga urinish `403` qaytaradi (ownership tekshiruvi backendda, frontendda alohida qilish shart emas).

**Javob — `200 OK`:** `{ "data": <Order obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 401 | token yo'q |
| 403 | bu buyurtma boshqa foydalanuvchiga tegishli |
| 404 | buyurtma topilmadi |

---

### 8.4 Barcha buyurtmalar — admin navbati

```
GET /api/v1/orders/admin
```

🔒 **Faqat admin**

- Barcha foydalanuvchilarning barcha buyurtmalari (admin panelidagi "buyurtmalar" jadvali uchun).

**Javob — `200 OK`:** `{ "data": [ <Order obyekti>, ... ] }`

**Xatoliklar:** `401`, `403`

---

### 8.5 To'lov holatini o'zgartirish

```
PATCH /api/v1/orders/{id}/payment-status
```

🔒 **Faqat admin**

**Body:**
```json
{ "status": "paid" }
```

- `status` — `"unpaid"` yoki `"paid"` bo'lishi shart, boshqa qiymat `400` qaytaradi.

**Javob — `200 OK`:** `{ "data": <yangilangan Order obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `status` `unpaid`/`paid`dan boshqa qiymat |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | buyurtma topilmadi |

---

### 8.6 Yetkazib berish holatini o'zgartirish

```
PATCH /api/v1/orders/{id}/delivery-status
```

🔒 **Faqat admin**

**Body:**
```json
{ "status": "handed_to_courier" }
```

- `status` — `"preparing"`, `"handed_to_courier"` yoki `"delivered"` bo'lishi shart.
- Faqat **oldinga** qarab o'zgarishi mumkin (masalan `delivered`dagi buyurtmani qayta `preparing`ga qaytarib bo'lmaydi) — orqaga urinish `409` qaytaradi.

**Javob — `200 OK`:** `{ "data": <yangilangan Order obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `status` yaroqsiz qiymat |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | buyurtma topilmadi |
| 409 | holatni orqaga qaytarishga urinish |

---

### 8.7 Qo'lda buyurtma yaratish (offline savdo)

```
POST /api/v1/admin/orders
```

🔒 **Faqat admin**

- Cart'siz, to'g'ridan-to'g'ri buyurtma yaratadi — masalan telefon orqali qabul qilingan yoki do'kondagi offline savdo uchun. Checkout bilan bir xil mantiq: Catalog'dan nom/narxni "surat" qiladi, `stock`ni kamaytiradi.

**Body:**
```json
{
  "address": "Toshkent sh., ...",
  "phone": "+998901234567",
  "note": "Telefon orqali qabul qilindi",
  "items": [
    { "product_id": "uuid", "quantity": 2 },
    { "product_id": "uuid-2", "quantity": 1 }
  ],
  "payment_status": "paid"
}
```

- `items` — bo'sh bo'lishi mumkin emas. Bir xil `product_id` bir nechta marta kelsa, backend ularni **birlashtirib** (quantity'larini qo'shib) stock'dan bittagina marta kamaytiradi.
- `payment_status` — `"unpaid"` yoki `"paid"` (offline savdo ko'pincha joyida to'langan bo'ladi, shu sabab bu yerda majburiy va boshlang'ich qiymat frontend tomonidan tanlanadi).

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <Order obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | manzil/telefon/`payment_status` noto'g'ri, `items` bo'sh |
| 401 | token yo'q |
| 403 | admin emas |
| 404 | `items` ichidagi mahsulotlardan biri topilmadi |
| 409 | mahsulotlardan biri uchun yetarli `stock` yo'q |
| 500 | server xatosi |

---

### 8.8 Buyurtmani bekor qilish

```
PATCH /api/v1/admin/orders/{id}/cancel
```

🔒 **Faqat admin**

- Buyurtmani bekor qiladi (`delivery_status: "cancelled"`) va checkout/qo'lda-yaratish paytida kamaytirilgan `stock`ni Catalog'ga **qaytaradi**.
- Faqat **`preparing`** bosqichidagi buyurtmalar bekor qilinishi mumkin — courier'ga topshirilgan (`handed_to_courier`) yoki yetkazilgan (`delivered`) buyurtmani bu yo'l bilan bekor qilib bo'lmaydi.

**Javob — `200 OK`:** `{ "data": <bekor qilingan Order obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 401 | token yo'q |
| 403 | admin emas |
| 404 | buyurtma topilmadi |
| 409 | buyurtma `preparing` bosqichida emas (allaqachon courier'ga topshirilgan/yetkazilgan) |

---

## 9. Sharhlar (Reviews)

Mahsulotga sharh (izoh + reyting) qoldirish. **Uchta alohida prefiks**: `/api/v1/reviews` (yaratish), `/api/v1/products/{id}/reviews` (o'qish), `/api/v1/admin/reviews` (moderatsiya) — sababi 8-bo'limdagidek, bare `/api/v1` Catalog'ga tegishli.

> Muhim cheklov: foydalanuvchi sharh qoldirishi uchun shu mahsulotni sotib olib, buyurtmasi **yetkazilgan (`delivered`)** bo'lishi shart, va har bir foydalanuvchi bitta mahsulotga **faqat bitta marta** sharh qoldirishi mumkin. Frontendda "sharh qoldirish" tugmasini faqat shu shartlar bajarilganda ko'rsating (masalan foydalanuvchining shu mahsulot bo'yicha yetkazilgan buyurtmasi bor-yo'qligini `GET /orders`dan tekshirib).

### Review obyekti

```json
{
  "id": "uuid",
  "user_id": "uuid",
  "product_id": "uuid",
  "order_id": "uuid",
  "rating": 5,
  "comment": "Juda chiroyli gullar, tavsiya qilaman!",
  "created_at": "2026-09-05T10:00:00Z"
}
```

- `rating` — 1 dan 5 gacha butun son.
- `comment` — ixtiyoriy (`omitempty`).
- `order_id` — sharh qaysi (yetkazilgan) buyurtma asosida qoldirilganini ko'rsatadi.

---

### 9.1 Mahsulotga sharh qoldirish

```
POST /api/v1/reviews
```

🔒 Autentifikatsiya talab qilinadi.

**Body:**
```json
{ "product_id": "uuid", "rating": 5, "comment": "Juda chiroyli gullar!" }
```

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <Review obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `rating` 1–5 oralig'ida emas, `product_id` bo'sh |
| 401 | token yo'q |
| 403 | bu mahsulot uchun sharh qoldirish huquqi yo'q (yetkazilgan buyurtma yo'q) |
| 409 | bu mahsulot uchun sharh allaqachon qoldirilgan |
| 500 | server xatosi |

---

### 9.2 Mahsulot sharhlari

```
GET /api/v1/products/{id}/reviews?page=<son>&page_size=<son>
```

- Auth talab qilinmaydi (ochiq).
- `page` — ixtiyoriy, default `1`.
- `page_size` — ixtiyoriy, default `20`, maksimal `100`.

**Javob — `200 OK`:**
```json
{
  "data": {
    "items": [ <Review obyekti>, ... ],
    "pagination": { "page": 1, "page_size": 20, "total_items": 12, "total_pages": 1 }
  }
}
```

> Bu — [4.1-bo'limdagi](#41-mahsulotlar-royxatini-olish-public) `{ "items": [...], "pagination": {...} }` shakli bilan bir xil umumiy pagination formati (Gallery'da ham xuddi shunday, [10.1](#101-galereya-postlarini-olish)ga qarang).

**Xatoliklar:** `500`

---

### 9.3 Sharhni o'chirish (moderatsiya)

```
DELETE /api/v1/admin/reviews/{id}
```

🔒 **Faqat admin**

- Istalgan foydalanuvchi yozgan sharhni butunlay o'chiradi. Egalik tekshiruvi yo'q — bu endpoint faqat admin uchun ochiq, shu bilan cheklanadi.

**Javob — `200 OK`:** `{ "data": "Sharh o'chirildi" }`

**Xatoliklar:** `401`, `403`, `404` — sharh topilmadi, `500`

---

## 10. Galereya (Gallery)

Bosh sahifadagi ilhom/portfolio postlari (masalan "Bizning ishlarimiz" bo'limi) — pricing/stock bilan bog'liq emas, eng oddiy context. Ikkita prefiks: `/api/v1/gallery` (ochiq o'qish), `/api/v1/admin/gallery` (yaratish/o'chirish).

### GalleryPost obyekti

```json
{
  "id": "uuid",
  "image_urls": [
    "https://your-bucket.s3.your-region.amazonaws.com/gallery-images/....jpg"
  ],
  "description": "Bahorgi buketlar to'plami",
  "created_at": "2026-09-05T10:00:00Z"
}
```

- `image_urls` — eng ko'pi bilan **3 ta** rasm, bo'sh massiv ham bo'lishi mumkin.
- `description` — ixtiyoriy, bo'sh bo'lishi mumkin (`omitempty`).

---

### 10.1 Galereya postlarini olish

```
GET /api/v1/gallery?page=<son>&page_size=<son>
```

- Auth talab qilinmaydi.
- `page` — ixtiyoriy, default `1`. `page_size` — ixtiyoriy, default `20`, maksimal `100`.

**Javob — `200 OK`:**
```json
{
  "data": {
    "items": [ <GalleryPost obyekti>, ... ],
    "pagination": { "page": 1, "page_size": 20, "total_items": 8, "total_pages": 1 }
  }
}
```

**Xatoliklar:** `500`

---

### 10.2 Yangi post yaratish

```
POST /api/v1/admin/gallery
```

🔒 **Faqat admin**

**Content-Type:** `multipart/form-data`

| Maydon | Turi | Majburiymi | Izoh |
|---|---|---|---|
| `images` | file (bir nechta) | ❌ yo'q | Eng ko'pi bilan **3 ta** rasm — `FormData.append('images', file)` bir necha marta |
| `description` | string | ❌ yo'q | Post tavsifi |

**Rasm cheklovi:** har biri maks. 3MB, formatlar `jpeg`/`png`/`webp` ([12-bo'lim](#12-fayl-yuklash-haqida-umumiy-qoidalar)ga qarang).

**Muvaffaqiyatli javob — `201 Created`:** `{ "data": <GalleryPost obyekti> }`

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | 3 tadan ortiq rasm, rasm formati/hajmi noto'g'ri |
| 401 | token yo'q |
| 403 | admin emas |
| 500 | server xatosi |

---

### 10.3 Postni o'chirish

```
DELETE /api/v1/admin/gallery/{id}
```

🔒 **Faqat admin**

- Postni bazadan o'chiradi va rasmlarini S3'dan ham tozalaydi (soft delete emas — bu context'da alohida "o'chirilganlarni ko'rish" admin ro'yxati yo'q).

**Javob — `200 OK`:** `{ "data": "Post o'chirildi" }`

**Xatoliklar:** `401`, `403`, `404` — post topilmadi, `500`

---

## 11. Admin Dashboard

Faqat-o'qish (read-only) statistika/hisobot paneli. Bounded context emas (biznes qoidasi yo'q, faqat agregatsiya) — shu sabab hech qanday `create`/`update`/`delete` yo'q, faqat 3 ta `GET`. Prefiks: **`/api/v1/admin/dashboard`**

🔒 **Uchala endpoint ham faqat admin uchun.**

> Diqqat: barcha summalar **bitta valyutada** hisoblanadi deb faraz qilinadi (loyihada amalda har doim `UZS`) — agar kelajakda ko'p valyuta qo'llab-quvvatlansa, bu endpointlar valyutalarni aralashtirib yuborishi mumkin (bu haqida `CLAUDE.md`da ochiq muammo sifatida qayd etilgan).

### 11.1 Umumiy statistika

```
GET /api/v1/admin/dashboard/summary
```

**Javob — `200 OK`:**
```json
{
  "data": {
    "total_products": 42,
    "total_units_sold": 310,
    "total_revenue": 45600000
  }
}
```

- `total_products` — o'chirilmagan mahsulotlar soni.
- `total_units_sold` / `total_revenue` — faqat **to'langan va bekor qilinmagan** (`payment_status = "paid" AND delivery_status <> "cancelled"`) buyurtmalar bo'yicha hisoblanadi.

**Xatoliklar:** `401`, `403`, `500`

---

### 11.2 Daromad tarixi

```
GET /api/v1/admin/dashboard/revenue-history?period=<day|month>
```

- `period` — **majburiy**, faqat `"day"` yoki `"month"` qiymatini qabul qiladi. Boshqa qiymat `400` qaytaradi.

**Javob — `200 OK`:**
```json
{
  "data": [
    { "period": "2026-09-01T00:00:00Z", "revenue": 1200000 },
    { "period": "2026-09-02T00:00:00Z", "revenue": 3400000 }
  ]
}
```

- Grafik chizish uchun mo'ljallangan (masalan bar/line chart) — `summary`dagi bilan bir xil "to'langan va bekor qilinmagan" filtri qo'llaniladi.

**Xatoliklar:**
| Status | Sabab |
|---|---|
| 400 | `period` `day`/`month`dan boshqa qiymat (yoki berilmagan) |
| 401 | token yo'q |
| 403 | admin emas |
| 500 | server xatosi |

---

### 11.3 Kam zaxirali mahsulotlar

```
GET /api/v1/admin/dashboard/low-stock
```

- Eng kam `stock`ga ega (o'chirilmagan) **top-5** mahsulot — admin panelida "tez orada tugaydi" ogohlantirishi uchun.

**Javob — `200 OK`:**
```json
{
  "data": [
    { "id": "uuid", "name_uz": "51 ta qizil atirgul", "stock": 2 },
    { "id": "uuid-2", "name_uz": "Bahorgi buket", "stock": 3 }
  ]
}
```

**Xatoliklar:** `401`, `403`, `500`

---

## 12. Fayl yuklash haqida umumiy qoidalar

Categories, Products, Events va Gallery — barchasida rasm quyidagi qoidalarga bo'ysunadi:

- **Maksimal hajm:** 3 MB (`multipart` form umumiy hajm chegarasi ham 10MB, lekin rasm faylining o'zi 3MB dan oshmasligi kerak)
- **Ruxsat etilgan formatlar:** `image/jpeg`, `image/png`, `image/webp`
- Rasmlar **AWS S3**'ga yuklanadi, qaytadigan URL — to'liq S3 havolasi (frontendda to'g'ridan-to'g'ri `<img src>` sifatida ishlatavering).

---

## 13. Rollar (Roles)

| Rol | Qanday beriladi | Nima qila oladi |
|---|---|---|
| `customer` | Har bir yangi `register` shu rolda yaratiladi (default) | Public GET endpointlar (`/categories`, `/events`, `/products`, `/gallery`, `/auth/me`) + o'zining `/wishlist`'i, `/cart`'i, `/orders`'i + `/checkout` + `/reviews` (sharh qoldirish) |
| `admin` | Faqat DB orqali qo'lda beriladi, frontendda tanlash yo'q | Category/Product/Event/Gallery yaratish-o'chirish-yangilash, `/categories/admin`, `/events/admin`, `/products/admin`, buyurtmalarni boshqarish (`/orders/admin`, `/admin/orders`, holat o'zgartirish), sharhlarni o'chirish (`/admin/reviews/{id}`), `/admin/dashboard`, + o'zining `/wishlist`'i va `/cart`'i |

Frontendda: login qilingandan keyin `GET /auth/me` chaqirib, javobdagi `role` maydoniga qarab admin panelni ko'rsatish/yashirishni belgilang.

---

## 14. Tezkor cheat-sheet

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
| `/api/v1/categories/{id}/image` | PUT | ✅ | admin |
| `/api/v1/categories/{id}` | DELETE | ✅ | admin |
| `/api/v1/categories/admin` | GET | ✅ | admin |
| `/api/v1/categories/{id}/products` | GET | ❌ | — |
| `/api/v1/products` | GET | ❌ | — |
| `/api/v1/products/{id}` | GET | ❌ | — |
| `/api/v1/products/slug/{slug}` | GET | ❌ | — |
| `/api/v1/products` | POST | ✅ | admin |
| `/api/v1/products/{id}` | PUT | ✅ | admin |
| `/api/v1/products/{id}/images` | POST | ✅ | admin |
| `/api/v1/products/{id}/images/{index}` | PUT | ✅ | admin |
| `/api/v1/products/{id}/images/{index}` | DELETE | ✅ | admin |
| `/api/v1/products/{id}` | DELETE | ✅ | admin |
| `/api/v1/products/admin` | GET | ✅ | admin |
| `/api/v1/events` | GET | ❌ | — |
| `/api/v1/events/{id}` | GET | ❌ | — |
| `/api/v1/events` | POST | ✅ | admin |
| `/api/v1/events/{id}` | PUT | ✅ | admin |
| `/api/v1/events/{id}/image` | PUT | ✅ | admin |
| `/api/v1/events/{id}` | DELETE | ✅ | admin |
| `/api/v1/events/admin` | GET | ✅ | admin |
| `/api/v1/wishlist` | GET | ✅ | har qanday |
| `/api/v1/wishlist/items/{product_id}` | POST | ✅ | har qanday |
| `/api/v1/wishlist/items/{product_id}` | DELETE | ✅ | har qanday |
| `/api/v1/cart` | GET | ✅ | har qanday |
| `/api/v1/cart/items` | POST | ✅ | har qanday |
| `/api/v1/cart/items/{product_id}` | PUT | ✅ | har qanday |
| `/api/v1/cart/items/{product_id}` | DELETE | ✅ | har qanday |
| `/api/v1/checkout` | POST | ✅ | har qanday |
| `/api/v1/orders` | GET | ✅ | har qanday |
| `/api/v1/orders/{id}` | GET | ✅ | egasi yoki admin |
| `/api/v1/orders/admin` | GET | ✅ | admin |
| `/api/v1/orders/{id}/payment-status` | PATCH | ✅ | admin |
| `/api/v1/orders/{id}/delivery-status` | PATCH | ✅ | admin |
| `/api/v1/admin/orders` | POST | ✅ | admin |
| `/api/v1/admin/orders/{id}/cancel` | PATCH | ✅ | admin |
| `/api/v1/reviews` | POST | ✅ | har qanday |
| `/api/v1/products/{id}/reviews` | GET | ❌ | — |
| `/api/v1/admin/reviews/{id}` | DELETE | ✅ | admin |
| `/api/v1/gallery` | GET | ❌ | — |
| `/api/v1/admin/gallery` | POST | ✅ | admin |
| `/api/v1/admin/gallery/{id}` | DELETE | ✅ | admin |
| `/api/v1/admin/dashboard/summary` | GET | ✅ | admin |
| `/api/v1/admin/dashboard/revenue-history` | GET | ✅ | admin |
| `/api/v1/admin/dashboard/low-stock` | GET | ✅ | admin |

---

## 15. Hali tayyor bo'lmagan (backendda yo'q) narsalar

Frontend ishini rejalashtirishda hisobga oling:

- ❌ Kategoriya ierarxiyasi (parent/child daraxti) — kategoriyada `parent_id` degan maydon umuman yo'q, barcha kategoriyalar "flat" ro'yxat
- ❌ Parolni tiklash / o'zgartirish, logout endpointi
- ❌ "Mening sharhlarim" (foydalanuvchining o'zi yozgan barcha sharhlari) — backendda use case bor (`GetUserReviewsUseCase`), lekin hali HTTP route ochilmagan
- ❌ Gallery uchun tahrirlash (update) — faqat yaratish/ro'yxat/o'chirish bor, tavsifni yoki bitta rasmni almashtirish endpointi yo'q
- ❌ Qidiruv va sahifalash `GET /categories`, `GET /products/admin`, `GET /categories/admin`, `GET /events`/`GET /events/admin` uchun to'liq emas — ba'zi admin ro'yxatlarida hali `page`/`page_size` yo'q (faqat Products/Reviews/Gallery ro'yxatlarida sahifalash bor)
- ❌ Redis keshlash hali ishlatilmayapti (real trafik kutilmoqda)

> Eslatma: endi tayyor bo'lgan narsalar — rasmlarni yangilash (`PUT /categories/{id}/image` — 3.6, `POST/PUT /products/{id}/images...` — 4.8/4.9, `PUT /events/{id}/image` — 5.6), Sevimlilar ([6-bo'lim](#6-sevimlilar-wishlist)), **Savat** ([7-bo'lim](#7-savat-cart)), **Buyurtma/checkout** ([8-bo'lim](#8-buyurtmalar-ordering)), **Mahsulot sharhlari** ([9-bo'lim](#9-sharhlar-reviews)), **Galereya** ([10-bo'lim](#10-galereya-gallery)) va **Admin statistika paneli** ([11-bo'lim](#11-admin-dashboard)).
