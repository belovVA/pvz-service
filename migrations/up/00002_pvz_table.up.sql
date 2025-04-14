CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS pvz (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
   city TEXT NOT NULL
    );

CREATE INDEX IF NOT EXISTS idx_pvz_city ON pvz(city);
CREATE INDEX IF NOT EXISTS idx_pvz_registration_date ON pvz (registration_date);