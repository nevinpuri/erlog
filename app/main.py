import duckdb
from fastapi import FastAPI, HTTPException
from fastapi import Request
import json
from fastapi.middleware.cors import CORSMiddleware
from util import flatten
from models import ErLog
from app.query import QBuilder
import os
from async_tail import atail
import asyncio
import structlog
from structlog import get_logger
import ujson

structlog.configure(processors=[structlog.processors.JSONRenderer()])


def insert_log(log):
    try:
        if log == "":
            return False

        l = ujson.loads(log)
        flattened = flatten(l)
        erlog = ErLog(log)
        erlog.parse_log(flattened)

        conn.execute(
            "INSERT INTO erlogs VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?, ?, ?)",
            [
                erlog._timestamp,
                erlog._string_keys,
                erlog._string_values,
                erlog._bool_keys,
                erlog._bool_values,
                erlog._number_keys,
                erlog._number_values,
                erlog._raw_log,
            ],
        )
    except Exception as e:
        return


async def read_from_file():
    f = os.environ["LOGS"]
    files = f.split(" ")
    async for line in atail("file1.txt"):
        # todo, get file name with it
        insert_log(str(line[0]))


conn = duckdb.connect("./logs.db")
conn.execute(
    "CREATE TABLE IF NOT EXISTS erlogs (id UUID primary key, timestamp DOUBLE, string_keys string[], string_values string[], bool_keys string[], bool_values bool[], number_keys string[], number_values double[], raw_log string);"
)

app = FastAPI()


@app.on_event("startup")
async def read_logs():
    if not "LOGS" in os.environ:
        return

    loop = asyncio.get_event_loop()
    loop.create_task(read_from_file())


origins = ["http://localhost", "http://localhost:59971", "*"]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.post("/search")
async def root(request: Request):
    body = await request.json()
    if isinstance(body, str):
        raise HTTPException(status_code=400, detail="Invalid json")

    user_query = body["query"]

    q = QBuilder()
    q.parse(user_query)
    query, params = q.query, q.params

    l = conn.execute(query, params).fetchall()

    logs = []
    for log in l:
        logs.append({"id": log[0], "timestamp": log[1], "log": log[2]})

    return logs


@app.post("/get")
async def get_log(request: Request):
    body = await request.json()

    if isinstance(body, str):
        raise HTTPException(status_code=400, detail="Invalid json")

    id = body["id"]

    if id is None:
        raise HTTPException(status_code=400, detail="Invalid json")

    h = conn.execute("SELECT id, timestamp, raw_log from erlogs WHERE id = ?", [id])
    log = h.fetchone()

    return {"id": log[0], "timestamp": log[1], "log": log[2]}


@app.post("/")
async def log(request: Request):
    body = await request.json()
    erlog = ErLog(json.dumps(body))

    if isinstance(body, str):
        raise HTTPException(status_code=400, detail="Invalid json")

    flattened = flatten(body)
    erlog.parse_log(flattened)

    # todo, use appender or add tis to a batch
    conn.execute(
        "INSERT INTO erlogs VALUES (gen_random_uuid(), ?, ?, ?, ?, ?, ?, ?, ?)",
        [
            erlog._timestamp,
            erlog._string_keys,
            erlog._string_values,
            erlog._bool_keys,
            erlog._bool_values,
            erlog._number_keys,
            erlog._number_values,
            erlog._raw_log,
        ],
    )
    return {"status": "OK"}
