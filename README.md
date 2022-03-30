# Solana Ping API 

## Purpose
- execute solana ping command
- provide http API service
- generate a report and send to slack
- store results of solana ping
- provide frequently datapoint for Solan Explorer
- active monitoring confirmation loss and send an alert to slack
## Server Setup
### PingService
This is similar to  "solana ping" tool in solana tool but can do concurrent rpc query.
It send transactions to rpc endpoint and wait for transactions is confirmed. 
Use `NoPingService: false` to turn on. The default is On. 

### RetensionService
Use `RetensionService: true` to turn on. Default is Off.
Clean database data periodically.

### SlackReportService
Use `NoSlackReportService: false` to turn on. Default is On.
Fetch Report Data and a send summary to a channel periodically.

### SlackAlertService
Use `NoSlackAlertService: false` to turn on. Default is On.
If Loss is greater than a thredhold, send alert to a channel

### Example:Run only API Query Server
In config.yaml ServerSetup: 
```
 NoPingService: true           
 RetensionService: false        
 NoSlackReportService: true
 NoSlackAlertService: true    
```

## Installation
- download executable file 
- or build from source
    - Install golang 
    - clone from github.com/solana-labs/solana-ping-api
    - go mod tidy to download packages
    - go build 
- mkdir ~/.config/ping-api
- put config.yaml in ~/.config/ping-api/config.yaml

## sugguested setup
- mkdir ~/ping-api-server
- cp scripts in script to ~/ping-api-server
- make solana-ping-api system service 
    - create a /etc/systemd/system/solana-ping-api.service
    - remember to di ```sudo systemctl daemon-reload```

```
[Unit]
Description=Solana Ping API Service
After=network.target
StartLimitIntervalSec=1

[Service]
Type=simple
Restart=always
RestartSec=30
User=sol
LogRateLimitIntervalSec=0
ExecStart=/home/sol/ping-api-server/solana-ping-restart.sh

[Install]
WantedBy=multi-user.target

```

- put executable file in ~/ping-api-server
- cp config.yaml.samle to ~/ping-api-server/config.yaml and modify it 
- use cp-to-real-config.sh to copy config.yaml to ~/.config/ping-api/config.yaml
- start service by sudo sysmtemctl start solana-ping-api.service
- you can check log by ```sudo tail -f /var/log/syslog | grep ping-api```
