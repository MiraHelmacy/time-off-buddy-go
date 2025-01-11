# time-off-buddy-go

## Usage
```
Calculate when you can take your next dream vacation based on your current time off and time off earned each pay period.

Usage:
  tobuddy [flags]

Flags:
      --ch int        Number of Hours accrued.
      --cm int        Number of Minutes accrued.
      --eh int        Number of hours earned per pay period
      --em int        Number of minutes earned per pay period
  -h, --help          help for tobuddy
  -i, --interactive   Start tobuddy in interactive mode.
  -t, --target int    Total time off time required in hours. (default 40)
  -v, --verbose       Print verbose output.
      --version       version for tobuddy
```

## Build
In the same folder as the one cloned run the following command:
```
go build
```