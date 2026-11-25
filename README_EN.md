# DDNS-IPv6-GO

Monitors the local IPv4/IPv6 address in real time, notifies via Webhook when the IPv6 address changes, and uses certbot to automatically issue certificates for the new IPv6 address, replace the old certificate, and restart Nginx to make it effective.

## Features

- Real-time detection of the local IPv4 / IPv6 address, with a configurable check frequency (Web UI or the `-f` flag)
- Obtain IP addresses via interface, network card, or command; IPv6 supports regular expression matching
- Trigger a Webhook notification when the IPv6 address changes, with customizable URL, request body and headers
- Webhook content supports the `#{ipv4Addr}`, `#{ipv6Addr}`, `#{oldIpv6Addr}` and `#{timestamp}` variables
- certbot integration: automatically issue an IP certificate for the new IPv6 address and replace the old one
- Certificates are deployed to a fixed directory (default `/etc/ssl/ddns-ipv6`) for Nginx on port 8443, and Nginx is restarted automatically after replacement
- Daily scheduled renewal via certbot, and manual certificate renewal from the Web UI
- Web UI status panel: current IP, certificate address and expiration time, latest task result, next renewal time
- Web UI configuration for IP detection, Webhook, certbot and account, plus log viewing
- Runs as a system service; supports Windows / Linux / macOS
- Chinese and English UI and logs
