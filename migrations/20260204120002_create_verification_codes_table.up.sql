-- Create verification_codes table
CREATE TABLE IF NOT EXISTS verification_codes (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255),
    phone VARCHAR(20),
    code VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_verification_email ON verification_codes(email, type);
CREATE INDEX idx_verification_phone ON verification_codes(phone, type);
CREATE INDEX idx_verification_expires ON verification_codes(expires_at);
