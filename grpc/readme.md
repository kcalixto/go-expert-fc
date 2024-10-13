```
protoc --go_out=. --go-grpc_out=. proto/course_category.proto
```

```
sqlite3 db.sqlite
CREATE TABLE categories (id string, name string, description string);
```

```
brew tap ktr0731/evans
brew install evans
```

```
evans -r repl
```