CREATE TABLE payroll_results (
  id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  payroll_period_id BIGINT NOT NULL REFERENCES payroll_periods(id),
  employee_id       BIGINT NOT NULL REFERENCES users(id),

  base_salary       BIGINT NOT NULL,
  working_days      INT NOT NULL,
  attendance_days   INT NOT NULL,
  take_home_pay     BIGINT NOT NULL,

  created_at        TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),

  UNIQUE (payroll_period_id, employee_id)
);