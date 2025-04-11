CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE reception (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           date_time TIMESTAMP NOT NULL DEFAULT NOW(),
                           is_closed BOOLEAN NOT NULL DEFAULT FALSE,
                           pvz_id UUID NOT NULL,
                           CONSTRAINT fk_pvz_id FOREIGN KEY (pvz_id) REFERENCES pvz(id)
);
