CREATE TABLE attendance (
  id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_id      BIGINT NOT NULL REFERENCES users(id),
  att_date     DATE   NOT NULL,                     -- <— date type, not varchar
  checkin_at   TIMESTAMPTZ(0),
  checkout_at  TIMESTAMPTZ(0),
    
  created_at   TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  created_by   BIGINT,
  updated_by   BIGINT,
  UNIQUE (user_id, att_date),
  CONSTRAINT fk_att_created_by FOREIGN KEY (created_by) REFERENCES users(id),
  CONSTRAINT fk_att_updated_by FOREIGN KEY (updated_by) REFERENCES users(id),
  CONSTRAINT chk_att_weekday CHECK (EXTRACT(DOW FROM att_date) BETWEEN 2 AND 6)
);