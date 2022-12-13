## Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [0.5.0] - 2022-12-13
### Added
- The function `DeleteFlightLog` to interact with the flight log service to
  delete a specific flight log.
- The `DELETE /flight-log/-` endpoint to handle the deletion of a specific
  flight log.

## [0.4.0] - 2022-12-13
### Added
- The function `UpdateFlightLog` to interact with the flight log microservice
  to update a specific flight log's data.
- The `PUT /flight-log/-` endpoint to handle the update of a flight log's data.

## [0.3.0] - 2022-12-13
### Added
- The function `CreateFlightLog` to interact with the microservice to create a
  new flight log in the flight log service.
- The `POST /flight-log` endpoint to handle the creation of a new flight log.

## [0.2.0] - 2022-12-10
### Added
- The function `GetFlightLog` to get a specific flight log of a user.
- The function `GetFlightLogs` to get all of a user's flight logs.
- The `GET /flight-log/-` endpoint handle the exchange to get a specific flight
  log for a user.
- The `GET /flight-log` endpoint to handle the exchange to get all of a user's
  flight logs.

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

