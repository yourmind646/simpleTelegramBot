CREATE TABLE IF NOT EXISTS users (
    user_id BIGINT PRIMARY KEY NOT NULL,
    username VARCHAR(33) NULL,
    fullname VARCHAR(129) NULL,
    register_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS admins (
    user_id BIGINT PRIMARY KEY NOT NULL REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TYPE file_type_enum AS ENUM (
  'photo',
  'document',
  'video',
  'audio',
  'voice',
  'sticker',
  'animation',
  'undefined'
);

CREATE TABLE IF NOT EXISTS files (
  file_id     TEXT PRIMARY KEY,
  file_key    VARCHAR(50) NOT NULL UNIQUE,
  file_type   file_type_enum NOT NULL,
  uploaded_by BIGINT NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
  uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS heroes (
    hero_id UUID NOT NULL UNIQUE,
    user_id BIGINT PRIMARY KEY NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    hp INT NOT NULL DEFAULT(100),
    energy INT NOT NULL DEFAULT(100),
    hunger INT NOT NULL DEFAULT(0),
    thirst INT NOT NULL DEFAULT(0),
    radiation INT NOT NULL DEFAULT(0)
);

-- CREATE TYPE item_category_enum AS ENUM ('food','liquid','medicine','materials','weapon','armor');

CREATE TABLE IF NOT EXISTS item_defs (
    item_def_id   BIGSERIAL PRIMARY KEY,
    code          VARCHAR(64) NOT NULL UNIQUE, -- "bandage", "ak74"
    name          VARCHAR(128) NOT NULL,
    category      item_category_enum NOT NULL,
    stackable     BOOLEAN NOT NULL DEFAULT TRUE,
    base_props    JSONB NOT NULL DEFAULT '{}'::jsonb,  -- weight, heal, dmg_base...
    icon_file_key  VARCHAR(50) NULL REFERENCES files(file_key) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS item_instances (
    item_instance_id BIGSERIAL PRIMARY KEY,
    item_def_id      BIGINT NOT NULL REFERENCES item_defs(item_def_id) ON DELETE RESTRICT,
    props            JSONB NOT NULL DEFAULT '{}'::jsonb, -- пер-инстанс: durability, roll, mods
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS inventory_items (
    user_id          BIGINT NOT NULL REFERENCES heroes(user_id) ON DELETE CASCADE, -- владелец (герой)
    item_instance_id BIGINT NOT NULL REFERENCES item_instances(item_instance_id) ON DELETE CASCADE,
    qty              INT NOT NULL DEFAULT 1 CHECK (qty > 0),
    location         VARCHAR(32) NOT NULL DEFAULT 'bag', -- bag/equip/stash...
    PRIMARY KEY (user_id, item_instance_id)
);

INSERT INTO item_defs (code, name, category, stackable, base_props)
VALUES
  ('water_bottle', 'Бутылка воды', 'liquid', true, '{"thirst_restore":25,"volume_ml":500,"weight":0.5}'::jsonb),
  ('canned_food', 'Банка консерв', 'food', true, '{"hunger_restore":20,"weight":0.4}'::jsonb),
  ('first_aid_kit', 'Набор первой помощи', 'medicine', false, '{"hp_restore":40,"weight":1.0}'::jsonb),
  ('bow_simple', 'Простой лук', 'weapon', false, '{"damage_min":8,"damage_max":14,"durability_max":100,"attack_cost_energy":4}'::jsonb)
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name,
    category = EXCLUDED.category,
    stackable = EXCLUDED.stackable,
    base_props = EXCLUDED.base_props;
