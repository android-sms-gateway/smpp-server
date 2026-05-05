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
- [Usage](#usage)
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
* **Authentication** - Username/password authentication via SMS Gateway REST API (JWT-based)
* **Dynamic Webhooks** - Automatic webhook registration per session for real-time delivery receipts
* **Delivery Receipts** - DELIVER_SM PDUs sent to ESME clients when message status changes
* **Error Mapping** - Comprehensive HTTP API to SMPP error code mapping
* **Metrics** - Built-in metrics for auth, server, and webhook operations

### Architecture

```
┌─────────────────┐  SMPP      ┌─────────────────┐
│ SMS Aggregator  │ ─────────▶│ SMPP Server    │
│ (External ESME) │ 2775/2776 │ (This Service)  │
└─────────────────┘            └───────┬────────┘
                                          │
                                 HTTP API (client-go)
                                          ▼
                                 ┌────────────────────────┐
                                 │ SMS Gateway REST API  │
                                 │ /api/3rdparty/v1/  │
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

WIP

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

WIP

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

### MVP Scope

- [ ] SMPP Server (ESME) implementation
- [ ] BIND_TRANSMITTER/RECEIVER/TRANSCEIVER commands
- [ ] SUBMIT_SM message submission
- [ ] QUERY_SM message status query
- [ ] DELIVER_SM delivery receipts
- [ ] UNBIND session management
- [ ] ENQUIRE_LINK keep-alive
- [ ] Authentication via SMS Gateway REST API (JWT)
- [ ] Dynamic webhook registration per session
- [ ] HTTP to SMPP error code mapping
- [ ] TLS/SSL support (port 2776)
- [ ] Built-in metrics

### Future Enhancements

- [ ] Rate limiting per client
- [ ] Message concatenation (multi-part SMS)
- [ ] Binary message support (UCS2 encoding)
- [ ] Alternative SMPP libraries evaluation
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
