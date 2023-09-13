MTK (MySQL Toolkit)
===================

A toolkit for exporting, sanitizing and packaging MySQL database.

## The Goal of this Project

To provide developers with the tools to share MySQL database dumps in a safe and repeatable manner.

## Example

The following example scenario will showcase how developers can use mtk to dump and package a MySQL image database.

### Dump a Sanitized Database with MTK

**Configuration File**

Prior to packaging, developers should assess the data that is stored within the database and determine what data should be sanitized or dropped.

The following mtk configuration file cover some of common Drupal 7/8 data which should be sanitized or dropped.

```bash
$ cat mtk.yml

---
rewrite:
  # Drupal 8.
  users_field_data:
    mail: concat(uid, "@localhost")
    # Quoting here denotes an explicit string rather than mysql expression. 
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
```

**Dump the database**

Now that you have a configuration file, it is time to dump the database.

The following command will dump a MySQL database using the configuration file created in the previous step.

```bash
$ export MTK_HOSTNAME=127.0.0.1
$ export MTK_USERNAME=root
$ export MTK_PASSWORD=password
$ export MTK_CONFIG=mtk.yml

$ mtk-dump test > db.sql
```

### Build a Database Image using Docker

Next we can build a database image using the database dump created by the step prior.

First, create a Dockerfile which will import the sanitized database dump.

In this example we are using MySQL images from our [image repository](https://github.com/skpr/image-mysql).

```dockerfile
FROM docker.io/skpr/mysql:8.x-v3-latest

ADD db.sql /tmp/db.sql
RUN database-import local local local /tmp/db.sql
```

Next, build the image with Docker.

```bash
docker build -t docker.io/my/database:latest .
```

Hooray! You have successfully packaged a sanitized MySQL database image!

### Integrate

The database image can then be integrated into local development workflows using Docker Compose (or similar).

Below is an example of how this can be configured.

```bash
$ cat docker-compose.yml

---
version: "3"

services:

  # Services used as part of the local development environment.
  mysql:
    image: docker.io/my/database:latest
```
