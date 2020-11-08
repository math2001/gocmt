# Assumptions made during development

- in conf.d file, we don't load yml files in subdirectories.
    Loading from subdirectory is easy to do if needed
- in conf.d, the yml files are loaded in alphabetical order, and NOT numerically aware.
    (2-foo.yml comes AFTER 10-bar.yml).
    DO: 002-foo.yml and 010-bar.yml
    DON't 2-foo.yml and 10-bar.yml

    However, numeric aware sort is trivial to do if needed
- in conf.d, if there is an invalid file, the error is printed to stdout, and ignored