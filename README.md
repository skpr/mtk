MySQL Toolkit
=============

Toolkit for exporting, sanitizing and packaging MySQL database.

## Scenario

The following example scenario showcase multiple _mtk_ tools:

* [dump](/dump)
* [package](/package)

In the following scenario we will be dumping a sanitized version of the database and packaging the image into a container image which can then be consumed by developers using Docker Compose.

The benefits for this approach are:

* **Safe** - The database are sanitized for safety.
* **Repeatable** - Images can be recreated very quickly given they are packaged as an image.
* **Easy** - Integrates with Docker Compose.

### Dump

**View the rules**

The following rules cover some common Drupal 7/8 scenarios were data should be sanitized or dropped.

```bash
$ cat mtk.yml

sanitize:
  tables:
    # Drupal 7
    - name: users
      fields:
        - name: mail
          value: "SANITIZED_MAIL"
        - name: pass
          value: "SANITIZED_PASSWORD"
    # Drupal 8
    - name: users_field_data
      fields:
        - name: mail
          value: "SANITIZED_MAIL"
        - name: pass
          value: "SANITIZED_PASSWORD"

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

The following command will dump a sanitized version of the database using the below environment variables.

```bash
$ export MTK_DUMP_CONFIG=mtk.yml
$ export MTK_DUMP_HOSTNAME=127.0.0.1
$ export MTK_DUMP_USERNAME=root
$ export MTK_DUMP_PASSWORD=password
$ export MTK_DUMP_DATABASE=test

$ mtk-dump > db.sql
```

### Package

```bash
$ docker run -it -v $HOME/.docker:/kaniko/.docker \
                 -v $(pwd):/workspace \
                 skpr/mtk-build:latest --context=/workspace \
                                       --dockerfile=/Dockerfile \
                                       --single-snapshot \
                                       --verbosity fatal \
                                       --destination=docker.io/my/image:latest \
                                       --destination=docker.io/my/image:$(date +%F)
```

### Integrate

```bash
$ cat docker-compose.yml

version: "3"

services:

  # Services used as part of the local development environment.
  mysql:
    image: docker.io/my/image:latest
```