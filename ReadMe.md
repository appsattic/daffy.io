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

## Renaming things ##

Once you have a copy of the project locally, you should rename everything that is Daffy related to your own name. This
includes:

* daffy.io -> example.com
* io-daffy -> com-example
* daffy -> example (or at least the 'project name')
* Daffy -> Example
* DAFFY -> EXAMPLE

This is true for both filenames and for the content within any file. The following commands may help you:

* `find . -name '*daffy.io*'`
* `find . -name '*daffy-io*'`
* `find . -name '*io-daffy*'`
* `find . -name '*daffy*'`
* `ack daffy.io` # Note: . is a wildcard so will pick up both dot and dash.
* `ack io.daffy` # Note: . is a wildcard so will pick up both dot and dash.
* `ack Daffy`
* `ack daffy`

## Author ##

[Andrew Chilton](https://chilts.org), [@andychilton](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com), [@AppsAttic](https://twitter.com/AppsAttic).

## License ##

[MIT](https://publish.li/mit-license-CPdxXSZb).

## Credits ##

Logo by [Ema Dimitrova](https://thenounproject.com/term/duck/152370/).

(Ends)
