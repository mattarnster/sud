# sud - Server upload/download

This is a small multi-platform application for downloading files from remote servers.

### Usage

To use this, you must create a config file in the same directory as the executable.

Example config.json:

```json
{
  "Protocol": "sftp",
  "Type": "download",
  "AllFiles": true,
  "Files": ["examplefile.txt"],
  "Hosts": ["yourserver:22"],
  "User": "shelluser",
  "KeyPath": "/Users/matt/.ssh/id_rsa",
  "SourceDirectory": "source", 
  "DestinationDirectory": "/Users/matt/destination",
  "DeleteOnRetrieve": true
}
```

Bear in mind, you cannot use "AllFiles:" and "Files:" together, pick one or the other.
This application currently supports "sftp" for the Protocol and "download" for the Type.
