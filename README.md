# Mekong
A lightweight api gateway written in golang.

### Config
### Top Level Config
| Parameter     | Description                                          | Required |
|---------------|------------------------------------------------------|----------|
| `listenPorts` | A list of strings that listens for incoming requests | `true`   |
| `routes`      | A list of routes                                     | `true`   |
| `heartbeat`    A list of heartbeat routes                            | `true`   |


### Routes
| Parameter                      | Description                                                            | Required |
|--------------------------------|------------------------------------------------------------------------|----------|
| `path`                         | String: the path, it supports wildcard paths like this: `/api/*/store` | true     |
| `headers`                      | List: of strings containing required headers                           | false    |
| `backendHost`                  | String: the host of the backend to be redirected to                    | true     |
| `basicAuthentication`          | Keys: basicAuthentication                                              | false    |
| `basicAuthentication.username` | Boolean: if the request must not have or must have query strings       | false    |
| `rules.hasBody`                | Boolean: if the request must not have or must have a body              | false    |
| `methods`                      | List: allowed http methods                                             | true     |
| `rules`                        | Keys: the rules for a request                                          | false    |
| `rules.hasQueryString`         | Boolean: if the request must not have or must have query strings       | false    |
| `rules.hasBody`                | Boolean: if the request must not have or must have a body              | false    |

### Heartbeat
| Parameter       | Description                                                            | Required |
|-----------------|------------------------------------------------------------------------|----------|
| `path`          | String: the path, it supports wildcard paths like this: `/heartbeat/*` | true     |
| `response_code` | Int: the response code, default is 200                                 | false    |
| `message`       | String: the body of the response, default is empty                     | false    |


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
