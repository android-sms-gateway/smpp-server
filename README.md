<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache License][license-shield]][license-url]
[![Go][go-shield]][go-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/android-sms-gateway/smpp-server">
    <img src="https://raw.githubusercontent.com/golang-samples/gopher-vector/master/gopher.png" alt="Logo" width="120" height="120">
  </a>

<h3 align="center">SMS Gateway SMPP Server</h3>

  <p align="center">
    Standalone SMPP Server (ESME) that allows SMS aggregators to submit SMS and receive delivery receipts through the SMS Gateway REST API.
    <br />
    <a href="https://github.com/android-sms-gateway/smpp-server"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/android-sms-gateway/smpp-server/issues">Report Bug</a>
    ·
    <a href="https://github.com/android-sms-gateway/smpp-server/issues">Request Feature</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
- [About The Project](#about-the-project)
  - [Key Features](#key-features)
  - [Architecture](#architecture)
  - [Built With](#built-with)
- [Getting Started](#getting-started)
  - [Option A: Download Release Binary](#option-a-download-release-binary)
  - [Option B: Docker Image](#option-b-docker-image)
  - [Option C: Build from Source](#option-c-build-from-source)
- [Usage](#usage)
  - [Connecting](#connecting)
  - [Supported Operations](#supported-operations)
  - [Delivery Receipts](#delivery-receipts)
- [Configuration](#configuration)
- [Roadmap](#roadmap)
  - [MVP Scope](#mvp-scope)
  - [Future Enhancements](#future-enhancements)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgments](#acknowledgments)



<!-- ABOUT THE PROJECT -->
## About The Project

This is a **standalone SMPP Server (ESME)** that allows SMS aggregators to submit SMS and receive delivery receipts through the SMS Gateway via the REST API.

### Key Features

* **SMPP Protocol Support** - Full implementation of SMPP v3.4 protocol on ports 2775 (SMPP) and 2776 (SMPPs/TLS)
* **Session Management** - Full SMPP session lifecycle with thread-safe state machine
* **SMPP PDU Handling** - Bind (TX/RX/TRX), SubmitSM, QuerySM, Unbind, EnquireLink with comprehensive error codes
* **SMS Gateway Client** - Per-session authenticated HTTP client with automatic webhook registration
* **Authentication** - Username/password authentication via SMS Gateway REST API (JWT-based)
* **Dynamic Webhooks** - Automatic webhook registration per session for real-time delivery receipts
* **Error Mapping** - Comprehensive HTTP API to SMPP error code mapping

### Architecture

```
┌─────────────────┐  SMPP      ┌──────────────────────────────────────────┐
│ SMS Aggregator  │ ─────────▶ │ SMPP Server                              │
│ (External ESME) │ 2775/2776  │                                          │
└─────────────────┘            │  ┌─────────────────────────────────────┐ │
                               │  │ SMPP Service (Listener)             │ │
                               │  └──────────────────┬──────────────────┘ │
                               │                     │                    │
                               │  ┌──────────────────▼──────────────────┐ │
                               │  │ Sessions Manager                    │ │
                               │  │ (State machine, PDU routing)        │ │
                               │  └──────────────────┬──────────────────┘ │
                               │                     │                    │
                               │  ┌──────────────────▼──────────────────┐ │
                               │  │ SMS Gateway Client                  │ │
                               │  │ (HTTP API + Webhooks)               │ │
                               │  └─────────────────────────────────────┘ │
                               └──────────────┬───────────────────────────┘
                                              │ HTTP API (client-go)
                                              ▼
                                       ┌────────────────────────┐
                                       │ SMS Gateway REST API   │
                                       │ /api/3rdparty/v1/      │
                                       └────────────────────────┘
                                              │
                                      Android Devices
                                              │
                                 DELIVER_SM (via webhooks)
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>


### Built With

* [![Go][go-shield]][go-url]
* [![go-smpp][gosmpp-shield]][gosmpp-url]
* [![Uber Fx][fx-shield]][fx-url]
* [![Fiber][fiber-shield]][fiber-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

### Option A: Download Release Binary

Binaries are built with GoReleaser and published as GitHub Release artifacts for Linux, macOS, and Windows.

```bash
# Download the latest release
curl -LO https://github.com/android-sms-gateway/smpp-server/releases/latest/download/smpp-server_Linux_x86_64.tar.gz
tar -xzf smpp-server_Linux_x86_64.tar.gz
./smpp-server
```

### Option B: Docker Image

The server is available as a Docker image from GitHub Container Registry (GHCR).

```bash
docker run -p 2775:2775 ghcr.io/android-sms-gateway/smpp-server:latest
```

### Option C: Build from Source

```bash
git clone https://github.com/android-sms-gateway/smpp-server.git
cd smpp-server
go build -o smpp-server .
./smpp-server
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

The SMPP Server is designed to integrate with any SMPP v3.4 compatible client (SMS aggregators, gateways, or other software).

### Connecting

Connect your SMPP client to the server using the configured bind address (default: `127.0.0.1:2775`). Authentication uses your SMS Gateway credentials:

- **Bind Type**: `BIND_TRANSMITTER`, `BIND_RECEIVER`, or `BIND_TRANSCEIVER`
- **system_id**: Your SMS Gateway username
- **password**: Your SMS Gateway password

### Supported Operations

| Operation    | PDU              | Description                        |
| ------------ | ---------------- | ---------------------------------- |
| Bind         | `BIND_TX/RX/TRX` | Authenticate and establish session |
| Submit SMS   | `SUBMIT_SM`      | Send SMS via the gateway           |
| Query Status | `QUERY_SM`       | Check message delivery state       |
| Unbind       | `UNBIND`         | Gracefully close session           |
| Keep-Alive   | `ENQUIRE_LINK`   | Maintain connection health         |

### Delivery Receipts

When a `BIND_RECEIVER` or `BIND_TRANSCEIVER` is used, the server automatically registers a webhook with the SMS Gateway to receive delivery receipts. Delivery receipts are forwarded to the ESME client via `DELIVER_SM` PDUs (currently WIP — webhook is registered but forwarding logic is not yet implemented).

> **Note**: The `gateway.webhook_base_url` configuration must be set to a publicly reachable URL so the SMS Gateway can deliver webhook callbacks.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONFIGURATION -->
## Configuration

The server is configured via environment variables or a `.env` file using the `koanf` library.

| Config Key                 | Env Var                     | Default                                | Description                                            |
| -------------------------- | --------------------------- | -------------------------------------- | ------------------------------------------------------ |
| `http.address`             | `HTTP__ADDRESS`             | `127.0.0.1:3000`                       | HTTP API bind address (health, metrics, OpenAPI)       |
| `smpp.bind_address`        | `SMPP__BIND_ADDRESS`        | `127.0.0.1:2775`                       | SMPP server bind address                               |
| `smpp.tls_cert`            | `SMPP__TLS_CERT`            |                                        | TLS certificate file path (enables SMPPs on port 2776) |
| `smpp.tls_key`             | `SMPP__TLS_KEY`             |                                        | TLS private key file path                              |
| `gateway.api_base_url`     | `GATEWAY__API_BASE_URL`     | `https://api.sms-gate.app/3rdparty/v1` | SMS Gateway REST API base URL                          |
| `gateway.webhook_base_url` | `GATEWAY__WEBHOOK_BASE_URL` |                                        | Public webhook callback URL for delivery receipts      |
| `gateway.timeout`          | `GATEWAY__TIMEOUT`          | `60s`                                  | HTTP client timeout for gateway requests               |

> **Note**: `gateway.webhook_base_url` must be publicly reachable by the SMS Gateway for webhook-based delivery receipts to work.


<!-- ROADMAP -->
## Roadmap

### MVP Scope

- [x] SMPP Server (ESME) implementation
- [x] BIND_TRANSMITTER/RECEIVER/TRANSCEIVER commands
- [x] SUBMIT_SM message submission
- [x] QUERY_SM message status query
- [x] UNBIND session management
- [x] ENQUIRE_LINK keep-alive
- [x] Authentication via SMS Gateway REST API (JWT)
- [x] Dynamic webhook registration per session
- [x] HTTP to SMPP error code mapping
- [x] TLS/SSL support (port 2776)
- [ ] DELIVER_SM delivery receipts (webhook registered, forwarding WIP)
- [ ] Built-in metrics

### Future Enhancements

- [ ] Delivery receipt forwarding (DELIVER_SM to ESME)
- [ ] Rate limiting per client
- [ ] Message concatenation (multi-part SMS)
- [ ] Binary message support (UCS2 encoding)
- [ ] SubmitMulti support
- [ ] DataSM support
- [ ] CancelSM/ReplaceSM support
- [ ] CI release pipeline hardening

See the [open issues](https://github.com/android-sms-gateway/smpp-server/issues) for a full list of proposed features (and known issues).

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the Apache 2.0 License. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Maintainer: [@capcom6](https://github.com/capcom6)

Project Link: [https://github.com/android-sms-gateway/smpp-server](https://github.com/android-sms-gateway/smpp-server)

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [go-smpp](https://github.com/fiorix/go-smpp) - SMPP protocol library
* [Go Fiber](https://github.com/gofiber/fiber) - HTTP framework
* [Uber Fx](https://github.com/uber-go/fx) - Dependency injection
* [Shields.io](https://shields.io) - Badges

<p align="right">(<a href="#readme-top">back to top</a>)</p>




<!-- MARKDOWN LINKS & IMAGES -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/smpp-server.svg?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/smpp-server/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/smpp-server.svg?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/smpp-server/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/smpp-server.svg?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/smpp-server/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/smpp-server.svg?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/smpp-server/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/smpp-server.svg?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/smpp-server/blob/master/LICENSE
[go-shield]: https://img.shields.io/badge/go-1.25%2B-00ADD8?style=for-the-badge&logo=go
[go-url]: https://go.dev/
[gosmpp-shield]: https://img.shields.io/badge/go--smpp-v2-00ADD8?style=for-the-badge
[gosmpp-url]: https://github.com/fiorix/go-smpp
[fiber-shield]: https://img.shields.io/badge/Fiber-v2-00b894?style=for-the-badge
[fiber-url]: https://github.com/gofiber/fiber
[fx-shield]: https://img.shields.io/badge/Uber%20Fx-DI-6f42c1?style=for-the-badge
[fx-url]: https://github.com/uber-go/fx
