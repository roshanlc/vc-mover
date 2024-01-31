# Vericred-Mover
- A basic file watcher tool that watches `Downloads` directory in Linux based distros for vericred files and moves them to their respective directory.

> Example: 
`vericred_some_name_20230101.json` will be moved to `vericred_some_name` directory.


## Supports
- Auto manage pre-existing vericred files
- Watch for new vericred files

## Todo
- Add support for deleting old files 
- Add config support


## Installation 
1. Install as a binary only
```bash
# make sure you have go installed on your machine.
go install github.com/roshanlc/vericred-mover@latest

# Run: (Logs are available at ~/.cache/vericred-mover/*.log)
vericred-mover
```

2. Install as a systemd service
> In-progress


## Author
Roshan Lamichhane