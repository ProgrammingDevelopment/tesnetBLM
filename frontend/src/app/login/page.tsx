"use client";

import React, { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import MathCaptcha, { MathCaptchaChallenge } from "@/components/MathCaptcha";
import ImageCaptcha, { ImageCaptchaChallenge } from "@/components/ImageCaptcha";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE || "http://localhost:8080";

export default function LoginPage() {
  const router = useRouter();
  const [formData, setFormData] = useState({
    identifier: "",
    password: "",
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
        const mathData = (await mathRes.json()) as MathCaptchaChallenge;
        const imageData = (await imageRes.json()) as ImageCaptchaChallenge;
        setMathCaptcha(mathData);
        setImageCaptcha(imageData);
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
      const res = await fetch(`${API_BASE}/api/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          identifier: formData.identifier,
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
        setError(data.message || "Login gagal");
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
          <div className="col-md-5 col-lg-4">
            <div className="card shadow-lg border-0" style={{ borderRadius: "1rem" }}>
              <div className="card-body p-5">
                <div className="text-center mb-4">
                  <h2 className="fw-bold" style={{ color: "#764ba2" }}>
                    Login
                  </h2>
                  <p className="text-muted">Antrean Butik Emas LM</p>
                </div>

                {error && <div className="alert alert-danger">{error}</div>}
                {captchaError && <div className="alert alert-warning">{captchaError}</div>}

                <form onSubmit={handleSubmit}>
                  <div className="mb-3">
                    <label className="form-label fw-semibold">Email / No. WhatsApp</label>
                    <input
                      type="text"
                      name="identifier"
                      className="form-control form-control-lg"
                      placeholder="email@example.com atau 08xxx"
                      value={formData.identifier}
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
                    {loading ? "Memproses..." : "MASUK"}
                  </button>
                </form>

                <div className="text-center mt-4">
                  <span className="text-muted">Belum punya akun? </span>
                  <a
                    href="/register"
                    className="text-decoration-none fw-semibold"
                    style={{ color: "#764ba2" }}
                  >
                    Daftar di sini
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