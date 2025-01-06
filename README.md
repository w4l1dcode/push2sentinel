# pushsecurity-sdk-go

Go SDK for the [PushSecurity](https://pushsecurity.com) public API. PushSecurity is a browser-based Identity Threat Detection & Response (ITDR) platform.

## Running

First create a yaml file, such as `config.yml`:
```yaml
log:
  level: DEBUG

push:
  api_token: ""
  lookback_hours: 1
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
INFO[0000] set log level                                 fields.level=debug
INFO[0000] Retrieving employees                         
INFO[0000] Retrieving findings                          
INFO[0000] Retrieving apps                              
INFO[0000] Retrieving accounts                          
INFO[0000] Retrieving browsers            
```

Or binary:
```shell
% pushsecurity -config=config.yml
```

## Building

```shell
% make build
```
