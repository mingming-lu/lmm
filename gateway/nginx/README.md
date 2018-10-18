### Example command to generate ssl cert:

```
openssl req \
-newkey rsa:2048 \
-x509 \
-nodes \
-keyout key.key \
-new \
-out crt.crt \
-subj /CN=hostname \
-reqexts SAN \
-extensions SAN \
-config <(cat /etc/ssl/openssl.cnf \
    <(printf '[SAN]\nsubjectAltName=DNS:hostname,DNS:other_hostname,IP:127.0.0.1')) \
-sha256 \
-days 3650
```
