CREATE TABLE IF NOT EXISTS people (
    id          INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        VARCHAR(128) NOT NULL,
    patronymic  VARCHAR(128),
    surname     VARCHAR(128) NOT NULL,
    age         INTEGER,
    gender      VARCHAR(8),
    nationality VARCHAR(128),
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_gender_type CHECK (gender IN ('male', 'female'))
    );
