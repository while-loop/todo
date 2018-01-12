todo
=======

<p align="center" style="font-family: verdana, serif; font-size:14pt; font-style:italic">
    <a href="https://godoc.org/github.com/while-loop/todo"><img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/while-loop/todo"><img src="https://img.shields.io/travis/while-loop/todo.svg?style=flat-square"></a>
    <a href="https://github.com/while-loop/todo/releases"><img src="https://img.shields.io/github/release/while-loop/todo.svg?style=flat-square"></a>
    <a href="https://coveralls.io/github/while-loop/todo"><img src="https://img.shields.io/coveralls/while-loop/todo.svg?style=flat-square"></a>
    <a href="LICENSE"><img src="https://img.shields.io/badge/license-Apache 2.0-blue.svg?style=flat-square"></a>
</p>

Auto-generate issues through TODOs in code using your favorite issue tracking
software and version control repository hosting service

Installation
------------

```
$ go get github.com/while-loop/todo/cmd/...
```

Running
-------

#### Command line

```bash
$ todod
```

#### Docker

```bash
$ docker run toyotasupra/todo
```

Example TODOs
-------------

Go:

```go
func Doer() error {
    // TODO(while-loop) Create Doer in main +label1 +feature/Doer @user1
    panic("implement me")
}
```

Python:

```python
def get_homepage(url):
    page = None
    try:
        page = urllib.request.urlopen(url).read()
    except Exception, e:
        # TODO(while-loop) Handle retries when retrieving homepage +api
        print e

    return page
}
```

Supported Services
------------------

#### Repository Hosting Service
- Github
- Gitlab

#### Issue Tracking
- Github
- Jira

#### Publishing Service
- Email (todo)
- RocketChat (todo)

Determining a Complete TODO
---------------------------

To automatically close issues, several methods of completion
are implemented.

Issues are closed when:

- A commit into master referencing the issue ID
- A changelog entry referencing the issue ID

Changelog
---------

The format is based on [Keep a Changelog](http://keepachangelog.com/) 
and this project adheres to [Semantic Versioning](http://semver.org/).

[CHANGELOG.md](CHANGELOG.md)

License
-------
todo is licensed under the Apache 2.0 License. See LICENSE for details.

Author
------

Anthony Alves
