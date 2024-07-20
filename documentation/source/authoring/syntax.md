---
title: Template Syntax
tags: [theming, authoring]
---

# Introduction

Templates follow the Go `text/template` syntax, with some added functions.

# Example

```markdown
{{`
    {{ if .Post }}
        This post is titled {{ .Post.Title }}, the relative link to it is {{ .Post.Href }}.
        It was created at {{ ftime .Post.CreatedAt }} and updated at {{ ftime .Post.UpdatedAt }}.

        This website is titled {{ .Site.Title }}. It uses a Go time format of {{ .Site.TimeFormat }}.

        To learn about creating a post, see {{ plink "authoring/create-first-post" }}.
    {{ else }}
        There is no post.
    {{end}}
`}}
```

Will output the following:

```text
    {{ if .Post }}
        This post is titled {{ .Post.Title }}, the relative link to it is {{ .Post.Href }}.
        It was created at {{ ftime .Post.CreatedAt }} and updated at {{ ftime .Post.UpdatedAt }}.

        This website is titled {{ .Site.Title }}. It uses a Go time format of {{ .Site.TimeFormat }}.

        To learn about creating a post, see {{ plink "authoring/create-first-post" }}.
    {{ else }}
        There is no post.
    {{end}}
```