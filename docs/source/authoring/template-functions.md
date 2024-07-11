---
title: Template Functions
tags: [authoring, theming]
---

These functions can be accessed via templates using the [Go template syntax](https://github.com/gofiber/template/blob/html/v2.1.0/html/TEMPLATES_CHEATSHEET.md#template-functions).

# Theme and source functions

## `add(int1, int2) -> int3`

Adds `int1` and `int2` together.

Example: `The sum is {{ add 10 1 }}`

Output: `The sum is 11`.

## `sub(int1, int2) -> int3`

Subtracts `int2` from `int1`.

## `ftime(time, <optional> format) string`

Takes a `time` object and optional format string. If no format string is specified it uses the `time_format` from `config.yaml`.

Example: `This post was created at {{ `{{ ftime .CreatedAt }}` }}`.

Output: `This post was created at 2024-07-10 12:30:00`.

Specifying a format string changes the output.

Example: `This post was created at {{ `{{ ftime .CreatedAt "06/02/01" }}` }}.`

Output: `This post was created at 24/07/10`.

# Source functions

These functions only work for posts in `source/`.

## `plink(path/to/post) string`

Returns the HTML anchor link to a given post.

Example: `{{`Link to a post: {{ plink "authoring/template-functions" }}.`}}`

Output: `Link to a post: {{ plink "authoring/template-functions" }}.`
