# VC-MOVER
- A basic file watcher tool that watches `Downloads` directory in Linux based distros for "vc" files and moves them to their respective directory.

> Example: 
`vc_some_name_20230101.json` will be moved to `vc_some_name` directory.


## Supports
- Auto manage pre-existing "vc" files
- Watch for new "vc" files

## Todo
- Add support for deleting old files 
- Add config support


## Installation 
1. Install as a binary only
```bash
# make sure you have go installed on your machine.
go install github.com/roshanlc/vc-mover@latest

# Run: (Logs are available at ~/.cache/vc-mover/*.log)
vc-mover
```

2. Install as a systemd service
> In-progress


## Author
Roshan Lamichhane