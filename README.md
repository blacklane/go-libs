# go-libs

Blacklane's [multi-module repository](https://github.com/golang/go/wiki/Modules#faqs--multi-module-repositories) for Go libraries :)

It's a work in progress, everyone is welcome to help. We use 
[GitHub Issues](https://github.com/blacklane/go-libs/issues) for bug report and requests.
 
## Repository structure

Every folder is a library, currently we have:
 - [logger](logger): a wrapper around [zerolog](https://github.com/rs/zerolog)
 and compliant with our [standards](http://handbook.int.blacklane.io/monitoring/kiev.html)
 - [tracking](tracking): adds UUID and other things to the [context](https://golang.org/pkg/context/)
 - [x](x): Experimental, they are experiments or work in progress which are not ready to be recommended
 as a standard
   - [events](https://github.com/blacklane/go-libs/tree/events/x/events): An event producer and consumer abstraction and a kafka implementation. Currently a work in progress on the branch [`events`](https://github.com/blacklane/go-libs/tree/events)

The maintainers of each library is defines in the CODEOWNERS file. If the library does not provide one,
consider the top-level [CODEOWNERS](CODEOWNERS) as the maintainers.

## Releases and versioning

As defined by [Go modules](https://github.com/golang/go/wiki/Modules) we use semantic versioning
as defined by [semver](https://semver.org/).

The tag for version 1.2.3 of module "my-repo/foo/rop" is "foo/rop/v1.2.3".

## Commit messages, Issues and Pull Requests

They should be prefixed by the affected module, e.g.:
 - `logger: release vx.y.z`
 - `tracking: add functions to read/set Gopher name in the context`


For commit messages we follow the [Go project standard](https://tip.golang.org/doc/contribute.html#commit_messages)
as much as possible.

>  A rule of thumb is that it should be written so to complete the sentence 
> "This change modifies X to _____." 
> That means it does not start with a capital letter, 
> is not a complete sentence, and actually summarizes the result of the change. 

## How to fetch Go private Modules
Run `GOPRIVATE=github.com/blacklane/* go get -u github.com/blacklane/go-libs/x/events`<br>
If the above command doesn't work follow one of the steps bellow:
* Set Github token from your Blacklane account in the env var `$GITHUB_TOKEN`. If you don't have a token yet follow this instructions: [create a token](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/creating-a-personal-access-token)
* Set git config to use ssh instead of https: `git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"`
