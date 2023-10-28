## Erlog

#### A Log Platform which runs on a $4 VPS

![img1](./assets/1.png)

## Features

- Ingest as many logs as you want from an http endpoint
- Query logs with a query language
- Runs on extremely low spec hardware

## Sending Logs

Just send a POST request to erlog with JSON

```
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

## Viewing Logs

![img2](./assets/2.png)

## Querying

Erql is extremely simple. Here are some examples

Querying a field
`field.bar = 'hi'`
`field.bar = 3.0`
`field.bar = false`

And Statements
`field.x = 3 and field.y = false`

Or Statements
`field.x = 3 or field.y = false`

Array index (this is getting improved)
`field.arr.1 = 3 or field.arr.2 = false`

## todo

- stop log on 'reset_erlog' message
- use https://github.com/jurismarches/luqum as query language
- get support for traces using `parentId` and `duration` or `start` `end` in ms
- default field shown will be event, and all other data will be shown in key=param
