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

### Example upstream conf:

#### For MacOS
```conf
upstream app {
    server host.docker.internal:8001;
}

upstream api {
    server host.docker.internal:8002;
}

upstream asset {
    server host.docker.internal:8003;
}

upstream manager {
    server host.docker.internal:8004;
}

upstream rabbitmq {
    server host.docker.internal:15672;
}
```
