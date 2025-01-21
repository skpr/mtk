---
sidebar_position: 1
---

# Tutorial

Let's discover **MTK in less than 5 minutes**.

## Scenario

You're a developer with an existing workflow that relies on the `mysqldump` command to export database images. However, there's a critical issue... **passwords and other sensitive information are leaking** into the build pipeline, non-production environments, and even every developer's local setup!

### It's time to implement MTK!

This tutorial demonstrates how developers can use MTK to fix this pipeline issue and export, santize and package a MySQL databases efficiently.

## Step 1 - Update the mysqldump comand

MTK serves as a seamless drop-in replacement for the `mysqldump` command.

The example below illustrates how you can easily replace an existing `mysqldump` command with MTK.

**Before**

```bash
$ mysqldump -u USERNAME -pPASSWORD -h HOSTNAME DATABASE > export.sql
```

**After**

```bash
$ mtk dump -u USERNAME -pPASSWORD -h HOSTNAME DATABASE > export.sql
```

## Configuration File

Now that the command is replaced you can begin to apply rules to santize and boost the speed of the export by skipping content eg. cache tables.

The configuration file, by default, is loaded in the same directory where the command is being executed.

```
mtk.yml
```

The configuration file (below) provides 3 manipulation types:

* **rewrite** - Dynamically modify content being exported
* **nodata** - Only export the structure of a table and skip all the content
* **ingore** - Completely skip the table during the export

```bash
rewrite:
  users_field_data:
    mail: concat(uid, "@localhost") <---- Update all email addresses to the format "uid@localhost"
    pass: '"password"'              <---- Set all passwords as "password"

nodata:
  - cache*           <---- Remove all cache tables to improve performance and storage size
  - captcha_sessions
  - history
  - flood
  - batch
  - queue            <---- Remove tasks which don't need to be queued on other environments
  - sessions
  - semaphore
  - search_api_task
  - search_dataset
  - search_index
  - search_total

ignore:
  - __ACQUIA_MONITORING__
```

## Environment Variables

The MTK dump command can also be configured using environment variables.

```bash
$ export MTK_HOSTNAME=127.0.0.1
$ export MTK_USERNAME=root
$ export MTK_PASSWORD=password
$ export MTK_CONFIG=mtk.yml

$ mtk dump DATABASE > export.sql
```