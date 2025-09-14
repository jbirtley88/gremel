# Gremel - "Dremel" In Go
Gremel is a utility for interrogating and extracting data from structured sources:

    - JSON
    - CSV
    - Apache CLF
    - Apache combined log
    - Syslog
    - Excel spreadsheets
    - SQL databases

Conceptually and functionally similar to Apache Drill, but considerably more lightweight and simpler to extend the functionality to accommodate whatever exotic data sources you have.

Currently in active development, but Gremel will enable you to do things like this:
```
    $ ls datafiles/
    foo.csv
    bar.csv
    other.csv

    $ gremel datafiles/*csv
    gremel> SELECT foo.name, bar.something FROM foo, bar WHERE foo.name LIKE 'a%' AND foo.id = bar.id
    foo.name      bar.something
    John Smith              111
    Micheal Mouse           123
```

## Why Is That Useful?
Many many times, a lot of the data you need to deal with has different pieces in different places and formats.  Gremel allows you to combine these disparate data and extract the information you need.

As a simple example, imagine you've got 2 datasources:

- `weblogs.log`
  An Apache combined logfile of web requests

  `192.168.143.149 - - [06/Sep/2025:03:46:43 +0100] "GET /api/v1/items HTTP/1.1" 301 1135 "https://docs.example.com/guide" "python-requests/2.31.0" 472`

- `ipaddresses.csv`
  An CSV file containing IP address and datacenter

  `ip,datacenter`

  `192.168.143.149 datacenter1`

Some of the web requests are showing high latency retrieving `/api/foo`, where "*high latency*" means "*more than 2000ms*".

You have a hunch that it might be related to cross-datacenter requests to a database, or it might be a single datacenter with the problem, so you need to see if the high latency is associated with particular datacenters.

To try to do this by hand would be massively complicated and time-consuming, even if your scripting skills are God-tier.

Since Gremel allows you to treat structured files as SQLite tables, and uses SQLite SQL syntax, you can do this easily in Gremel:

```
$ gremel weblogs.log ipaddresses.xlsx
gremel> SELECT
  i.datacenter,
  COUNT(DISTINCT CASE WHEN CAST(w.latency AS INTEGER) > 2000 THEN i.ip END) AS "latency>2000"
FROM ipaddresses AS i
LEFT JOIN weblogs AS w
  ON w.host = i.ip
WHERE w.request LIKE  'GET /api/foo%'
GROUP BY i.datacenter
ORDER BY i.datacenter;

datacenter   latency>2000
----------   ------------
datacenter1           237
datacenter2             0
datacenter3             0
datacenter4             0
```

So it's pretty obvious that whatever the problem is, it's to do with fielding requests in datacenter1.

Gremel is still very much a work-in-progress, please feel free to contribute.

# Parsing your own structured data
Compose from the `BaseParser` in `data/base_parser.go` and:

    - implement the `Parse()` method to convert your input data (an `io.Reader` into a `[]data.Row`).

    - implement the `GetHeaders()` method to return the names of the columns you'll be dealing with

Each entry in the `[]data.Row` is conceptually the same as a SQL `Row`.

# Parser TODO List

# Licences
    Gremel uses the very permissive BSD-3-Clause licence.  There are zero restrictions on how you use it, other than you acknowledge using it (and please give it a star on github!)

