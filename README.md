# nagios-api-go
## Simple nagios api

This is my golang learning project. However the first features already works which is read a service satatus.

TODO
- Implement authentication using jwt. Currently I ran use private network only controled by ec2 security group.
- Add all api call to read data (currently I only need service status for my work but better to implement them all)
- Maybe - add the execute nagios command (like downtime, etc..)

## Quick start

- Checkout this git repo in the host running nagios
- Change to the root directory and run go build 

```
cd nagios-api-go
go build
```

- Verify `config.json` is matching with current system
- Start the program interactively `./nagios-api-go`
- Try to curl the endpoint. Replace `{nagios_host}` and `{nagios_service}` with the real one you want to query

```
curl http://localhost:8000/{nagios_host}/service/{nagios_service}
```

- If it works, you can run it as a daemon using option `-d`

```
./nagios-api-go -d
```

More to come.
