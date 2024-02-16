# Golang Example

# Requirement
- Golang > go version go1.21.4

# Install
```bash
make install
```

# Run
```bash
make run
```

# Note
- `configs` folder to store config
  - config.json for application config file
  - dbconfig.yml for database migration
  - sqlc.yaml for config sql (sql to go generator)
- `db` to store sql file
  - factories to store init table etc
  - migrations to store function or view migration
  - queries to store query that converted to golang
  - seeds to store table seeder
- `generated` to store generated file
- internal to store codebase of internal application
- pkg to store util file and third party that can be imported to other folder