# Diskcheck
A Nagios check intended to replace default disk checks. The disk checks bundled with Nagios are parameterized, and expect prior knowledge of a host's disks.

`check_disks` is modeled after Linux utility `df`, and will monitor and alert for all file systems currently mounted.