# microsoft365 Security score
Prometheus exporter for various Microsoft Defender ATP metrics, taken from REST API

## Temporary notice

> https://docs.microsoft.com/en-gb/microsoft-365/security/mtp/microsoft-secure-score-whats-new?view=o365-worldwide#incompatibility-with-identity-secure-score-and-graph-api

## Usage

1. Clone this repository
2. `go build ./cmd/microsoft365-security-center-exporter`
3. Set required env var
```
export AZURE_TENANT_ID="..."
export AZURE_CLIENT_ID="..."
export AZURE_CLIENT_SECRET="..."
```
Check documentation to create the application : https://docs.microsoft.com/en-us/graph/auth-v2-service?context=graph%2Fapi%2F1.0&view=graph-rest-1.0
4. `./microsoft365-security-center-exporter`
5. `curl localhost:8080/metrics`


## Metrics outputs

```
# HELP microsoft365_security_active_user_count Number of active user
# TYPE microsoft365_security_active_user_count gauge
microsoft365_security_active_user_count{tenant_id="<YOUR_TENANT_ID>"} 29
# HELP microsoft365_security_licensed_user_count Number of licensed user
# TYPE microsoft365_security_licensed_user_count gauge
microsoft365_security_licensed_user_count{tenant_id="<YOUR_TENANT_ID>"} 130
# HELP microsoft365_security_secure_score_current Secure score
# TYPE microsoft365_security_secure_score_current gauge
microsoft365_security_secure_score_current{tenant_id="<YOUR_TENANT_ID>"} 251
# HELP microsoft365_security_secure_score_max Max Secure score
# TYPE microsoft365_security_secure_score_max gauge
microsoft365_security_secure_score_max{tenant_id="<YOUR_TENANT_ID>"} 223
```