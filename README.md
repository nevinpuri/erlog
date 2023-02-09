## Erlog

Erlog is a service which lets you forward json logs from your existing libraries. It's meant to be really lightweight, easy to run, and work with your current logging setup.

## How does it work?

Erlog is just a web framework which batch inserts into an sqlite3 server. We then provide many tutorials on how to forward your logs to erlog **from withing your code**. There's no collector or external dependencies required.

## Is it Fast?

No official benchmarks yet, but I can currently get around 9000 log insertions per second. It isn't a ton but it should be enough for most small scale projects.

## Future Goals

I don't really know
I want to rewrite erlog in either rust or elixir just for the fun of it (and maybe minor performance gains), as well as have better observability, a web ui, and a hosted service. If you have problems with logging you'd like me to solve with this product, please email me at me@nevin.cc. I'd love to hear your problems.

## Otel Support

Right now I think otel is too much of a pain to use in small/medium sized projects. There's insane competition in the space, and I feel the implementation is still really bad. My main goal with this project is solving 80% of logging problems with 20% of the work, not providing complete application observability if it means spending 4 hours on otel. If otel becomes easier to use, or enough people ask for otel support, then I'll add support for it.

## Todo

- make docs on using erlog with specific libraries
- Add Discord invite link
- Add events which fire when logs aren't as usual, or logs deviate from a norm
- use Server Sent Messages (or whatever they're called, openai uses them) to forward logs to the user as they arrive
- Parse logs when they're sent in (probably bind them to a struct and save that instead of just validating to remove random chars at the end of the json)
- Rewrite in elixir (I really want to try out live view + the concurrency)
- Otel might be the play (maybe)
- Not important: make a color system where each specific tag gets a specific color (eg: level=debug gets green or whatever color it gets, message= get a different one, but this is all on the fly (doesn't have elif and all that preconfigured colors) and prefferably stays the same after reloading the app so you can get used to the colors. Probably colordb or something like that
