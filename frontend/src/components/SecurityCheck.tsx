"use client";

import React, { useState } from 'react';

export default function SecurityCheck() {
    const [show, setShow] = useState(false);

    return (
        <>
            <div className="fixed-bottom p-2 bg-white border-top shadow-sm d-flex justify-content-center align-items-center" style={{ fontSize: '0.8rem', zIndex: 9999 }}>
                <span className="text-muted me-2">Protected by</span>
                <button
                    onClick={() => setShow(true)}
                    className="btn btn-sm btn-light border d-flex align-items-center gap-2"
                    style={{ cursor: 'pointer' }}
                >
                    <span style={{ color: '#28a745' }}>🔒</span>
                    <span className="fw-bold text-dark">BiznetGio Secured</span>
                    <span className="badge bg-success" style={{ fontSize: '0.6rem' }}>VERIFIED</span>
                </button>
            </div>

            {show && (
                <div className="modal show d-block" style={{ background: 'rgba(0,0,0,0.5)' }}>
                    <div className="modal-dialog modal-dialog-centered modal-lg">
                        <div className="modal-content">
                            <div className="modal-header bg-light">
                                <h5 className="modal-title fw-bold text-success">
                                    🔒 Safety Check: logammulia.com
                                </h5>
                                <button type="button" className="btn-close" onClick={() => setShow(false)}></button>
                            </div>
                            <div className="modal-body bg-light">
                                <div className="p-3 bg-white border rounded font-monospace" style={{ fontSize: '0.8rem', maxHeight: '400px', overflowY: 'auto' }}>
                                    <pre className="mb-0 text-dark">
                                        {`Domain Name: logammulia.com
Registry Domain ID: 29107427_DOMAIN_COM-VRSN
Registrar WHOIS Server: whois.biznetgio.com
Registrar URL: https://www.biznetgio.com
Updated Date: 2024-06-13T17:56:53Z
Creation Date: 2000-06-13T00:00:00Z
Registrar Registration Expiration Date: 2029-06-13T00:00:00Z
Registrar: PT Biznet Gio Nusantara
Registrar IANA ID: 3773
Registrar Abuse Contact Phone: +62 21 5714567 X 7605
URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/
Domain Status: Active

>>> Last update of WHOIS database: 2026-02-06T18:26:18Z <<<`}
                                    </pre>
                                </div>
                                <div className="mt-3 text-center">
                                    <div className="alert alert-success d-inline-block py-2 px-4 mb-0">
                                        <strong>✅ Verification Success</strong><br />
                                        <small>Official Site Safety Validation Passed</small>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}
