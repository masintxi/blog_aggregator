# RSS Feeds aggregator

This is a command line feed manager supporting different users in a local environment. Each user can add a new feed to the database and the aggregator will store the posts extracted from the feed. The users will be able to follow feeds and then browse the different posts amongst the followed feeds.

The aggregator is written in GO 1.23.5 and uses PostgreSQL as a database and server. 

# Setup and configuration

First of all, Go and Postgres must be installed to run the application. Then a simple configuration file must be created to set the database connection parameters.

## Go installation 

If Go is not yet installed on the system, you can install the latest Go toolchain from the [official website](https://go.dev/doc/install). Follow the instructions and you are ready to go.

## Postgres installation

For this project Postgres v15 or later is required

1. Postgres can be installed on most Linux distributions using the package manager:

```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
```

Ensure that installation is successful by running the following command:

```bash
psql --version
```

If the command returns the version of Postgres, you are ready to go.

2. In linux systems the postgres password has to be updated. Be sure to *remember the password* you set (note this is the system user's password, not the database password):

```bash
sudo passwd postgres
```

3. Start the Postgres service to run in the background:

```bash
sudo service postgresql start
```

4. Connect to the Postgres database using the default client (psql):

```bash
sudo -u postgres psql
```

You should see a new prompt like this:

```bash
postgres=#
```

5. Create a new database. In this example, we will use the name "gator" but you can choose any name you like:

```bash
CREATE DATABASE gator;
```

6. Connect to the new database just created:

```bash
\c gator
```

The prompt should now show the name of the database:

```bash
gator=#
```

7. This new database needs its own password. Again, Be sure to *remember the password* you set:

```bash
ALTER USER postgres PASSWORD 'your_password';
```

8. The installation is now complete. From here SQL queries can be run against the new database.

To exit the Postgres client, type _exit_ or `\q` and press Enter.

9. Get your *connection string* (that's just a URL with all of the information needed to connect to a database). The format is:

```bash
postgres://username:password@host:port/database
```

Test the connection string by running the following command:

```bash
sudo -u postgres psql "postgres://username:password@host:port/database"
```

It should connect directly to the new database. If it does, note the connection string and exit the Postgres client. We are finished with the postgres installation.

## Configuration file

Create a new file named `gatorconfig.json` in your home directory.

```bash
touch ~/.gatorconfig.json
```

The file should contain the *connection string* to the database in this format:

```json
{
   "db_url": "protocol://username:password@host:port/database?sslmode=disable"
}
```

This _json file_ will keep track of the current logged in user and the database connection parameters.

# Usage

This are the commands available:

| Command | Usage | Description |
|---------|-------|-------------|
| register | resgister "user" | Register a new user to the system |
| login | login "user" | Login to the system |
| users | users | List all the users |
| addfeed | addfeed "feed name" "feed URL" | Add a new feed to the database |
| delfeed | delfeed "feed URL" | Delete a feed from the database |
| feeds | feeds | List all the feeds stored on the database |
| follow | follow "feed URL" | Follow a feed (current user) |
| unfollow | unfollow "feed URL" | Unfollow a feed (current user) |
| following | following | List the feeds followed by the current user |
| browse | browse "limit" | Browse the posts from the feeds followed by the current user. The limit parameter is optional and defaults to 2 posts |
| agg | agg "time interval" | Start the aggregation process (end it with Ctrl+C). The time interval determines how often the aggregator will check for new posts and has to be in the format: 20s, 2m, 1h, etc. This time interval is optional and defaults to 30s |
