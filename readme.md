### TCP OpenTibia Login server
This is an experimental version of Golang implementation of an OpenTibia TCP login server.

What is it a must-have to grow to a non-experimental?
- add an IP rate limiter to wrong login tries
- add support to other database schemas (Nostalrius, OTX2, TFS)
- the current version is tested on a 7.72 game version, so to add a configurable support to other protocol versions is definitelly a must-have

Other features that is a nice-to-have:
- add support to gameservers which have cast-system
- add support to gameservers which have cam-system
- add full support to multiple words
- add support to proxies

### To use, you should:

- fill your database credentials in .env file (this repo has .env.example file with the needed fields)
- fill config.yaml with hostname and IP addresses

