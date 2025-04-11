CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE pvz (
                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                    registration_date TIMESTAMP NOT NULL DEFAULT NOW(),
                    city TEXT NOT NULL
);

