# Maze Rats Chargen - TODO

## Sooner

* Change ability generation to +2, +1, +0? Or keep random?
* Add dokku pre-flight CHECKS - see http://dokku.viewdocs.io/dokku/deployment/zero-downtime-deploys/
* Complete json data file
* Hook up spell slot generation to character sheet

## Later

* 100% test coverage
* Split codebase into a few package files
* Convert raw "Manner" data to adjectives (e.g. "Anecdotes" -> "Anecdotal")
* Multi-character, local saving support (splash screen)
* Export/import character as json
* Get nginx/dokku to cache static assets with a longer ttl

## Done

### 2017-08-02 (and earlier)

* Add Google Analytics
* Display containers with "handrawn" bordering (as per charsheet)
* Display abilities, details with flexbox
* Render static assets (no caching, do that later with nginx)
* Improve background sentence with "naturalLanguageSplice" (more than just a or an) (needs further work)

### 2017-06-19

* Description field "a or an"