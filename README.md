## Erlog

A self hosted server to send and view all your application logs.

asdf

## Benchmarks

- No official benchmarks, I can currently get around 9000 log insertions per second. It isn't a ton but it should be enough for most small scale projects.

## Log Forwarding

TODO

## Otel Support

Right now I think otel is too much of a pain to use in small/medium sized projects. If otel becomes easier to use, or enough people ask for otel support, then I'll consider adding support for it.

## Todo

- make docs on using erlog with specific libraries
- make queue flush (queue side) instead of continuing to append to the array (maybe it already does, you just have to check the flush print statements)
- Add Discord invite link
- Add events which fire when logs aren't as usual, or logs deviate from a norm
- use Server Sent Messages (or whatever they're called, openai uses them) to forward logs to the user as they arrive
