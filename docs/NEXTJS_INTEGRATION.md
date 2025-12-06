# Panduan Integrasi Next.js (Frontend) - SimNikah

Panduan ini menggunakan **Next.js 14+ (App Router)**.

## 1. Setup Environment
Buat file `.env.local` di root project Next.js:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 2. Daftar Halaman yang Harus Dibuat

### A. Halaman Publik
| Path | Deskripsi | API Utama |
|------|-----------|-----------|
| `/` | Landing Page & Cek Jadwal | `GET /simnikah/kalender-ketersediaan` |
| `/login` | Halaman Login | `POST /login` |
| `/register` | Halaman Daftar (Catin) | `POST /register` |

### B. Dashboard Catin (Calon Pengantin)
| Path | Deskripsi | API Utama |
|------|-----------|-----------|
| `/dashboard/catin` | Status Pendaftaran | `GET /simnikah/pendaftaran/status` |
| `/dashboard/catin/daftar` | Form Pendaftaran | `POST /simnikah/pendaftaran` |

### C. Dashboard Staff
| Path | Deskripsi | API Utama |
|------|-----------|-----------|
| `/dashboard/staff` | Overview Dashboard | `GET /simnikah/dashboard/staff` |
| `/dashboard/staff/verifikasi` | List Pendaftaran | `GET /simnikah/staff/pendaftaran` |
| `/dashboard/staff/verifikasi/[id]` | Detail & Action | `POST /verify-formulir`, `POST /approve` |

### D. Dashboard Kepala KUA
| Path | Deskripsi | API Utama |
|------|-----------|-----------|
| `/dashboard/kepala-kua` | Statistik | `GET /simnikah/dashboard/kepala-kua` |
| `/dashboard/kepala-kua/assign` | Penugasan Penghulu | `GET /available-penghulu`, `POST /assign-penghulu` |

### E. Dashboard Penghulu
| Path | Deskripsi | API Utama |
|------|-----------|-----------|
| `/dashboard/penghulu` | Jadwal Hari Ini | `GET /simnikah/penghulu/today-schedule` |
| `/dashboard/penghulu/tugas` | Riwayat Tugas | `GET /simnikah/penghulu/assigned-registrations` |

---

## 3. Contoh Code Implementasi

### A. Helper API (lib/api.ts)
Gunakan axios atau fetch wrapper agar token otomatis terkirim.

```typescript
// lib/api.ts
import axios from 'axios';
import Cookies from 'js-cookie'; // Install: npm i js-cookie

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add Token Interceptor
api.interceptors.request.use((config) => {
  const token = Cookies.get('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
```

### B. Fetch Data di Server Component (Next.js)
Untuk halaman dashboard, gunakan fetch langsung di component (Server Component).

```typescript
// app/dashboard/catin/page.tsx
import { cookies } from 'next/headers';

async function getStatusPendaftaran() {
  const token = cookies().get('token')?.value;
  
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/simnikah/pendaftaran/status`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
    cache: 'no-store', // Selalu ambil data terbaru
  });
  
  return res.json();
}

export default async function CatinDashboard() {
  const { data } = await getStatusPendaftaran();

  if (!data.has_registration) {
    return <div>Anda belum mendaftar. <a href="/dashboard/catin/daftar">Daftar Sekarang</a></div>;
  }

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold">Status Pendaftaran</h1>
      <div className="card bg-white p-4 mt-4 shadow">
        <p>Nomor: {data.registration.nomor_pendaftaran}</p>
        <p>Status: <span className="badge">{data.registration.status_pendaftaran}</span></p>
        <p>Tanggal: {data.registration.tanggal_nikah}</p>
      </div>
    </div>
  );
}
```

### C. Submit Form (Client Component)
Untuk form pendaftaran, gunakan Client Component (`use client`).

```typescript
// app/dashboard/catin/daftar/page.tsx
'use client';

import { useState } from 'react';
import api from '@/lib/api';
import { useRouter } from 'next/navigation';

export default function FormPendaftaran() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    
    const formData = new FormData(e.currentTarget);
    const payload = {
      calon_laki_laki: {
        nama_dan_bin: formData.get('suami_nama'),
        // ... field lainnya
      },
      // ... struktur body sesuai dokumentasi
    };

    try {
      await api.post('/simnikah/pendaftaran', payload);
      alert('Pendaftaran Berhasil!');
      router.push('/dashboard/catin');
    } catch (error) {
      alert('Gagal mendaftar');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* Input fields */}
      <button type="submit" disabled={loading}>
        {loading ? 'Mengirim...' : 'Daftar Nikah'}
      </button>
    </form>
  );
}
```

### D. Middleware (Proteksi Route)
Buat file `middleware.ts` di root project.

```typescript
// middleware.ts
import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function middleware(request: NextRequest) {
  const token = request.cookies.get('token');
  const { pathname } = request.nextUrl;

  // Redirect ke login jika akses dashboard tanpa token
  if (pathname.startsWith('/dashboard') && !token) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // Redirect ke dashboard jika sudah login tapi buka login page
  if (pathname === '/login' && token) {
    // Anda bisa decode token JWT di sini untuk redirect sesuai role
    return NextResponse.redirect(new URL('/dashboard', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/dashboard/:path*', '/login', '/register'],
};
```
