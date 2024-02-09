# VC-MOVER
- A basic file watcher tool that watches `Downloads` directory in Linux based distros for "vc" files and moves them to their respective directory.

> Example: 
`vc_some_name_20230101.json` will be moved to `vc_some_name` directory.


## Supports
- Auto manage pre-existing "vc" files
- Watch for new "vc" files

## Todo
- Add config support
- Add support for deleting old files 


## Installation 
1. Install using one-liner

    Installs the binary and sets it auto-run on login
```bash
# make sure you have go installed on your machine.
bash -c "$(wget -qO-  https://github.com/roshanlc/vc-mover/blob/master/scripts/setup.sh)"
```

2. Install as a binary only
```bash
# make sure you have go installed on your machine.
go install github.com/roshanlc/vc-mover@latest

# Run: (Logs are available at ~/.cache/vc-mover/*.log)
vc-mover
```

## Author
Roshan Lamichhane