# copi

```sh
$ copi -h
```

```
copi [source] [destination]

Copies files and folders from [source] to [destination]

Features:
- Can backup [destination] to other location. (Default keeps 3 backups)
- Can ignore the files and folders described in the list.
- Can transform the files described in the list.

Usage:
  copi [flags]

Flags:
  -b, --backup string      filesystem path to backup folder
  -h, --help               help for copi
  -s, --ignore string      filesystem path to list of files and folders to ignore
  -k, --keep int           number of backups to keep (default 3)
  -r, --remove             remove destination contents (default true)
  -t, --transform string   filesystem path to list of files to transform
```

__Example list of files and folders to ignore__
```ignore.json:```
```json
{
  ".well-known/": {},
  "TemplateFiles/": {},
  "Logs/": {},
  "wwwroot/files/": {},
  "appsettings.json": {},
  "web.config": {}
}
```

In the example above, paths ending in "/" are folders, these folders and the contents will not be copied to the destination and deleted from the destination. (Same thing applies for the file paths.)

__Example list of files to transform__
```transform.json:```
```json
{
  "appsettings.release.json": "appsettings.json",
  "web.release.config": "web.config"
}
```
In the example above, "appsettings.release.json" and "web.release.config" are relative file paths, these files will be moved to corresponding file paths, in that case "appsettings.json" and "web.config". If there is a file with the same name as the destination file, it will be overwritten without a warning.
