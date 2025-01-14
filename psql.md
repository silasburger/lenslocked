# Create tables

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  age INT,
  first_name TEXT,
  last_name TEXT,
  email TEXT UNIQUE NOT NULL
);
```

# Update query

```sql
UPDATE users
SET first_name = 'Anonymous', last_name = 'Teenager'
WHERE age < 20 AND age > 12;
```

This will update every single record
```sql
UPDATE users
SET first_name = 'Jon';
```

# Delete record

```sql
DELETE FROM users
WHERE id = 1;
```