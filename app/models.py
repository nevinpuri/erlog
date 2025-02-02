from datetime import datetime, timezone
from pydantic import BaseModel
import uuid


class User(BaseModel):
    username: str
    email: str | None = None


class DBUser(User):
    hashed_password: str


def decode_token(token):
    return User(username=token + "fakecoded", email="hi@example.com")


def get_current_user(token):
    user = decode_token(token)
    return user


class ErLog:
    def __init__(self, raw_log):
        self._id = None
        self._timestamp = 0.00

        self._string_keys = []
        self._string_values = []

        self._number_keys = []
        self._number_values = []

        self._bool_keys = []
        self._bool_values = []

        self._raw_log = raw_log
        self._parent_id = None
        self._child_logs = 0

    def parse_log(self, log):
        for k, v in log.items():
            if k == "timestamp":
                self._timestamp = float(v)
                continue

            if k == "parentId" or k == "parent_id":
                self._parent_id = str(v)
                continue

            if k == "id":
                # uuid parse
                self._id = str(v)
                continue

            if isinstance(v, str):
                self._string_keys.append(k)
                self._string_values.append(v)

            elif (
                isinstance(v, int)
                and not isinstance(v, bool)
                or isinstance(v, float)
                and not isinstance(v, bool)
            ):
                self._number_keys.append(k)
                self._number_values.append(v)

            elif isinstance(v, bool):
                self._bool_keys.append(k)
                self._bool_values.append(v)

        if self._timestamp == 0.00:
            self._timestamp = datetime.now().timestamp()
