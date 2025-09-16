# Gremel - "Dremel" In Go
Gremel is a utility for interrogating and extracting data from structured sources:

- JSON
- CSV
- Apache CLF
- Apache combined log
- Syslog
- Excel spreadsheets

Conceptually and functionally similar to Apache Drill, but considerably more lightweight and simpler to extend the functionality to accommodate whatever exotic data sources you have.

Currently in active development, but Gremel will enable you to do things like this (you can copy/paste these commands - they refer to the files in `test_resources/...`):
```
    $ ls datafiles/
    foo.csv
    bar.csv
    other.csv

    $ ./gremel
    Type '.help' or '?' for help
    gremel> -- Mount the 'accounts.csv' to the 'accounts' table
    gremel> .mount accounts test_resources/accounts.csv
    gremel> -- Mount the 'people.csv' to the 'people' table
    gremel> .mount people test_resources/people.csv
    gremel> -- We want the full name (in people.csv) and email (in people.csv and accounts.csv)
gremel> SELECT
    ...>   people.id,
    ...>   accounts.username,
    ...>   people.fullname
    ...> FROM
    ...>   people,
    ...>   accounts
    ...> WHERE
    ...>   accounts.email = people.email
    ...> LIMIT 10;
id    username         fullname                
--    --------         --------                
1     krawnsley0       Marcellina Benedicto    
2     eatmore1         Aubert Akers            
3     wshirtcliffe2    Felicle Paynton         
4     walywin3         Ashla Palatini          
5     mbowdler4        Babbette Tratton        
6     supshall5        Mendy Doiley            
7     dclapp6          Nedi Grattan            
8     awillarton7      Liva Eagell             
9     sagnew8          Celestyn Perchard       
10    zmarkussen9      Devonna Stedmond        
10 rows
```
You can also one-shot things in a single command-line (which can be handy for scripting):
```sh
$ echo -e ".mount accounts test_resources/accounts.json\nSELECT username FROM accounts;" | ./gremel --silent
klunbech0           
rpoyle1             
agartery2           
mallsup3            
isim4
```
The `--silent` (or `-q`) removes all non-data output, so it is incredibly useful for scipting as part of a pipeline.

Gremel is still very much a work-in-progress, please feel free to contribute.

## Why Is That Useful?
Many many times, a lot of the data you need to deal with has different pieces in different places and formats.  Gremel allows you to combine these disparate data and extract the information you need.

As a simple example, imagine you've got a REST API which is showing very high latency for some requests for `/api/foo` (where "*high latency*" means "*more than 2000ms*").

You have a hunch that it might be related to cross-datacenter requests to a database, or it might be a single datacenter with the problem, so you need to see if the high latency is associated with particular datacenters.

You've got 2 datasources - the combined weblog file and an Excel spreadsheet of IP address and which datacenter that handles that IP address (because request routing is based on hashing the source IP, so it is consistent):

- `weblogs.log`
  An Apache combined logfile of web requests

  `192.168.143.149 - - [06/Sep/2025:03:46:43 +0100] "GET /api/v1/items HTTP/1.1" 301 1135 "https://docs.example.com/guide" "python-requests/2.31.0" 472`

- `ipaddresses.xslx`
  An Excel file containing IP address and datacenter

  `ip,datacenter`

  `192.168.143.149 datacenter1`

You somehow need to:

- find the requests to `/api/foo` which have > 2000ms latency from the weblog
- grab the client IP address for that high latency request from the weblog

No problem so far, you can do that with a single line of `awk(1)`.

But here's where it gets tricky:

- you need look up which datacenter is handling each IP address in the Excel file
- figure out if your hunch is right - i.e. the high latency is associated with only 1 datacenter

To try to do this by hand would be massively complicated and time-consuming, even if your scripting skills are God-tier.  Sure, for this example you could export the Excel file to CSV and use a bit of `awk(1)` and `join(1)` but there could be many different data sources in many different formats for the problem at hand.

Since Gremel allows you to treat structured files as SQLite tables, and uses SQLite SQL syntax, you can do this easily:

```
$ ./gremel
gremel> -- Mount the ipaddresses.xlsx
gremel> .mount ipaddresses test_resources/ipaddresses.xlsx
2025/09/14 15:36:11 Creating table ipaddresses with SQL:
DROP TABLE IF EXISTS ipaddresses;
CREATE TABLE ipaddresses (
    ip TEXT,
    datacenter TEXT
);

gremel> -- Mount the web server logs
gremel> .mount weblogs test_resources/weblogs.log
2025/09/14 15:35:53 Creating table weblogs with SQL:
DROP TABLE IF EXISTS weblogs;
CREATE TABLE weblogs (
    host TEXT,
    ident TEXT,
    status INTEGER,
    request TEXT,
    latency INTEGER,
    user TEXT,
    size INTEGER,
    referer TEXT,
    time INTEGER,
    useragent TEXT
);

gremel> -- Run the SQL query which groups high latency for /api/foo by datacenter
gremel> SELECT
    ...> i.datacenter,
    ...> COUNT(DISTINCT CASE WHEN CAST(w.latency AS INTEGER) > 2000 THEN i.ip END) AS "latency>2000"
    ...> FROM ipaddresses AS i
    ...> LEFT JOIN weblogs AS w ON w.host = i.ip
    ...> WHERE w.request LIKE  'GET /api/foo%'
    ...> GROUP BY i.datacenter
    ...> ORDER BY i.datacenter 
    ...> ;
Executing SQL: SELECT i.datacenter, COUNT(DISTINCT CASE WHEN CAST(w.latency AS INTEGER) > 2000 THEN i.ip END) AS "latency>2000" FROM ipaddresses AS i LEFT JOIN weblogs AS w ON w.host = i.ip WHERE w.request LIKE  'GET /api/foo%' GROUP BY i.datacenter ORDER BY i.datacenter
datacenter     latency>2000    
----------     ------------    
datacenter1    237             
datacenter2    0               
datacenter3    0               
datacenter4    0               
4 rows
```

So it's pretty obvious that whatever the problem is, it's to do with fielding requests in datacenter1.

You can fire up gremel and copy/paste the above commands to try it out for yourself.

# Parsing your own structured data
Compose from the `BaseParser` in `data/base_parser.go` and:

- implement the `Parse()` method to convert your input data (an `io.Reader` into a `[]data.Row`).

- implement the `GetHeaders()` method to return the names of the columns you'll be dealing with

Each entry in the `[]data.Row` is conceptually the same as a SQL `Row`.

# Parser TODO List

# Gremel TODO
- REST api
- Allow the mounting of files via HTTP File upload
- Allow mounting of datasources by url, e.g. `https://myapi.com/api/v1/foo`
- Allow exporting to files (particularly Excel)

# Licences
Gremel uses the very permissive BSD-3-Clause licence.  There are zero restrictions on how you use it, other than you acknowledge using it (and please give it a star on github!)

