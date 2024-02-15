# Makego

Makego is a customizable code generator that sets up the basics of an API framework in Go
to let you quickly launch an api.

## Running the app

```
Usage:
  makego [flags] [package_name]

Flags:
      --config string      config file (default is $HOME/makego.yaml)
      --copyright string   copyright holder (and contact if desired)
      --database string    database type to use (mysql, mariadb, postgres, etc) (default "postgres")
  -d, --docker             whether to use docker
      --envprefix string   how to expect env variables to be prefixed
      --folder string      application folder, can be left blank for no folder
  -a, --header             whether to show copyright headers on most files
  -h, --help               help for makego
      --license string     license, can be left blank for proprietary code
      --name string        application name
      --orm string         ORM to use for models (defaults to gorm) (default "gorm")
      --router string      router to use (echo, gin, http, mux) (default "gin")
  -s, --sentry             whether to use sentry
```

`[package_name]` is required if the `go.mod` file is not already set up.

WARNING: Already existing files with the same name will be overwritten. Care should be taken to backup files before running this.

Optionally, a config file can be used with the above flags. If the $HOME/makego.yaml exists, it will be used, so you can use that to cut down on the amount of flags you need to use, especially if you set the same flags consistently.

The config file also includes a `templates` section, where you can specify additional files to create (see below for an example). Templates are given in the form of `filepath: contents`. Where filepath is both relative and regulated to project folder.

### Example config.yml file

```
name: Example App                  # Application Name (setting this in $HOME/makego.yaml is not recommended)
folder: application                # Application folder, can be left out or blank for no folder
license: mit                       # License, can be left blank for proprietary code
copyright: user <user@example.com> # The copyright holder (and contact if desired)
router: gin                        # The router to use (echo, gin, mux, etc)
orm: gorm                          # The orm to use (currently only gorm is supported)
database: postgres                 # The database type to use (mysql, mariadb, postgres, etc)
sentry: true                       # Whether or not to use Sentry
header: true                       # Whether or not to add copyright header to code files
docker: true                       # Whether or not to use Docker
envprefix: app                     # How to expect environment variables to be prefixed, can be left out or blank for no prefix
templates:
  application/example.txt: |       # File name (including path from base folder)
    this
    is
    a
    file
  CONTRIBUTING.md: |               # File name (including path from base folder)
    # Rules

    Here are some rules for contributing to this project
```

If you wish to include the license header in your template, put {{ template "header.template" . }} at the beginning of the file.
