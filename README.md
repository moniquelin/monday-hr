# 💼 Monday HR

HR app to track employee attendance and calculate payroll. Built with **Go**.

---

## ✨ Features
### 🖥️ System
- Health check (`GET /v1/health`)
- User can log in as admin or employee  (`POST /v1/auth/login`)

### 📅 Attendance
- User (employees) can record check in (`POST /v1/attendance/checkin`)
- User (employees) can record check out (`POST /v1/attendance/checkout`)

### 💸 Payroll
- User (admin) can create payroll periods (`POST /v1/payroll/period`)
- User (admin) can run payroll (`PUT /v1/payroll/run-payroll`)

--- 

## 📑 Seeding the Data
Set the environment variable `MONDAY_DB_ALLOW_SEED = true` and hit the following endpoints
- (`POST /v1/seed/users`) to seed 100 employee into the database 
- (`POST /v1/seed/attendance`) to seed attendance into the database
    - JSON input: `start_date` and `end_date` (format: "DDDD-MM-YY")
The attendance seed endpoint should result in
- 70% probability of employee checking in and out
- 20% probability of employee only checking in
- 10% probability of employee being absent

---

## 🏗️ Tech Stack

- **Language:** Go  
- **Database:** PostgreSQL  
- **Authentication:** JWT tokens

---

## 💡 Planned

- Payslip generation  
- Payroll summary generation for admin-side
- Testing  
- Reimbursement feature
- Overtime feature