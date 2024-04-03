# HTTP Multiplexer

This project is a HTTP server with a single handler. The handler accepts a POST request with a list of URLs in JSON format. The server then fetches data from all these URLs and returns the result to the client in JSON format.

### Usage

Using docker-compose:

`docker compose -f docker-compose.yml -p multiplexer up`

### What to improve

- More customization through config / consider moving config to YAML/.env.
- Increase coverage of the project.
