## untsc - extractor for TSComp files

A proprietary format by The Stirling Group. These files 
are DCL imploded with some basic header.

#### Requires

To build:
 - golang 1.16

To install binary release:
 - See releases for builds for most major OS and architectures.

#### Usage

```shell
usage: untsc [-h] -e FILENAME [-d PATH]

Extract TSComp files.

optional arguments:
  -h, --help            show this help message and exit
  -e FILENAME, --extract FILENAME
                        The TSC file to extract.
  -d PATH, --destination PATH
                        An optional output folder.

Please file bugs on the GitHub Issues. Thanks
```

### Format Info
Reverse engineering of the format I found:

- 6 byte file signature: "\x65\x5D\x13\x8C\x08\x01"
- 8 bytes unknown
- Entry header: "\x12"
- 4 byte integer - Compressed data size
- 10 bytes unknown
- 1 byte filename length
- n bytes Filename (null terminated)
- remainder - DCL Imploded payload

### References
Archive format information: http://fileformats.archiveteam.org/wiki/TSComp
DOS Archiver: https://www.sac.sk/download/pack/tscomp.zip

#### Author

- Kris Hunt (@ctfkris)