# Maze Rats Character Generator

http://maze-rats-chargen.cloud.maclure.info/

https://www.drivethrurpg.com/product/197158/Maze-Rats

## Deployment Notes

```
# on host
dokku postgres:create maze-rats-chargen

# local repo
git remote add dokku dokku@maze-rats-chargen.cloud.maclure.info:maze-rats-chargen

# push to deploy
git push dokku master

```