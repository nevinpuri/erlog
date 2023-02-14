create table er_logs (
    Id              UUID,
    Timestamp       Int64,
    ServiceName     String,
    StringKeys     Array(String),
    StringValues   Array(String),
    NumberKeys     Array(String),
    NumberValues   Array(Float64),
    BoolNames      Array(String),
    BoolValues     Array(UInt8)
) Engine MergeTree()
PARTITION BY toDate(Timestamp)
ORDER BY toUnixTimestamp(Timestamp)