-- name: GrantItemByCode :one
WITH def AS (
  SELECT item_def_id, stackable
  FROM item_defs
  WHERE code = $2
),
existing AS (
  SELECT ii.item_instance_id
  FROM inventory_items ii
  JOIN item_instances ins ON ins.item_instance_id = ii.item_instance_id
  WHERE ii.user_id = $1
    AND ii.location = $4
    AND ins.item_def_id = (SELECT item_def_id FROM def)
    AND ins.props = $5::jsonb
  LIMIT 1
),
upd AS (
  UPDATE inventory_items AS ii
  SET qty = ii.qty + $3
  WHERE ii.user_id = $1
    AND ii.item_instance_id = (SELECT item_instance_id FROM existing)
    AND (SELECT stackable FROM def) = true
  RETURNING ii.item_instance_id
),
new_inst AS (
  INSERT INTO item_instances (item_def_id, props)
  SELECT item_def_id, $5::jsonb
  FROM def
  WHERE NOT EXISTS (SELECT 1 FROM upd)
  RETURNING item_instance_id
),
ins_inv AS (
  INSERT INTO inventory_items (user_id, item_instance_id, qty, location)
  SELECT $1, item_instance_id, $3, $4
  FROM new_inst
  RETURNING item_instance_id
)
SELECT item_instance_id FROM upd
UNION ALL
SELECT item_instance_id FROM ins_inv;

-- name: GetInventoryItemForUse :one
SELECT
  ii.user_id,
  ii.item_instance_id,
  ii.qty,
  ii.location,
  d.item_def_id,
  d.code,
  d.name,
  d.category,
  d.stackable,
  d.base_props,
  d.icon_file_key,
  ins.props AS instance_props,
  ins.created_at
FROM inventory_items ii
JOIN item_instances ins ON ins.item_instance_id = ii.item_instance_id
JOIN item_defs d ON d.item_def_id = ins.item_def_id
WHERE ii.user_id = sqlc.arg(user_id)
  AND ii.item_instance_id = sqlc.arg(item_instance_id)
FOR UPDATE OF ii, ins;

-- name: ConsumeInventoryItem :one
WITH cur AS (
  SELECT ii.user_id, ii.item_instance_id, ii.qty
  FROM inventory_items ii
  WHERE ii.user_id = sqlc.arg(user_id)
    AND ii.item_instance_id = sqlc.arg(item_instance_id)
    AND ii.qty > 0
  FOR UPDATE OF ii
),
dec AS (
  UPDATE inventory_items ii
  SET qty = ii.qty - 1
  WHERE ii.user_id = (SELECT user_id FROM cur)
    AND ii.item_instance_id = (SELECT item_instance_id FROM cur)
    AND (SELECT qty FROM cur) > 1
  RETURNING ii.qty AS qty_after
),
del_inst AS (
  DELETE FROM item_instances ins
  WHERE ins.item_instance_id = (SELECT item_instance_id FROM cur)
    AND (SELECT qty FROM cur) = 1
  RETURNING 1 AS removed
)
SELECT
  COALESCE((SELECT qty_after FROM dec), 0)::int AS qty_after,
  EXISTS(SELECT 1 FROM del_inst) AS removed
FROM cur;

-- name: DecrementDurabilityAndMaybeBreak :one
WITH cur AS (
  SELECT ins.item_instance_id
  FROM item_instances ins
  JOIN inventory_items ii ON ii.item_instance_id = ins.item_instance_id
  WHERE ii.user_id = sqlc.arg(user_id)
    AND ins.item_instance_id = sqlc.arg(item_instance_id)
  FOR UPDATE OF ins
),
del_inst AS (
  DELETE FROM item_instances ins
  WHERE ins.item_instance_id = (SELECT item_instance_id FROM cur)
    AND COALESCE((ins.props->>'durability')::int, 0) <= sqlc.arg(damage)::int
  RETURNING 1 AS broken
),
upd AS (
  UPDATE item_instances ins
  SET props = jsonb_set(
    ins.props,
    '{durability}',
    to_jsonb(GREATEST(COALESCE((ins.props->>'durability')::int, 0) - sqlc.arg(damage)::int, 0)),
    true
  )
  WHERE ins.item_instance_id = (SELECT item_instance_id FROM cur)
    AND NOT EXISTS (SELECT 1 FROM del_inst)
  RETURNING (ins.props->>'durability')::int AS durability_after
)
SELECT
  COALESCE((SELECT durability_after FROM upd), 0)::int AS durability_after,
  EXISTS(SELECT 1 FROM del_inst) AS broken
FROM cur;

-- name: GetInventoryByUser :many
SELECT
  ii.user_id,
  ii.item_instance_id,
  ii.qty,
  ii.location,
  d.item_def_id,
  d.code,
  d.name,
  d.category,
  d.stackable,
  d.base_props,
  d.icon_file_key,
  ins.props AS instance_props,
  ins.created_at
FROM inventory_items ii
JOIN item_instances ins ON ins.item_instance_id = ii.item_instance_id
JOIN item_defs d ON d.item_def_id = ins.item_def_id
WHERE ii.user_id = $1
ORDER BY d.category, d.name;
