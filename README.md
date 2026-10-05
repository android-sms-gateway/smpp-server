# 📱 SMPP Server

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

A standalone SMPP v3.4 server that lets SMS aggregators submit SMS through the SMSGate REST API. Delivery receipts are received via per-session gateway webhook registration; forwarding them to SMPP clients with `DELIVER_SM` is not implemented yet. See [SMPP Server](https://docs.sms-gate.app/services/smpp-server/) for the full guide.

## 📖 About

The SMPP Server implements the ESME-facing side of the SMPP v3.4 protocol and bridges it to the SMSGate platform: aggregator sessions authenticate against the SMSGate REST API, submitted messages are relayed through it, and delivery receipts flow back via per-session webhooks.

## 📚 Table of Contents

- [📖 About](#-about)
- [📚 Table of Contents](#-table-of-contents)
- [⭐ Features](#-features)
- [📦 Prerequisites](#-prerequisites)
- [🚀 Getting Started](#-getting-started)
- [⚙️ Configuration](#️-configuration)
- [📦 Deployment](#-deployment)
- [🔌 API Overview](#-api-overview)
- [📚 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## ⭐ Features

- SMPP v3.4: `BIND_TX/RX/TRX`, `SUBMIT_SM`, `QUERY_SM`, `UNBIND`, `ENQUIRE_LINK`
- Authentication against the SMSGate REST API (per-session, JWT-based)
- Automatic webhook registration per receiver/transceiver session for delivery receipts
- Optional SMPPs (TLS) on the bind address when cert and key are configured
- ESME error mapping (`ESME_RINVPASWD` vs `ESME_RBINDFAIL`)
- Health and Prometheus metrics endpoints

## 📦 Prerequisites

- Go 1.25+ for building from source
- SMSGate account credentials (username/password)
- Public webhook URL for delivery receipts (receiver mode)
- Optional: TLS certificate and key for SMPPs

## 🚀 Getting Started

```bash
go build -o smpp-server .
GATEWAY__WEBHOOK_BASE_URL=https://your-host.example ./smpp-server
```

Point any SMPP v3.4 client at `127.0.0.1:2775` and bind with your SMSGate username (`system_id`) and password.

## ⚙️ Configuration

Configuration is loaded from environment variables (double-underscore nesting) or an optional YAML file via `CONFIG_PATH`.

| Variable                  | Default                                | Description                                  |
| ------------------------- | -------------------------------------- | -------------------------------------------- |
| `HTTP__ADDRESS`           | `127.0.0.1:3000`                       | HTTP API listen address (health, metrics)    |
| `SMPP__BIND_ADDRESS`      | `127.0.0.1:2775`                       | SMPP listen address                          |
| `SMPP__TLS_CERT`          | *(empty)*                              | TLS cert path; enables SMPPs on bind address |
| `SMPP__TLS_KEY`           | *(empty)*                              | TLS key path                                 |
| `SMPP__SOURCE_TON`        | `1`                                    | Source type of number                        |
| `SMPP__SOURCE_NPI`        | `1`                                    | Source numbering plan indicator              |
| `SMPP__DEST_TON`          | `1`                                    | Destination type of number                   |
| `SMPP__DEST_NPI`          | `1`                                    | Destination numbering plan indicator         |
| `GATEWAY__API_BASE_URL`   | `https://api.sms-gate.app/3rdparty/v1` | SMSGate REST API base URL                    |
| `GATEWAY__WEBHOOK_BASE_URL` | *(empty)*                            | Public webhook URL for delivery receipts     |
| `GATEWAY__TIMEOUT`        | `60s`                                  | Gateway HTTP client timeout                  |
| `CONFIG_PATH`             | *(empty)*                              | Optional YAML config file path               |

Full reference: [SMPP Server](https://docs.sms-gate.app/services/smpp-server/).

## 📦 Deployment

```bash
docker run -d --name smpp-server \
  -e GATEWAY__API_BASE_URL=https://api.sms-gate.app/3rdparty/v1 \
  -e GATEWAY__WEBHOOK_BASE_URL=https://your-host.example \
  -p 2775:2775 \
  -e SMPP__BIND_ADDRESS=0.0.0.0:2775 \
  ghcr.io/android-sms-gateway/smpp-server:latest
```

The SMPP listener defaults to `127.0.0.1` inside the container, so the example sets `SMPP__BIND_ADDRESS=0.0.0.0:2775` to make the published port reachable. The HTTP endpoints listen on `127.0.0.1:3000` and are not published; expose them through a suitably restricted reverse proxy instead of publishing port 3000 directly.

## 🔌 API Overview

| Method | Path          | Description                       |
| ------ | ------------- | --------------------------------- |
| GET    | `/health`     | Health probes (`live`, `ready`)   |
| GET    | `/metrics`    | Prometheus metrics                |
| GET    | `/api/v1/docs`| Swagger/OpenAPI documentation     |

## 📚 Documentation

- [SMPP Server](https://docs.sms-gate.app/services/smpp-server/)
- [SMPP integration](https://docs.sms-gate.app/integration/smpp/)
- [Central docs](https://docs.sms-gate.app/)

## 🤝 Contributing

Contributions are welcome. Open an issue or pull request and follow the repo's `make fmt` / `make lint` / `make test` checks.

## 📄 License

Apache-2.0. See [LICENSE](LICENSE).

<!-- Reference-style badge URLs: style=for-the-badge is mandatory -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/smpp-server?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/smpp-server/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/smpp-server?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/smpp-server/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/smpp-server?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/smpp-server/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/smpp-server?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/smpp-server/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/smpp-server?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/smpp-server/blob/master/LICENSE
