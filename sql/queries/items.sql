-- name: UpsertItemDef :exec
INSERT INTO item_defs (code, name, category, stackable, base_props, icon_file_key)
VALUES ($1, $2, $3, $4, $5::jsonb, $6)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    category = EXCLUDED.category,
    stackable = EXCLUDED.stackable,
    base_props = EXCLUDED.base_props,
    icon_file_key = EXCLUDED.icon_file_key;

-- name: GetItemDefByCode :one
SELECT *
FROM item_defs
WHERE code = $1;

-- name: ListItemDefs :many
SELECT *
FROM item_defs
ORDER BY category, name;

-- name: CreateItemInstance :one
INSERT INTO item_instances (item_def_id, props)
VALUES ($1, $2::jsonb)
RETURNING *;

-- name: CreateItemInstanceByCode :one
INSERT INTO item_instances (item_def_id, props)
SELECT item_def_id, $2::jsonb
FROM item_defs
WHERE code = $1
RETURNING *;
