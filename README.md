# 💼 Monday HR

HR app to track employee attendance, overtime, reimbursement, and generate payroll. Built with **Go**.

---

## ✨ Features

- Health check (GET `/v1/health`)
- User (admin & employees) can log in with JWT authentication (POST `/v1/auth/login`)
- User (employees) can record attendance check in (POST `/v1/attendance/checkin`)

---

## 🏗️ Tech Stack

- **Language:** Go  
- **Database:** PostgreSQL  

---

## 💡 Planned

- Attendance endpoints  
- Payroll period creation & processing  
- Payslip and payroll summary generation  
- Audit logs & request tracing  
- Docker setup for local development  
- Automated testing  
- Deployment setup
- Overtime endpoints
- Reimbursement endpoints

---

## 📝 Development Log

- **5 Nov 2025:** Initial project setup with working health endpoint
- **8 Nov 2025:** Prepared database connection
- **17 Nov 2025:** Improved project layout for better readability
- **30 Nov 2025:** Added seeding function
- **3 Dec 2025:** Implemented user login with JWT authentication
- **9 Dec 2025:** Implemented attendance check in