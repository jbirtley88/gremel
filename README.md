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