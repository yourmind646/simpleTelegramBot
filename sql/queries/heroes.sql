-- name: IsHeroExists :one
SELECT hero_id FROM heroes
WHERE user_id = $1;

-- name: CreateHero :exec
INSERT INTO heroes (hero_id, user_id)
VALUES ($1, $2);

-- name: GetHeroByUser :one
SELECT * FROM heroes
WHERE user_id = $1;

-- name: ApplyHeroDelta :one
UPDATE heroes
SET hp = LEAST(100, GREATEST(0, hp + $2)),
    energy = LEAST(100, GREATEST(0, energy + $3)),
    hunger = LEAST(100, GREATEST(0, hunger + $4)),
    thirst = LEAST(100, GREATEST(0, thirst + $5)),
    radiation = LEAST(100, GREATEST(0, radiation + $6))
WHERE user_id = $1
RETURNING *;
