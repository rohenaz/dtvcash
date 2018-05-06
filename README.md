# Memo

### Prerequisites

- MySQL
- Memcache
- Bitcoin node (ABC, Unlimited, etc)

### Setup

- Get repo
    ```sh
    go get github.com/memocash/memo/...
    ```

- Create config.yaml

    ```yaml
    MYSQL_HOST: 127.0.0.1
    MYSQL_USER: memo_user
    MYSQL_PASS: memo_password
    MYSQL_DB: memo
    
    MEMCACHE_HOST: 127.0.0.1
    MEMCACHE_PORT: 11211
    
    BITCOIN_NODE_HOST: 127.0.0.1
    BITCOIN_NODE_PORT: 8333
    ```

### Running

```sh
cd memo
go build

# Run node
./memo main-node

# Separately run web server
./memo web --insecure
```

### View

Visit `http://127.0.0.1:8261` in your browser
