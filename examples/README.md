# Examples

This folder contains example configuration files in four
supported formats. You can use any format you want for
the config file, just give it the appropriate suffix for
the format. A JSON file should end with `.json`, and
YAML with `.yaml`. The default format is always TOML and
may have any _other_ suffix.

# Kubernetes

There are two files for Kubernetes deployment examples.
Feel free to use them as you see fit.
Please make sure to the delete all comments before
deploying and make sure to fill in with correct values.

# Notes

When adding new content to this folder, **DO NOT MAKE NEW FOLDERS**,
it will break `make install` on macOS (used for homebrew).
