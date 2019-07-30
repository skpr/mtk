mtk-dump
========

Dumps a sanitized database based on a ruleset given by the user.

## Example

```bash
# View the rules which will be used.
$ cat mtk.yml

.......

# Configure the command. NOTE: We also accept myqldump flags.
$ export MTK_DUMP_CONFIG=mtk.yml
$ export MTK_DUMP_HOSTNAME=127.0.0.1
$ export MTK_DUMP_USERNAME=root
$ export MTK_DUMP_PASSWORD=password
$ export MTK_DUMP_DATABASE=test

# Run the command!
$ mtk-dump > sanitized.sql
```