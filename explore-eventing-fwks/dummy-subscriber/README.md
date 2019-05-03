# Overview

A dummy subscriber that can be used for debugging and troubleshooting event delivery

It exposes a **REST API** `POST /v1/events` for receiving the events.

The REST API responds with the `HTTP 201` and logs the full HTTP request.