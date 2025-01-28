# Introduction
A CLI tool to perform tasks such as uploading files, purging cache, and more using Bunny.net APIs.

# Usage
```

Usage:
  bunny-cli [command]

Available Commands:
  completion       Generate the autocompletion script for the specified shell
  gen-docs         Generate CLI documentation
  help             Help about any command
  purge-cache-full Purge the full pull zone cache
  purge-cache-url  Purge cache for URLs listed in a file
  upload-file      Upload a single file
  upload-folder    Upload a folder concurrently

Flags:
  -h, --help   help for bunny-cli


Required Environment Variables(For uploading files):
  STORAGE_ZONE_NAME
  STORAGE_ACCESS_KEY

Required Environment Variables(For purging cache):
  BUNNYCDN_API_KEY

Optional Environment Variables:
  STORAGE_ZONE_HOSTNAME  (Default: sg.storage.bunnycdn.com)

Use "bunny-cli [command] --help" for more information about a command
```