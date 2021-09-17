# Simple app to export definitions from a RabbitMQ vhost

This application is meant to allow users to export definitions from a single RabbitMQ vhost. RabbitMQ Management API endpoints (`/api/definitions` and `/api/definitions/vhost`) require admin privileges and therefore cannot be used by regular users. This app is meant to be used as a "proxy" or a way to gain temporary admin permissions for the purpose of exporting definitions from a vhost.

**Warning**: once deployed, this application performs no user authentication. Therefore anyone with access to this application will be able to export definitions from any vhost. This is likely harmless, because the applications uses `/api/definitions/vhost` endpoint, which does not export users nor permissions but make sure it's acceptable in your environment.

# Installation

This application is meant to be deployed to Cloud Foundry. It should be bound to a user-provided service that grants admin permissions.

1. Create a user-provided service
```
cf cups shared-instance-admin -p 'URL, username, password'
URL> https://pivotal-rabbitmq.<SYS-DOMAIN>:443
username> admin
password> admin
```
The name of the service (`shared-instance-admin`) is hardcoded in the application so if you want to changed it, you need to edit `main.go` as well.

Also, don't forget about the port specification (`:443`). The default RabbitMQ Management API port is different (15672/15671) so the application will not work with this URL, unless the port is correct.

2. Deploy the application
```
cf push
```

Since the application is written in Go, it can be pushed using the `go_buildpack` or built locally and deployed using `binary_buildpack`.

# Usage

Send a GET request to the application and specify the `vhost` parameter with the value you want:

```
curl 'https://export-vhost.apps.rosebonbon.cf-app.com/?vhost=my_vhost'
```

Of course you can use a browser or any other HTTP client instead.

The response should be a JSON that you can now import to any other RabbitMQ instance. Import can be performed from the Management UI or command line:

```
rabbitmqadmin import file.json
```
