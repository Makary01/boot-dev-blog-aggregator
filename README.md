# Setup

## 1. Install PostgreSQL

Install PostgreSQL and create a database named `gator`.

```sql
CREATE DATABASE gator;
```

---

## 2. Run the database migrations

From the `sql/schema` directory, run:

```sh
goose postgres "<protocol>://<username>:<password>@<host>:<port>/gator" up
```

For example:

```sh
goose postgres "postgres://postgres:password@localhost:5432/gator?sslmode=disable" up
```

---

## 3. Create the configuration file

Create a file named `.gatorconfig.json` in your home directory (`~/.gatorconfig.json`):

```json
{
  "db_url": "postgres://<username>:<password>@<host>:<port>/gator?sslmode=disable",
  "current_user_name": "<username>"
}
```

The `current_user_name` field is optional and will be set automatically after registering or logging in.

---

## 4. Build the application

From the project root, build the executable:

```sh
go build -o gator .
```

---

## 5. Install the executable

Move the binary to a directory on your `PATH`, for example:

```sh
sudo mv gator /usr/local/bin/
```

You should now be able to run the application from anywhere:

```sh
gator register <username>
```

# Commands

## `register <username>`

Creates a new user and automatically logs in as that user.

**Usage**

```sh
gator register alice
```

---

## `login <username>`

Logs in as an existing user.

**Usage**

```sh
gator login alice
```

---

## `users`

Lists all registered users. The currently logged-in user is marked with `(current)`.

**Usage**

```sh
gator users
```

---

## `reset`

Deletes all users from the database.

**Usage**

```sh
gator reset
```

---

## `agg <duration>`

Starts the RSS feed aggregator, fetching feeds at the specified interval.

**Usage**

```sh
gator agg 30s
gator agg 5m
gator agg 1h
```

---

## `addfeed <name> <url>`

Adds a new RSS feed and automatically follows it.

**Usage**

```sh
gator addfeed "Boot.dev Blog" https://blog.boot.dev/index.xml
```

---

## `feeds`

Lists all registered feeds and the user who added each one.

**Usage**

```sh
gator feeds
```

---

## `follow <feed_url>`

Follows an existing feed.

**Usage**

```sh
gator follow https://blog.boot.dev/index.xml
```

---

## `following`

Lists all feeds followed by the currently logged-in user.

**Usage**

```sh
gator following
```

---

## `unfollow <feed_url>`

Stops following a feed.

**Usage**

```sh
gator unfollow https://blog.boot.dev/index.xml
```

---

## `browse [limit]`

Displays the latest posts from feeds you follow. If `limit` is omitted, the default is `2`.

**Usage**

```sh
gator browse
gator browse 10
```
