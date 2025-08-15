# Gator, the Aggregator
## About

**Gator**, an RSS feed aggregator, is a CLI tool written in *golang*, with use of goose and . It allows users to:
* add RSS feeds from URL to be collected;
* store the collected posts in PostgreSQL database;
* follow and unfollow RSS feeds added in the database by users;
* view summaries of the aggregated posts in the terminal.

### Requirements

You'll need to have [Go](https://go.dev/doc/install) and [PostgreSQL](https://www.postgresql.org/download/) to run gator.

## Installation

```
go install https://github.com/ar3ty/gator
```

### Database configuration

1. Make sure you're on version 15+ of Postgres:

```
psql --version
```
2. (Linux only) Update postgres password:

```
sudo passwd postgres
```
Enter a password, and be sure you won't forget it. You can just use something easy like `postgres`.

3. Start the Postgres server in the background

    * Mac: `brew services start postgresql@15`
    * Linux: `sudo service postgresql start`

4. Enter the psql shell:

    * Mac: `psql postgres`
    * Linux: `sudo -u postgres psql`

5. Create a new database (`gator` is good enough):

```
CREATE DATABASE gator;
```

6. Connect to db:

```
\c gator
```

7. Set the user password (Linux only). I used `postgres` by default:

```
ALTER USER postgres PASSWORD 'postgres';
```

8. You can type `exit` to leave the **psql** shell.

### Configure file

To begin with it is necessary to set up a config file in your home directory, `~/.gatorconfig.json`, with the following content:

```
{
  "db_url": "postgres://username:password@localhost:5432/database?sslmode=disable",
}
```

(It will further also contain a current user logged in **gator**, but the program will set it up automatically.)


`username` - is username you set previously

`password` - is password you set previously

`@localhost` - by default it supposed that you're running it locally. You may not, if you want so

`5432` - is default `port` for SQL databases, and for **PostgreSQL** too

`/database` - is your name of database, set previously

`?sslmode=disable` - query for app, it doesn't supposed to try to use SSL locally

## Usage

#### Login
`gator login <username>` Sets the current user as *username*, if it is registered

#### Register a new user
`gator register <username>` Adds new user and sets as current.

#### List users
`gator users` Lists all the users and indicates which one is currently logged in

#### Agg
`gator agg <time_between_requests>` Fetches RSS feeds, parses and stores them

#### Add feed
`gator addfeed <name> <url>` Add a new RSS feed, titled by *name* from *url* address

#### List feeds
`gator feeds` List all feeds to the console

#### Follow
`gator follow <url>` It takes a single *url* argument and creates a new feed follow record for the current user

#### Following
`gator following` Prints all the names of the feeds the current user is following

#### Unfollow
`gator unfollow <url>` Accepts an *url* as an argument and unfollows it for the current user

#### Browse
`gator browse (<limit>)` Prints posts for current user in terminal. Takes an optional *limit* parameter, defaults to **2**