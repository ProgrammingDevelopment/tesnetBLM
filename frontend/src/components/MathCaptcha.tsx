"use client";

import React from "react";

export type MathCaptchaChallenge = {
  a: number;
  b: number;
  token: string;
  expires_at: number;
};

type MathCaptchaProps = {
  challenge: MathCaptchaChallenge | null;
  value: string;
  onChange: (value: string) => void;
  onReload: () => void;
  loading: boolean;
};

export default function MathCaptcha({
  challenge,
  value,
  onChange,
  onReload,
  loading,
}: MathCaptchaProps) {
  return (
    <div className="captcha-card">
      <div className="d-flex justify-content-between align-items-center mb-2">
        <span className="fw-semibold">Verifikasi Matematika</span>
        <button
          type="button"
          className="btn btn-sm btn-outline-secondary"
          onClick={onReload}
          disabled={loading}
        >
          Muat Ulang
        </button>
      </div>
      {challenge ? (
        <>
          <label className="form-label">
            {challenge.a} + {challenge.b} = ?
          </label>
          <input
            type="number"
            className="form-control form-control-lg"
            value={value}
            onChange={(e) => onChange(e.target.value)}
            required
            placeholder="Jawaban"
            disabled={loading}
          />
        </>
      ) : (
        <div className="text-muted small">Captcha belum siap.</div>
      )}
    </div>
  );
}
