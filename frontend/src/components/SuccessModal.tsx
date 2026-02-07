"use client";

import React from 'react';

interface SuccessModalProps {
    ticketNumber?: string;
    onViewDetails?: () => void;
}

export default function SuccessModal({ ticketNumber = "A-1234", onViewDetails }: SuccessModalProps) {
    return (
        <div className="card card-custom p-5 text-center animate-fade-in state-success border-success">
            <div className="card-body">
                <div className="mb-4 text-success">
                    <svg xmlns="http://www.w3.org/2000/svg" width="80" height="80" fill="currentColor" className="bi bi-check-circle-fill" viewBox="0 0 16 16">
                        <path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zm-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417 5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z" />
                    </svg>
                </div>

                <h2 className="display-6 fw-bold text-success mb-2">Selamat! Slot Teramankan.</h2>
                <p className="lead fw-bold text-dark">Posisi kuota Anda telah kami kunci.</p>

                <div className="card bg-white bg-opacity-75 border-0 my-4 p-3 shadow-sm">
                    <p className="text-muted mb-1 small text-uppercase fw-bold">Nomor Antrean Anda</p>
                    <div className="display-4 fw-bold text-dark">{ticketNumber}</div>
                </div>

                <p className="text-muted mb-4">
                    Sistem sedang memfinalisasi data pendaftaran Anda.
                    Mohon jangan tinggalkan halaman ini sampai tiket Anda muncul.
                </p>

                <button
                    className="btn btn-success btn-lg w-100 shadow-sm"
                    onClick={onViewDetails}
                >
                    Lihat Detail Tiket
                </button>
            </div>
        </div>
    );
}
