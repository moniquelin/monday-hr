CREATE TABLE payroll_results (
  id                BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  payroll_period_id BIGINT NOT NULL REFERENCES payroll_periods(id),
  employee_id       BIGINT NOT NULL REFERENCES users(id),

  base_salary       BIGINT NOT NULL,
  working_days      INT NOT NULL,
  attendance_days   INT NOT NULL,
  take_home_pay     BIGINT NOT NULL,

  created_at        TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),

  CONSTRAINT chk_payroll_salary_positive CHECK (
    base_salary >= 0 AND net_salary >= 0
  ),
  CONSTRAINT chk_payroll_days CHECK (
    working_days > 0
    AND attendance_days >= 0
    AND attendance_days <= working_days
  ),

  UNIQUE (payroll_period_id, employee_id)
);