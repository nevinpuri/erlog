from datetime import datetime, timezone
import uuid


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

    def parse_log(self, log):
        for k, v in log.items():
            if k == "timestamp":
                self._timestamp = float(v)
                continue

            if k == "parentId" or k == "parent_id":
                self._parent_id = uuid.UUID(v)
                continue

            if k == "id":
                # uuid parse
                self._id = uuid.UUID(v)
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
