# sud - Server upload/download

This is a small multi-platform application for downloading files from remote servers.

### Usage

To use this, you must create a config file in the same directory as the executable.
Example config.json:

```json
{
  "Protocol": "sftp", // Only supports sftp right now
  "Type": "download", // Only supports download right now
  
  // Pick from AllFiles or Files:
  "AllFiles": true,
  // or
  "Files": ["examplefile.txt"],

  "Hosts": ["yourserver:22"],
  "User": "shelluser",
  "KeyPath": "/Users/matt/.ssh/id_rsa",
  "SourceDirectory": "source", // Path to the directory on the server where the files are placed
  "DestinationDirectory": "/Users/matt/destination",
  "DeleteOnRetrieve": true // Should sud delete the files once it has retrieved them?
}
```
