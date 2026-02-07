"use client";

import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

interface GoldStock {
  size: string;
  available: number;
  limited: boolean;
}

interface Location {
  id: string;
  name: string;
  stock: GoldStock[];
  quota: number;
  openTime: string;
}

type StoredUser = {
  id: string;
  nik: string;
  nama: string;
  whatsapp: string;
};

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

const readStoredUser = (): StoredUser | null => {
  if (typeof window === "undefined") return null;
  const stored = localStorage.getItem("user");
  if (!stored) return null;
  try {
    const parsed = JSON.parse(stored) as Partial<StoredUser>;
    if (
      !parsed ||
      typeof parsed.id !== "string" ||
      typeof parsed.nik !== "string" ||
      typeof parsed.nama !== "string"
    ) {
      return null;
    }
    return parsed as StoredUser;
  } catch {
    return null;
  }
};

const MOCK_LOCATIONS: Location[] = [
  {
    id: "graha-dipta",
    name: "Butik Emas LM - Graha Dipta",
    quota: 30,
    openTime: "07:00 - 08:00 WIB",
    stock: [
      { size: "0,5 gr", available: 50, limited: false },
      { size: "1 gr", available: 40, limited: false },
      { size: "2 gr", available: 30, limited: false },
      { size: "3 gr", available: 25, limited: false },
      { size: "5 gr", available: 20, limited: false },
      { size: "10 gr", available: 15, limited: false },
      { size: "25gr", available: 5, limited: true },
    ],
  },
  {
    id: "juanda",
    name: "Butik Emas LM - Juanda",
    quota: 25,
    openTime: "07:00 - 08:00 WIB",
    stock: [
      { size: "0,5 gr", available: 45, limited: false },
      { size: "1 gr", available: 35, limited: false },
      { size: "2 gr", available: 25, limited: false },
      { size: "3 gr", available: 20, limited: false },
      { size: "5 gr", available: 15, limited: false },
      { size: "10 gr", available: 10, limited: true },
      { size: "25gr", available: 3, limited: true },
    ],
  },
  {
    id: "gedung-antam",
    name: "Butik Emas LM - Gedung Antam",
    quota: 40,
    openTime: "07:00 - 08:00 WIB",
    stock: [
      { size: "0,5 gr", available: 60, limited: false },
      { size: "1 gr", available: 50, limited: false },
      { size: "2 gr", available: 35, limited: false },
      { size: "3 gr", available: 30, limited: false },
      { size: "5 gr", available: 25, limited: false },
      { size: "10 gr", available: 20, limited: false },
      { size: "25gr", available: 8, limited: true },
    ],
  },
  {
    id: "setiabudi-one",
    name: "Butik Emas LM - Setiabudi One",
    quota: 20,
    openTime: "07:00 - 08:00 WIB",
    stock: [
      { size: "0,5 gr", available: 30, limited: false },
      { size: "1 gr", available: 25, limited: false },
      { size: "2 gr", available: 20, limited: false },
      { size: "3 gr", available: 15, limited: false },
      { size: "5 gr", available: 10, limited: false },
      { size: "10 gr", available: 5, limited: true },
      { size: "25gr", available: 2, limited: true },
    ],
  },
];

const TIME_SLOTS = [
  "08:30:00 - 09:00:00 WIB",
  "09:00:00 - 09:30:00 WIB",
  "09:30:00 - 10:00:00 WIB",
  "10:00:00 - 10:30:00 WIB",
  "10:30:00 - 11:00:00 WIB",
];

export default function AntreanPage() {
  const router = useRouter();
  const [user] = useState<StoredUser | null>(() => readStoredUser());
  const [selectedLocation, setSelectedLocation] = useState<Location | null>(null);
  const [selectedTime, setSelectedTime] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!user) {
      router.push("/login");
    }
  }, [router, user]);

  const handleSubmit = async () => {
    if (!user || !selectedLocation || !selectedTime) return;

    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/api/ticket`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          user_id: user.id,
          location_id: selectedLocation.id,
          time_slot: selectedTime,
        }),
      });
      const data = await res.json();
      if (data.status === "success") {
        localStorage.setItem("ticket", JSON.stringify(data.ticket));
        router.push("/tiket");
      } else {
        alert(data.message || "Gagal mengambil antrean");
      }
    } catch {
      alert("Koneksi ke server gagal");
    }
    setLoading(false);
  };

  const handleLogout = () => {
    localStorage.removeItem("user");
    router.push("/login");
  };

  const today = new Date();
  const dateStr = today.toLocaleDateString("id-ID", {
    weekday: "long",
    day: "numeric",
    month: "long",
    year: "numeric",
  });

  if (!user) return null;

  return (
    <div className="min-vh-100" style={{ background: "#f8f9fa" }}>
      {/* Header */}
      <nav
        className="navbar"
        style={{ background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)" }}
      >
        <div className="container">
          <a href="/antrean" className="text-white text-decoration-none">
            <small>Profile</small>
          </a>
          <span className="navbar-brand text-white fw-bold mx-auto">Antrean BELM</span>
          <button className="btn btn-link text-white text-decoration-none" onClick={handleLogout}>
            Logout
          </button>
        </div>
      </nav>

      <div className="container py-4">
        {/* User Info */}
        <div className="alert alert-info mb-4">
          <strong>Selamat datang, {user.nama}</strong> | NIK: ***{user.nik.slice(-4)}
        </div>

        {/* Location Selection */}
        <div className="text-center mb-4">
          <h4 className="fw-bold" style={{ color: "#764ba2" }}>
            Kuota Tersedia
          </h4>
          <select
            className="form-select form-select-lg mx-auto"
            style={{ maxWidth: "400px", borderColor: "#764ba2" }}
            onChange={(e) =>
              setSelectedLocation(MOCK_LOCATIONS.find((l) => l.id === e.target.value) || null)
            }
            value={selectedLocation?.id || ""}
          >
            <option value="">-- Pilih Lokasi Butik --</option>
            {MOCK_LOCATIONS.map((loc) => (
              <option key={loc.id} value={loc.id}>
                {loc.name}
              </option>
            ))}
          </select>
        </div>

        {selectedLocation && (
          <>
            {/* Quota Info */}
            <div className="text-center mb-3">
              <p className="mb-1">{selectedLocation.name}</p>
              <span
                className="badge rounded-pill fs-5"
                style={{ background: "#28a745", padding: "8px 20px" }}
              >
                Sisa: {selectedLocation.quota}
              </span>
              <p className="mt-2 text-muted">
                Sesi waktu ambil antrean: Pukul {selectedLocation.openTime}
              </p>
            </div>

            {/* Stock Info Box */}
            <div className="card mb-4" style={{ borderColor: "#e9d8a6", background: "#fefae0" }}>
              <div className="card-body">
                <p className="fw-bold mb-2 text-uppercase" style={{ fontSize: "0.8rem" }}>
                  INFORMASI STOK {dateStr.toUpperCase()}:
                </p>
                <p className="mb-2">
                  {selectedLocation.stock.map((s, i) => (
                    <span key={i}>
                      {s.size}
                      {s.limited ? <span className="text-danger"> (Limited)</span> : ""}
                      {i < selectedLocation.stock.length - 1 ? " | " : ""}
                    </span>
                  ))}
                </p>
                <ul className="mb-0 small" style={{ color: "#666" }}>
                  <li>1 KTP hanya bisa bertransaksi 2x dalam 1 bulan.</li>
                  <li>Pendaftaran ANTREAN tidak menjamin ketersediaan stok</li>
                  <li>Mohon datang sesuai jam kedatangan di nomor antrian.</li>
                </ul>
              </div>
            </div>

            {/* Time Slot Selection */}
            <div className="mb-4">
              <select
                className="form-select form-select-lg"
                value={selectedTime}
                onChange={(e) => setSelectedTime(e.target.value)}
              >
                <option value="">--Pilih Waktu Kedatangan--</option>
                {TIME_SLOTS.map((slot, i) => (
                  <option key={i} value={slot}>
                    {slot}
                  </option>
                ))}
              </select>
            </div>

            {/* Submit Button */}
            <button
              className="btn btn-lg w-100 text-white fw-bold py-3"
              style={{ background: "linear-gradient(135deg, #667eea 0%, #764ba2 100%)" }}
              disabled={!selectedTime || loading}
              onClick={handleSubmit}
            >
              {loading ? "Memproses..." : "AMBIL ANTREAN"}
            </button>
          </>
        )}
      </div>
    </div>
  );
}