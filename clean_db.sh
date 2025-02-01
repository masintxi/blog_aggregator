cd sql/schema/
goose postgres postgres://postgres:postgres@localhost:5432/gator down
echo table deleted
goose postgres postgres://postgres:postgres@localhost:5432/gator up
echo table created