from luqum.parser import parser
from luqum.tree import SearchField, Word, Phrase, From, To


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
        print(repr(f))
        print(f.children[0].include)
        print(dir(f.children[0].children[0]))
        print(dir(f.children[0]))
        self.parse_class(f)

        self.query += " ORDER BY timestamp DESC"
        print(self.query, self.params)

    def parse_class(self, f):
        if isinstance(f, From) or isinstance(f, To):
            self.parse_op(f)
        if isinstance(f, Word):
            self.parse_word(f)
        elif isinstance(f, Phrase):
            self.parse_phrase(f)
        elif isinstance(f, SearchField):
            self.parse_searchfield(f)
            pass

        pass

    def parse_op(self, o):
        """
        Will return inner word
        """
        if isinstance(o, To):
            op = "<"
        else:
            op = ">"

        if o.include:
            op += "="

        print(repr(o.children[0]))
        return op, o.children[0]

    def parse_searchfield(self, s, op="="):
        fname = s.name
        fval = s.children[0]

        if isinstance(fval, To) or isinstance(fval, From):
            op, fval = self.parse_op(fval)
            print(op)
        else:
            op = "="

        kf, kv, val = self.parse_value(fval.value.replace("'", "").replace('"', ""))
        if self.added == False:
            self.query += " WHERE "
            self.added = True

        self.query += "{}[list_indexof({}, ?)] {} ?".format(kv, kf, op)
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
            kf = "bool_keys"
            kv = "bool_values"
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
            if val.lower() == "true":
                return True
            else:
                return False
            # val = bool(vala
            # return val
        except Exception:
            return None


if __name__ == "__main__":
    QBuilder().parse("hey.whatever:<=100")
# print(type(f))
# print(dir(f))
# print(f.value)
# print(repr(parser.parse('title:"foo bar"')))
