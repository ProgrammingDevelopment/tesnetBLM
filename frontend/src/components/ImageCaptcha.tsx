"use client";

import React from "react";

export type ImageCaptchaChallenge = {
  prompt: string;
  options: string[];
  token: string;
  expires_at: number;
};

type ImageCaptchaProps = {
  challenge: ImageCaptchaChallenge | null;
  value: string;
  onSelect: (value: string) => void;
  onReload: () => void;
  loading: boolean;
};

const IconGoldbar = () => (
  <svg viewBox="0 0 64 40" className="captcha-icon">
    <rect x="6" y="10" width="52" height="22" rx="4" fill="#d6a500" stroke="#8a6d00" />
    <rect x="12" y="14" width="40" height="14" rx="3" fill="#f0c419" />
  </svg>
);

const IconCoin = () => (
  <svg viewBox="0 0 64 64" className="captcha-icon">
    <circle cx="32" cy="32" r="22" fill="#e0b13c" stroke="#8a6d00" />
    <circle cx="32" cy="32" r="14" fill="#f3d27d" />
  </svg>
);

const IconRing = () => (
  <svg viewBox="0 0 64 64" className="captcha-icon">
    <circle cx="32" cy="32" r="20" fill="#d4af37" />
    <circle cx="32" cy="32" r="12" fill="#ffffff" />
  </svg>
);

const IconWallet = () => (
  <svg viewBox="0 0 64 44" className="captcha-icon">
    <rect x="6" y="8" width="52" height="28" rx="6" fill="#4b5563" />
    <rect x="10" y="12" width="28" height="20" rx="4" fill="#6b7280" />
    <rect x="40" y="18" width="14" height="8" rx="2" fill="#9ca3af" />
  </svg>
);

const iconMap: Record<string, React.ReactNode> = {
  goldbar: <IconGoldbar />,
  coin: <IconCoin />,
  ring: <IconRing />,
  wallet: <IconWallet />,
};

const labelMap: Record<string, string> = {
  goldbar: "Emas batangan",
  coin: "Koin",
  ring: "Cincin",
  wallet: "Dompet",
};

export default function ImageCaptcha({
  challenge,
  value,
  onSelect,
  onReload,
  loading,
}: ImageCaptchaProps) {
  return (
    <div className="captcha-card">
      <div className="d-flex justify-content-between align-items-center mb-2">
        <span className="fw-semibold">Verifikasi Gambar</span>
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
          <div className="mb-2 text-muted small">{challenge.prompt}</div>
          <div className="captcha-grid">
            {challenge.options.map((option) => (
              <button
                key={option}
                type="button"
                className={`captcha-tile ${value === option ? "active" : ""}`}
                onClick={() => onSelect(option)}
                disabled={loading}
              >
                <div className="captcha-icon-wrap">{iconMap[option]}</div>
                <div className="small">{labelMap[option] || option}</div>
              </button>
            ))}
          </div>
        </>
      ) : (
        <div className="text-muted small">Captcha belum siap.</div>
      )}
    </div>
  );
}
