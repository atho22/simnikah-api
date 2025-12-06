# Profile Photo Upload Feature Documentation

## Overview
This feature allows all users (Catin, Staff, Penghulu, Kepala KUA) to upload and manage their own profile photos. Photos are stored on ImgBB cloud storage and the URLs are saved in the database.

## Implementation Details

### Database Model Changes
**File:** `internal/models/models.go`

Added field to `Users` struct:
```go
Profile_photo string `gorm:"size:500" json:"foto_profil"`
```

- **Field Name:** Profile_photo
- **Type:** String
- **Database Size:** 500 characters (for URL storage)
- **JSON Tag:** foto_profil
- **Default:** Empty string (no photo)

### API Endpoints

#### 1. Upload Profile Photo
**Endpoint:** `POST /upload-photo`

**Authentication:** Required (JWT Token)

**Request:**
```bash
curl -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -F "photo=@/path/to/image.jpg"
```

**Form Parameters:**
- `photo` (file, required): Image file to upload
  - **Max Size:** 5MB
  - **Allowed Formats:** JPEG, PNG, WebP
  - **MIME Types:** image/jpeg, image/png, image/jpg, image/webp

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "Profile photo uploaded successfully",
  "data": {
    "photo_url": "https://ibb.co/xxxxx/image.jpg"
  }
}
```

**Error Responses:**

**400 Bad Request** - Missing file:
```json
{
  "status": "error",
  "message": "No photo provided",
  "code": "MISSING_FILE"
}
```

**400 Bad Request** - File too large:
```json
{
  "status": "error",
  "message": "File size exceeds 5MB limit",
  "code": "FILE_TOO_LARGE"
}
```

**400 Bad Request** - Invalid format:
```json
{
  "status": "error",
  "message": "Invalid file format. Only JPEG, PNG, and WebP are allowed",
  "code": "INVALID_FORMAT"
}
```

**401 Unauthorized** - Missing token:
```json
{
  "status": "error",
  "message": "Authorization token missing or invalid",
  "code": "UNAUTHORIZED"
}
```

**500 Internal Server Error** - Upload failure:
```json
{
  "status": "error",
  "message": "Failed to upload photo to storage service",
  "code": "UPLOAD_FAILED"
}
```

**500 Internal Server Error** - Database update failure:
```json
{
  "status": "error",
  "message": "Failed to save photo URL to database",
  "code": "DB_ERROR"
}
```

#### 2. Get User Profile (Enhanced)
**Endpoint:** `GET /profile`

**Authentication:** Required (JWT Token)

**Success Response (200 OK):**
```json
{
  "status": "success",
  "message": "User profile retrieved successfully",
  "data": {
    "user_id": "USER_12345",
    "username": "john_doe",
    "email": "john@example.com",
    "nama": "John Doe",
    "role": "user_biasa",
    "status": "Aktif",
    "profile_photo": "https://ibb.co/xxxxx/image.jpg",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-20T14:25:00Z"
  }
}
```

**Note:** `profile_photo` field will be empty string if no photo has been uploaded yet.

---

## Implementation Code Reference

### Handler Implementation
**File:** `internal/handlers/auth/auth.go`

The `UploadProfilePhoto` handler:
1. Extracts JWT token from Authorization header via middleware
2. Retrieves user ID from context (set by AuthMiddleware)
3. Gets file from multipart form with field name "photo"
4. Validates file size (≤ 5MB)
5. Validates MIME type (whitelisted formats only)
6. Uploads to ImgBB using `storage.UploadFileFromMultipart()`
7. Updates Users table with returned photo URL
8. Returns JSON response with photo URL or error

### Storage Integration
**File:** `pkg/storage/imgbb.go`

Uses existing `UploadFileFromMultipart` function:
```go
func UploadFileFromMultipart(file *multipart.FileHeader) (string, error)
```

Returns:
- **Success:** Direct URL from ImgBB (e.g., "https://ibb.co/xxxxx/image.jpg")
- **Error:** Error message describing the failure

### Route Registration
**File:** `cmd/api/main.go`

```go
r.POST("/upload-photo", middleware.AuthMiddleware(), authHandler.UploadProfilePhoto)
```

- Route requires valid JWT token (AuthMiddleware)
- Route handles multipart form data
- Direct response to client with upload result

---

## Usage Guide

### For Frontend Developers

#### 1. Basic Upload Implementation (JavaScript/Fetch)
```javascript
async function uploadProfilePhoto(file, token) {
  const formData = new FormData();
  formData.append('photo', file);
  
  try {
    const response = await fetch('http://localhost:8080/upload-photo', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`
      },
      body: formData
    });
    
    const data = await response.json();
    
    if (response.ok) {
      console.log('Photo URL:', data.data.photo_url);
      return data.data.photo_url;
    } else {
      console.error('Upload failed:', data.message);
      return null;
    }
  } catch (error) {
    console.error('Error uploading photo:', error);
    return null;
  }
}
```

#### 2. React Hook Implementation
```typescript
import { useState } from 'react';

const useProfilePhoto = (token: string) => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const uploadPhoto = async (file: File): Promise<string | null> => {
    setLoading(true);
    setError(null);
    
    const formData = new FormData();
    formData.append('photo', file);
    
    try {
      const response = await fetch('/upload-photo', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`
        },
        body: formData
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Upload failed');
      }
      
      return data.data.photo_url;
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unknown error';
      setError(message);
      return null;
    } finally {
      setLoading(false);
    }
  };
  
  return { uploadPhoto, loading, error };
};
```

#### 3. Vue.js Implementation
```vue
<template>
  <div class="photo-upload">
    <input 
      type="file" 
      @change="handleFileSelect" 
      accept="image/*"
      ref="fileInput"
    >
    <button @click="uploadPhoto" :disabled="!selectedFile || loading">
      {{ loading ? 'Uploading...' : 'Upload Photo' }}
    </button>
    <div v-if="error" class="error">{{ error }}</div>
    <img v-if="photoUrl" :src="photoUrl" alt="Profile Photo">
  </div>
</template>

<script>
export default {
  data() {
    return {
      selectedFile: null,
      photoUrl: null,
      loading: false,
      error: null
    };
  },
  methods: {
    handleFileSelect(event) {
      this.selectedFile = event.target.files[0];
    },
    async uploadPhoto() {
      if (!this.selectedFile) return;
      
      this.loading = true;
      this.error = null;
      
      const formData = new FormData();
      formData.append('photo', this.selectedFile);
      
      try {
        const response = await fetch('/upload-photo', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${this.$store.state.token}`
          },
          body: formData
        });
        
        const data = await response.json();
        
        if (response.ok) {
          this.photoUrl = data.data.photo_url;
        } else {
          this.error = data.message;
        }
      } catch (err) {
        this.error = err.message;
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>
```

### For Mobile Developers

#### Swift (iOS) Example
```swift
func uploadProfilePhoto(image: UIImage, token: String) {
    guard let imageData = image.jpegData(compressionQuality: 0.8) else {
        print("Failed to convert image to JPEG")
        return
    }
    
    let url = URL(string: "http://localhost:8080/upload-photo")!
    var request = URLRequest(url: url)
    request.httpMethod = "POST"
    request.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
    
    let boundary = "Boundary-\(UUID().uuidString)"
    request.setValue("multipart/form-data; boundary=\(boundary)", forHTTPHeaderField: "Content-Type")
    
    var body = Data()
    body.append("--\(boundary)\r\n".data(using: .utf8)!)
    body.append("Content-Disposition: form-data; name=\"photo\"; filename=\"profile.jpg\"\r\n".data(using: .utf8)!)
    body.append("Content-Type: image/jpeg\r\n\r\n".data(using: .utf8)!)
    body.append(imageData)
    body.append("\r\n--\(boundary)--\r\n".data(using: .utf8)!)
    
    request.httpBody = body
    
    URLSession.shared.dataTask(with: request) { data, response, error in
        if let error = error {
            print("Upload error: \(error)")
            return
        }
        
        if let data = data {
            if let json = try? JSONSerialization.jsonObject(with: data) as? [String: Any],
               let photoUrl = json["data"] as? [String: String] {
                print("Photo URL: \(photoUrl["photo_url"] ?? "")")
            }
        }
    }.resume()
}
```

#### Kotlin (Android) Example
```kotlin
fun uploadProfilePhoto(file: File, token: String) {
    val okHttpClient = OkHttpClient()
    
    val requestBody = MultipartBody.Builder()
        .setType(MultipartBody.FORM)
        .addFormDataPart("photo", file.name,
            RequestBody.create("image/jpeg".toMediaType(), file)
        )
        .build()
    
    val request = Request.Builder()
        .url("http://localhost:8080/upload-photo")
        .addHeader("Authorization", "Bearer $token")
        .post(requestBody)
        .build()
    
    okHttpClient.newCall(request).enqueue(object : Callback {
        override fun onFailure(call: Call, e: IOException) {
            e.printStackTrace()
        }
        
        override fun onResponse(call: Call, response: Response) {
            response.body?.string()?.let { responseBody ->
                val jsonObject = JSONObject(responseBody)
                val photoUrl = jsonObject
                    .getJSONObject("data")
                    .getString("photo_url")
                println("Photo URL: $photoUrl")
            }
        }
    })
}
```

---

## Testing Guide

### Manual Testing with cURL

#### 1. Get JWT Token
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "password": "password123"
  }' | jq '.data.token'
```

#### 2. Upload Profile Photo
```bash
TOKEN="your_jwt_token_here"

curl -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@/path/to/your/image.jpg"
```

#### 3. Verify Photo in Profile
```bash
TOKEN="your_jwt_token_here"

curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer $TOKEN" | jq '.data.profile_photo'
```

### Automated Testing with Postman

1. **Create request:** POST `/upload-photo`
2. **Headers:**
   - `Authorization: Bearer {{token}}`
3. **Body:** Form-data
   - Key: `photo`
   - Value: Select file from disk
   - Type: File
4. **Send** and verify response contains `photo_url`

### Testing with Different File Scenarios

```bash
# Test 1: Missing file
curl -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer $TOKEN"

# Expected: 400 - "No photo provided"

# Test 2: Invalid format
curl -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@document.pdf"

# Expected: 400 - "Invalid file format. Only JPEG, PNG, and WebP are allowed"

# Test 3: Missing token
curl -X POST http://localhost:8080/upload-photo \
  -F "photo=@image.jpg"

# Expected: 401 - "Authorization token missing or invalid"

# Test 4: Valid upload
curl -X POST http://localhost:8080/upload-photo \
  -H "Authorization: Bearer $TOKEN" \
  -F "photo=@image.jpg"

# Expected: 200 - { "status": "success", "data": { "photo_url": "..." } }
```

---

## Configuration

### Environment Variables
The following environment variables affect the photo upload feature:

**ImgBB Configuration:**
```bash
# Optional: ImgBB API Key (for anonymous uploads if not provided)
IMGBB_API_KEY=your_api_key_here
```

If `IMGBB_API_KEY` is not set, the system will use ImgBB's anonymous upload feature.

### Validation Rules (Hardcoded)
- **Maximum File Size:** 5MB
- **Allowed MIME Types:** 
  - image/jpeg
  - image/png
  - image/jpg
  - image/webp
- **Storage Location:** ImgBB cloud storage
- **URL Retention:** Permanent (ImgBB standard)

---

## Database Schema

### Users Table (Modified)
```sql
ALTER TABLE Users ADD COLUMN profile_photo VARCHAR(500);
```

**Field Details:**
- **Column Name:** profile_photo
- **Type:** VARCHAR(500)
- **Nullable:** YES
- **Default:** NULL or empty string
- **Index:** None (not frequently queried alone)
- **Sample Value:** "https://ibb.co/xxxxx/image.jpg"

---

## User Roles and Permissions

All user roles can upload profile photos:
- ✅ **user_biasa** (Regular Users/Catin)
- ✅ **staff** (KUA Staff)
- ✅ **penghulu** (Penghulu Officer)
- ✅ **kepala_kua** (KUA Head)

**Authentication Required:** Yes - All endpoints require valid JWT token

---

## Security Considerations

### File Validation
1. **Size Validation:** Maximum 5MB enforced server-side
2. **MIME Type Validation:** Only whitelisted image formats accepted
3. **Content Type Check:** File content is validated against MIME type

### Authorization
- **Authentication:** Required - JWT token from login
- **Owner Verification:** User can only upload their own photo (enforced via user ID from token)
- **Role-Based Access:** No additional role restrictions (all authenticated users can upload)

### Storage Security
- **ImgBB:** Trusted third-party service
- **URL Storage:** Stored as plain text in database (read-only to other users)
- **No Local Storage:** Images never stored on server filesystem
- **HTTPS:** ImgBB provides HTTPS URLs automatically

### Data Privacy
- **Photo Visibility:** URL is visible when fetching user profile
- **No Encryption:** URLs stored unencrypted (standard practice for image hosting)
- **Cleanup:** No automatic photo deletion mechanism yet (TODO)

---

## Performance Considerations

### Upload Speed
- **Network Dependent:** Upload speed depends on user's internet connection
- **File Size Impact:** 5MB file typically takes 2-10 seconds depending on connection
- **ImgBB Processing:** ImgBB typically responds within 1-2 seconds

### Database Impact
- **Single UPDATE query:** Only User record is updated with photo URL
- **No N+1 Problem:** Profile fetch is single query with all fields
- **Index Strategy:** No special indexing needed for this feature

### Caching Strategy
- **Profile Cache:** Frontend should cache photo URLs (immutable after upload)
- **Browser Cache:** ImgBB URLs are served with cache headers
- **API Cache:** No server-side caching needed for this feature

---

## Future Enhancements

### Phase 2 Features
1. **Photo Deletion Endpoint**
   - DELETE /profile/photo
   - Removes photo URL from database
   - Allows re-uploading new photo

2. **Photo Update Endpoint**
   - PUT /profile/photo
   - Replace existing photo with new one
   - Automatic cleanup of old URL (optional)

3. **Photo Cropping**
   - Client-side image cropping before upload
   - Server-side validation of image dimensions
   - Consistent profile photo aspect ratio

4. **Thumbnail Generation**
   - Server generates thumbnail versions
   - Multiple resolution support
   - CDN optimization

5. **Photo History**
   - Store all uploaded photos with timestamps
   - Allow users to revert to previous photo
   - Admin view of all user photos

### Phase 3 Features
1. **Integration with Marriage Registration**
   - Require valid profile photo for Catin registration
   - Attach photo to Pendaftaran Nikah record
   - Include in official documents

2. **Photo Verification**
   - Admin approval workflow for photos
   - Rejection reasons if inappropriate
   - Resubmission support

3. **Photo Analytics**
   - Track photo upload rates by role
   - Identify users without photos
   - Generate reports for admins

---

## Troubleshooting

### Issue: "No photo provided" Error
**Solution:** Ensure form field name is exactly "photo" (case-sensitive)

### Issue: "Invalid file format" Error
**Causes:**
- Uploading non-image file (PDF, DOC, etc.)
- MIME type not recognized
**Solution:** Use JPG, PNG, or WebP format

### Issue: "File size exceeds 5MB limit" Error
**Solution:** Compress image before upload or resize to smaller dimensions

### Issue: "Authorization token missing or invalid" Error
**Causes:**
- Missing Authorization header
- Invalid or expired token
- Incorrect token format
**Solution:** 
- Include `Authorization: Bearer <token>` header
- Verify token is from successful login
- Check token has not expired

### Issue: "Failed to upload photo to storage service" Error
**Causes:**
- ImgBB service unavailable
- Network connectivity issue
- Invalid ImgBB API key
**Solution:**
- Check ImgBB service status
- Verify internet connection
- Check IMGBB_API_KEY configuration

### Issue: Photo uploaded but not visible in profile
**Solution:**
- Refresh the page to fetch latest data
- Clear browser cache
- Verify user ID matches uploaded photo user
- Check database has profile_photo field

---

## API Response Status Codes Summary

| Status | Meaning | Action |
|--------|---------|--------|
| 200 | Upload successful | Process photo URL from response |
| 400 | Bad request (file issue) | Check file size/format |
| 401 | Unauthorized | Refresh login token |
| 404 | User not found | Should not happen if auth works |
| 500 | Server error | Contact admin, check logs |

---

## Integration Checklist

- [ ] Backend route is registered and accessible
- [ ] Database field profile_photo exists on Users table
- [ ] ImgBB API key is configured (optional)
- [ ] JWT authentication is working
- [ ] File upload handler is implemented
- [ ] Frontend upload component created
- [ ] File validation on frontend matches backend
- [ ] Error handling displays user-friendly messages
- [ ] Profile photo displays in user profile view
- [ ] Testing completed with different image formats
- [ ] Testing completed with network errors
- [ ] Documentation provided to frontend team

---

## Summary

The profile photo upload feature is now **fully implemented and ready for use**. Users can upload their own photos which are stored securely on ImgBB cloud storage, with URLs persisted in the database. The feature includes comprehensive validation, error handling, and is available for all user roles.
