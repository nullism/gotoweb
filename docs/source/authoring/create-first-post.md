---
title: Creating Your First Post
tags: [authoring]
created_at: 2024-01-02 16:00:00
updated_at: 2024-07-08 12:00:00
---

# Creating a post

A post consists of a single markdown file that is rendered against the theme's `post.html` template.

## Create hello.md

In the editor of your choice, create a file at `source/hello.md` with the following contents:

```markdown
---
title: Hello!
---

This is my *first post*!
```

### Post config

Everything between `---` and `---` at the start of the post is part of the post configuration.

The syntax of the configuration is in `yaml`.

_Note: to be considered configuration, the very first line of the file must start with `---`._

#### Post config properties

* `title` - Sets the post title.
* `blurb` - An optional summary for previews (such as the search page).
* `tags` - A list of tags for this post (typically shown on search page).
  * Example: `tags: [my-first-post, hello]`
* `skip_publish` - If this is set to true, the post will not be published or added to the search index.
* `skip_index` - If true the post will not be searchable.

## Run build

Run `gotoweb build` from the directory containing `config.yaml`.

_Note: You may also run this command in child directories_.

```text
you@host /path/to/my-site/$ gotoweb build
...
[06:32:47] INFO built!
  took (s):     0.014175595
```