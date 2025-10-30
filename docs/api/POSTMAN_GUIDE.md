# 📮 SimNikah API - Postman Testing Guide

## 🚀 Quick Setup

### 1. Import to Postman

1. Open Postman
2. Click **Import** button
3. Select **Raw text** tab
4. Copy-paste the collection below
5. Click **Import**

### 2. Setup Environment Variables

Create a new environment in Postman with these variables:

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| `base_url` | `http://localhost:8080` | (same) |
| `token` | (empty) | (will be auto-filled after login) |
| `user_id` | (empty) | (will be auto-filled after login) |
| `pendaftaran_id` | (empty) | (will be set after creating registration) |

---

## 📋 Postman Collection

### Collection: SimNikah API

```json
{
  "info": {
    "name": "SimNikah API",
    "description": "Complete API collection for SimNikah - Marriage Registration System",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "base_url",
      "value": "http://localhost:8080",
      "type": "string"
    },
    {
      "key": "token",
      "value": "",
      "type": "string"
    }
  ],
  "auth": {
    "type": "bearer",
    "bearer": [
      {
        "key": "token",
        "value": "{{token}}",
        "type": "string"
      }
    ]
  },
  "item": [
    {
      "name": "1. Authentication",
      "item": [
        {
          "name": "Register User",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "if (pm.response.code === 201) {",
                  "    const jsonData = pm.response.json();",
                  "    pm.environment.set('user_id', jsonData.user.user_id);",
                  "}"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"username\": \"testuser\",\n  \"email\": \"test@example.com\",\n  \"password\": \"SecurePass123!\",\n  \"nama\": \"Test User\",\n  \"role\": \"user_biasa\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/register",
              "host": ["{{base_url}}"],
              "path": ["register"]
            }
          }
        },
        {
          "name": "Login",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "if (pm.response.code === 200) {",
                  "    const jsonData = pm.response.json();",
                  "    pm.environment.set('token', jsonData.token);",
                  "    pm.environment.set('user_id', jsonData.user.user_id);",
                  "    console.log('Token saved:', jsonData.token);",
                  "}"
                ]
              }
            }
          ],
          "request": {
            "auth": {
              "type": "noauth"
            },
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"username\": \"testuser\",\n  \"password\": \"SecurePass123!\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/login",
              "host": ["{{base_url}}"],
              "path": ["login"]
            }
          }
        },
        {
          "name": "Get Profile",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/profile",
              "host": ["{{base_url}}"],
              "path": ["profile"]
            }
          }
        }
      ]
    },
    {
      "name": "2. Marriage Registration",
      "item": [
        {
          "name": "Check Registration Status",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/pendaftaran/status",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "pendaftaran", "status"]
            }
          }
        },
        {
          "name": "Create Complete Registration",
          "event": [
            {
              "listen": "test",
              "script": {
                "exec": [
                  "if (pm.response.code === 201) {",
                  "    const jsonData = pm.response.json();",
                  "    pm.environment.set('pendaftaran_id', jsonData.data.pendaftaran.id);",
                  "    console.log('Registration ID:', jsonData.data.pendaftaran.id);",
                  "}"
                ]
              }
            }
          ],
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"scheduleAndLocation\": {\n    \"weddingLocation\": \"Di KUA\",\n    \"weddingAddress\": \"KUA Kecamatan Banjarmasin Tengah\",\n    \"weddingDate\": \"2025-12-25\",\n    \"weddingTime\": \"09:00\",\n    \"dispensationNumber\": \"\"\n  },\n  \"groom\": {\n    \"groomFullName\": \"Ahmad Fauzan\",\n    \"groomNik\": \"6371012501950001\",\n    \"groomCitizenship\": \"WNI\",\n    \"groomPassportNumber\": \"\",\n    \"groomPlaceOfBirth\": \"Banjarmasin\",\n    \"groomDateOfBirth\": \"1995-01-25\",\n    \"groomStatus\": \"Belum Kawin\",\n    \"groomReligion\": \"Islam\",\n    \"groomEducation\": \"S1\",\n    \"groomOccupation\": \"Pegawai Swasta\",\n    \"groomOccupationDescription\": \"Staff IT\",\n    \"groomPhoneNumber\": \"081234567890\",\n    \"groomEmail\": \"ahmad@example.com\",\n    \"groomAddress\": \"Jl. Lambung Mangkurat No. 45\"\n  },\n  \"bride\": {\n    \"brideFullName\": \"Siti Aminah\",\n    \"brideNik\": \"6371016702980001\",\n    \"brideCitizenship\": \"WNI\",\n    \"bridePassportNumber\": \"\",\n    \"bridePlaceOfBirth\": \"Banjarmasin\",\n    \"brideDateOfBirth\": \"1998-02-27\",\n    \"brideStatus\": \"Belum Kawin\",\n    \"brideReligion\": \"Islam\",\n    \"brideEducation\": \"S1\",\n    \"brideOccupation\": \"Guru\",\n    \"brideOccupationDescription\": \"Guru SD\",\n    \"bridePhoneNumber\": \"081298765432\",\n    \"brideEmail\": \"siti@example.com\",\n    \"brideAddress\": \"Jl. A. Yani Km 5 No. 12\"\n  },\n  \"groomParents\": {\n    \"groomFather\": {\n      \"groomFatherPresenceStatus\": \"Hidup\",\n      \"groomFatherName\": \"Muhammad Ali\",\n      \"groomFatherNik\": \"6371010101700001\",\n      \"groomFatherCitizenship\": \"WNI\",\n      \"groomFatherCountryOfOrigin\": \"Indonesia\",\n      \"groomFatherPlaceOfBirth\": \"Banjarmasin\",\n      \"groomFatherDateOfBirth\": \"1970-01-01\",\n      \"groomFatherReligion\": \"Islam\",\n      \"groomFatherOccupation\": \"Wiraswasta\",\n      \"groomFatherAddress\": \"Jl. Veteran No. 20\"\n    },\n    \"groomMother\": {\n      \"groomMotherPresenceStatus\": \"Hidup\",\n      \"groomMotherName\": \"Fatimah\",\n      \"groomMotherNik\": \"6371014501720001\",\n      \"groomMotherCitizenship\": \"WNI\",\n      \"groomMotherCountryOfOrigin\": \"Indonesia\",\n      \"groomMotherPlaceOfBirth\": \"Banjarmasin\",\n      \"groomMotherDateOfBirth\": \"1972-05-01\",\n      \"groomMotherReligion\": \"Islam\",\n      \"groomMotherOccupation\": \"Ibu Rumah Tangga\",\n      \"groomMotherAddress\": \"Jl. Veteran No. 20\"\n    }\n  },\n  \"brideParents\": {\n    \"brideFather\": {\n      \"brideFatherPresenceStatus\": \"Hidup\",\n      \"brideFatherName\": \"Abdullah Rahman\",\n      \"brideFatherNik\": \"6371015501680001\",\n      \"brideFatherCitizenship\": \"WNI\",\n      \"brideFatherCountryOfOrigin\": \"Indonesia\",\n      \"brideFatherPlaceOfBirth\": \"Banjarmasin\",\n      \"brideFatherDateOfBirth\": \"1968-05-15\",\n      \"brideFatherReligion\": \"Islam\",\n      \"brideFatherOccupation\": \"PNS\",\n      \"brideFatherAddress\": \"Jl. A. Yani Km 5 No. 12\"\n    },\n    \"brideMother\": {\n      \"brideMotherPresenceStatus\": \"Hidup\",\n      \"brideMotherName\": \"Khadijah\",\n      \"brideMotherNik\": \"6371012201700001\",\n      \"brideMotherCitizenship\": \"WNI\",\n      \"brideMotherCountryOfOrigin\": \"Indonesia\",\n      \"brideMotherPlaceOfBirth\": \"Banjarmasin\",\n      \"brideMotherDateOfBirth\": \"1970-01-22\",\n      \"brideMotherReligion\": \"Islam\",\n      \"brideMotherOccupation\": \"Ibu Rumah Tangga\",\n      \"brideMotherAddress\": \"Jl. A. Yani Km 5 No. 12\"\n    }\n  },\n  \"guardian\": {\n    \"guardianFullName\": \"Abdullah Rahman\",\n    \"guardianNik\": \"6371015501680001\",\n    \"guardianRelationship\": \"Ayah Kandung\",\n    \"guardianStatus\": \"Hidup\",\n    \"guardianReligion\": \"Islam\",\n    \"guardianAddress\": \"Jl. A. Yani Km 5 No. 12\",\n    \"guardianPhoneNumber\": \"081234567899\"\n  }\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/simnikah/pendaftaran/form-baru",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "pendaftaran", "form-baru"]
            }
          }
        },
        {
          "name": "Get Status Flow",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/pendaftaran/{{pendaftaran_id}}/status-flow",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "pendaftaran", "{{pendaftaran_id}}", "status-flow"]
            }
          }
        },
        {
          "name": "Get All Registrations (Staff)",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/pendaftaran?status=Menunggu Verifikasi&page=1&limit=10",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "pendaftaran"],
              "query": [
                {
                  "key": "status",
                  "value": "Menunggu Verifikasi"
                },
                {
                  "key": "page",
                  "value": "1"
                },
                {
                  "key": "limit",
                  "value": "10"
                }
              ]
            }
          }
        }
      ]
    },
    {
      "name": "3. Calendar & Schedule",
      "item": [
        {
          "name": "Get Availability Calendar",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/kalender-ketersediaan?bulan=12&tahun=2025",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "kalender-ketersediaan"],
              "query": [
                {
                  "key": "bulan",
                  "value": "12"
                },
                {
                  "key": "tahun",
                  "value": "2025"
                }
              ]
            }
          }
        },
        {
          "name": "Get Date Detail",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/kalender-tanggal-detail?tanggal=2025-12-25",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "kalender-tanggal-detail"],
              "query": [
                {
                  "key": "tanggal",
                  "value": "2025-12-25"
                }
              ]
            }
          }
        },
        {
          "name": "Get Penghulu Schedule",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/penghulu-jadwal/2025-12-25",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "penghulu-jadwal", "2025-12-25"]
            }
          }
        }
      ]
    },
    {
      "name": "4. Map & Location",
      "item": [
        {
          "name": "Geocoding (Address to Coordinates)",
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"address\": \"Jl. Lambung Mangkurat No. 45, Banjarmasin\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/simnikah/location/geocode",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "location", "geocode"]
            }
          }
        },
        {
          "name": "Reverse Geocoding",
          "request": {
            "method": "POST",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"latitude\": -3.3194374,\n  \"longitude\": 114.5900675\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/simnikah/location/reverse-geocode",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "location", "reverse-geocode"]
            }
          }
        },
        {
          "name": "Search Address",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/location/search?q=Lambung Mangkurat&limit=5",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "location", "search"],
              "query": [
                {
                  "key": "q",
                  "value": "Lambung Mangkurat"
                },
                {
                  "key": "limit",
                  "value": "5"
                }
              ]
            }
          }
        }
      ]
    },
    {
      "name": "5. Notifications",
      "item": [
        {
          "name": "Get User Notifications",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/notifikasi/user/{{user_id}}?status=Belum Dibaca",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "notifikasi", "user", "{{user_id}}"],
              "query": [
                {
                  "key": "status",
                  "value": "Belum Dibaca"
                }
              ]
            }
          }
        },
        {
          "name": "Get Notification Stats",
          "request": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{base_url}}/simnikah/notifikasi/user/{{user_id}}/stats",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "notifikasi", "user", "{{user_id}}", "stats"]
            }
          }
        },
        {
          "name": "Mark as Read",
          "request": {
            "method": "PUT",
            "header": [],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"status_baca\": \"Sudah Dibaca\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            },
            "url": {
              "raw": "{{base_url}}/simnikah/notifikasi/1/status",
              "host": ["{{base_url}}"],
              "path": ["simnikah", "notifikasi", "1", "status"]
            }
          }
        }
      ]
    }
  ]
}
```

---

## 🎯 Testing Workflow

### Step 1: Health Check
```bash
GET {{base_url}}/health
```

### Step 2: Register & Login
1. Run **Register User** (token akan auto-save ke environment)
2. Run **Login** (token akan auto-save ke environment)
3. Run **Get Profile** untuk verify token

### Step 3: Create Marriage Registration
1. Run **Check Registration Status** (pastikan belum ada)
2. Run **Create Complete Registration**
3. Verify response: `pendaftaran_id` akan auto-save

### Step 4: Check Status
1. Run **Get Status Flow**
2. Verify current status: `"Menunggu Verifikasi"`

### Step 5: Test Calendar
1. Run **Get Availability Calendar**
2. Run **Get Date Detail** untuk tanggal tertentu

### Step 6: Test Location Services
1. Run **Geocoding** (alamat → koordinat)
2. Run **Reverse Geocoding** (koordinat → alamat)
3. Run **Search Address** (autocomplete)

### Step 7: Test Notifications
1. Run **Get User Notifications**
2. Run **Get Notification Stats**
3. Run **Mark as Read**

---

## 📝 Pre-request Scripts

### Auto Token Setup
Add this to collection-level Pre-request Script:

```javascript
// Check if token exists
const token = pm.environment.get('token');
if (!token) {
    console.log('⚠️ No token found. Please login first.');
}
```

### Auto Test Scripts
Add this to collection-level Tests:

```javascript
// Log response time
console.log(`⏱️ Response time: ${pm.response.responseTime}ms`);

// Check status code
pm.test("Status code is 200 or 201", function () {
    pm.expect(pm.response.code).to.be.oneOf([200, 201]);
});

// Check response has data
pm.test("Response has data", function () {
    const jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('data').or.property('message');
});
```

---

## 🔧 Environment Variables Template

```json
{
  "name": "SimNikah Development",
  "values": [
    {
      "key": "base_url",
      "value": "http://localhost:8080",
      "enabled": true
    },
    {
      "key": "token",
      "value": "",
      "enabled": true
    },
    {
      "key": "user_id",
      "value": "",
      "enabled": true
    },
    {
      "key": "pendaftaran_id",
      "value": "",
      "enabled": true
    }
  ]
}
```

---

## 📌 Tips

1. **Always check token** - Jika request return 401, login ulang
2. **Use variables** - Gunakan `{{pendaftaran_id}}` untuk dynamic testing
3. **Check response time** - Should be < 1500ms for registration
4. **Test error cases** - Try invalid NIK, past dates, etc.
5. **Sequential testing** - Follow the workflow order

---

**Happy Testing! 🚀**

