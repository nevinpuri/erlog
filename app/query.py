from luqum.parser import parser
from luqum.tree import (
    SearchField,
    Word,
    Phrase,
    From,
    To,
    AndOperation,
    OrOperation,
    Group,
)
import uuid


class QBuilder:
    def __init__(self):
        self.query = "SELECT id, Timestamp, raw_log from erlogs"
        self.params = {}
        self.added = False

    def parse(self, q, page):
        if q == "":
            # WHERE parent_id ISNULL
            qu = str(uuid.uuid4())
            self.query += f" WHERE ISNULL(parent_id) ORDER BY Timestamp DESC LIMIT 50 OFFSET %({qu})s "
            self.params.update({qu: int(page * 50)})
            return

        f = parser.parse(q)
        print(repr(f))
        # print(f.children[0].include)
        # print(dir(f.children[0].children[0]))
        # print(dir(f.children[0]))
        self.query += " WHERE "
        self.parse_class(f)

        # nothing was parsed
        # if self.query.strip()[-5:] == "WHERE":
        #     pass

        qeu = str(uuid.uuid4())
        self.query += f" ORDER BY Timestamp DESC LIMIT 50 OFFSET %({qeu})s "
        self.params.update({qeu: int(page * 50)})
        # print(self.query, self.params)

    def parse_class(self, f):
        if isinstance(f, AndOperation):
            self.parse_and(f)
        if isinstance(f, OrOperation):
            self.parse_or(f)
        if isinstance(f, Group):
            self.parse_group(f)
        if isinstance(f, From) or isinstance(f, To):
            self.parse_op(f)
        if isinstance(f, Word):
            self.parse_word(f)
        elif isinstance(f, Phrase):
            self.parse_phrase(f)
        elif isinstance(f, SearchField):
            self.parse_searchfield(f)

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

    def parse_group(self, g):
        print("in group")
        print(dir(g))
        self.query += " ( "
        for child in g.children:
            self.parse_class(child)
        self.query += " ) "
        pass

    def parse_range(self, r):
        """
        TODO: implement
        """
        pass

    def parse_and(self, a):
        self.query += " ( "
        for i in range(len(a.children)):
            self.parse_class(a.children[i])
            if i != len(a.children) - 1:
                self.query += " AND "

        self.query += " ) "

    def parse_or(self, a):
        self.query += " ( "
        for i in range(len(a.children)):
            self.parse_class(a.children[i])
            if i != len(a.children) - 1:
                self.query += " OR "

        self.query += " ) "

    def parse_null(self, fname):
        pass

    def parse_searchfield(self, s):
        fname = s.name
        fval = s.children[0]

        if isinstance(fval, To) or isinstance(fval, From):
            op, fval = self.parse_op(fval)
        else:
            op = "="

        print("fval", fval)
        field = None
        if fval.value.lower() == "null":
            if fname == "id":
                field = "id"
            elif fname == "timestamp":
                field = "Timestamp"
            elif fname == "parent_id":
                field = "parent_id"

            if field != None:
                self.query += "ISNULL({})".format(field)
            else:
                su = str(uuid.uuid4())
                nu = str(uuid.uuid4())
                bu = str(uuid.uuid4())

                self.query += f" has(string_keys, %({su})s) == 0 AND has(number_keys, %({nu})s) == 0 AND has(bool_keys, %({bu})s) == 0"
                self.params.update({su: fname, bu: fname, nu: fname})
                # self.params.append(fname)
                # self.params.append(fname)
                # self.params.append(fname)
            return

        kf, kv, val = self.parse_value(fval.value.replace("'", "").replace('"', ""))

        if fname == "timestamp":
            tu = str(uuid.uuid4())
            self.query += f"Timestamp {op} %({tu})s"
            # self.params.append(val)
            self.params.update({tu: val})
            return

        # elif fname == "parentId"

        # also do this for parentId

        # if self.added == False:
        #     self.query += " WHERE "
        #     self.added = True

        bf = str(uuid.uuid4())
        bv = str(uuid.uuid4())
        self.query += f"{kv}[indexOf({kf}, %({bf})s)] {op} %({bv})s"
        self.params.update({bf: fname, bv: val})
        # self.params.append(fname)
        # self.params.append(val)
        # print(s.name, s.children)
        # print(dir(s))
        pass

    def parse_word(self, w):
        print(w.value)
        kf, kv, val = self.parse_value(w.value)
        # if self.added == False:
        #     self.query += " WHERE "
        #     self.added = True

        uv = str(uuid.uuid4())
        sv = str(uuid.uuid4())
        self.query += f"has({kf}, %({uv})s) OR has({kv}, %({sv})s)"
        self.params.update({uv: val, sv: val})
        # self.params.append(val)
        # self.params.append(val)

    def parse_phrase(self, p):
        print(dir(p))
        print(p.value)
        kf, kv, val = self.parse_value(p.value.replace("'", "").replace('"', ""))
        # if self.added == False:
        #     self.query += " WHERE "
        #     self.added = True

        fv = str(uuid.uuid4())
        sv = str(uuid.uuid4())
        self.query += f"has({kf}, %({fv})s) OR has({kv}, %({sv})s)"
        self.params.update({fv: val, sv: val})
        # self.params.append(val)
        # self.params.append(val)
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
    QBuilder().parse("event:null", 1)
# print(type(f))
# print(dir(f))
# print(f.value)
# print(repr(parser.parse('title:"foo bar"')))
