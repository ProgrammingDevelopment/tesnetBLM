"use client";

import React from 'react';

export default function FailureCard() {
    return (
        <div className="card card-custom state-failure border-danger p-4 text-center animate-fade-in">
            <div className="card-body">
                <div className="mb-3 text-danger">
                    <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" fill="currentColor" className="bi bi-x-circle-fill" viewBox="0 0 16 16">
                        <path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zM5.354 4.646a.5.5 0 1 0-.708.708L7.293 8l-2.647 2.646a.5.5 0 0 0 .708.708L8 8.707l2.646 2.647a.5.5 0 0 0 .708-.708L8.707 8l2.647-2.646a.5.5 0 0 0-.708-.708L8 7.293 5.354 4.646z" />
                    </svg>
                </div>

                <h3 className="fw-bold text-danger mb-3">Mohon Maaf, Kuota Habis.</h3>

                <p className="card-text text-dark mb-4">
                    Terima kasih atas antusiasme Anda. Karena tingginya permintaan,
                    seluruh kuota antrean untuk hari ini (Total: 5.000) telah terisi penuh.
                </p>

                <div className="alert alert-light border-0 shadow-sm mb-4">
                    <h6 className="fw-bold mb-1">Saran:</h6>
                    <p className="small mb-0">Silakan mencoba kembali besok pagi pukul <strong>07:00 WIB</strong>.</p>
                </div>

                <button className="btn btn-outline-danger w-100" onClick={() => window.location.reload()}>
                    Kembali ke Beranda
                </button>
            </div>
        </div>
    );
}
