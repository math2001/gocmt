# Assumptions made during development

## Configuration

- in conf.d file, we don't load yml files in subdirectories.
    Loading from subdirectory is easy to do if needed
- in conf.d, the yml files are loaded in alphabetical order, and NOT numerically aware.
    (2-foo.yml comes AFTER 10-bar.yml).
    DO: 002-foo.yml and 010-bar.yml
    DON't 2-foo.yml and 10-bar.yml

    However, numeric aware sort is trivial to do if needed
- in conf.d, if there is an invalid file, the error is printed to stdout, and ignored

- configuration lists are only merged at the top level. It should be easy to
    make a recursive alg which concatenates even nested lists.

### Improvments for the conf:

There are two different types of settings:

1. some for the framework
2. some for the tests

They should be clearly separated (espeically for general names like `urls`.
This could apply for both the framework and the checks).

Also, it seems safer to have settings *per check*, and avoid global settings
(where some check can access settings he doesn't care about).

So, here is what I propose (i could update the python implementation, to make
sure both tools are still compatible)

```yaml
framework_settings:
    probe_id: ...
    probe_group: ...

    graylog_udp_gelf_servers:
        ...
    
    graylog_http_gelf_servers:
        ...

    teams_channel:
        ...

    checks:
        ...

checks_settings:
    _globals:
        ... # accessible by all checks

    checka:
        ... # only accessible by checka
    checkb:
        ... # only accessible by checkb
```

Note that there is nothing preventing the framework from exposing some of his
settings to the tests, should it be needed (if you think of case, let me
know).

## Runner

## Reporter