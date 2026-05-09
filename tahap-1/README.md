# Phase 1 - Logic & Algorithm

This repository contains the solutions for the Phase 1 technical test focusing on problem-solving skills and logic optimization using the Go programming language.

## Folder Structure
```text
tahap-1/
├── shared/         # Helper functions (input reader, money formatter, stack, etc)
├── soal1/          # String matching algorithm
├── soal2/          # Cashier change calculation
├── soal3/          # Bracket structure validation
└── soal4/          # Employee leave quota calculation
```

---

## Problem 1: Matching Strings
Finding case-insensitive string matches and returning their indices.

### Execution
```bash
go run tahap-1/soal1/main.go
```

---

## Problem 2: Cashier Change
Calculating change denominations rounded down to the nearest 100.

### Execution
```bash
go run tahap-1/soal2/main.go
```

---

## Problem 3: Bracket Validation
Validating bracket strings (length 1-4096) without using Regular Expressions.

### Execution
```bash
go run tahap-1/soal3/main.go
```

---

## Problem 4: Employee Leave
Calculating annual leave quota based on join date, collective leave, and pro-rata rules.

### Execution
```bash
go run tahap-1/soal4/main.go
```
