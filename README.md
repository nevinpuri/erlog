## Erlog

ErLog is a minimalist log collection service. You can either forward logs from existing log libraries (eg: zerolog or winston), or use the collector to forward logs from stdout or stderr.

## Features (how can this help you)

- Simple log collection service which runs on a $4 vps
- Works with your existing libraries

## Setup

Run the docker image for the server.

## How does it work?

ErLog is just an api which batch inserts data efficiently into an sqlite3 server.

## Is it Fast?

No official benchmarks yet, but I can currently get around 9000 log insertions per second locally. It isn't a ton but it should be enough for most small scale projects.

## How Can I Scale it?

I'm not sure of that now, but I'm developing a hosted version which auto scales based on your usage. If this sounds interesting to you, email me at me@nevin.cc and I'll set it up for you.

## Otel Support

Right now I think otel is too much of a pain to use in small/medium sized projects. For the insane amount of adoption, I feel the implementation is still really bad. If otel becomes easier to use or enough people ask for otel support, I'll add support for it.

## Todo

- make docs on using erlog with specific libraries
- Add Discord invite link
- Add events which fire when logs aren't as usual, or logs deviate from a norm
- use Server Sent Messages (or whatever they're called, openai uses them) to forward logs to the user as they arrive
- Parse logs when they're sent in (probably bind them to a struct and save that instead of just validating to remove random chars at the end of the json)
- Rewrite in elixir (I really want to try out live view + the concurrency)
- Otel might be the play (maybe)
- Not important: make a color system where each specific tag gets a specific color (eg: level=debug gets green or whatever color it gets, message= get a different one, but this is all on the fly (doesn't have elif and all that preconfigured colors) and prefferably stays the same after reloading the app so you can get used to the colors. Probably colordb or something like that
