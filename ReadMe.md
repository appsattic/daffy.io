# daffy.io #

A simple skeleton webapp which does a number of things for you.

![Logo for daffy.io](https://raw.githubusercontent.com/appsattic/daffy.io/master/daffy-logo.png "daffy.io")

* runs a server (which check for which PORT to listen to)
* deals with encrypted sessions for logins
* allows login with the following social networks:
    * Twitter
    * GitHub
    * others easily added
* opens and uses a BoltDB key/value datastore
* stores all social IDs in a `social` table
* stores all users in a `user` table
* stores the mapping from social ID to user separately
* allows user to change their username
* allows user to see which social networks they have logged in with

This project is not designed to be deployed but instead to be cloned and changed as you will.

This is a [gb](https://getgb.io/) project, but probably can be converted to the vanilla go toolchain easily enough.

## Author ##

[Andrew Chilton](https://chilts.org), [@andychilton](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com), [@AppsAttic](https://twitter.com/AppsAttic).

## License ##

[MIT](https://appsattic.mit-license.org/2017/).

## Credits ##

Logo by [Ema Dimitrova](https://thenounproject.com/term/duck/152370/).

(Ends)
