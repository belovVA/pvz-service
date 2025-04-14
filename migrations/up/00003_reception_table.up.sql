CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS reception (
                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                         date_time TIMESTAMP NOT NULL DEFAULT NOW(),
                                         is_closed BOOLEAN NOT NULL DEFAULT FALSE,
                                         pvz_id UUID NOT NULL,
                                         CONSTRAINT fk_pvz_id FOREIGN KEY (pvz_id) REFERENCES pvz(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_reception_pvz_id ON reception(pvz_id);
CREATE INDEX IF NOT EXISTS idx_reception_date_time ON reception(date_time);
