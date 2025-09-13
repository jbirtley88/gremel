# Gremel - "Dremel" In Go
Gremel is a utility for interrogating and extracting data from structured sources:

    - JSON
    - CSV
    - Apache CLF
    - Apache combined log
    - Syslog
    - Excel spreadsheets
    - SQL databases

Conceptually and functionally similar to Apache Drill.

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

Some of the web requests are showing high latency retrieving `/api/foo`, where "*high latency*" means "*more than 1000ms*".

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

# Conditions Package (Deprecated)
The `conditions` package evaluates expressions to a `true` or `false` value.  Conceptually (and syntactically) similar to what you'd find in a `WHERE` clause in SQL.

It supports variable expansion, and nested / parenthesised expressions

## Operators
| Operator | Symbol | Syntax | Description |
| -------- | ------ | ------ | ----------- |
| AND | `AND` | `lhs AND rhs` | Will be `true` if both `lhs` and `rhs` are true |
| OR | `OR` | `lhs OR rhs` | Will be `true` if either `lhs` and `rhs` are true |
| EQ | `=` | `lhs = rhs` | Will be `true` if `lhs` is equal to `rhs` |
| NEQ | `!=` | `lhs != rhs` | Will be `true` if `lhs` is not equal to `rhs` |
| LT | `<` |  `lhs < rhs` | Will be `true` if `lhs` is less than`rhs` |
| LTE | `<=` | `lhs <= rhs` | Will be `true` if `lhs` is less than or equal to `rhs` |
| GT | `>` | `lhs > rhs` | Will be `true` if `lhs` is greater than to `rhs` |
| GTE | `>=` | `lhs >= rhs` | Will be `true` if `lhs` is greater than or equal to `rhs` |
| EREG | `=~` | `lhs =~ regex` | Will be `true` if `lhs` matches regular expression `regex` |
| NEREG | `!~` | `lhs !~ rhs` | Will be `true` if `lhs` does not match regular expression `regex` |
| IN | `IN` | `lhs IN [array, of, values]` | Will be `true` if `lhs` is in the array |
| NOTIN | `NOT IN` | `lhs NOT IN [array, of, values]` | Will be `true` if `lhs` is not in the array |
| CONTAINS | `CONTAINS` | `[array, of, values] CONTAINS rhs` | Will be `true` if `rhs` is in the array |
| CONTAINS | `NOT` | `[array, of, values] CONTAINS rhs` | Will be `true` if `rhs` is in the array |

## Examples
```
    package main
    
    import (
        "context"
        "fmt"
        "log"
        "strings"

        "github.com/jbirtley88/gremel/conditions"
    )

    func main() {
        // Variables can live in the context
        ctx := context.WithValue(
            context.Background(),
            "foo",
            123,
        )

        // Variables can also live in a map.
        // The map is checked first and ifthe variable is not found, then the context
        // is checked.
        vars := map[string]any{
            "bar": 999,
        }

	expression := `{foo} > 100 AND {bar} < 1000`
        p := conditions.NewParser(strings.NewReader(expression))
        expr, err := p.Parse()
        if err != nil {
            log.Fatal(err)
        }

        // Evaluate expression passing data for {vars}
        // viper.Set("config.undefined_is_zero", true)
        r, err := conditions.Evaluate(ctx, expr, vars)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("%s: %t\n", expression, r)
    }
```

# Credits
The conditions package is built on top of [https://github.com/oleksandr/conditions](https://github.com/oleksandr/conditions).

The main differences are

    - Changed the syntax of variables from `[foo]` to `{foo}`.
    - Added capability for deriving zero-value for unknown variables via `config.undefined_is_zero: true`
    - Added capability for variables to be stored (and retrieved) as context values
    - Added `CONTAINS`.
    - Added float comparison with epsilon error tolerance.
    - Optimized long array `IN`/`CONTAINS` operator.
    - Removed redundant RWMutex for better performance.

# Licences
    Gremel uses the very permissive BSD-3-Clause licence.  There are zero restrictions on how you use it, other than you acknowledge using it (and please give it a star on github!)
    https://github.com/oleksandr/conditions (MIT)

