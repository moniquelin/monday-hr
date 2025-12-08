CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
  id             BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  role           TEXT           NOT NULL DEFAULT 'employee',

  name           VARCHAR(255)   NOT NULL,
  email          CITEXT         NOT NULL UNIQUE,
  password_hash  BYTEA          NOT NULL,
  salary         BIGINT         NOT NULL,
  CHECK (salary >= 0),
    
  created_at     TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  created_by     BIGINT,
  updated_by     BIGINT,

  CONSTRAINT fk_users_created_by FOREIGN KEY (created_by) REFERENCES users(id),
  CONSTRAINT fk_users_updated_by FOREIGN KEY (updated_by) REFERENCES users(id)
);

