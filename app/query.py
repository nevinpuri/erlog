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
import structlog

logger = structlog.get_logger()

class QBuilder:
    def __init__(self):
        self.query = "SELECT id, timestamp, raw_log, child_logs FROM erlogs"
        self.params = {}
        self.where_conditions = []

    def parse(self, user_query, page, show_children, time_range="all"):
        try:
            # Add time filter
            time_filter = self.build_time_filter(time_range)
            if time_filter:
                self.where_conditions.append(time_filter)

            # Add parent filter
            if not show_children:
                self.where_conditions.append("parent_id = '00000000-0000-0000-0000-000000000000'")

            # Parse user query if present
            if user_query:
                try:
                    tree = parser.parse(user_query)
                    condition = self.parse_node(tree)
                    if condition:
                        self.where_conditions.append(condition)
                except Exception as e:
                    logger.error("Query parse error", error_message=str(e))

            # Combine WHERE conditions
            if self.where_conditions:
                self.query += " WHERE " + " AND ".join(self.where_conditions)

            # Add ordering and pagination
            self.query += " ORDER BY timestamp DESC"
            page_param = str(uuid.uuid4())
            self.query += f" LIMIT 50 OFFSET %({page_param})s"
            self.params[page_param] = int(page * 50)

            return self.query, self.params

        except Exception as e:
            logger.error("Query build error", error_message=str(e))
            raise

    def parse_node(self, node):
        if isinstance(node, AndOperation):
            left = self.parse_node(node.children[0])
            right = self.parse_node(node.children[1])
            return f"({left} AND {right})"
        
        if isinstance(node, OrOperation):
            left = self.parse_node(node.children[0])
            right = self.parse_node(node.children[1])
            return f"({left} OR {right})"
        
        if isinstance(node, Group):
            return self.parse_node(node.children[0])
        
        if isinstance(node, SearchField):
            return self.parse_field(node)
        
        if isinstance(node, Word):
            return self.parse_word(node.value)
        
        if isinstance(node, Phrase):
            return self.parse_word(node.value.strip('"\''))
        
        return None

    def parse_field(self, node):
        field = node.name
        value = node.children[0].value.strip('"\'')
        
        # Handle special fields
        if field in ['id', 'parent_id', 'timestamp']:
            param_id = str(uuid.uuid4())
            self.params[param_id] = value
            return f"{field} = %({param_id})s"
            
        # Handle array fields
        param_id = str(uuid.uuid4())
        self.params[param_id] = field
        value_param = str(uuid.uuid4())
        self.params[value_param] = value
        
        return f"""
            (arrayExists(x -> x = %({param_id})s, string_keys) AND 
             arrayFirst(i -> string_keys[i] = %({param_id})s, arrayEnumerate(string_keys)) > 0 AND 
             string_values[arrayFirst(i -> string_keys[i] = %({param_id})s, arrayEnumerate(string_keys))] = %({value_param})s)
            OR
            (arrayExists(x -> x = %({param_id})s, number_keys) AND 
             arrayFirst(i -> number_keys[i] = %({param_id})s, arrayEnumerate(number_keys)) > 0 AND 
             toString(number_values[arrayFirst(i -> number_keys[i] = %({param_id})s, arrayEnumerate(number_keys))]) = %({value_param})s)
            OR
            (arrayExists(x -> x = %({param_id})s, bool_keys) AND 
             arrayFirst(i -> bool_keys[i] = %({param_id})s, arrayEnumerate(bool_keys)) > 0 AND 
             toString(bool_values[arrayFirst(i -> bool_keys[i] = %({param_id})s, arrayEnumerate(bool_keys))]) = %({value_param})s)
        """

    def parse_word(self, value):
        param_id = str(uuid.uuid4())
        self.params[param_id] = f"%{value}%"
        
        return f"""
            arrayExists(x -> position(lower(x), lower(%({param_id})s)) > 0, string_values) OR
            arrayExists(x -> position(lower(x), lower(%({param_id})s)) > 0, string_keys) OR
            position(lower(raw_log), lower(%({param_id})s)) > 0
        """

    def build_time_filter(self, time_range):
        if time_range == "all":
            return None
            
        current_time = time.time()
        time_filters = {
            "1h": current_time - 3600,
            "24h": current_time - 86400,
            "7d": current_time - 604800,
            "30d": current_time - 2592000
        }
        
        if time_range in time_filters:
            param_id = str(uuid.uuid4())
            self.params[param_id] = time_filters[time_range]
            return f"timestamp >= %({param_id})s"
            
        return None


if __name__ == "__main__":
    QBuilder().parse("event:null", 1)
# print(type(f))
# print(dir(f))
# print(f.value)
# print(repr(parser.parse('title:"foo bar"')))
