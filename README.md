# push2sentinel

A Go program that exports [PushSecurity](https://pushsecurity.com/) (A browser-based Identity Threat Detection & Response (ITDR) platform) logs to Microsoft Sentinel SIEM. 
The table used; `PushSecurity_CL`.

## Running

First create a yaml file, such as `config.yml`:
```yaml
log:
  level: DEBUG

microsoft:
  app_id: 
  secret_key: 
  tenant_id: 
  subscription_id: 
  dcr:
    endpoint: 
    rule_id: 
    stream_name: 
  resource_group: 
  workspace_name: 
  retention_days: 

push:
  api_token: ""
  lookback_hours: 1
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
INFO[0000] set log level                                 fields.level=debug
INFO[0000] waiting for log ingestion to finish          
INFO[0000] Retrieving employees                         
INFO[0000] Retrieving findings                          
INFO[0000] Retrieving apps                              
INFO[0000] Retrieving accounts                          
INFO[0000] Retrieving browsers            
DEBU[0005] uploading logs                                module=sentinel_ingest total=9
DEBU[0005] successfully uploaded pushsecurity logs       module=sentinel_ingest total_logs=9
INFO[0005] shipped logs                                  module=sentinel_logs
INFO[0005] successfully sent logs to sentinel            total=319
```

Or binary:
```shell
% push2sentinel -config=config.yml
```

## Building

```shell
% make build
```
