CREATE TABLE zones (
    id char(6) NOT NULL,
    type VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    effective_date TIMESTAMPTZ NOT NULL,
    expiration_date TIMESTAMPTZ NOT NULL,
    state CHAR(2),
    PRIMARY KEY(id, type),
    CONSTRAINT fk_state FOREIGN KEY(state) REFERENCES areas(id)
);