-- name: IsFileExistsByKey :one
SELECT file_id FROM files
WHERE file_key = $1;

-- name: CreateFile :exec
INSERT INTO files (file_id, file_key, file_type, uploaded_by)
VALUES ($1, $2, $3, $4);

-- name: GetFileByKey :one
SELECT * FROM files
WHERE file_key = $1;

-- name: UpdateFile :exec
UPDATE files
SET file_id = $1,
    file_type = $2,
    uploaded_by = $3,
    uploaded_at = now()
WHERE file_key = $4;
