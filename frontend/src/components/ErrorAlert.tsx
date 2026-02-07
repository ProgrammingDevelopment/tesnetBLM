"use client";

import React from 'react';

interface ErrorAlertProps {
    onRetry: () => void;
    message?: string;
}

export default function ErrorAlert({ onRetry, message }: ErrorAlertProps) {
    return (
        <div className="alert alert-danger shadow-sm animate-fade-in" role="alert">
            <h4 className="alert-heading fw-bold">
                <i className="bi bi-wifi-off me-2"></i>
                Terjadi Kendala Koneksi
            </h4>
            <p>
                {message || "Kami mengalami kesulitan memproses permintaan Anda karena kepadatan lalu lintas jaringan yang sangat tinggi."}
            </p>
            <hr />
            <div className="d-flex justify-content-end">
                <button className="btn btn-danger" onClick={onRetry}>
                    <i className="bi bi-arrow-clockwise me-2"></i>
                    Coba Lagi
                </button>
            </div>
        </div>
    );
}
