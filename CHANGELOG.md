# Change Log

## 0.0.6 (20 February 2021)

- Fix `--selfupdate`, issue with quotes [#9](https://github.com/addshore/dockerit/issues/9)

## 0.0.5 (20 February 2021)

- Add `--magic` command to use some predefined options for various images
  - [composer](https://hub.docker.com/_/composer) with `--me`, `--pwd`, `--home`, `COMPOSER_HOME`
  - [node](https://hub.docker.com/_/node) with `--me`, `--pwd`, `--home`, `npm_config_cache`
- Add `--volume` flag for a custom volume mount [#8](https://github.com/addshore/dockerit/issues/8)
- `--port` flag can now be used multiple times
- `-e` is no longer a shorthand for `--env`

## 0.0.4 (7 February 2021)

- Add `--selfupdate` flag for easy updating of the command [#5](https://github.com/addshore/dockerit/issues/5)
- Fix `runtime error: slice bounds out of range` when outputing [#3](https://github.com/addshore/dockerit/issues/3)

## 0.0.3 (7 February 2021)

- Add `--env` flag [#2](https://github.com/addshore/dockerit/issues/2)
- Fix input issues (including copy and paste) [#4](https://github.com/addshore/dockerit/issues/4)

## 0.0.2 (10 January 2021)

- Added docker API version negotiation [#1](https://github.com/addshore/dockerit/pull/1)
- Fixed nil pointer due to multiple terminal resets [#1](https://github.com/addshore/dockerit/pull/1)

## 0.0.1 (10 January 2021)

Initial version...
