from luqum.parser import parser
from luqum.tree import SearchField, Word, Phrase


class QBuilder:
    def __init__(self):
        self.query = "SELECT id, timestamp, raw_log from erlogs"
        self.params = []
        self.added = False

    def parse(self, q):
        if q == "":
            self.query += " ORDER BY timestamp DESC"
            return

        f = parser.parse(q)
        print(type(f))

        if isinstance(f, Word):
            self.parse_word(f)
        elif isinstance(f, Phrase):
            self.parse_phrase(f)
        elif isinstance(f, SearchField):
            self.parse_searchfield(f)
            pass

        self.query += " ORDER BY timestamp DESC"
        print(self.query, self.params)

    def parse_searchfield(self, s):
        fname = s.name
        fval = s.children[0]

        kf, kv, val = self.parse_value(fval.value)
        if self.added == False:
            self.query += " WHERE "
            self.added = True

        self.query += "{}[list_indexof({}, ?)] {} ?".format(kv, kf, "=")
        self.params.append(fname)
        self.params.append(val)
        # print(s.name, s.children)
        # print(dir(s))
        pass

    def parse_word(self, w):
        print(w.value)
        kf, kv, val = self.parse_value(w.value)
        if self.added == False:
            self.query += " WHERE "
            self.added = True

        self.query += "list_contains({}, ?) OR list_contains({}, ?)".format(kf, kv)
        self.params.append(val)
        self.params.append(val)

    def parse_phrase(self, p):
        print(dir(p))
        print(p.value)
        kf, kv, val = self.parse_value(p.value.replace("'", "").replace('"', ""))
        if self.added == False:
            self.query += " WHERE "
            self.added = True

        self.query += "list_contains({}, ?) OR list_contains({}, ?)".format(kf, kv)
        self.params.append(val)
        self.params.append(val)
        # self.query += "list_contains({}, ?)".format(kf)
        # self.params.append(val)

    def parse_value(self, val):
        v = self.parse_float(val)
        if v != None:
            kf = "number_keys"
            kv = "number_values"
            return kf, kv, v

        v = self.parse_bool(val)
        if v != None:
            kf = "number_keys"
            kv = "number_values"
            return kf, kv, v

        v = val
        kf = "string_keys"
        kv = "string_values"
        return kf, kv, v

    def parse_float(self, val):
        try:
            val = float(val)
            return val
        except Exception:
            return None

    def parse_bool(self, val):
        if val.lower() != "true" and val.lower() != "false":
            return None
        try:
            val = bool(val)
            return val
        except Exception:
            return None


if __name__ == "__main__":
    QBuilder().parse("name.first:foo")
# print(type(f))
# print(dir(f))
# print(f.value)
# print(repr(parser.parse('title:"foo bar"')))
