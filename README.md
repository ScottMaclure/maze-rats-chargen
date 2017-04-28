# Maze Rats Character Generator

http://maze-rats-chargen.cloud.maclure.info/

## Deployment Notes

```
# on host
dokku postgres:create maze-rats-chargen

# local repo
git remote add dokku dokku@maze-rats-chargen.cloud.maclure.info:maze-rats-chargen

# push to deploy
git push dokku master

```