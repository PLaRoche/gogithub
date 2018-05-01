# gogithub
experiment creating an issue in GitHub via go


### build

`make bootstap`
`make build`

### run

Set up GITHUB_TOKEN env var in your local, or use [torus](https://torus.sh) and do it
the right way

`GITHUB_TOKEN=<cool_token_here> ./bin/gogoithub`
or
`torus run ./bin/gogithub`

### Args
`-o` owner of the repo (i.e PLaRoche)
`-r` repository name (i.e. gogithub)
`-v` version tag to use when making the ticket (i.e. v2.0.2)
`-f` filename of the template (optional: tmpl/issue.txt)

## Example

`GITHUB_TOKEN=<cool_token_here> ./bin/gogoithub -o PLaRoche -r gogithub -v v2.0.2 -f tmpl/issue.txt`
