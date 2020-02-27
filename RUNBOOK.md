# Zorrodatage

Here are a few command to help internally managing the backend application.

```console
docker build . -t rockside/timestamping-showcase
docker run -p 8080:8080 -e ROCKSIDE_API_KEY={} rockside/timestamping-showcase -identity {} -port 8080
```

To start the service under TLS and on the mainnet.

```console
docker run -p 443:443 -e ROCKSIDE_API_KEY={} rockside/timestamping-showcase -identity {} -mainnet true -port 443 -tls timestamping-showcase.rockside.io
```