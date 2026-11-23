ALTER TABLE calendar_events ADD COLUMN IF NOT EXISTS reminder_at TIMESTAMPTZ;

UPDATE calendar_events
SET reminder_at = event_at - reminder_hours * interval '1 hour'
WHERE reminder_hours > 0 AND reminder_at IS NULL;
