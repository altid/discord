# Altid - Discord

discord is an altid service allowing a user to connect to a discord network

[![Go Report Card](https://goreportcard.com/badge/github.com/altid/discord)](https://goreportcard.com/report/github.com/altid/discord) ![Tests](https://github.com/altid/discord/workflows/Tests/badge.svg) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org)

`go get github.com/altid/discord/cmd/discord@latest`

## Usage

*Currently, we use the password field and set that to our token. This isn't ideal, we will be moving to oauth2*

`discord [-d] [-m] [-s <srv>] [-a <address>]`
- `<srv>` service name to use. (Default `discord`)

## Configuration

### With -conf

[![asciicast](https://asciinema.org/a/68P7qa0h8ZpWIUXOIalrNwOpz.svg)](https://asciinema.org/a/68P7qa0h8ZpWIUXOIalrNwOpz)

### Manually

```ndb
# altid/config - Place this in your operating systems' default configuration directory

service=discord address=discordapp.com auth=password
    password=hunter2
    user=myloginemail
    log=/usr/halfwit/log
```

- service matches the given servicename (default "discord")

- address is currently ignored
- auth is the authentication method, one of password|factotum
- factotum uses a local factotm (Plan9, plan9port) to find your password
- if auth=password, a matching password= tuple is required
- user is your login email for Discord
- log is a location to store channel logs. A special value of `none` disables logging.

> See [altid configuration](https://altid.github.io/altid-configurations.html) for more information
