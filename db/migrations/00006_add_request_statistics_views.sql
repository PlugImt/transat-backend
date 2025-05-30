-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create a view showing average response time per endpoint
CREATE OR REPLACE VIEW endpoint_statistics AS
SELECT 
    endpoint,
    method,
    COUNT(*) AS request_count,
    AVG(duration_ms) AS avg_duration_ms,
    MIN(duration_ms) AS min_duration_ms,
    MAX(duration_ms) AS max_duration_ms,
    AVG(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 ELSE 0 END) * 100 AS success_rate_percent,
    MIN(request_received) AS first_request,
    MAX(request_received) AS last_request
FROM 
    request_statistics
GROUP BY 
    endpoint, method
ORDER BY 
    endpoint, method;

-- Create a view showing global average response time
CREATE OR REPLACE VIEW global_statistics AS
SELECT 
    COUNT(*) AS total_request_count,
    AVG(duration_ms) AS global_avg_duration_ms,
    MIN(duration_ms) AS global_min_duration_ms,
    MAX(duration_ms) AS global_max_duration_ms,
    AVG(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 ELSE 0 END) * 100 AS global_success_rate_percent,
    MIN(request_received) AS first_request,
    MAX(request_received) AS last_request
FROM 
    request_statistics;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP VIEW IF EXISTS endpoint_statistics;
DROP VIEW IF EXISTS global_statistics; 