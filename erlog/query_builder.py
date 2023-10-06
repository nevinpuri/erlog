from erlog_utils import Operation


class QueryBuilder:
    def __init__(self):
        self.query = "SELECT id, timestamp, raw_log from erlogs"
        self.params = []
        self.added = False

    def add(self, ar, keyword, limit=50):
        if ar[0] == None:
            self.query += " ORDER BY timestamp DESC"
            return

        for expr in ar:
            if isinstance(expr.val, float):
                key_field = "number_keys"
                val_field = "number_values"
            elif isinstance(expr.val, bool):
                key_field = "bool_keys"
                val_field = "bool_values"
            # string
            else:
                key_field = "string_keys"
                val_field = "string_values"

            if self.added == False:
                self.query += " WHERE "
                self.added = True
            else:
                # literally just change and to or
                self.query += " {} ".format(keyword.upper())

            op = self.get_op_from_operation(expr.operation)
            self.query += "{}[list_indexof({}, ?)] {} ?".format(
                val_field, key_field, op
            )

            self.params.append(expr.field)
            self.params.append(expr.val)

        self.query += " ORDER BY timestamp ASC"

    def get_op_from_operation(self, op):
        if op == Operation.Eq:
            return "="
        elif op == Operation.Gt:
            return ">"
        elif op == Operation.Gte:
            return ">="
        elif op == Operation.Lt:
            return "<"
        elif op == Operation.Lte:
            return "<="

    def get_query_and_params(self):
        return self.query, self.params
