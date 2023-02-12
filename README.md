[discord-shield]: https://img.shields.io/discord/1032709477054562334?label=Discord
[discord-url]: https://discord.com/invite/7FtTBwRdQK

[![Discord][discord-shield]][discord-url]

## Erlog

ErLog is a minimalist log collection service. You can either forward structured logs from existing log libraries (eg: zerolog or winston), or use the collector to forward structured logs from stdout or stderr (coming soon).

## Features

- Simple log collection service which runs on a $4 vps
- Works with your existing libraries

## Setup

### Docker:

`docker run -p 8080:8080 nevin1901/erlog:latest`

### From Binaries

Download binaries at [releases](https://github.com/Nevin1901/erlog/releases) and run ./erlog

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

Library specific examples coming soon

## Querying Logs

The ErLog DB is literally just an sqlite file, so you can open it with any sqlite3 client. You can use any of sqlite3's [JSON](https://www.sqlite.org/json1.html) functions to query the `data` column in the `er_log` table.

Example:

```bash
sqlite> select COUNT(*) from er_logs where json_extract(data, "$.level") = "error";
56
```

## Is it Fast?

No official benchmarks yet, but I can currently get around 9000 log insertions per second locally. It isn't a ton but it should be enough for most small scale projects.

## How Can I Scale it?

I'm not sure of that now, but I'm developing a hosted version which auto scales based on your usage. If this sounds interesting to you, email me at me@nevin.cc and I'll set it up for you.

## Otel Support

Right now I think otel is too much of a pain to use in small/medium sized projects. For the insane amount of adoption, I feel the implementation is still really bad. If otel becomes easier to use or enough people ask for otel support, I'll add support for it.

## Todo

- Implement log searching functions
- make docs on using erlog with specific libraries
- Add Discord invite link
- Add events which fire when logs aren't as usual, or logs deviate from a norm
- Rewrite in rust elixir (I really want to try out live view + the concurrency)
- Otel support (maybe)
- Not important: make a color system where each specific tag gets a specific color (eg: level=debug gets green or whatever color it gets, message= get a different one, but this is all on the fly (doesn't have elif and all that preconfigured colors) and prefferably stays the same after reloading the app so you can get used to the colors. Probably colordb or something like that
