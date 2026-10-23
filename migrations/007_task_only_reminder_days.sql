UPDATE homework h
SET reminder_days = COALESCE(s.reminder_days, ARRAY[7]::smallint[])
FROM user_settings s
WHERE s.chat_id = h.chat_id
  AND h.reminder_days IS NULL;

UPDATE homework
SET reminder_days = ARRAY[7]::smallint[]
WHERE reminder_days IS NULL;

ALTER TABLE user_settings
    DROP COLUMN IF EXISTS reminder_days;
