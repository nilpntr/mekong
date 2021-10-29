# Mekong
A lightweight api gateway written in golang.

### Config
### Top Level Config
| Parameter     | Description                                          | Required |
|---------------|------------------------------------------------------|----------|
| `listenPorts` | A list of strings that listens for incoming requests | `true`   |
| `sentry`      | A config object to use sentry                        | `false`  |
| `routes`      | A list of routes                                     | `true`   |
| `heartbeat`   | A list of heartbeat routes                           | `false`  |
| `server`      | A config object to configure some server parameters  | `false`  |

### Sentry
| Parameter     | Description                                                                                                     | Required |
|---------------|-----------------------------------------------------------------------------------------------------------------|----------|
| `dsn`         | String: The sentry dsn url                                                                                      | `true`   |
| `environment` | String: The environment to be sent with events.                                                                 | `false`  |
| `release`     | String: The release to be sent with events.                                                                     | `false`  |
| `debug`       | Boolean: In debug mode, the debug information is printed to stdout to help you understand what sentry is doing. | `false`  |

### Routes
| Parameter                        | Description                                                            | Required |
|----------------------------------|------------------------------------------------------------------------|----------|
| `path`                           | String: the path, it supports wildcard paths like this: `/api/*/store` | true     |
| `headers`                        | List: of strings containing required headers                           | false    |
| `backendHost`                    | String: the host of the backend to be redirected to                    | true     |
| `basicAuthentication`            | Keys: basicAuthentication                                              | false    |
| `basicAuthentication[].username` | String: basic auth username                                            | true     |
| `basicAuthentication[].password` | String: basic auth password                                            | true     |
| `rules.hasBody`                  | Boolean: if the request must not have or must have a body              | false    |
| `methods`                        | List: allowed http methods                                             | true     |
| `rules`                          | Keys: the rules for a request                                          | false    |
| `rules.hasQueryString`           | Boolean: if the request must not have or must have query strings       | false    |
| `rules.hasBody`                  | Boolean: if the request must not have or must have a body              | false    |

### Heartbeat
| Parameter       | Description                                                            | Required |
|-----------------|------------------------------------------------------------------------|----------|
| `path`          | String: the path, it supports wildcard paths like this: `/heartbeat/*` | true     |
| `response_code` | Int: the response code, default is 200                                 | false    |
| `message`       | String: the body of the response, default is empty                     | false    |

### Server
| Parameter    | Description                                  | Required |
|--------------|----------------------------------------------|----------|
| `timeout`    | Int: The reverse proxy timeout in seconds    | `false`  |
| `keep_alive` | Int: The reverse proxy keep alive in seconds | `false`  |

### Example
```
listenPorts:
  - :54321
  - :8080
routes:
  - path: "/api/*/store"
    backendHost: http://127.0.0.1
    methods:
      - POST
    rules:
      hasQueryString: true
      hasBody: true
  - path: "/api/push/fcm"
    backendHost: http://127.0.0.1
    methods:
      - POST
    basicAuthentication:
      username: admin
      password: kaasje
```

## 2. Run
To run it you could use the flag `--config-file` or env variable `MEKONG_CONFIG_FILE`
