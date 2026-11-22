## ADDED Requirements

### Requirement: Automatic subdomains
Every web service SHALL get an automatic subdomain under the platform's base domain, served over HTTPS with a wildcard certificate obtained through a DNS challenge (Cloudflare first).

#### Scenario: New service
- **WHEN** a service named `shop` is deployed with base domain `example.com`
- **THEN** it is reachable at `https://shop.example.com` with a valid certificate

### Requirement: No-domain mode
Without a base domain, the platform SHALL serve services on sslip.io hostnames derived from the server IP, with valid TLS certificates.

#### Scenario: Server without domain
- **WHEN** the server IP is 1.2.3.4 and no base domain is set
- **THEN** a service `shop` is reachable at `https://shop.1-2-3-4.sslip.io`

### Requirement: Custom domains
A service SHALL accept custom domains; the platform SHALL show the DNS record to create and SHALL request a certificate only after verifying that the domain resolves to the server.

#### Scenario: DNS not ready
- **WHEN** the admin adds `shop.com` but its DNS does not point to the server yet
- **THEN** the domain stays "pending verification" and no certificate request is made

#### Scenario: DNS ready
- **WHEN** verification succeeds
- **THEN** a certificate is issued and the domain serves the service over HTTPS

### Requirement: Redirects
All HTTP traffic SHALL be redirected to HTTPS; per service, the admin SHALL choose whether `www` redirects to the apex or the reverse, and MAY add custom redirect rules.

#### Scenario: Plain HTTP
- **WHEN** a client requests `http://shop.com`
- **THEN** it receives a permanent redirect to `https://shop.com`

### Requirement: Access restriction
The platform SHALL let the admin protect a service without code changes by HTTP basic auth and/or an allowlist of IPs or CIDR ranges.

#### Scenario: IP not allowed
- **WHEN** a service has an allowlist and a request comes from another IP
- **THEN** the request is rejected with 403

### Requirement: Routing independent of the platform process
Live traffic to services SHALL keep flowing while the platform process restarts or upgrades.

#### Scenario: Platform restart
- **WHEN** the platform container restarts
- **THEN** existing service routes keep answering requests during the restart

### Requirement: Panel exposure
The panel SHALL be served over HTTPS on its own subdomain; before a domain is configured, the setup wizard SHALL be served on a dedicated port with a self-signed certificate.

#### Scenario: Before domain setup
- **WHEN** the operator opens the server IP on the setup port right after install
- **THEN** the wizard loads over HTTPS with a self-signed certificate warning
