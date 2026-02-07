"use client";

import React, { useState, useEffect } from 'react';

export default function PreWarBanner({ targetTime }: { targetTime: Date }) {
  const [timeLeft, setTimeLeft] = useState<{hours: number, minutes: number, seconds: number} | null>(null);

  useEffect(() => {
    const timer = setInterval(() => {
      const now = new Date().getTime();
      const distance = targetTime.getTime() - now;

      if (distance < 0) {
        setTimeLeft(null);
        clearInterval(timer);
      } else {
        setTimeLeft({
          hours: Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)),
          minutes: Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60)),
          seconds: Math.floor((distance % (1000 * 60)) / 1000)
        });
      }
    }, 1000);

    return () => clearInterval(timer);
  }, [targetTime]);

  return (
    <div className="card card-custom p-4 mb-4 text-center animate-fade-in">
      <div className="card-body">
        <h2 className="display-6 fw-bold mb-3">Antrean Dibuka Pukul 07:00 WIB</h2>
        <p className="lead mb-4">
          Mohon pastikan Anda sudah login dan koneksi internet stabil. 
          Tombol &quot;Ambil Antrean&quot; akan aktif secara otomatis saat waktu menunjukkan pukul 07:00:00 WIB.
          <br/>
          <span className="text-danger fw-bold">Harap tidak melakukan refresh halaman secara berlebihan.</span>
        </p>
        
        {timeLeft && (
          <div className="display-4 fw-bold font-monospace text-primary my-4">
            {String(timeLeft.hours).padStart(2, '0')}:
            {String(timeLeft.minutes).padStart(2, '0')}:
            {String(timeLeft.seconds).padStart(2, '0')}
          </div>
        )}

        <button className="btn btn-secondary btn-lg w-100 disabled" disabled>
          <i className="bi bi-clock me-2"></i>
          Menunggu Waktu Pembukaan
        </button>
      </div>
    </div>
  );
}
