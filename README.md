# Maze Rats Character Generator

## Overview

Maze Rats RPG, by Ben Milton:

https://www.drivethrurpg.com/product/197158/Maze-Rats

Generate a new level 1 character, here:

http://maze-rats-chargen.cloud.maclure.info/

Generate in JSON format, to use in other programs:

http://maze-rats-chargen.cloud.maclure.info/json

## Deployment Notes

```
# on host
dokku postgres:create maze-rats-chargen

# local repo
git remote add dokku dokku@maze-rats-chargen.cloud.maclure.info:maze-rats-chargen

# push to deploy
git push dokku master

```