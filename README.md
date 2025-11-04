# macvm-arc-plugin

A lightweight **GitHub Actions Runner Controller (ARC) plugin** that provisions and tears down **ephemeral macOS virtual machines** using [macvmagt](https://github.com/changty97/macvmagt).

This project bridges the gap between **Kubernetes ARC** and your **on-prem Mac hosts**, allowing ARC to dynamically request macOS VMs for use as GitHub Actions runners — similar to how ARC manages cloud VMs (e.g., EC2, GCE), but purpose-built for Apple hardware.

---

## 🚀 Overview

**macvm-arc-plugin** exposes a simple HTTP API used by ARC or custom automation logic to:

- **Create** ephemeral macOS VMs via `macvmagt`
- **Delete** VMs after workflows complete
- Integrate into ARC’s runner lifecycle (without modifying ARC)

---

### 🧩 Architecture
```
+---------------------------------------------------------------+
| GitHub Actions Runner Controller (ARC) |
| • Manages runners via Kubernetes CRDs |
| • Triggers plugin to provision macOS runners |
+---------------------------+-----------------------------------+
|
| REST/gRPC calls
▼
+---------------------------------------------------------------+
| macvm-arc-plugin (Kubernetes Pod) |
| • Exposes /create and /delete endpoints |
| • Communicates with macvmagt over HTTP |
| • Runs as a Deployment + Service in the cluster |
+---------------------------+-----------------------------------+
|
| HTTP -> macvmagt API
▼
+---------------------------------------------------------------+
| macvmagt (macOS Host Daemon) |
| • Controls macOS VMs using Virtualization Framework |
| • Launches ephemeral GitHub Actions runner VMs |
| • Reports VM status and teardown completion |
+---------------------------------------------------------------+
```
---

## 🏗️ Setup

### 1. Build the plugin

```bash
docker build -t tylerchang97/macvm-arc-plugin .
docker push tylerchang97/macvm-arc-plugin
