import duckdb
import sys
from fastapi import FastAPI, HTTPException
from fastapi import Request
from pydantic import BaseModel
import json
import erlog_utils
from fastapi.middleware.cors import CORSMiddleware
from util import flatten, isint, isfloat
from query_builder import QueryBuilder
from models import ErLog

conn = duckdb.connect("./logs.db")
conn.execute(
    "CREATE TABLE IF NOT EXISTS erlogs (id UUID primary key, timestamp DOUBLE, string_keys string[], string_values string[], bool_keys string[], bool_values bool[], number_keys string[], number_values double[], raw_log string);"
)

app = FastAPI()

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

    q = QueryBuilder()
    try:
        p = erlog_utils.parse_input(user_query)
    except Exception:
        raise HTTPException(status_code=400, detail="Invalid query")

    if len(p.__getattribute__("and")) > 0:
        a = p.__getattribute__("and")
        keyword = "and"
    elif len(p.__getattribute__("or")) > 0:
        a = p.__getattribute__("or")
        keyword = "or"
    else:
        a = [p.none]
        keyword = ""

    q.add(a, keyword)

    query, params = q.get_query_and_params()
    print(query, params)

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
