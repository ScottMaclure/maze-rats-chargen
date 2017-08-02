# Maze Rats Character Generator

## Overview

Maze Rats RPG, by Ben Milton:

https://www.drivethrurpg.com/product/197158/Maze-Rats

Generate a new level 1 character, here:

http://maze-rats-chargen.cloud.maclure.info/

Generate in JSON format, to use in other programs:

http://maze-rats-chargen.cloud.maclure.info/json

## Development Notes

Setup:

```
go get -u github.com/kardianos/govendor
```

Tests:

```
# run with coverage report
go test github.com/scottmaclure/maze-rats-chargen -cover

# verbose mode, for logs
go test github.com/scottmaclure/maze-rats-chargen -test.v

# run a single test function
go test -test.run=Test_naturalLanguageSplice
```

Keeping deps up-to-date:

```
govendor list
github.com/jinzhu/inflection
```

## Deployment Notes

```
# on host
dokku postgres:create maze-rats-chargen

# local repo
git remote add dokku dokku@maze-rats-chargen.cloud.maclure.info:maze-rats-chargen

# push to deploy
git push dokku master

```
