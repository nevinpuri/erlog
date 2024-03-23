select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfYear(timestamp) as date, toHour(timestamp) as hour, toMinute(timestamp) as minute, COUNT(*) as count from metrics GROUP BY minute, hour, date, month, year ORDER BY year, month, date, hour, minute;

select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfYear(timestamp) as date, toHour(timestamp) as hour, COUNT(*) as count from metrics GROUP BY hour, date, month, year ORDER BY year, month, date, hour;

select toYear(timestamp) as year, toMonth(timestamp) as month, toDayOfMonth(timestamp) as date, toHour(timestamp) as hour, toMinute(timestamp) as minute, COUNT(*) as count from metrics GROUP BY minute, hour, date, month, year ORDER BY year, month, date, hour, minute;