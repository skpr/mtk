mtk-dump
========

Dumps a sanitized database based on a ruleset given by the user.

## Example

```bash
# View the rules which will be used.
$ cat mtk.yml

rewrite:
  # Drupal 8.
  users_field_data:
    mail: concat(uid, "@localhost")
    # Quoting here denotes an explicit string rather than mysql expression 
    pass: '"password"'
  # Drupal 7.
  users:
    mail: concat(uid, "@localhost")
    pass: '"password"'

nodata:
  - cache*
  - captcha_sessions
  - history
  - flood
  - batch
  - queue
  - sessions
  - semaphore
  - search_api_task
  - search_dataset
  - search_index
  - search_total

ignore:
  - __ACQUIA_MONITORING__
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
