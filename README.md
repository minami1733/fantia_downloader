# fantia_downloader
Download original images and other media data from fantia posts.

## usage
```
fantia_downloader_64.exe --key <key> [--output <path>] [--fanclub <id>] [--excludeFreePlan <bool>]
```
### key
You must specify a key as a required argument. The key is the value of the _session_id cookie that is stored in your browser when you are logged in to [fantia.jp](https://fantia.jp).

### output
Specify the path of the directory where the images will be downloaded.
If not specified, downloads will be created in the execution directory.
### fanclub
Specify the ID of the fan club.
If not specified, all participating fan clubs will be downloaded.
### excludeFreePlan
You can specify whether you want to exclude free plans or not.  
- `true`: Exclude free plans and download only those that are part of paid plans.  
- `false`: Download all participating plans, including free plans.  
If not specified, it is `false`.

## Roadmap
- Add year, month and date option
- Organizing logs