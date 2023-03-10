[discord-shield]: https://img.shields.io/discord/1032709477054562334?label=Discord
[discord-url]: https://discord.com/invite/7FtTBwRdQK

[![Discord][discord-shield]][discord-url]

## Erlog

ErLog is a minimalist log collection service. You can either forward structured logs from existing log libraries (eg: zerolog or winston), or use any log collector to forward logs to erlog.

## Features

- Simple log collection service which runs on a $4 vps
- Works with your existing libraries

## Setup

### Docker-Compose:

Download the latest release and run `docker-compose up`

## Sending Logs

ErLog supports json formatted logs. To send logs to erlog, send a post request to `localhost:8080` with whatever data you want to be saved.

Example request body

```json
{
  "timestamp": "1675955819",
  "level": "debug",
  "service": "my_service",
  "key": "value",
  "data": {
    "another_key": "another value"
  }
}
```

## Querying Logs

No current way to query logs. API's Coming soon.

## Is it Fast?

No idea, but it uses clickhouse and an optimized way of storing logs so it should be fast enough.

## How Can I Scale it?

Probably scale clickhouse. Or email me@nevin.cc and I'll set up a hosted service for you.

## OpenTelemtry Support

Coming later with the search feature.

## Known Limitations

- Does not support search on arrays of objects (well it does, but they get transformed weird)
