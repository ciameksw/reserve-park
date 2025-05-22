# Reserve-Park

## Overview

**Reserve-Park** is a modular backend system for managing parking reservations. It is built with a microservices architecture in **Golang**, featuring dedicated services for users, spots, and reservations, all orchestrated through a facade gateway. The system uses MongoDB for data storage and Docker Compose for easy local deployment.

---

## Architecture Overview

![Reserve-Park Architecture](docs/overview.jpg)

---

## Index

- [Overview](#overview)
- [Architecture Overview](#architecture-overview)
- [Index](#index)
- [Local Deployment](#local-deployment)
- [Usage of Reserve-Park Backend](#usage-of-reserve-park-backend)
- [Internal Services Endpoints](#internal-services-endpoints)

---

## Local Deployment

See the [local deployment guide](docs/local.md) for step-by-step instructions on running Reserve-Park locally, including restoring MongoDB data from a backup.

---

## Usage of Reserve-Park Backend

The main entry point for interacting with Reserve-Park is the **facade service**.  
See the [facade service endpoint documentation](docs/facade.md) for available API endpoints and usage examples.

---

## Internal Services Endpoints

Each internal service exposes its own set of endpoints.  
Refer to the following documentation for details:

- [User Service Endpoints](docs/user.md)
- [Spot Service Endpoints](docs/spot.md)
- [Reservation Service Endpoints](docs/reservation.md)

---