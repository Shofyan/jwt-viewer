// Tab switching
document.querySelectorAll('.tab-button').forEach(button => {
    button.addEventListener('click', () => {
        const tabName = button.getAttribute('data-tab');
        
        // Remove active class from all tabs and buttons
        document.querySelectorAll('.tab-button').forEach(btn => btn.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
        
        // Add active class to clicked button and corresponding content
        button.classList.add('active');
        document.getElementById(`${tabName}-tab`).classList.add('active');
    });
});

// ==================== DECODE TAB ====================
const jwtInput = document.getElementById('jwt-input');
let decodeTimeout;

jwtInput.addEventListener('input', () => {
    clearTimeout(decodeTimeout);
    decodeTimeout = setTimeout(() => {
        const token = jwtInput.value.trim();
        if (token) {
            decodeJWT(token);
        } else {
            clearDecodeOutput();
        }
    }, 300); // Debounce for 300ms
});

async function decodeJWT(token) {
    try {
        const response = await fetch('/api/decode', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ token }),
        });

        const data = await response.json();

        if (data.error) {
            showDecodeError(data.error);
            return;
        }

        hideDecodeError();
        displayDecodedJWT(data);
    } catch (error) {
        showDecodeError(`Failed to decode: ${error.message}`);
    }
}

function displayDecodedJWT(data) {
    // Display header
    const headerOutput = document.getElementById('header-output');
    headerOutput.textContent = JSON.stringify(data.header, null, 2);
    
    // Display algorithm badge
    const algBadge = document.getElementById('header-alg');
    algBadge.textContent = data.header?.alg || 'N/A';

    // Display payload
    const payloadOutput = document.getElementById('payload-output');
    payloadOutput.textContent = JSON.stringify(data.payload, null, 2);

    // Display signature
    const signatureOutput = document.getElementById('signature-output');
    signatureOutput.textContent = data.signature;

    // Display claim info
    const claimInfo = data.claim_info;
    const claimInfoDiv = document.getElementById('claim-info');
    
    if (claimInfo) {
        let claimHTML = '<div style="font-weight: 600; margin-bottom: 10px;">Standard Claims:</div>';
        
        if (claimInfo.exp_string) {
            const status = claimInfo.is_expired ? 'Expired' : 'Valid';
            const statusClass = claimInfo.is_expired ? 'status-expired' : 'status-valid';
            claimHTML += `<div><strong>Expires:</strong> ${claimInfo.exp_string} <span class="status-badge ${statusClass}">${status}</span></div>`;
            
            // Update exp status badge
            const expStatus = document.getElementById('exp-status');
            expStatus.textContent = status;
            expStatus.className = `status-badge ${statusClass}`;
        }
        
        if (claimInfo.iat_string) {
            claimHTML += `<div><strong>Issued At:</strong> ${claimInfo.iat_string}</div>`;
        }
        
        if (claimInfo.nbf_string) {
            claimHTML += `<div><strong>Not Before:</strong> ${claimInfo.nbf_string}</div>`;
        }
        
        if (claimInfo.iss) {
            claimHTML += `<div><strong>Issuer:</strong> ${claimInfo.iss}</div>`;
        }
        
        if (claimInfo.sub) {
            claimHTML += `<div><strong>Subject:</strong> ${claimInfo.sub}</div>`;
        }
        
        if (claimInfo.aud) {
            claimHTML += `<div><strong>Audience:</strong> ${claimInfo.aud}</div>`;
        }
        
        claimInfoDiv.innerHTML = claimHTML;
        claimInfoDiv.style.display = 'block';
    } else {
        claimInfoDiv.style.display = 'none';
    }
}

function clearDecodeOutput() {
    document.getElementById('header-output').textContent = '';
    document.getElementById('payload-output').textContent = '';
    document.getElementById('signature-output').textContent = '';
    document.getElementById('claim-info').style.display = 'none';
    document.getElementById('header-alg').textContent = '';
    document.getElementById('exp-status').textContent = '';
    hideDecodeError();
}

function showDecodeError(message) {
    const errorDiv = document.getElementById('decode-error');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    clearDecodeOutput();
}

function hideDecodeError() {
    document.getElementById('decode-error').style.display = 'none';
}

// ==================== ENCODE TAB ====================
const encodeButton = document.getElementById('encode-button');
const copyEncodeButton = document.getElementById('copy-encode-button');

encodeButton.addEventListener('click', async () => {
    const algorithm = document.getElementById('encode-algorithm').value;
    const headerText = document.getElementById('encode-header').value;
    const payloadText = document.getElementById('encode-payload').value;
    const secret = document.getElementById('encode-secret').value;

    // Parse JSON
    let header, payload;
    try {
        header = JSON.parse(headerText);
        payload = JSON.parse(payloadText);
    } catch (error) {
        showEncodeError(`Invalid JSON: ${error.message}`);
        return;
    }

    if (!secret) {
        showEncodeError('Secret is required');
        return;
    }

    try {
        const response = await fetch('/api/encode', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                header,
                payload,
                secret,
                algorithm,
            }),
        });

        const data = await response.json();

        if (data.error) {
            showEncodeError(data.error);
            return;
        }

        hideEncodeError();
        document.getElementById('encode-result').value = data.token;
        copyEncodeButton.style.display = 'block';
    } catch (error) {
        showEncodeError(`Failed to encode: ${error.message}`);
    }
});

copyEncodeButton.addEventListener('click', () => {
    const resultTextarea = document.getElementById('encode-result');
    resultTextarea.select();
    navigator.clipboard.writeText(resultTextarea.value);
    
    const originalText = copyEncodeButton.textContent;
    copyEncodeButton.textContent = '✓ Copied!';
    setTimeout(() => {
        copyEncodeButton.textContent = originalText;
    }, 2000);
});

function showEncodeError(message) {
    const errorDiv = document.getElementById('encode-error');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    document.getElementById('encode-result').value = '';
    copyEncodeButton.style.display = 'none';
}

function hideEncodeError() {
    document.getElementById('encode-error').style.display = 'none';
}

// ==================== VERIFY TAB ====================
const verifyButton = document.getElementById('verify-button');

verifyButton.addEventListener('click', async () => {
    const token = document.getElementById('verify-token').value.trim();
    const secret = document.getElementById('verify-secret').value;

    if (!token) {
        showVerifyError('Token is required');
        return;
    }

    if (!secret) {
        showVerifyError('Secret is required');
        return;
    }

    try {
        const response = await fetch('/api/verify', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                token,
                secret,
            }),
        });

        const data = await response.json();

        if (data.error && !data.message) {
            showVerifyError(data.error);
            return;
        }

        hideVerifyError();
        displayVerificationResult(data);
    } catch (error) {
        showVerifyError(`Failed to verify: ${error.message}`);
    }
});

function displayVerificationResult(data) {
    const resultDiv = document.getElementById('verify-result');
    const claimsDiv = document.getElementById('verify-claims');
    
    if (data.valid) {
        resultDiv.className = 'verification-result valid';
        resultDiv.innerHTML = `<strong>✓ Signature Verified</strong><div>${data.message}</div>`;
    } else {
        resultDiv.className = 'verification-result invalid';
        resultDiv.innerHTML = `<strong>✗ Verification Failed</strong><div>${data.message}</div>`;
    }

    // Display claims if available
    if (data.claims) {
        const claimsOutput = document.getElementById('verify-claims-output');
        claimsOutput.textContent = JSON.stringify(data.claims, null, 2);
        
        // Display claim info
        const claimInfo = data.claim_info;
        const claimInfoDiv = document.getElementById('verify-claim-info');
        
        if (claimInfo) {
            let claimHTML = '<div style="font-weight: 600; margin-bottom: 10px;">Standard Claims:</div>';
            
            if (claimInfo.exp_string) {
                const status = claimInfo.is_expired ? 'Expired' : 'Valid';
                const statusClass = claimInfo.is_expired ? 'status-expired' : 'status-valid';
                claimHTML += `<div><strong>Expires:</strong> ${claimInfo.exp_string} <span class="status-badge ${statusClass}">${status}</span></div>`;
            }
            
            if (claimInfo.iat_string) {
                claimHTML += `<div><strong>Issued At:</strong> ${claimInfo.iat_string}</div>`;
            }
            
            if (claimInfo.nbf_string) {
                claimHTML += `<div><strong>Not Before:</strong> ${claimInfo.nbf_string}</div>`;
            }
            
            if (claimInfo.iss) {
                claimHTML += `<div><strong>Issuer:</strong> ${claimInfo.iss}</div>`;
            }
            
            if (claimInfo.sub) {
                claimHTML += `<div><strong>Subject:</strong> ${claimInfo.sub}</div>`;
            }
            
            if (claimInfo.aud) {
                claimHTML += `<div><strong>Audience:</strong> ${claimInfo.aud}</div>`;
            }
            
            claimInfoDiv.innerHTML = claimHTML;
            claimInfoDiv.style.display = 'block';
        }
        
        claimsDiv.style.display = 'block';
    } else {
        claimsDiv.style.display = 'none';
    }
}

function showVerifyError(message) {
    const errorDiv = document.getElementById('verify-error');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
    document.getElementById('verify-result').innerHTML = '';
    document.getElementById('verify-claims').style.display = 'none';
}

function hideVerifyError() {
    document.getElementById('verify-error').style.display = 'none';
}

// Auto-populate verify tab from decode tab (convenience feature)
document.querySelector('[data-tab="verify"]').addEventListener('click', () => {
    const decodedToken = jwtInput.value.trim();
    if (decodedToken) {
        document.getElementById('verify-token').value = decodedToken;
    }
});
