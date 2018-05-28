# snapshot

The purpose of this microservice is to take in a URL as input and return a string containing the fully rendered website. (After the DOM has been built by all of the HTML, CSS, Javascript, etc)

## Environment Variables

* LISTEN_PORT - service will listen on localhost:LISTEN_PORT. Default 9999
* CHROME_PATH - Path to the chrome binary. Default is /usr/bin/chromium-browser
* SLEEP_TIME - the integer number of seconds to pause to let chrome load the page. Default 2

## Usage

There's only one parameter: url.

Send a get request to localhost:LISTEN_PORT like so:

```
curl localhost:9999/?url=https://foundirl.com
```

## Returns

JSON

```
{
  "result":
  {
    "type": "string",
    "value": "<Entire contents of DOM>"
  }
}
```

## Daemonize

Run this in an ubuntu OS after installing chromium-browser by adding this systemd file:

/etc/systemd/system/snapshot.service

```
[Unit]
Description=snapshot - create a string from a specific web location
After=network.target

[Service]
Environment="LISTEN_PORT=9999"
Environment="CHROME_PATH=/usr/bin/chromium-browser"
Environment="SLEEP_TIME=2"
User=ubuntu
ExecStart=/home/ubuntu/web/snapshot/snapshot
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=snapshot

[Install]
WantedBy=multi-user.target
```
