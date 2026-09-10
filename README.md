# 🛡️ Nexus-AST

> **Autonomous AST-Driven Authorization & BOLA Linter for Cloud-Native APIs**

[![Go Report Card](https://goreportcard.com/badge/github.com/realsandeep1271-ui/nexus-ast)](https://goreportcard.com/report/github.com/realsandeep1271-ui/nexus-ast)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GitHub Action](https://img.shields.io/badge/CI%2FCD-GitHub%20Action-green)](.github/workflows/nexus-audit.yml)

Traditional static analysis tools (SAST) and dynamic scanners (DAST) identify syntax-level vulnerabilities like SQL injection or XSS, but remain **completely blind to Broken Object Level Authorization (BOLA / IDOR) and Business Logic flaws**.

**Nexus-AST** bridges this gap. It parses source code Abstract Syntax Trees (AST) and OpenAPI 3.0 specifications to detect unauthorized multi-tenant data leaks and missing object-level access controls before code reaches production.

---

## ⚡ Key Features

- **Static AST Route Inspection**: Parses Go and REST controllers using native AST parsing. Detects controllers querying database objects via unvalidated route parameters without session tenant verification.
- **Dual-Token Dynamic Verifier**: Simulates cross-tenant interactions (Tenant A vs Tenant B) against live staging APIs with zero false positives.
- **CI/CD Native (GitHub Action)**: Automatically audits Pull Requests and blocks merges when critical authorization bugs are discovered.
- **Zero Heavy Dependencies**: Written purely in high-performance Go with minimal memory footprint.

---

## 🚀 Quick Start

### Installation

```bash
git clone https://github.com/realsandeep1271-ui/nexus-ast.git
cd nexus-ast
go build -o nexus-ast main.go
```

### Static Source Code Scan

```bash
# Scan a Go backend repository for missing authorization checks
./nexus-ast scan --repo ./path-to-repo
```

### Dynamic Cross-Tenant API Verification

```bash
# Verify live endpoints against OpenAPI specification using dual tokens
./nexus-ast scan \
  --spec ./swagger.json \
  --url https://staging.api.example.com \
  --token-a "TENANT_A_TOKEN" \
  --token-b "TENANT_B_TOKEN" \
  --output report.json
```

---

## 🤖 GitHub Action Integration

Add Nexus-AST directly to your repository's CI pipeline (`.github/workflows/nexus.yml`):

```yaml
name: Security Audit
on: [pull_request]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Nexus-AST Authorization Scan
        uses: realsandeep1271-ui/nexus-ast@main
        with:
          repo-path: '.'
          fail-on-critical: 'true'
```

---

## 📋 Architecture

```
[Developer Code / PR]
       │
       ▼
 [AST Inspector] ─── (Parses Route Handlers & DB Queries)
       │
       ▼
 [Policy Matrix] ─── (Checks: Does Query Assert session.TenantID?)
       │
       ├────────► [Pass: PR Check Green ✔]
       │
       └────────► [Fail: Highlight Line in Diff + Post PR Review Comment 🚨]
```

---

## 👤 Author

**Sandeep Yadav**
- LinkedIn: [in/sandeep-yadav-1892a136b](https://linkedin.com/in/sandeep-yadav-1892a136b)
- HackerOne: [sandeep_yadav345](https://hackerone.com/sandeep_yadav345)
- Medium: [@realsandeep1271](https://medium.com/@realsandeep1271)

Licensed under the Apache License 2.0.
