MailSlurper
===========

Simple mail SMTP server that slurps mail into oblivion! MailSlurper Server is designed to run on a small network server for multiple developers to use for development and debugging of mail functionality in their applications. When a mail is received the mail item is stored in a database and display in a web-based application.

This application uses a lot of libraries.

* [Gorilla Mux](http://www.gorillatoolkit.org/pkg/mux)
* [Gorilla Context](http://www.gorillatoolkit.org/pkg/context)
* [Alice](https://github.com/justinas/alice)
* [GoHttpService](https://github.com/adampresley/GoHttpService)
* [Logging](https://github.com/adampresley/logging)
* [Bootstrap](http://getbootstrap.com/)
* [Font Awesome](http://fortawesome.github.io/Font-Awesome/)
* [Promiscuous](https://github.com/RubenVerborgh/promiscuous)
* [jQuery](http://jquery.com/) - MIT
* [Moment.js](http://momentjs.com) - MIT
* [RequireJS](http://requirejs.org) - MIT

Compiling
---------
The instructions below assume you have the following tools.

* NodeJS/NPM
* Bower
* Go 1.4.2 (or higher)

```bash
$ bower install
$ npm install
$ cd www/assets/promiscuous
$ node ./build/build.js
$ cd ../../../
$ go get
$ go build
```

License
-------
The MIT License (MIT)

Copyright (c) 2015 Adam Presley

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
