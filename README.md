# Discordfs

discordfs is an altid service allowing a user to connect to a discord network

[![Go Report Card](https://goreportcard.com/badge/github.com/altid/discordfs)](https://goreportcard.com/report/github.com/altid/discordfs)

`go get github.com/altid/discordfs`

## Usage


`discordfs [-p <dir>] [-s <srv>]`

- `<dir>` fileserver path. Will default to /tmp/altid if none is given
- `<srv>` service name to use. (Default `discord`)

## Configuration

```ndb
# altid/config - Place this in your operating systems' default configuration directory

service=discord address=discordapp.com auth=pass=hunter2
    user=myloginemail
    log=/usr/halfwit/log
    #listen_address=192.168.0.4
```

- service matches the given servicename (default "discord")

- address is currently ignored
- auth is the authentication method, one of password|factotum|none
- factotum uses a local factotm (Plan9, plan9port) to find your password
- if auth=password, a matching password= tuple is required
- user is your login email for Discord
- log is a location to store channel logs. A special value of `none` disables logging.
- listen_address is a more advanced topic, explained here: [Using listen_address](https://altid.github.io/using-listen-address.html)

> See [altid configuration](https://altid.github.io/altid-configurations.html) for more information
