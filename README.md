# gogithub
experiment creating an issue in GitHub via go


### build

`make bootstap`
`make build`

### run

Set up GITHUB_TOKEN in your local, or use [torus](https://torus.sh) and do it
the right way

`torus run ./bin/gogithub`

### Args
`-o` owner (i.e plaroche)
`-r` repository
`-v` version tag to use when making the ticket
