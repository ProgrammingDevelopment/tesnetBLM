"use client";

import React from 'react';
import { Modal, Spinner } from 'react-bootstrap';

interface LoadingModalProps {
    show: boolean;
}

export default function LoadingModal({ show }: LoadingModalProps) {
    return (
        <Modal show={show} backdrop="static" keyboard={false} centered>
            <Modal.Body className="text-center p-5">
                <div className="mb-4">
                    <Spinner animation="border" variant="primary" className="spinner-xl" role="status">
                        <span className="visually-hidden">Loading...</span>
                    </Spinner>
                </div>
                <h3 className="fw-bold mb-3">Sedang Memproses Permintaan...</h3>
                <p className="text-muted mb-4">
                    Sistem sedang mengecek ketersediaan kuota untuk Anda.
                    Mohon tunggu sejenak, Anda sedang berada dalam antrean virtual.
                </p>
                <div className="alert alert-warning border-warning d-flex align-items-center" role="alert">
                    <i className="bi bi-exclamation-triangle-fill fs-4 me-3"></i>
                    <div className="text-start">
                        <strong>MOHON JANGAN TUTUP ATAU REFRESH HALAMAN INI.</strong>
                        <br />
                        Melakukan refresh akan membuat Anda kehilangan posisi antrean.
                    </div>
                </div>
            </Modal.Body>
        </Modal>
    );
}
