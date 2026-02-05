CREATE TABLE attendance_snapshots (
  id                BIGINT         GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  payroll_period_id BIGINT         NOT NULL REFERENCES payroll_periods(id),
  employee_id       BIGINT         NOT NULL REFERENCES users(id),
  att_date          DATE           NOT NULL,
  checkin_at        TIMESTAMPTZ(0) NOT NULL,
  checkout_at       TIMESTAMPTZ(0),

  created_by   BIGINT         REFERENCES users(id),
  created_at   TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
  UNIQUE (payroll_period_id, employee_id, att_date)
);
