-- UUID生成関数のために拡張機能を有効化
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Devices Table
CREATE TABLE IF NOT EXISTS devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hardware_id VARCHAR(255) NOT NULL UNIQUE, -- 物理ID
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'UNREGISTERED', -- UNREGISTERED, ACTIVE, REVOKED
    metadata JSONB DEFAULT '{}', -- 柔軟な属性情報
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Certificates Table (履歴管理)
CREATE TABLE IF NOT EXISTS certificates (
    serial_number BIGINT PRIMARY KEY, -- CAが発行したシリアル
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    fingerprint VARCHAR(255) NOT NULL,
    pem_raw TEXT NOT NULL,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    valid_to TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_certs_device_id ON certificates(device_id);

-- Enrollment Tokens (初期登録用)
CREATE TABLE IF NOT EXISTS enrollment_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id UUID NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Audit Logs (操作ログ)
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    target_device_id UUID REFERENCES devices(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL, -- "REVOKE_CERT", "UPDATE_FIRMWARE"
    details TEXT,
    actor VARCHAR(100) DEFAULT 'system',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
