# OPDS 2.0 Library

This Go library is meant for applications that want to manipulate OPDS 1.x or 2.0 feeds.

## Using the app

In addition to libraries, this project can be compiled into a binary that converts OPDS 1.x into OPDS 2.0.

The converter simply takes an OPDS 1.X URI as an argument and prints an OPDS 2.0 feed.

Example : ./libopds2-go http://www.feedbooks.com/store/recent.atom

## Features

- [x] OPDS 2.0 model
- [x] OPDS 1.x model
- [x] Parsing OPDS 1.x
- [x] Generating OPDS 2.0
- [ ] Parsing OPDS 2.0
- [ ] Helpers for OPDS 2.0
