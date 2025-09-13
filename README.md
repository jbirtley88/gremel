# Gremel - "Dremel" In Glo
Gremel is a utility for interrogating and extracting data from structured sources:

    - JSON
    - CSV
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

Gremel is still very much a work-in-progress.  One part of it that is fairly mature is the `conditions` package.

# Parsing your own structured data
Compose from the `BaseParser` in `data/base_parser.go` and:

    - implement the `Parse()` method to convert your input data (an `io.Reader` into a `[]map[string]any`).

    - implement the `GetHeaders()` method to return the names of the columns you'll be dealing with

Each entry in the `[]map[string]any` is conceptually the same as a SQL `Row`.

The default `BaseSelector` in `helper/selector.go` will then operate quite happily on the `[]map[string]any`:

    - `Select(needle string, haystack []map[string]any) ([]map[string]any, error)`

    - `Where(rows []map[string]any, whereClause string) ([]map[string]any, error)`

    - `Order(input []map[string]any, by string) ([]map[string]any, error)`

The `BaseSelector` implementation of `Where()` delegates row matching to the `conditions` package.

# Conditions Package
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
