INSERT INTO recipes (name, url, ingredients, steps)
VALUES (?, ?, ?, ?)
RETURNING id;