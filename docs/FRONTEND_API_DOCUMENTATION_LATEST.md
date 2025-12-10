# SimNikah API — Frontend Reference (Latest)

This document summarizes the API surface frontend teams need to integrate with SimNikah backend (including the new profile photo upload feature). It includes authentication rules, headers, main endpoints, example requests/responses, file upload guidance, CORS, and useful tips.

Base URL
- Local development: `http://localhost:8080`
- Production: Set by deployment (Railway) — use provided domain

Authentication
- Type: Bearer JWT
- Header: `Authorization: Bearer <token>`
- How to get token: `POST /login` (see Login endpoint)
- Token expiry: 24 hours (by default)

Common headers
- `Content-Type: application/json` for JSON bodies
- `Accept: application/json`
- For file uploads: do NOT set `Content-Type` manually; use `multipart/form-data` with boundary handled by browser/HTTP client

Rate limiting
- Auth endpoints use stricter rate limiting
- Global rate limiter: ~100 req/min per IP

CORS
- The server reads `ALLOWED_ORIGINS`; defaults include `http://localhost:3000`, `http://localhost:5173`, and the production frontend `https://kua-ku.vercel.app`.
- Ensure your frontend origin is added to `ALLOWED_ORIGINS` in production if different.

Environment variables relevant to frontend integration
- `PORT` — server port (default 8080)
- `ALLOWED_ORIGINS` — comma-separated CORS whitelist
- `IMGBB_API_KEY` — optional ImgBB API key used for profile photo uploads

Response format conventions
- Success responses typically:
  {
    "success": true,
    "message": "...",
    "data": { ... }
  }
- Error responses typically:
  {
    "success": false,
    "message": "...",
    "error": "..."
  }

--

## Authentication

### POST /register
- Description: Register a new user
- Auth: Public
- Body (JSON):
  {
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "nama": "John Doe",
    "role": "user_biasa"  // one of: user_biasa, staff, penghulu, kepala_kua
  }
- Response: 201 Created with created user summary

### POST /login
- Description: Login and receive JWT token
- Auth: Public
- Body (JSON):
  {
    "username": "john_doe",
    "password": "password123"
  }
- Success Response (200):
  {
    "success": true,
    "message": "Login berhasil",
    "token": "<JWT_TOKEN>",
    "user": { "user_id":"USR...", "username":"john_doe", "email":"...", "role":"user_biasa", "nama":"John Doe" },
    "data": { "token":"<JWT_TOKEN>", "user": {...} }
  }
- Use `token` for authenticated requests.

--

## User profile

### GET /profile
- Description: Get current user profile
- Auth: Required (Bearer token)
- Response (200):
  {
    "success": true,
    "message": "Profile berhasil diambil",
    "data": {
      "user_id": "USR...",
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user_biasa",
      "nama": "John Doe",
      "status": "Aktif",
      "profile_photo": "https://..." // empty string if none
    }
  }

### POST /upload-photo
- Description: Upload profile photo for currently authenticated user
- Auth: Required (Bearer token)
- Content-Type: `multipart/form-data`
- Form field name: `photo` (file)
- Validation:
  - Max file size: 5MB
  - Allowed formats/MIME types: `image/jpeg`, `image/png`, `image/jpg`, `image/webp`
- Success Response (200):
  {
    "success": true,
    "message": "Foto profil berhasil diupload",
    "data": { "profile_photo": "<URL>", "user_id":"...", "username":"..." }
  }

cURL example:
```bash
TOKEN="<JWT_TOKEN>"
curl -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@/path/to/photo.jpg"
```

Fetch example (browser):
```javascript
const form = new FormData();
form.append('photo', fileInput.files[0]);
fetch('http://localhost:8080/upload-photo', {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: form
})
.then(r => r.json()).then(console.log)
```

Notes:
- Do not set `Content-Type` header manually when sending `FormData` from browser — the browser will add the correct boundary.
- The backend saves the image in ImgBB (cloud) and stores only the returned URL in `profile_photo`.

--

## SimNikah (application) endpoints — main subsets
All application routes are under `/simnikah` group. Below are the most relevant endpoints the frontend usually consumes.

### Registration (Pendaftaran)
- POST /simnikah/pendaftaran
  - Create a registration (Catin)
  - Auth: Required
  - Body: JSON with applicant data (see backend schema) — large form; forward files/documents as required by specific endpoints
- GET /simnikah/pendaftaran/status
  - Get current user's registration status
  - Auth: Required
- GET /simnikah/pendaftaran
  - List registrations (staff/kepala_kua)
  - Auth: Required + Role middleware (staff, kepala_kua)

### Feedback
- POST /simnikah/feedback-pernikahan
  - Create wedding feedback
  - Auth: Required

### Calendar & Availability
- GET /simnikah/kalender-ketersediaan
- GET /simnikah/ketersediaan-jam
- GET /simnikah/pernikahan-tanggal
  - Mostly public or auth-protected depending on endpoint

### Staff management (Kepala KUA)
- GET /simnikah/staff
  - List staff (role: kepala_kua)
- PUT /simnikah/staff/:id
  - Update staff
- POST /simnikah/staff/pendaftaran
  - Create registration for user (staff privilege)

### Penghulu (officer) endpoints
- GET /simnikah/penghulu
- PUT /simnikah/penghulu/:id
- POST /simnikah/penghulu/verify-documents/:id
- GET /simnikah/penghulu/assigned-registrations
- GET /simnikah/penghulu/today-schedule
- POST /simnikah/penghulu/complete-marriage/:id

### Dashboard (role-specific)
- GET /simnikah/dashboard/kepala-kua
  - Role: kepala_kua
- GET /simnikah/dashboard/staff
  - Role: staff
- GET /simnikah/dashboard/statistik-pernikahan
  - Roles: staff, kepala_kua

### Location services
- POST /simnikah/location/geocode
- POST /simnikah/location/reverse-geocode
- GET /simnikah/location/search
- PUT /simnikah/pendaftaran/:id/location
- GET /simnikah/pendaftaran/:id/location

### Notification system
- POST /simnikah/notifikasi (create)
  - Roles: staff, kepala_kua
- GET /simnikah/notifikasi/user/:user_id
- GET /simnikah/notifikasi/:id
- PUT /simnikah/notifikasi/:id/status
- PUT /simnikah/notifikasi/user/:user_id/mark-all-read
- DELETE /simnikah/notifikasi/:id

--

## Error handling and common errors
- 400 Bad Request: Invalid payload, missing required fields, file errors
- 401 Unauthorized: Missing/invalid token or user inactive
- 403 Forbidden: Role-based access blocked
- 404 Not Found: Resource missing
- 429 Too Many Requests: Rate-limited
- 500 Internal Server Error: Unexpected backend error

Frontend should display `message` and optionally `error` for troubleshooting.

--

## Tips for frontend integration
- Always include `Authorization` header for protected routes.
- Perform client-side validation on forms to reduce server load (but do not rely on it; server-side validates again).
- For file upload, show progress and compress client-side if needed to keep under 5MB.
- Cache profile photo URLs in UI components; when updating a photo, refresh profile data to get new URL.
- Respect rate limits: implement exponential backoff for retries.

## Postman / Collection
- If you want, I can export a Postman collection (or OpenAPI/Swagger) with all endpoints and example requests — tell me which format you prefer.

--

## Where to find more detailed backend docs
- Full API details and examples: `docs/API_DOCUMENTATION_COMPLETE.md` (in repo root `docs/`)
- Developer-level docs and architecture: `DEVELOPER_REFERENCE.md`, `PROJECT_SUMMARY.md`
- Profile-photo-specific doc: `docs/PROFILE_PHOTO_FEATURE.md`

--

If you want, I can:
- Generate a Postman collection or OpenAPI (Swagger) spec from the code
- Add example responses for the large registration endpoints
- Add a small React component snippet to handle upload with preview and immediate profile refresh

Tell me which of the above you want next and I will implement it.
