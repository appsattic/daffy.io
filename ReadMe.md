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

## Using Daffy ##

You have two options depending on whether you're starting a new project or you have an existing one:

1. clone the repo, remove the `.git` directory, and do the renames below
2. `cp -r` the `etc`, `scripts`, `src`, `static`, `templates`, and `vendor` directories (and the `Makefile` too)

## Renaming things ##

Once you have a copy of the project locally, you should rename everything that is Daffy related to your own name. This
includes:

* daffy.io -> blah.xyz
* io-daffy -> xyz-blah
* daffy -> blah (or at least the 'project name')
* Daffy -> Blah
* DAFFY -> BLAH

This is true for both filenames and for the content within any file. The following commands may help you:

* `find . -name '*daffy.io*'`
* (there is no daffy-io)
* `find . -name '*io-daffy*'`
* `find . -name '*daffy*'`
* `ack daffy.io` # Note: . is a wildcard so will pick up both dot and dash.
* `ack io.daffy` # Note: . is a wildcard so will pick up both dot and dash.
* `ack Daffy`
* `ack daffy`

Try the following:

```
# clone the repo
cd /tmp
git clone https://github.com/appsattic/daffy.io.git blah.xyz
cd blah.xyz

# do some renames
rm daffy-logo.png # make your own :-p
mv ./etc/supervisor/conf.d/{io-daffy,xyz-blah}.conf.m4
mv ./etc/caddy/vhosts/{io.daffy,xyz.blah}.conf.m4

# change all the daffy stuff to blah
find . -type f | xargs perl -pi -e 's{daffy\.io}{blah.xyz}gxms'
git commit -a -m 'Rename all daffy.io to blah.xyz'

find . -type f | xargs perl -pi -e 's{io-daffy}{xyz-blah}gxms'
git commit -a -m 'Rename all io-daffy to xyz-blah'

find . -type f | xargs perl -pi -e 's{daffy}{blah}gxms'
git commit -a -m 'Rename all daffy to blah'

find . -type f | xargs perl -pi -e 's{Daffy}{Blah}gxms'
git commit -a -m 'Rename all Daffy to Blah'

find . -type f | xargs perl -pi -e 's{DAFFY}{BLAH}gxms'
git commit -a -m 'Rename all DAFFY to BLAH'

cp set-env-example.sh set-env-dev.sh
vi set-env-dev.sh
```

## Author ##

[Andrew Chilton](https://chilts.org), [@andychilton](https://twitter.com/andychilton).

For [AppsAttic](https://appsattic.com), [@AppsAttic](https://twitter.com/AppsAttic).

## License ##

[MIT](https://publish.li/mit-license-CPdxXSZb).

## Credits ##

Logo by [Ema Dimitrova](https://thenounproject.com/term/duck/152370/).

(Ends)
