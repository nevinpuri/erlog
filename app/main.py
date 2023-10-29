import duckdb
from fastapi import FastAPI, HTTPException
from fastapi import Request
import json
from fastapi.middleware.cors import CORSMiddleware
from util import flatten
from models import ErLog
from query import QBuilder
import os
from async_tail import atail
import asyncio
import structlog
from structlog import get_logger
import ujson
import uuid

structlog.configure(processors=[structlog.processors.JSONRenderer()])


def insert_log(log):
    # logger = get_logger()
    # id = str(uuid.uuid4())
    # logger.info("Inserting log", id=id)
    if log == "":
        # logger.error("Log is none, skipping", parent_id=id)
        return False

    try:
        l = ujson.loads(log)
    except Exception as e:
        # logger.error("Failed parsing log json", parent_id=id, e=str(e))
        return False

    try:
        flattened = flatten(l)
    except Exception as e:
        # logger.error("Failed flattening json", parent_id=id, e=str(e))
        return False

    try:
        erlog = ErLog(log)
        erlog.parse_log(flattened)

        if erlog._id == None:
            erlog._id = str(uuid.uuid4())
    except Exception as e:
        # logger.error("Failed parsing log", parent_id=id, e=str(e))
        return False

    try:
        # logger.info("Inserting into table", parent_id=id)
        conn.execute(
            "INSERT INTO erlogs VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
            [
                erlog._id,
                erlog._parent_id,
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
        # logger.error("Failed inserting into erlogs table", parent_id=id, e=str(e))
        return False

    return True


async def read_from_file():
    f = os.environ["LOGS"]
    files = f.split(" ")
    async for line in atail("file1.txt"):
        # todo, get file name with it
        insert_log(str(line[0]))


conn = duckdb.connect("./logs.db")
conn.execute(
    "CREATE TABLE IF NOT EXISTS erlogs (id UUID primary key, parent_id UUID, timestamp DOUBLE, string_keys string[], string_values string[], bool_keys string[], bool_values bool[], number_keys string[], number_values double[], raw_log string);"
)
conn.execute(
    "CREATE TABLE IF NOT EXISTS etest (id UUID, parent_id UUID, timestamp DOUBLE)"
)

app = FastAPI()


@app.on_event("startup")
async def read_logs():
    if not "LOGS" in os.environ:
        logger = get_logger()
        logger.info("No logs in os.environ", l=len(os.environ))
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
    logger = get_logger()
    id = str(uuid.uuid4())
    logger.info("search request", id=id)
    body = await request.json()
    if not isinstance(body, object) or isinstance(body, str):
        logger.error("Invalid json", body=body, parent_id=id)
        raise HTTPException(status_code=400, detail="Invalid json")

    user_query = body["query"]
    page = body["page"]

    if user_query == None:
        user_query = ""
    if page == None:
        page = 0

    try:
        p = int(page)
    except Exception:
        logger.error("Invalid page", status_code=400, page=page, parent_id=id)
        raise HTTPException(status_code=400, detail="Page is invalid")

    try:
        logger.info("building query", user_query=user_query, p=p, parent_id=id)
        q = QBuilder()
        q.parse(user_query, p)
        query, params = q.query, q.params
    except Exception as e:
        logger.error("Failed building query", parent_id=id)
        raise HTTPException(status_code=400, detail="Failed building query")

    # maybe put try catch

    try:
        logger.info("executing query", query=query, parent_id=id)
        l = conn.execute(query, params).fetchall()
    except Exception as e:
        logger.error("Failed executing query", query=query, parent_id=id, err=str(e))
        raise HTTPException(status_code=400, detail="Failed executing query")

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

    h = conn.execute(
        "SELECT id, parent_id, timestamp, raw_log from erlogs WHERE id = ?", [id]
    )
    log = h.fetchone()

    # print(log[])
    # if log[1] != None:
    c = conn.execute(
        "SELECT id, parent_id, timestamp, raw_log from erlogs WHERE parent_id = ? ORDER BY timestamp ASC",
        [id],
    )

    children = []
    clogs = c.fetchall()
    for c in clogs:
        children.append({"id": c[0], "parent_id": c[1], "timestamp": c[2], "log": c[3]})
        # print("hi")
        # print(len(child))

    print(children)

    return {"id": log[0], "timestamp": log[2], "log": log[3], "children": children}


@app.post("/")
async def log(request: Request):
    body = await request.json()
    s = ujson.dumps(body)
    # erlog = ErLog(s)
    # print(type(s))
    print(s)
    status = insert_log(s)

    # if isinstance(body, str):
    #     raise HTTPException(status_code=400, detail="Invalid json")

    # flattened = flatten(body)
    # erlog.parse_log(flattened)

    # # todo, use appender or add tis to a batch
    # if erlog._id == None:
    #     erlog._id = uuid.uuid4()

    # conn.execute(
    #     "INSERT INTO erlogs VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
    #     [
    #         erlog._id,
    #         erlog._timestamp,
    #         erlog._parent_id,
    #         erlog._string_keys,
    #         erlog._string_values,
    #         erlog._bool_keys,
    #         erlog._bool_values,
    #         erlog._number_keys,
    #         erlog._number_values,
    #         erlog._raw_log,
    #     ],
    # )
    return {"status": status}


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app)
