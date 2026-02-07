"use client";

import React, { useState, useEffect, useMemo } from 'react';
import PreWarBanner from '@/components/PreWarBanner';
import LoadingModal from '@/components/LoadingModal';
import SuccessModal from '@/components/SuccessModal';
import FailureCard from '@/components/FailureCard';
import ErrorAlert from '@/components/ErrorAlert';

type AppState = 'PRE_WAR' | 'IDLE' | 'LOADING' | 'SUCCESS' | 'FAILURE' | 'ERROR';
const API_BASE = process.env.NEXT_PUBLIC_API_BASE || 'http://localhost:8080';

export default function Home() {
  const [appState, setAppState] = useState<AppState>('IDLE');
  const [ticketNumber, setTicketNumber] = useState<string>('');
  const [quota, setQuota] = useState<number>(0);
  const [backendStatus, setBackendStatus] = useState<'online' | 'offline'>('offline');

  const targetTime = useMemo(() => {
    const now = new Date();
    const nextSeven = new Date();
    nextSeven.setHours(7, 0, 0, 0);
    if (now.getHours() >= 7) {
      nextSeven.setDate(nextSeven.getDate() + 1);
    }
    return nextSeven;
  }, []);

  useEffect(() => {
    // Fetch quota from backend
    const fetchQuota = async () => {
      try {
        const res = await fetch(`${API_BASE}/api/status`);
        const data = await res.json();
        setQuota(data.quota_remaining || 0);
        setBackendStatus('online');
      } catch {
        setBackendStatus('offline');
      }
    };
    fetchQuota();
    const interval = setInterval(fetchQuota, 3000); // Refresh every 3s
    return () => clearInterval(interval);
  }, []);

  const handleWarClick = async () => {
    setAppState('LOADING');

    try {
      // Real API call to Go backend
      const response = await fetch(`${API_BASE}/api/war`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: `user_${Date.now()}`,
          name: 'Test User'
        }),
      });

      const data = await response.json();

      if (data.status === 'success') {
        setTicketNumber(data.ticket_number);
        setAppState('SUCCESS');
      } else if (data.status === 'failed') {
        setAppState('FAILURE');
      } else {
        setAppState('ERROR');
      }
    } catch (err) {
      console.error('API Error:', err);
      setAppState('ERROR');
    }
  };

  return (
    <main className="min-h-screen container py-5">
      <div className="row justify-content-center">
        <div className="col-md-8 col-lg-6">

          <div className="text-center mb-5">
            <h1 className="fw-bold text-primary">War Tiket System</h1>
            <p className="text-muted">High Performance Queuing Engine Demo</p>
          </div>

          {/* MAIN CONTENT AREA */}
          <div className="content-area">

            {appState === 'PRE_WAR' && (
              <PreWarBanner targetTime={targetTime} />
            )}

            {appState === 'IDLE' && (
              <div className="card card-custom p-5 text-center animate-fade-in">
                <div className="card-body">
                  <h3 className="mb-4 fw-bold">Antrean Dibuka!</h3>
                  <p className="mb-4">
                    Kuota tersedia: <span className={`badge ${quota > 0 ? 'bg-success' : 'bg-danger'}`}>{quota.toLocaleString()}</span>
                    <span className={`ms-2 badge ${backendStatus === 'online' ? 'bg-primary' : 'bg-secondary'}`}>
                      {backendStatus === 'online' ? 'Backend Online' : 'Backend Offline'}
                    </span>
                  </p>
                  <button
                    className="btn btn-primary btn-lg w-100 py-3 fw-bold fs-5 shadow-sm"
                    onClick={handleWarClick}
                  >
                    AMBIL ANTREAN SEKARANG
                  </button>
                </div>
              </div>
            )}

            {appState === 'FAILURE' && <FailureCard />}

            {appState === 'SUCCESS' && (
              <SuccessModal
                ticketNumber={ticketNumber}
                onViewDetails={() => alert('Redirecting to ticket details...')}
              />
            )}

            {appState === 'ERROR' && (
              <ErrorAlert onRetry={handleWarClick} />
            )}
          </div>

          <LoadingModal show={appState === 'LOADING'} />

          {/* DEBUG CONTROLS (Remove in production) */}
          <div className="mt-5 p-3 border rounded bg-light opacity-50 small">
            <p className="mb-2 fw-bold">Dev Controls:</p>
            <div className="btn-group btn-group-sm">
              <button className="btn btn-outline-secondary" onClick={() => setAppState('PRE_WAR')}>Pre-War</button>
              <button className="btn btn-outline-secondary" onClick={() => setAppState('IDLE')}>Idle</button>
              <button className="btn btn-outline-secondary" onClick={() => setAppState('LOADING')}>Loading</button>
              <button className="btn btn-outline-secondary" onClick={() => setAppState('SUCCESS')}>Success</button>
              <button className="btn btn-outline-secondary" onClick={() => setAppState('FAILURE')}>Failure</button>
              <button className="btn btn-outline-secondary" onClick={() => setAppState('ERROR')}>Error</button>
            </div>
          </div>

        </div>
      </div>
    </main>
  );
}
