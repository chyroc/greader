# greader

RSS service, providing api similar to google reader.

## Deploy

### By Docker

- **Run**

```bash
docker run -d \
  -e MYSQL_HOST=xx \
  -e MYSQL_USERNAME=xx \
  -e MYSQL_PASSWORD=xx \
  -e MYSQL_DATABASE=xx \
  -p 8081:8081 \
  chyroc/greader:latest
```

### By Docker Compose

- **Create `docker-compose.yml`**

```bash
curl -fsSL https://raw.githubusercontent.com/chyroc/greader/main/docker-compose.yml > docker-compose.yml
``
```

- **Run**

```bash
docker-compose up -d
```

### By Binary

- **Install Binary**

```bash
go install github.com/chyroc/greader@latest
```

- **Run**

```bash
MYSQL_HOST=xx \
  MYSQL_USERNAME=xx \
  MYSQL_PASSWORD=xx \
  MYSQL_DATABASE=xx \
  greader
```

## Usage

### RSS API

```text
http://<address>:8082/api/greader
```

![](./screenshot/list.png)

## Ref

- https://ranchero.com/downloads/GoogleReaderAPI-2009.pdf
- https://github.com/Ranchero-Software/NetNewsWire/blob/mac-6.1.1b1/Account/Sources/Account/ReaderAPI/ReaderAPICaller.swift
- https://github.com/FreshRSS/FreshRSS/blob/1.20.2/p/api/greader.php
