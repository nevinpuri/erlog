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
import time


class QBuilder:
    def __init__(self):
        self.query = "SELECT id, timestamp, raw_log, child_logs from erlogs"
        self.params = {}
        self.added = False

    def parse(self, q, page, show_children, time_range="all"):
        where_conditions = []
        
        # Add time filter first
        time_filter = self.build_time_filter(time_range)
        if time_filter:
            where_conditions.append(time_filter)

        if q == "":
            if show_children == False:
                where_conditions.append("parent_id = '00000000-0000-0000-0000-000000000000'")
        else:
            self.query += " WHERE "
            self.parse_class(f)
            # Remove the last WHERE if nothing was parsed
            if self.query.strip()[-5:] == "WHERE":
                self.query = self.query.strip()[:-5]
            else:
                where_conditions.append(self.query.split("WHERE", 1)[1])
                # Reset query to before WHERE clause
                self.query = self.query.split("WHERE", 1)[0]

        if show_children == False and q != "":
            where_conditions.append("parent_id = '00000000-0000-0000-0000-000000000000'")

        # Combine all WHERE conditions
        if where_conditions:
            self.query += " WHERE " + " AND ".join(where_conditions)

        # Add limit and offset
        u = str(uuid.uuid4())
        self.query += f" ORDER BY timestamp DESC LIMIT 50 OFFSET %({u})s "
        self.params.update({u: int(page * 50)})

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
                field = "timestamp"
            elif fname == "parent_id" or fname == "parentId":
                field = "parent_id"

            if field != None:
                self.query += "{} ISNULL".format(field)
            else:
                str_keyid = str(uuid.uuid4())
                num_keyid = str(uuid.uuid4())
                bool_keyid = str(uuid.uuid4())
                self.query += f" has(string_keys, %({str_keyid})s) == false AND has(number_keys, %({num_keyid})s) == false AND has(bool_keys, %({bool_keyid})s) == false "
                self.params.update({str_keyid: fname})
                self.params.update({num_keyid: fname})
                self.params.update({bool_keyid: fname})
            return

        kf, kv, val = self.parse_value(fval.value.replace("'", "").replace('"', ""))

        if fname == "timestamp":
            v_id = str(uuid.uuid4())
            self.query += f" timestamp {op} %({v_id})s "
            self.params.update({v_id: val})
            return

        # also for id
        if fname == "parent_id" or fname == "parentId":
            v_id = str(uuid.uuid4())
            self.query += f" parent_id {op} %({v_id})s "
            self.params.update({v_id: val})
            return

        if fname == "id":
            v_id = str(uuid.uuid4())
            self.query += f" id {op} %({v_id})s "
            self.params.update({v_id: val})
            return

        # elif fname == "parentId"

        # also do this for parentId

        # if self.added == False:
        #     self.query += " WHERE "
        #     self.added = True

        fname_id = str(uuid.uuid4())
        val_id = str(uuid.uuid4())
        self.query += f"{kv}[indexOf({kf}, %({fname_id})s)] {op} %({val_id})s"
        self.params.update({fname_id: fname})
        self.params.update({val_id: val})
        # print(s.name, s.children)
        # print(dir(s))
        pass

    def parse_word(self, w):
        kf, kv, val = self.parse_value(w.value)
        if kf == "string_keys":
            val_id = str(uuid.uuid4())
            val_id2 = str(uuid.uuid4())
            # self.query += f"has({kf}, %({val_id})s) OR has({kv}, %({val_id2})s)"
            self.query += f" arrayExists(x -> lower(x) LIKE %({val_id})s, {kf}) or arrayExists(x -> lower(x) LIKE %({val_id2})s, {kv})"
            self.params.update({val_id: f"%{val}%"})
            self.params.update({val_id2: f"%{val}%"})
            return
        # if self.added == False:
        #     self.query += " WHERE "
        #     self.added = True

        val_id = str(uuid.uuid4())
        val_id2 = str(uuid.uuid4())
        self.query += f"has({kf}, %({val_id})s) OR has({kv}, %({val_id2})s)"
        self.params.update({val_id: val})
        self.params.update({val_id2: val})

    def parse_phrase(self, p):
        print(dir(p))
        print(p.value)
        kf, kv, val = self.parse_value(p.value.replace("'", "").replace('"', ""))
        # if self.added == False:
        #     self.query += " WHERE "
        #     self.added = True

        val_id = str(uuid.uuid4())
        val_id2 = str(uuid.uuid4())
        self.query += f"has({kf}, %({val_id})s) OR has({kv}, %({val_id2})s)"
        self.params.update({val_id: val})
        self.params.update({val_id2: val})
        # self.query += "has({}, ?)".format(kf)
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

    def build_time_filter(self, time_range):
        if time_range == "all":
            return ""
            
        current_time = time.time()
        
        time_filters = {
            "1h": current_time - 3600,
            "24h": current_time - 86400,
            "7d": current_time - 604800,
            "30d": current_time - 2592000
        }
        
        if time_range in time_filters:
            return f"timestamp >= {time_filters[time_range]}"
            
        return ""


if __name__ == "__main__":
    QBuilder().parse("event:null", 1)
# print(type(f))
# print(dir(f))
# print(f.value)
# print(repr(parser.parse('title:"foo bar"')))
