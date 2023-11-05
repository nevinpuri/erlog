CREATE TABLE IF NOT EXISTS erlogs (
    id UUID,
    parent_id Nullable(UUID),
    Timestamp Float64,
    string_keys Array(String),
    string_values Array(String),
    bool_keys Array(String),
    bool_values Array(String),
    number_keys Array(String),
    number_values Array(Float64),
    raw_log String
)
Engine = MergeTree()
ORDER BY toUnixTimestamp(Timestamp)
PARTITION BY toDate(Timestamp);
