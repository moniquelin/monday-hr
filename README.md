# 💼 Monday HR

HR app to track employee attendance, overtime, reimbursement, and generate payroll. Built with **Go**.

---

## ✨ Features
### 🖥️ System
- Health check (GET `/v1/health`)
- User can log in as admin or employee  (`POST /v1/auth/login`)

### 📅 Attendance
- User (employees) can record check in (`POST /v1/attendance/checkin`)
- User (employees) can record check out (`POST /v1/attendance/checkout`)

### 💸 Payroll
-  User (admin) can create payroll periods (`POST /v1/payroll/period`)

---

## 🏗️ Tech Stack

- **Language:** Go  
- **Database:** PostgreSQL  
- **Authentication:** JWT tokens

---

## 💡 Planned

- Payroll processing  
- Payslip generation  
- Payroll summary generation for admin-side
- Docker setup for local development  
- Testing  
- Deployment setup
- Front end (React & Vite)
- Audit logs & request tracing  

---

## 📝 Development Log

- **5 Nov 2025:** Initial project setup with working health endpoint
- **8 Nov 2025:** Prepared database connection
- **17 Nov 2025:** Improved project layout for better readability
- **30 Nov 2025:** Added seeding function
- **3 Dec 2025:** Implemented user login with JWT authentication
- **9 Dec 2025:** Implemented attendance check-in & check-out feature
- **15 Dec 2025:** Users can now create payroll periods