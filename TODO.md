* ~~sacct --json for time range < works with account specification~~
    * ~~wrap it in timer for excessive accounts (e.g. eta > 40GB json)~~
    * ~~configurable time-range with cmdline switch?~~
* ~~config file reading from 1. /etc/sc.conf , 2. ~/sc.conf~~
* job templates readong from 1. /etc/sc/templates , 2. ~/sc/templates , 3. templates_config from config file if it's set
* build process, handling diferent openapi.json versions
* table sorting, e.g. ctrl-1,2,3 (sort by 1st, 2nd, 3rd,...)
* statistics: add counting per user, per group/account?
* ~~UI table responsiveness when len(rows)>1++k~~
* terminal columns ~150 breaks stats windows
* on low number of rows (<60), on jobtab, stats and info overlap, messing up screen (laptops)
* ~~make sorting of boxes static! otherwise flapping in display happens on same number items~~
