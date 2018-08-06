gemer
=====
[![GitHub release](http://img.shields.io/github/release/shuheiktgw/gemer.svg?style=flat-square)](release)
[![Build Status](https://travis-ci.org/shuheiktgw/gemer.svg?branch=master)](https://travis-ci.org/shuheiktgw/gemer)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

`gemer` is a CLI tool to release your private Ruby gem with one command.

## Demo
Coming soon...

## Usage
In order to use gemer, first you need to set GitHub personal access token (See [GitHub personal access token section](#how-to-get-a-github-personal-access-token) for more information). Then, simply run the command below.

```
gemer [options]
```

`gemer` command actually does the following stuff for you, to prepare your private Ruby gem to release.

1. Creates a new Pull Request which increments `Version` constant in `version.rb`
2. Drafts a new Release with a new version tag

After running the command above, the last things you need to do is to merge the Pull Request and publish the Release!

### How to get a GitHub personal access token
gemer needs a GitHub personal access token with enough permission to release your gem. If you are not familiar with the access token, [GitHub Help page](https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/) guides you though how to create one.

Please be aware that, for a public repository, you just need `public_repo` scope, and for a private repository, you need whole `repo` scope.

### How to set a GitHub personal access token
Currently, there are two ways to specify your GitHub personal access token.

1. Environment variable
```
$ export GITHUB_TOKEN="Your GitHub personal access token"
```

2. `-t` or `-token` option

```
gemer -t="Your GitHub personal access token" [Other Options]
```


## Example

```
gemer -username='shuheiktgw' -repository='some-gem' -path='lib/github-api-test/version.rb'
```

## Options

You can set these options below:

```bash
$ gemer \
    -t or -token \        # Set a GitHub personal access token
    -u or -username \     # Set a GitHub username
    -r or -repository \   # Set a GitHub repository name
    -b or -branch \       # Set a GitHub branch name your release is based on, default is master
    -p or -path \         # Set a path to version.rb file in your gem, default is lib/[repo name]/version.rb
    -v or -version \      # Return a current version of gemer
    -d or -dry-run \      # Return a current version of gemer
    -major \              # Increments a major version of your gem
    -minor \              # Increments a minor version of your gem
    -patch \              # Increments a patch version of your gem (default)
```


## Author

[Shuhei Kitagawa](https://github.com/shuheiktgw)






