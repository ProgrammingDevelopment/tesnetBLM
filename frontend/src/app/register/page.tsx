"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import MathCaptcha, { MathCaptchaChallenge } from "@/components/MathCaptcha";
import ImageCaptcha, { ImageCaptchaChallenge } from "@/components/ImageCaptcha";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export default function RegisterPage() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    nik: "",
    nama: "",
    whatsapp: "",
    email: "",
    password: "",
    confirmPassword: "",
    captchaMathAnswer: "",
    captchaImageAnswer: "",
  });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const [captchaLoading, setCaptchaLoading] = useState(false);
  const [captchaError, setCaptchaError] = useState("");
  const [mathCaptcha, setMathCaptcha] = useState<MathCaptchaChallenge | null>(null);
  const [imageCaptcha, setImageCaptcha] = useState<ImageCaptchaChallenge | null>(null);

  useEffect(() => {
    const loadCaptchas = async () => {
      setCaptchaLoading(true);
      setCaptchaError("");
      try {
        const [mathRes, imageRes] = await Promise.all([
          fetch(`${API_BASE}/api/captcha/math`),
          fetch(`${API_BASE}/api/captcha/image`),
        ]);
        if (!mathRes.ok || !imageRes.ok) {
          throw new Error("Captcha fetch failed");
        }
        setMathCaptcha((await mathRes.json()) as MathCaptchaChallenge);
        setImageCaptcha((await imageRes.json()) as ImageCaptchaChallenge);
      } catch {
        setCaptchaError("Gagal memuat captcha. Coba muat ulang.");
      } finally {
        setCaptchaLoading(false);
      }
    };

    loadCaptchas();
  }, []);

  const reloadCaptchas = async () => {
    setFormData((prev) => ({ ...prev, captchaMathAnswer: "", captchaImageAnswer: "" }));
    setError("");
    setCaptchaError("");
    try {
      setCaptchaLoading(true);
      const [mathRes, imageRes] = await Promise.all([
        fetch(`${API_BASE}/api/captcha/math`),
        fetch(`${API_BASE}/api/captcha/image`),
      ]);
      if (!mathRes.ok || !imageRes.ok) {
        throw new Error("Captcha fetch failed");
      }
      setMathCaptcha((await mathRes.json()) as MathCaptchaChallenge);
      setImageCaptcha((await imageRes.json()) as ImageCaptchaChallenge);
    } catch {
      setCaptchaError("Gagal memuat captcha. Coba muat ulang.");
    } finally {
      setCaptchaLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.nik.length !== 16) {
      setError("NIK harus 16 digit");
      return;
    }
    if (formData.password !== formData.confirmPassword) {
      setError("Password tidak cocok");
      return;
    }
    if (!mathCaptcha || !imageCaptcha) {
      setError("Captcha belum siap.");
      return;
    }
    if (!formData.captchaMathAnswer || !formData.captchaImageAnswer) {
      setError("Captcha wajib diisi.");
      return;
    }

    setLoading(true);
    try {
      const res = await fetch(`${API_BASE}/api/register`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          nik: formData.nik,
          nama: formData.nama,
          whatsapp: formData.whatsapp,
          email: formData.email,
          password: formData.password,
          captcha_math_token: mathCaptcha.token,
          captcha_math_answer: formData.captchaMathAnswer,
          captcha_image_token: imageCaptcha.token,
          captcha_image_answer: formData.captchaImageAnswer,
        }),
      });
      const data = await res.json();
      if (data.status === "success") {
        localStorage.setItem("user", JSON.stringify(data.user));
        router.push("/antrean");
      } else {
        setError(data.message || "Registrasi gagal");
        await reloadCaptchas();
      }
    } catch {
      setError("Koneksi ke server gagal");
    }
    setLoading(false);
  };

  return (
    <div className="min-vh-100 d-flex align-items-center" style={{ background: "#667eea" }}>
      <div className="container">
        <div className="row justify-content-center">
          <div className="col-md-6 col-lg-5">
            <div className="card shadow-lg border-0" style={{ borderRadius: "1rem" }}>
              <div className="card-body p-5">
                <div className="text-center mb-4">
                  <h2 className="fw-bold" style={{ color: "#764ba2" }}>
                    Daftar Antrean
                  </h2>
                  <p className="text-muted">Butik Emas Logam Mulia</p>
                </div>

                {error && <div className="alert alert-danger">{error}</div>}
                {captchaError && <div className="alert alert-warning">{captchaError}</div>}

                <form onSubmit={handleSubmit}>
                  <div className="mb-3">
                    <label className="form-label fw-semibold">NIK (16 digit)</label>
                    <input
                      type="text"
                      name="nik"
                      className="form-control form-control-lg"
                      placeholder="3201234567890001"
                      maxLength={16}
                      value={formData.nik}
                      onChange={handleChange}
                      required
                    />
                  </div>

                  <div className="mb-3">
                    <label className="form-label fw-semibold">Nama Lengkap (Sesuai KTP)</label>
                    <input
                      type="text"
                      name="nama"
                      className="form-control form-control-lg"
                      placeholder="NAMA LENGKAP"
                      value={formData.nama}
                      onChange={handleChange}
                      required
                      style={{ textTransform: "uppercase" }}
                    />
                  </div>

                  <div className="mb-3">
                    <label className="form-label fw-semibold">No. WhatsApp</label>
                    <input
                      type="tel"
                      name="whatsapp"
                      className="form-control form-control-lg"
                      placeholder="08123456789"
                      value={formData.whatsapp}
                      onChange={handleChange}
                      required
                    />
                  </div>

                  <div className="mb-3">
                    <label className="form-label fw-semibold">Email</label>
                    <input
                      type="email"
                      name="email"
                      className="form-control form-control-lg"
                      placeholder="email@example.com"
                      value={formData.email}
                      onChange={handleChange}
                      required
                    />
                  </div>

                  <div className="mb-3">
                    <label className="form-label fw-semibold">Password</label>
                    <input
                      type="password"
                      name="password"
                      className="form-control form-control-lg"
                      value={formData.password}
                      onChange={handleChange}
                      required
                      minLength={6}
                    />
                  </div>

                  <div className="mb-3">
                    <label className="form-label fw-semibold">Konfirmasi Password</label>
                    <input
                      type="password"
                      name="confirmPassword"
                      className="form-control form-control-lg"
                      value={formData.confirmPassword}
                      onChange={handleChange}
                      required
                    />
                  </div>

                  <MathCaptcha
                    challenge={mathCaptcha}
                    value={formData.captchaMathAnswer}
                    onChange={(value) => setFormData((prev) => ({ ...prev, captchaMathAnswer: value }))}
                    onReload={reloadCaptchas}
                    loading={captchaLoading}
                  />

                  <ImageCaptcha
                    challenge={imageCaptcha}
                    value={formData.captchaImageAnswer}
                    onSelect={(value) => setFormData((prev) => ({ ...prev, captchaImageAnswer: value }))}
                    onReload={reloadCaptchas}
                    loading={captchaLoading}
                  />

                  <button
                    type="submit"
                    className="btn btn-lg w-100 text-white fw-bold"
                    style={{ background: "#764ba2" }}
                    disabled={loading || captchaLoading}
                  >
                    {loading ? "Memproses..." : "DAFTAR SEKARANG"}
                  </button>
                </form>

                <div className="text-center mt-4">
                  <span className="text-muted">Sudah punya akun? </span>
                  <a
                    href="/login"
                    className="text-decoration-none fw-semibold"
                    style={{ color: "#764ba2" }}
                  >
                    Login di sini
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}