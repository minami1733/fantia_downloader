# fantia_downloader
Download original images and other media data from fantia posts.

## usage
```
fantia_downloader_64.exe --key <key> [--output <path>] [--fanclub <id>] [--excludeFreePlan <bool>] [--date <yyyy-mm-dd>]
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
### date
Specifies a date in `yyyy-mm-dd` format.
If specified, only posts made after the specified date will be downloaded.
## Roadmap
- Organizing logs
- Organize fantia packages

## Copyright and Disclaimer
This software is free software and is made for personal use. Please use it freely only for personal use.
The copyright of this software is held by the author, "あくあ".

The author is in no way affiliated with とらのあな, Fantia, or Fantia Development.
You may use this software at your own risk, provided that you do not violate the Fantia Terms of Use.

I'm not responsible for any failure, damage, or defect caused by the use of this software.
Please use the software at your own risk, as neither I, nor any of my associates, nor any group or organization to which I belong will be liable for any damage, loss, or defect caused by the use of the software.
Please use this software at your own risk.
