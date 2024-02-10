BEGIN;

CREATE TABLE flashcards (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    question text NOT NULL,
    answer text NOT NULL
);

COMMIT;