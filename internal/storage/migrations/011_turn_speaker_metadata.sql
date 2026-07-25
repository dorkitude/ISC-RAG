ALTER TABLE turns ADD COLUMN IF NOT EXISTS speaker_name TEXT;
ALTER TABLE turns ADD COLUMN IF NOT EXISTS speaker_role TEXT;

CREATE INDEX IF NOT EXISTS turns_speaker_name_idx ON turns (speaker_name);
CREATE INDEX IF NOT EXISTS turns_speaker_role_idx ON turns (speaker_role);
