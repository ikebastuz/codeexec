-- name: InsertCodeExecutionResult :one
INSERT INTO code_execution_results (
    code, language, encoded_code, stdout, stderr, error, build_duration, exec_duration
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, created_at;

-- name: GetCodeExecutionResult :one
SELECT * FROM code_execution_results
WHERE encoded_code = $1 AND language = $2
ORDER BY created_at DESC
LIMIT 1; 