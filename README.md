# Vintage Jazz Records Go/Gin Tutorial
https://go.dev/doc/tutorial/web-service-gin


# DB Setup:

```
$ psql -U postgres
```

```
$ CREATE USER myuser WITH PASSWORD 'mypassword';
```

```
$ CREATE DATABASE mydb;
```

```
$ GRANT ALL PRIVILEGES ON DATABASE mydb TO myuser;
```

Navigate into specific DB
```
$ \c vintagejazzrecord
```

Check Table Existence
```
$ \d+ albums
```

Drop Table
```
$ DROP TABLE IF EXISTS albums;
```