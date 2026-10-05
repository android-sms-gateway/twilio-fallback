# 📱 Twilio Fallback Service

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

A bridge that intercepts Twilio message status callbacks and resends failed messages via the SMSGate API, providing a fallback channel without changing existing Twilio-based software. See [Twilio Fallback](https://docs.sms-gate.app/services/twilio-fallback/) for the full guide.

## 📖 About

The Twilio Fallback Service listens for Twilio outbound message status callbacks, validates the `X-Twilio-Signature` header, and for failed or blocked messages automatically resends them through SMSGate. Applications keep using Twilio as their primary channel; SMSGate becomes the delivery fallback.

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

- Twilio message status callback webhook at `/api/twilio`
- Twilio request signature validation
- Automatic resend of failed messages via the SMSGate API
- Health and Prometheus metrics endpoints

## 📦 Prerequisites

- Linux server with a public IP address (for Twilio webhooks)
- A Twilio account with message status callbacks enabled
- An SMSGate account
- Domain name and SSL/TLS certificate (recommended for webhook security)

## 🚀 Getting Started

```bash
make build
cp .env.example .env   # fill in Twilio and SMSGate credentials
./bin/twilio-fallback
```

Configure your Twilio account to send message status callbacks to `https://your-host/api/twilio`.

## ⚙️ Configuration

Configuration is loaded from environment variables (double-underscore nesting) or an optional YAML file. The YAML file is selected via `CONFIG_PATH`, which must be set in the process environment before startup — a value placed only in `.env` cannot select the YAML file. See [.env.example](.env.example).

| Variable                | Default                                | Description                                  |
| ----------------------- | -------------------------------------- | -------------------------------------------- |
| `HTTP__ADDRESS`         | `127.0.0.1:3000`                       | HTTP listen address                          |
| `HTTP__PROXY_HEADER`    | `X-Forwarded-For`                      | Trusted proxy header                         |
| `HTTP__PROXIES`         | *(empty)*                              | Comma-separated trusted proxies              |
| `TWILIO__ACCOUNT_SID`   | *(required)*                           | Twilio account SID                           |
| `TWILIO__AUTH_TOKEN`    | *(required)*                           | Twilio auth token                            |
| `TWILIO__CALLBACK_URL`  | *(empty)*                              | Public callback URL (defaults from Host)     |
| `SMSGATE__BASE_URL`     | `https://api.sms-gate.app/3rdparty/v1` | SMSGate API base URL                         |
| `SMSGATE__USERNAME`     | *(required)*                           | SMSGate username                             |
| `SMSGATE__PASSWORD`     | *(required)*                           | SMSGate password                             |
| `SMSGATE__TIMEOUT`      | `1s`                                   | SMSGate API timeout                          |
| `CONFIG_PATH`           | *(empty)*                              | Optional YAML config file path               |

Full reference: [Twilio Fallback](https://docs.sms-gate.app/services/twilio-fallback/).

## 📦 Deployment

```bash
docker compose up -d
```

Or with Docker directly:

```bash
docker run -d --name twilio-fallback \
  -p 3000:3000 \
  --env-file .env \
  -e HTTP__ADDRESS=:3000 \
  ghcr.io/android-sms-gateway/twilio-fallback:latest
```

See [compose.yaml](compose.yaml) for the full service definition.

## 🔌 API Overview

| Method | Path            | Description                             |
| ------ | --------------- | --------------------------------------- |
| POST   | `/api/twilio`   | Twilio message status callback webhook  |
| GET    | `/health`       | Liveness probe                          |
| GET    | `/health/ready` | Readiness probe                         |
| GET    | `/metrics`      | Prometheus metrics                      |

## 📚 Documentation

- [Twilio Fallback](https://docs.sms-gate.app/services/twilio-fallback/)
- [Central docs](https://docs.sms-gate.app/)

## 🤝 Contributing

Contributions are welcome. Open an issue or pull request; follow the repo's `make fmt` / `make lint` / `make test` checks and the [contribution guide](https://docs.sms-gate.app/contributing/).

## 📄 License

Apache-2.0. See [LICENSE](LICENSE).

## 📜 Legal Notice

This project is not affiliated with, endorsed by, or sponsored by Twilio. It is an independent project that utilizes the Twilio API. Twilio is a trademark of Twilio Inc.

<!-- Reference-style badge URLs: style=for-the-badge is mandatory -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/twilio-fallback?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/twilio-fallback/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/twilio-fallback?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/twilio-fallback/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/twilio-fallback?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/twilio-fallback/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/twilio-fallback?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/twilio-fallback/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/twilio-fallback?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/twilio-fallback/blob/master/LICENSE
