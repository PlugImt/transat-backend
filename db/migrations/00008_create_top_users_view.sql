-- Create view for top users statistics
CREATE OR REPLACE VIEW top_users_statistics AS
SELECT 
    user_email AS email,
    COUNT(*) AS request_count,
    AVG(duration_ms) AS avg_duration_ms,
    MIN(duration_ms) AS min_duration_ms,
    MAX(duration_ms) AS max_duration_ms,
    (100.0 * COUNT(CASE WHEN status_code < 400 THEN 1 END) / COUNT(*)) AS success_rate_percent,
    MIN(request_received) AS first_request,
    MAX(request_received) AS last_request,
    COUNT(CASE WHEN status_code < 400 THEN 1 END) AS success_count,
    COUNT(CASE WHEN status_code >= 400 THEN 1 END) AS error_count
FROM 
    request_statistics
WHERE 
    user_email IS NOT NULL -- Only consider authenticated users
GROUP BY 
    user_email
ORDER BY 
    request_count DESC
LIMIT 10; -- Get only top 10 users 