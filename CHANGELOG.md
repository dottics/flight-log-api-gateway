## Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [0.1.0] - 2022-12-10
### Added
- The function `GetAircraftTypes` to interact with the microservice to fetch all
  possible aircraft types.
- The endpoint `GET /aircraft-types` to get all the aircraft types from the
  microservice.

## [0.0.0] - 2022-12-10
### Added
- The initial api setup.
- Copied over the basics from another repo, such as:
  - `POST /login` for login.
  - `DELETE /logout` for logout.
  - `POST /forgot-password` for forgot password.
  - `POST /reset-password` for reset password.
  - `POST /contact-us` contact us.
  - Note: These are still to be revised and updated as needed.
- Added the `flightserv` MSP.

