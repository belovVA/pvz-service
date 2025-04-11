CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE product
(
    id           UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    date_time    TIMESTAMP    NOT NULL DEFAULT NOW(),
    type_product VARCHAR(255) NOT NULL,
    reception_id UUID         NOT NULL,
    CONSTRAINT fk_reception_id FOREIGN KEY (reception_id) REFERENCES reception(id)
);