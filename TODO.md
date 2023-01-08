## TIPS:

- Serve local package docs:
  - `godoc -http=:6060` and navigate to /pkg/github.com/

## BUGS:

- Check if -o has .pdf and append if necessary
  - already does this, update the README - apparently this doesn't actually work? - it does actually work?
  - Need to have proper (integration) test cases for this

## TODO:

- bin sizes with different flag implementations:
  - kingpin:      8.4M
  - pflag:        6.0M
  - builtin flag: 5.6M

- Per golang naming conventions, acronyms (like PDF) should be all-caps
- Tweak test layout
  - https://github.com/juju/utils/blob/master/uuid_test.go
  - https://medium.com/@mwhudson/a-patronizing-guide-to-what-makes-a-good-unit-test-fe895951ce35#.mxnfuax4f
- Profile execution
  - https://www.youtube.com/watch?v=uBjoTxosSys
- Replace gofpdf import with my forked copy
- Replace .md files with [reStructuredText](http://docutils.sourceforge.net/rst.html)
- Drop kingpin for the std flags package
- Reimplement Generate.Write to use io.Writer
  - gofpdf's Write() requires a full io.Writer, so is there any point if we have to construct and handle io.Writer stuff in main()?
  - io.Writer is just a warpper around a single Write() function
- Document exports
  - appease `golint`
- Makefile, particularly for test commands
  - move test commands out of .travis.yml
  - https://medium.com/@joshroppo/setting-go-1-5-variables-at-compile-time-for-versioning-5b30a965d33e#.rbvyijh6m
    - https://github.com/mholt/caddy/blob/master/main.go#L232
    - https://github.com/mholt/caddy/blob/master/build.bash
  - https://github.com/ailispaw/talk2docker/blob/master/Makefile
  - https://gist.github.com/isaacs/62a2d1825d04437c6f08
- Sort files array ignoring scanlation group's tag
  - http://adampresley.com/2015/09/06/sorting-inventory-items-in-go.html
  - Use https://github.com/erengy/anitomy (C++ project) for folder title parsing? Or just do a naive brackets strip & look for numbers?
    - Maybe just strip first set of '[]' and sort everything after
- Rewrite Generate() to use our own PDF writer implementation

- reorganize into flat structure
  - ~~note that we're doing dot imports on the `fs` and `generator` packages, so isolating these as packages doesn't actually
  get us anything in terms of encapsulation~~
  - writing a initialization method for Generator has solved the dot imports issue, so leave as-is for now?
