from fastapi import FastAPI, HTTPException, Depends
from fastapi import Request, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from app.util import flatten
from app.models import ErLog, User, get_current_user
from app.query import QBuilder
import structlog
from structlog import get_logger
import ujson
import uuid
import secrets
from typing import Annotated
from fastapi.security import OAuth2PasswordBearer
from clickhouse_driver import Client

api_keys = ["ek-QldMOqfEWSpG_u6VCJv3ng_OD97OiXPDh5Luqvc"]

# todo: use flask_httpauth
# https://stackoverflow.com/questions/817882/unique-session-id-in-python/6092448#6092448
# https://medium.com/@anubabajide/rest-api-authentication-in-flask-481518a7479b
# instead of token generation use this qhwer u query the db

structlog.configure(processors=[structlog.processors.JSONRenderer()])


def insert_log(log):
    logger = get_logger()
    id = str(uuid.uuid4())
    if log == "":
        return False

    try:
        l = ujson.loads(log)
    except Exception as e:
        logger.error("Failed parsing log json", parent_id=id, e=str(e))
        return False

    try:
        flattened = flatten(l)
    except Exception as e:
        logger.error("Failed flattening json", parent_id=id, e=str(e))
        return False

    erlog = ErLog(log)
    erlog.parse_log(flattened)

    if erlog._id == None:
        erlog._id = str(uuid.uuid4())

    if erlog._parent_id == None:
        erlog._parent_id = str("00000000-0000-0000-0000-000000000000")
    else:
        print("updating child logs")
        client.execute(
            "ALTER TABLE erlogs UPDATE child_logs = child_logs + 1 WHERE id = %(parent_id)s",
            {"parent_id": str(erlog._parent_id)},
        )

    client.execute(
        "INSERT INTO erlogs VALUES",
        [
            [
                str(erlog._id),
                str(erlog._parent_id),
                erlog._timestamp,
                erlog._string_keys,
                erlog._string_values,
                erlog._bool_keys,
                erlog._bool_values,
                erlog._number_keys,
                erlog._number_values,
                erlog._raw_log,
                erlog._child_logs,
            ]
        ],
    )

    return erlog._id


client = Client(host="localhost", password="test123")

print("Creating tables..")
# client.execute("CREATE DATABASE IF NOT EXISTS e ENGINE = Atomic;")
res = client.execute(
    "CREATE TABLE IF NOT EXISTS erlogs (id UUID primary key, parent_id UUID, timestamp DOUBLE, string_keys Array(String), string_values Array(String), bool_keys Array(String), bool_values Array(Boolean), number_keys Array(String), number_values Array(Double), raw_log String, child_logs UInt32) Engine = MergeTree;"
)
client.execute(
    "CREATE TABLE IF NOT EXISTS users (id UUID primary key, email String, hashed_password String) Engine = MergeTree;"
)

client.execute(
    "CREATE TABLE IF NOT EXISTS sessions (id UUID primary key, session_id String, user_id Int32) Engine = MergeTree;"
)

print("Finished creating tables")

fake_users_db = {
    "johndoe": {
        "username": "johndoe",
        "full_name": "John Doe",
        "email": "johndoe@example.com",
        "hashed_password": "fakehashedsecret",
        "disabled": False,
    },
    "alice": {
        "username": "alice",
        "full_name": "Alice Wonderson",
        "email": "alice@example.com",
        "hashed_password": "fakehashedsecret2",
        "disabled": True,
    },
}


def hash_password(password: str):
    return "hashed" + password


def get_user(session_id: str):
    user_ids = client.execute(
        "SELECT user_id FROM sessions WHERE session_id = %(session_id)s;",
        {"session_id": session_id},
    )
    print(user_ids)
    pass


app = FastAPI()

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="token")


@app.get("/users/me")
async def read_me(current_user: Annotated[User, Depends(get_current_user)]):
    return current_user


@app.get("/items/")
async def items(token: Annotated[str, Depends(oauth2_scheme)]):
    return {"token": token}


# @app.on_event("startup")
# async def read_logs():
#     if not "LOGS" in os.environ:
#         logger = get_logger()
#         logger.info("No logs in os.environ", l=len(os.environ))
#         return

#     loop = asyncio.get_event_loop()
#     loop.create_task(read_from_file())


origins = ["http://localhost", "http://localhost:59971", "*"]

app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.post("/api_key")
async def gen_api_key(request: Request):
    key = "ek-" + secrets.token_urlsafe(29)
    api_keys.append(key)
    return key


@app.post("/metrics")
async def metrics(request: Request):
    return "ok!"


@app.post("/test")
async def test_log(request: Request):
    logger = get_logger()
    logger.info("hello!")
    return "ok"


@app.post("/search")
async def root(request: Request):
    # fine
    logger = get_logger()

    # create a library which wraps structlog but also gives you access to uuid so
    # id = create_id()
    # logger.info()

    id = str(uuid.uuid4())
    # so each of these will send a post request to the service with the id
    logger.info("search request", id=id)

    body = await request.json()
    if not isinstance(body, object) or isinstance(body, str):

        # if the thing already exists with an id in the database
        logger.error("Invalid json", body=body, parent_id=id)
        raise HTTPException(status_code=400, detail="Invalid json")

    user_query = body["query"]
    page = body["page"]
    show_children = body["showChildren"]
    if show_children == "true":
        show_children = True
    else:
        show_children = False
    print(show_children)
    print("SHOW CHILDREN")

    if user_query == None:
        user_query = ""
    if page == None:
        page = 0

    try:
        p = int(page)
    except Exception:
        logger.error("Invalid page", status_code=400, page=page, parent_id=id)
        raise HTTPException(status_code=400, detail="Page is invalid")

    # try:
    logger.info("building query", user_query=user_query, p=p, parent_id=id)
    q = QBuilder()
    q.parse(user_query, p, show_children)
    query, params = q.query, q.params
    # except Exception as e:
    #     print(e)
    #     logger.error("Failed building query", parent_id=id)
    #     raise HTTPException(status_code=400, detail="Failed building query")

    # maybe put try catch

    # try:
    logger.info("executing query", query=query, parent_id=id)
    print(query, params)
    l = client.execute(query, params)
    # except Exception as e:
    # logger.error("Failed executing query", query=query, parent_id=id, err=str(e))
    # raise HTTPException(status_code=400, detail="Failed executing query")

    logs = []
    for log in l:
        logs.append(
            {"id": log[0], "timestamp": log[1], "log": log[2], "child_logs": log[3]}
        )

    return logs


@app.post("/get")
async def get_log(request: Request):
    body = await request.json()

    if isinstance(body, str):
        raise HTTPException(status_code=400, detail="Invalid json")

    id = body["id"]

    if id is None:
        raise HTTPException(status_code=400, detail="Invalid json")

    h = client.execute(
        "SELECT id, parent_id, timestamp, raw_log from erlogs WHERE id = %(s)s",
        {"s": id},
    )

    # TODO: check if less than one
    log = h[0]

    # print(log[])
    # if log[1] != None:
    c = client.execute(
        "SELECT id, parent_id, timestamp, raw_log from erlogs WHERE parent_id = %(s)s ORDER BY timestamp ASC",
        {"s": id},
    )

    children = []
    clogs = c
    for c in clogs:
        children.append({"id": c[0], "parent_id": c[1], "timestamp": c[2], "log": c[3]})
        # print("hi")
        # print(len(child))

    return {"id": log[0], "timestamp": log[2], "log": log[3], "children": children}


@app.post("/")
async def log(request: Request):
    # if not "Authorization" in request.headers:
    #     raise HTTPException(401, "Unauthorized")

    # key = request.headers["Authorization"]
    # key = key.replace("Bearer ", "")

    # if not key in api_keys:
    #     raise HTTPException(401, "Unauthorized")

    body = await request.json()
    s = ujson.dumps(body)
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


@app.get("/children/{log_id}")
async def get_children(log_id: str):
    children = client.execute(
        "SELECT id, timestamp, raw_log, child_logs FROM erlogs WHERE parent_id = %(parent_id)s ORDER BY timestamp ASC",
        {"parent_id": log_id}
    )
    
    return [
        {
            "id": child[0],
            "timestamp": child[1],
            "log": child[2],
            "child_logs": child[3]
        }
        for child in children
    ]


@app.get("/preview/{log_id}")
async def get_log_preview(log_id: str):
    log = client.execute(
        """
        SELECT id, timestamp, raw_log, child_logs 
        FROM erlogs 
        WHERE id = %(log_id)s
        LIMIT 1
        """,
        {"log_id": log_id}
    )
    
    if not log:
        raise HTTPException(status_code=404, detail="Log not found")
        
    return {
        "id": log[0][0],
        "timestamp": log[0][1],
        "log": log[0][2],
        "child_logs": log[0][3]
    }


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app)
