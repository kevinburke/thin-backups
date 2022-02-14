# backupdeleter

Given a directory and a set of flags, prints out files that should be deleted to
match the specified preferences.

- `--extension='.zst'` - Only check files matching this file extension.

- `--everything-until=7d`: Preserve everything in the directory that's less than
  7 days old.

- `--daily-until=14d`: Preserve a maximum of one file per day that's less than
  14 days old. If `--everything-until` is present this will preserve 1 file per
  day between the everything-until date and the daily-until date.

- `--weekly-until=2mo`: Preserve a maximum of one file per week that's less than
2 months old. If `--daily-until` is present this will preserve 1 file per week
between the daily-until date and the weekly-until date. Specify `2mo` or `2
month` or `2 months` for months.

- `--monthly-until=2y`: Preserve a maximum of one file per month (30 days)
that's less than 2 years old. If `--weekly-until` is present this will preserve
1 file per week between the daily-until date and the weekly-until date. Specify
`2mo` or `2 month` or `2 months` for months. For our purposes a month is 30
days.
