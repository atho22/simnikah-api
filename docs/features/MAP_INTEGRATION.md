# 🗺️ Map Integration - SimNikah (100% GRATIS dengan OpenStreetMap)

## 📋 Pendahuluan

Fitur Map Integration memungkinkan pengguna dan penghulu untuk:
- **Input alamat** dengan validasi koordinat GPS
- **Melihat lokasi** nikah di peta menggunakan **OpenStreetMap (GRATIS)**
- **Mendapatkan rute** ke lokasi pernikahan
- **Autocomplete alamat** untuk kemudahan input

Khusus untuk **pernikahan di luar KUA**, sistem akan mencatat koordinat GPS (latitude, longitude) untuk memudahkan penghulu menemukan lokasi.

> 🎉 **100% GRATIS** - Menggunakan **OpenStreetMap Nominatim API** tanpa perlu API Key atau Billing!

---

## 🚀 Fitur Utama

### 1. Geocoding (Alamat → Koordinat)
Convert alamat menjadi koordinat GPS menggunakan OpenStreetMap Nominatim API

### 2. Reverse Geocoding (Koordinat → Alamat)  
Convert koordinat GPS menjadi alamat menggunakan OpenStreetMap Nominatim API

### 3. Address Autocomplete (Search)
Search dan suggest alamat saat user mengetik (minimal 3 karakter)

### 4. Map Display dengan Leaflet.js
- **Leaflet.js** - Library JavaScript gratis untuk OpenStreetMap
- **OpenStreetMap tiles** - Peta gratis tanpa batas
- **Marker dan popup** - Tandai lokasi nikah

### 5. Navigation Links untuk Penghulu
- Link ke **OpenStreetMap** (view di browser)
- Link ke **Google Maps** (hanya untuk navigasi, tidak perlu API key)
- Link ke **Waze** (navigasi mobile)

---

## 📡 API Endpoints

### Base URL
```
http://localhost:8080/simnikah
```

### Authentication
Semua endpoint memerlukan authentication token di header:
```
Authorization: Bearer <your_token>
```

---

## 1. Geocoding - Alamat ke Koordinat

**Endpoint:** `POST /location/geocode`

**Deskripsi:** Mendapatkan koordinat GPS dari alamat.

**Request Body:**
```json
{
  "alamat": "Jl. Pangeran Antasari No.1, Banjarmasin, Kalimantan Selatan"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Koordinat berhasil ditemukan",
  "data": {
    "alamat": "Jl. Pangeran Antasari No.1, Banjarmasin, Kalimantan Selatan",
    "latitude": -3.3149,
    "longitude": 114.5925,
    "map_url": "https://www.google.com/maps?q=-3.3149,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3149&mlon=114.5925&zoom=16"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "success": false,
  "message": "Lokasi tidak ditemukan",
  "error": "Alamat tidak dapat ditemukan di peta. Silakan periksa kembali alamat yang dimasukkan."
}
```

**Use Case:**
- User mengetik alamat lengkap
- Sistem convert ke koordinat untuk disimpan di database
- Tampilkan pin di map

---

## 2. Reverse Geocoding - Koordinat ke Alamat

**Endpoint:** `POST /location/reverse-geocode`

**Deskripsi:** Mendapatkan alamat dari koordinat GPS (saat user pilih lokasi di map).

**Request Body:**
```json
{
  "latitude": -3.3149,
  "longitude": 114.5925
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Alamat berhasil ditemukan",
  "data": {
    "latitude": -3.3149,
    "longitude": 114.5925,
    "alamat": "Jalan Pangeran Antasari, Sungai Baru, Banjarmasin Tengah, Kota Banjarmasin, Kalimantan Selatan, 70123, Indonesia",
    "detail": {
      "road": "Jalan Pangeran Antasari",
      "suburb": "Sungai Baru",
      "city_district": "Banjarmasin Tengah",
      "city": "Kota Banjarmasin",
      "state": "Kalimantan Selatan",
      "postcode": "70123",
      "country": "Indonesia"
    },
    "map_url": "https://www.google.com/maps?q=-3.3149,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3149&mlon=114.5925&zoom=16"
  }
}
```

**Use Case:**
- User klik lokasi di map
- Sistem otomatis isi field alamat
- User konfirmasi atau edit

---

## 3. Search Address (Autocomplete)

**Endpoint:** `GET /location/search?q={query}`

**Deskripsi:** Search alamat untuk autocomplete (minimal 3 karakter).

**Request:**
```
GET /simnikah/location/search?q=Banjarmasin
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Hasil pencarian alamat",
  "data": {
    "query": "Banjarmasin",
    "results": [
      {
        "display_name": "Banjarmasin, Kalimantan Selatan, Indonesia",
        "latitude": "-3.3149",
        "longitude": "114.5925",
        "address": {
          "city": "Banjarmasin",
          "state": "Kalimantan Selatan",
          "country": "Indonesia"
        }
      },
      {
        "display_name": "Banjarmasin Utara, Banjarmasin, Kalimantan Selatan, Indonesia",
        "latitude": "-3.2896",
        "longitude": "114.5972",
        "address": {
          "city_district": "Banjarmasin Utara",
          "city": "Banjarmasin",
          "state": "Kalimantan Selatan",
          "country": "Indonesia"
        }
      }
    ],
    "count": 2
  }
}
```

**Use Case:**
- User mulai mengetik alamat
- Tampilkan suggestion dropdown
- User pilih dari list

---

## 4. Update Lokasi Nikah dengan Koordinat

**Endpoint:** `PUT /pendaftaran/:id/location`

**Deskripsi:** Update alamat dan koordinat lokasi nikah.

**Request Body:**
```json
{
  "alamat_akad": "Jl. Pangeran Antasari No.1, Banjarmasin",
  "latitude": -3.3149,
  "longitude": 114.5925
}
```

**Note:** Jika `latitude` dan `longitude` tidak disediakan, sistem akan otomatis geocoding dari `alamat_akad`.

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Lokasi nikah berhasil diupdate",
  "data": {
    "pendaftaran_id": 123,
    "nomor_pendaftaran": "NIK20250127001",
    "alamat_akad": "Jl. Pangeran Antasari No.1, Banjarmasin",
    "tempat_nikah": "Di Luar KUA",
    "latitude": -3.3149,
    "longitude": 114.5925,
    "map_url": "https://www.google.com/maps?q=-3.3149,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3149&mlon=114.5925&zoom=16",
    "updated_at": "2025-01-27T10:00:00Z"
  }
}
```

---

## 5. Get Detail Lokasi untuk Penghulu

**Endpoint:** `GET /pendaftaran/:id/location`

**Deskripsi:** Mendapatkan detail lokasi nikah dengan semua link navigasi (untuk penghulu).

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Detail lokasi nikah berhasil diambil",
  "data": {
    "pendaftaran_id": 123,
    "nomor_pendaftaran": "NIK20250127001",
    "tanggal_nikah": "2025-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Pangeran Antasari No.1, Banjarmasin, Kalimantan Selatan",
    "latitude": -3.3149,
    "longitude": 114.5925,
    "has_coordinates": true,
    "is_outside_kua": true,
    "note": "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi.",
    "map_url": "https://www.google.com/maps?q=-3.3149,114.5925",
    "google_maps_url": "https://www.google.com/maps/search/?api=1&query=-3.3149,114.5925",
    "google_maps_directions_url": "https://www.google.com/maps/dir/?api=1&destination=-3.3149,114.5925",
    "osm_url": "https://www.openstreetmap.org/?mlat=-3.3149&mlon=114.5925&zoom=16",
    "waze_url": "https://www.waze.com/ul?ll=-3.3149,114.5925&navigate=yes"
  }
}
```

**Tanpa Koordinat (jika belum di-set):**
```json
{
  "success": true,
  "message": "Detail lokasi nikah berhasil diambil",
  "data": {
    "pendaftaran_id": 123,
    "nomor_pendaftaran": "NIK20250127001",
    "tanggal_nikah": "2025-02-14",
    "waktu_nikah": "09:00",
    "tempat_nikah": "Di Luar KUA",
    "alamat_akad": "Jl. Pangeran Antasari No.1, Banjarmasin",
    "has_coordinates": false,
    "message": "Koordinat belum tersedia untuk lokasi ini",
    "is_outside_kua": true,
    "note": "Pernikahan dilaksanakan di luar KUA. Penghulu perlu datang ke lokasi."
  }
}
```

---

## 🎨 Frontend Integration

### React Example - Geocoding Alamat

```javascript
// Component untuk input alamat dengan geocoding
const AddressInput = () => {
  const [alamat, setAlamat] = useState('');
  const [coordinates, setCoordinates] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleGeocoding = async () => {
    if (alamat.length < 10) {
      alert('Masukkan alamat lengkap minimal 10 karakter');
      return;
    }

    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/simnikah/location/geocode', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ alamat })
      });

      const data = await response.json();
      
      if (data.success) {
        setCoordinates({
          lat: data.data.latitude,
          lng: data.data.longitude
        });
        console.log('Koordinat:', data.data);
        // Tampilkan di map
      } else {
        alert(data.error);
      }
    } catch (error) {
      console.error('Error:', error);
      alert('Gagal mendapatkan koordinat');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <input
        type="text"
        placeholder="Masukkan alamat lengkap..."
        value={alamat}
        onChange={(e) => setAlamat(e.target.value)}
      />
      <button onClick={handleGeocoding} disabled={loading}>
        {loading ? 'Loading...' : 'Cari Koordinat'}
      </button>
      
      {coordinates && (
        <div>
          <p>Lat: {coordinates.lat}, Lng: {coordinates.lng}</p>
          <MapComponent center={coordinates} zoom={16} />
        </div>
      )}
    </div>
  );
};
```

### React Example - Address Autocomplete

```javascript
import { useState, useEffect } from 'react';
import { debounce } from 'lodash';

const AddressAutocomplete = ({ onSelect }) => {
  const [query, setQuery] = useState('');
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);

  // Debounced search function
  const searchAddress = debounce(async (searchQuery) => {
    if (searchQuery.length < 3) {
      setSuggestions([]);
      return;
    }

    setLoading(true);
    try {
      const response = await fetch(
        `http://localhost:8080/simnikah/location/search?q=${encodeURIComponent(searchQuery)}`,
        {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      );

      const data = await response.json();
      
      if (data.success) {
        setSuggestions(data.data.results);
      }
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  }, 500); // Wait 500ms after user stops typing

  useEffect(() => {
    searchAddress(query);
  }, [query]);

  const handleSelect = (suggestion) => {
    setQuery(suggestion.display_name);
    setSuggestions([]);
    onSelect({
      alamat: suggestion.display_name,
      latitude: parseFloat(suggestion.latitude),
      longitude: parseFloat(suggestion.longitude)
    });
  };

  return (
    <div className="autocomplete">
      <input
        type="text"
        placeholder="Ketik alamat (min 3 karakter)..."
        value={query}
        onChange={(e) => setQuery(e.target.value)}
      />
      
      {loading && <div>Mencari...</div>}
      
      {suggestions.length > 0 && (
        <ul className="suggestions">
          {suggestions.map((suggestion, index) => (
            <li key={index} onClick={() => handleSelect(suggestion)}>
              {suggestion.display_name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};
```

### React Example - Leaflet + OpenStreetMap (GRATIS)

**Install dependencies:**
```bash
npm install leaflet react-leaflet
# atau
yarn add leaflet react-leaflet
```

**Import CSS di `index.html` atau `App.js`:**
```html
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
```

**Component:**
```javascript
import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';

// Fix default marker icon (Leaflet issue with bundlers)
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
});

const MapComponent = ({ latitude, longitude, alamat }) => {
  const position = [latitude, longitude];

  return (
    <MapContainer 
      center={position} 
      zoom={16} 
      style={{ height: '400px', width: '100%' }}
    >
      {/* OpenStreetMap tiles (GRATIS) */}
      <TileLayer
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
      />
      
      {/* Marker lokasi nikah */}
      <Marker position={position}>
        <Popup>
          <strong>Lokasi Nikah</strong><br/>
          {alamat}
        </Popup>
      </Marker>
    </MapContainer>
  );
};
```

### React Example - View untuk Penghulu

```javascript
const PenghuluLocationView = ({ registrationId }) => {
  const [location, setLocation] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchLocation();
  }, [registrationId]);

  const fetchLocation = async () => {
    try {
      const response = await fetch(
        `http://localhost:8080/simnikah/pendaftaran/${registrationId}/location`,
        {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        }
      );

      const data = await response.json();
      
      if (data.success) {
        setLocation(data.data);
      }
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Loading...</div>;
  if (!location) return <div>Data tidak ditemukan</div>;

  return (
    <div className="location-detail">
      <h2>Detail Lokasi Pernikahan</h2>
      
      <div className="info">
        <p><strong>Nomor:</strong> {location.nomor_pendaftaran}</p>
        <p><strong>Tanggal:</strong> {location.tanggal_nikah}</p>
        <p><strong>Waktu:</strong> {location.waktu_nikah}</p>
        <p><strong>Lokasi:</strong> {location.tempat_nikah}</p>
        <p><strong>Alamat:</strong> {location.alamat_akad}</p>
      </div>

      {location.has_coordinates ? (
        <>
          <div className="map">
            <MapComponent
              latitude={location.latitude}
              longitude={location.longitude}
              alamat={location.alamat_akad}
            />
          </div>

          <div className="navigation-buttons">
            <a 
              href={location.google_maps_directions_url} 
              target="_blank" 
              className="btn btn-primary"
            >
              🗺️ Buka di Google Maps
            </a>
            <a 
              href={location.waze_url} 
              target="_blank" 
              className="btn btn-info"
            >
              🚗 Navigasi dengan Waze
            </a>
            <a 
              href={location.osm_url} 
              target="_blank" 
              className="btn btn-secondary"
            >
              🌍 Lihat di OpenStreetMap
            </a>
          </div>
        </>
      ) : (
        <div className="alert alert-warning">
          ⚠️ Koordinat belum tersedia untuk lokasi ini.
        </div>
      )}

      {location.is_outside_kua && (
        <div className="alert alert-info">
          ℹ️ {location.note}
        </div>
      )}
    </div>
  );
};
```

---

## 🌍 Map Services (Semua GRATIS!)

### 1. OpenStreetMap Nominatim API (Backend Geocoding)
✅ **100% FREE** - Tidak perlu API key  
✅ Geocoding & Reverse Geocoding  
✅ Address search  

**Requirements:**
- Rate limit: **1 request/second** (gunakan debounce di frontend)
- **User-Agent header required**: `SimNikah-KUA-Banjarmasin/1.0`
- Filter Indonesia: `countrycodes=id`

**Dokumentasi:** https://nominatim.org/release-docs/latest/api/Overview/

**Sudah terimplementasi di:**
- `helper/utils.go` - `GetCoordinatesFromAddress()`
- `catin/location.go` - Semua endpoint

---

### 2. Leaflet.js + OpenStreetMap Tiles (Frontend Display)
✅ **100% FREE** - Tidak perlu API key  
✅ Display peta interaktif  
✅ Marker, popup, zoom, pan  

**Install:**
```bash
npm install leaflet react-leaflet
```

**Tiles URL:** `https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png`

**Dokumentasi:** 
- https://leafletjs.com/
- https://react-leaflet.js.org/

---

### 3. Google Maps & Waze (Navigation Links Only)
✅ **TIDAK PERLU API KEY** - Hanya deep links untuk navigasi

**Google Maps Navigation:**
```javascript
// Link langsung ke Google Maps (browser)
const googleMapsUrl = `https://www.google.com/maps?q=${lat},${lon}`;

// Link untuk navigasi
const directionsUrl = `https://www.google.com/maps/dir/?api=1&destination=${lat},${lon}`;
```

**Waze Navigation:**
```javascript
const wazeUrl = `https://www.waze.com/ul?ll=${lat},${lon}&navigate=yes`;
```

> ⚠️ **CATATAN:** Kita TIDAK menggunakan Google Maps JavaScript API (yang berbayar). Kita hanya pakai deep links untuk navigasi (GRATIS).

---

## 💡 Best Practices

### 1. Input Alamat di Frontend

**DO:**
✅ Gunakan autocomplete untuk suggest alamat  
✅ Validasi minimal 10 karakter  
✅ Tampilkan preview map setelah geocoding  
✅ Beri opsi untuk adjust pin di map  
✅ Save koordinat bersamaan dengan alamat  

**DON'T:**
❌ Kirim request geocoding setiap keystroke (pakai debounce!)  
❌ Skip validasi alamat  
❌ Simpan alamat tanpa koordinat  

### 2. Tampilan untuk Penghulu

**DO:**
✅ Tampilkan map dengan marker lokasi  
✅ Berikan button "Buka di Google Maps" dan "Waze"  
✅ Tampilkan jarak dari KUA ke lokasi (optional)  
✅ Berikan info lengkap (tanggal, waktu, alamat)  

### 3. Error Handling

```javascript
try {
  // Call geocoding API
} catch (error) {
  if (error.status === 404) {
    alert('Alamat tidak ditemukan. Silakan cek kembali.');
  } else if (error.status === 429) {
    alert('Terlalu banyak request. Tunggu sebentar.');
  } else {
    alert('Terjadi kesalahan. Coba lagi nanti.');
  }
}
```

### 4. Rate Limiting

Nominatim API has 1 request/second limit:
```javascript
// Debounce untuk autocomplete
const debouncedSearch = debounce(searchAddress, 500);

// Atau gunakan rate limiter
import pThrottle from 'p-throttle';
const throttle = pThrottle({ limit: 1, interval: 1000 });
const geocode = throttle(async (address) => { /* ... */ });
```

---

## 🔧 Testing

### Test Geocoding
```bash
curl -X POST http://localhost:8080/simnikah/location/geocode \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "alamat": "Jl. Pangeran Antasari No.1, Banjarmasin, Kalimantan Selatan"
  }'
```

### Test Reverse Geocoding
```bash
curl -X POST http://localhost:8080/simnikah/location/reverse-geocode \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "latitude": -3.3149,
    "longitude": 114.5925
  }'
```

### Test Search Address
```bash
curl -X GET "http://localhost:8080/simnikah/location/search?q=Banjarmasin" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Test Get Location Detail
```bash
curl -X GET http://localhost:8080/simnikah/pendaftaran/123/location \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## 🐛 Troubleshooting

### Problem: Koordinat tidak ditemukan
**Solution:**
- Pastikan alamat lengkap (jalan, kelurahan, kota)
- Tambahkan landmark terkenal
- Gunakan format: "Jalan, Kelurahan, Kecamatan, Kota, Provinsi"

### Problem: Rate limit exceeded (429)
**Solution:**
- Implementasi debounce di frontend
- Tunggu 1 detik antara request
- Pertimbangkan cache hasil geocoding

### Problem: Map tidak muncul di frontend (Leaflet)
**Solution:**
- Pastikan CSS Leaflet sudah di-import:
  ```javascript
  import 'leaflet/dist/leaflet.css';
  ```
- Pastikan icon fix sudah ditambahkan (lihat contoh di atas)
- Check console browser untuk error
- Pastikan container punya height (misal: `style={{ height: '400px' }}`)
- Check koneksi internet untuk loading tiles

---

## 📚 Referensi

### Backend (Go)
- [OpenStreetMap Nominatim API](https://nominatim.org/release-docs/latest/api/Overview/) - Geocoding API (GRATIS)
- [Nominatim Usage Policy](https://operations.osmfoundation.org/policies/nominatim/) - Rate limits & best practices

### Frontend (React)
- [Leaflet.js Documentation](https://leafletjs.com/) - JavaScript library untuk maps (GRATIS)
- [React Leaflet](https://react-leaflet.js.org/) - React wrapper untuk Leaflet
- [OpenStreetMap Tiles](https://wiki.openstreetmap.org/wiki/Tiles) - Map tiles gratis

### Navigation Deep Links
- [Google Maps URLs](https://developers.google.com/maps/documentation/urls/get-started) - Deep linking (GRATIS, tidak perlu API key)
- [Waze Deep Links](https://developers.google.com/waze/deeplinks) - Waze navigation

---

## 🎯 Kesimpulan

✅ **Backend**: OpenStreetMap Nominatim API (GRATIS)  
✅ **Frontend**: Leaflet.js + OpenStreetMap tiles (GRATIS)  
✅ **Navigation**: Deep links ke Google Maps & Waze (GRATIS)  

**Tidak perlu API key, tidak perlu billing, tidak ada biaya apapun!** 🎉

---

**Terakhir diupdate:** 2025-01-27  
**Versi:** 1.0 (100% FREE with OpenStreetMap)

