<p align="center"><img src="logo/horizontal.png" alt="mailslurper" height="100px"></p>

MailSlurper
===========

MailSlurper is a small SMTP mail server that slurps mail into oblivion! MailSlurper is perfect for individual developers or small teams writing mail-enabled applications that wish to test email functionality without the risk or hassle of installing and configuring a full blown email server. It's simple to use! Simply setup MailSlurper, configure your code and/or application server to send mail through the address where MailSlurper is running, and start sending emails! MailSlurper will capture those emails into a database for you to view at your leisure.

Compiling
---------
The following are general instructions for compiling MailSlurper. Your details may vary a bit here and there. The below example is based on a Unix-style system, such as Ubuntu or OSX. Furthermore for instructional purposes it is assumed that your GOPATH is set to *~/code/go*, and that you have a folder in your source directory called **github.com**. Your setup may vary. The instructions below also assume you have the following already installed.

* Go 1.10 (or higher)
* Git

```bash
$ cd ~/code/go/src/github.com
$ mkdir mailslurper
$ cd mailslurper
$ git clone https://github.com/mailslurper/mailslurper.git
$ go get github.com/mjibson/esc
$ cd mailslurper/cmd/mailslurper
$ go get
$ go generate
$ go build
```

Quickstart With Docker
----------------------

```bash
# Build container image (adjust repo location as necessary)
docker build -t mailslurper 'https://github.com/mailslurper/mailslurper#master'
# Run a temporary container. Note that upon shutdown, all stored messages will be lost when using this config.
docker run -it --rm --name mailslurper -p 8080:8080 -p 8085:8085 -p 2500:2500 mailslurper
```

Library and Framework Credits
-----------------------------
This application uses a lot of great open source libraries.

* [BlockUI](http://jquery.malsup.com/block/) - MIT
* [bluemonday](https://github.com/microcosm-cc/bluemonday) - BSD 3-Clause
* [Bootstrap](http://getbootstrap.com/) - MIT
* [bootstrap-dialog](https://github.com/nakupanda/bootstrap3-dialog) - MIT
* [bootstrap-growl](https://github.com/ifightcrime/bootstrap-growl) - MIT
* [Date Range Picker for Bootstrap](http://www.daterangepicker.com) - MIT
* [Gorilla Context](https://github.com/gorilla/context) - BSD 3-Clause
* [Gorilla Secure Cookie](https://github.com/gorilla/securecookie) - BSD 3-Clause
* [Gorilla Sessions](https://github.com/gorilla/sessions) - BSD 3-Clause
* [Copier](https://github.com/jinzhu/copier) - MIT
* [Echo](https://github.com/labstack/echo) - MIT
* [errors](https://github.com/pkg/errors) - BSD 2-Clause
* [esc](https://github.com/mjibson/esc) - MIT
* [Font Awesome](http://fortawesome.github.io/Font-Awesome/) - Fonts under OFL License, CSS under MIT license
* [GoUUID](https://github.com/nu7hatch/gouuid) - MIT
* [go-cache](https://github.com/patrickmn/go-cache) - MIT
* [go-mssqldb](https://github.com/denisenkom/go-mssqldb)
* [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql) - Mozilla Public License Version 2.0
* [go-sqlite3](https://github.com/mattn/go-sqlite3) - MIT
* [Handlebars](http://handlebarsjs.com) - MIT
* [jQuery](http://jquery.com/) - MIT
* [jwt-go](https://github.com/dgrijalva/jwt-go) - MIT
* [lightbox2](http://lokeshdhakar.com/projects/lightbox2/) - MIT
* [Logrus](https://github.com/sirupsen/logrus) - MIT
* [Moment.js](http://momentjs.com) - MIT
* [NPO](https://github.com/getify/native-promise-only) - MIT
* [open](https://github.com/skratchdot/open-golang) - MIT

Themes by Thomas Park at [Bootswatch](http://bootswatch.com/).

License
-------
The MIT License (MIT)

Copyright (c) 2013-2018 Adam Presley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
