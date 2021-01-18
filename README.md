# tkdo

aglio -i ./docs/tkdo.apib -o ./docs/index.html

NEEDS UPDATING

- uses make for building
- uses aglio to generate html docs

make goals:
- localBuild
- lint
- database
- tearDownDb
- build
- dredd
- postman
- run
- kill
- apiDocs
- delAdmin
- help

requires env vars:
- TKDO_HOST
- TKDO_PORT
- TKDO_USER
- TKDO_PASSWORD
- TKDO_DBNAME
