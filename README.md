# Backend for the Lightwell Network UI

## Local development

### Requirements

- Go 1.26

### Create your configuration

Create a config file from the example: 

```sh
cp configs/config.yaml.example configs/config.yaml
```

### Linting

```sh
make lint
```

### Unit tests

```sh
make test-unit
```

### Run the server 

```sh
make run
```

### Hit the API

```sh
curl http://localhost:8000/ping
```