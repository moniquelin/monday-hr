# 💼 Monday HR

HR app to track employee attendance, overtime, reimbursement, and generate payroll. Built with **Go**.

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

## 🏗️ Tech Stack

- **Language:** Go  
- **Database:** PostgreSQL  
- **Authentication:** JWT tokens

---

## 💡 Planned

- Payslip generation  
- Payroll summary generation for admin-side
- Docker setup for local development  
- Testing  
- Deployment setup