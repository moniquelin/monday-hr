CREATE TYPE payroll_status AS ENUM ('draft', 'processed');

CREATE TABLE payroll_periods (
  id           BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  start_date   DATE NOT NULL,
  end_date     DATE NOT NULL,
  status       payroll_status NOT NULL DEFAULT 'draft',

  processed_at TIMESTAMPTZ(0),
  processed_by BIGINT REFERENCES users(id),

  created_by   BIGINT,
  updated_by   BIGINT,
  created_at   TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_period_date_order CHECK (end_date >= start_date),

  CONSTRAINT payroll_periods_prevent_date_overlap EXCLUDE USING GIST (
      daterange(start_date, end_date, '[]') WITH &&
    )
);
