## Erlog

**Stop Checking your logs to see when stuff go wrong. ErLog automatically does that and lets you trace through your services**

Erlog basically monitors your logs and sends you events when things are whack.

## Benchmarks

- No official benchmarks, I can currently get around 9000 log insertions per second. It isn't a ton but it should be enough for most small scale projects.

## Log Forwarding

TODO

## Otel Support

Right now I think otel is too much of a pain to use in small/medium sized projects. If otel becomes easier to use, or enough people ask for otel support, then I'll consider adding support for it.

## Todo

- make docs on using erlog with specific libraries
- Add Discord invite link
- Add events which fire when logs aren't as usual, or logs deviate from a norm
- use Server Sent Messages (or whatever they're called, openai uses them) to forward logs to the user as they arrive
- Parse logs when they're sent in (probably bind them to a struct and save that instead of just validating to remove random chars at the end of the json)
- Rewrite in elixir (I really want to try out live view + the concurrency)
- Otel might be the play (maybe)
- Not important: make a color system where each specific tag gets a specific color (eg: level=debug gets green or whatever color it gets, message= get a different one, but this is all on the fly (doesn't have elif and all that preconfigured colors) and prefferably stays the same after reloading the app so you can get used to the colors. Probably colordb or something like that
