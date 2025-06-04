# 🚀 Twilio™ Fallback Service

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=for-the-badge)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://img.shields.io/badge/Go%20Report-A%2B-brightgreen.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/android-sms-gateway/twilio-fallback)
[![GitHub Actions Status](https://img.shields.io/badge/CI-passing-brightgreen.svg?style=for-the-badge)](https://github.com/android-sms-gateway/twilio-fallback/actions)

A service that provides a fallback mechanism for Twilio SMS messages using SMSGate.

## 📖 Overview

Many companies using Twilio's SMS API encounter challenges maintaining optimal deliverability rates, especially concerning evolving industry requirements. Existing software built for Twilio often lacks the flexibility to integrate with alternative SMS gateways without significant code modifications.

This Twilio Fallback Service provides a solution by acting as a bridge between Twilio and SMSGate. It intercepts Twilio failure webhooks, extracts the necessary message details, and reroutes the messages via SMSGate, enabling seamless fallback without requiring any changes to the existing Twilio-based software.

## 📝 Description

The Twilio Fallback Service handles Twilio message failure callbacks and resends failed messages via the SMSGate API. This service ensures that critical messages are delivered even if the primary Twilio channel fails.

## 📖 Table of Contents

- [🚀 Twilio™ Fallback Service](#-twilio-fallback-service)
  - [📖 Overview](#-overview)
  - [📝 Description](#-description)
  - [📖 Table of Contents](#-table-of-contents)
  - [🚀 Getting Started](#-getting-started)
    - [⚙️ Prerequisites](#️-prerequisites)
    - [📦 Installation](#-installation)
      - [1. Pre-built Binaries](#1-pre-built-binaries)
      - [2. Docker Image](#2-docker-image)
  - [⚙️ Configuration](#️-configuration)
  - [🚀 Deployment](#-deployment)
    - [Local Deployment](#local-deployment)
    - [Docker Deployment](#docker-deployment)
  - [📄 License](#-license)
  - [⚖️ Legal Notice](#️-legal-notice)

## 🚀 Getting Started

### ⚙️ Prerequisites

- Linux VPS with a public IP address
- Docker (if using Docker deployment)
- Domain name and SSL certificate ([strongly recommended](https://www.twilio.com/docs/usage/webhooks/webhooks-security#httpstls))

### 📦 Installation

Choose one of the following installation methods:

#### 1. Pre-built Binaries

Download the appropriate binary for your system from the [GitHub Releases page](https://github.com/android-sms-gateway/twilio-fallback/releases/latest).

#### 2. Docker Image

Alternatively, you can use the Docker image hosted on `ghcr.io`:

```sh
docker pull ghcr.io/android-sms-gateway/twilio-fallback:latest
```

## ⚙️ Configuration

Create a `.env` file in the root directory of the project. Copy the contents from `.env.example` and fill in the required values.

To enable the fallback mechanism, configure your Twilio account to send message status callbacks to the service's endpoint: `[YOUR_SERVICE_URL]/api/twilio`. This endpoint will receive webhooks for failed or blocked messages and automatically resend them via SMSGate. For more information on setting up Twilio callbacks, see the [Twilio documentation](https://www.twilio.com/docs/usage/webhooks/messaging-webhooks#outbound-message-status-callback).

The following environment variables are available to configure the service:

| Variable               | Description                                  | Default                                |
| ---------------------- | -------------------------------------------- | -------------------------------------- |
| `HTTP__ADDRESS`        | HTTP server address                          | `127.0.0.1:3000`                       |
| `HTTP__PROXY_HEADER`   | HTTP proxy header                            | *empty*                                |
| `HTTP__PROXIES`        | Comma separated list of trusted proxies      | *empty*                                |
| `TWILIO__ACCOUNT_SID`  | Twilio account identifier                    | **required**                           |
| `TWILIO__AUTH_TOKEN`   | Twilio authentication token                  | **required**                           |
| `TWILIO__CALLBACK_URL` | Publicly accessible URL for Twilio callbacks | Dynamic based on `Host` header         |
| `SMSGATE__BASE_URL`    | SMSGate API endpoint                         | `https://api.sms-gate.app/3rdparty/v1` |
| `SMSGATE__USERNAME`    | SMSGate API username                         | **required**                           |
| `SMSGATE__PASSWORD`    | SMSGate API password                         | **required**                           |
| `SMSGATE__TIMEOUT`     | SMSGate API timeout                          | `1s`                                   |

## 🚀 Deployment

### Local Deployment

To run the service locally:

```sh
./twilio-fallback
```

### Docker Deployment

Run the Docker container:

```sh
docker run -p 3000:3000 --env-file .env ghcr.io/android-sms-gateway/twilio-fallback:latest
```

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## ⚖️ Legal Notice

This project is not affiliated with, endorsed by, or sponsored by Twilio. It is an independent project that utilizes the Twilio API.
