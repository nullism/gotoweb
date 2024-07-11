---
title: Creating a Theme
tags: [theming]
---

# Creating a theme

> [!NOTE] This is for advanced users.

A theme is composed of Gotoweb HTML templates and static files.

## Folder structure

```
  /my-theme/
    dist/
    tpl/
    post.html
    post-list.html
    post-preview.html
    search.html
    theme.yaml
```

### dist/

The `dist/` folder contains static assets use by the theme. Examples include CSS, Javascript, and images.

This folder will be copied to the public directory of the site.

### tpl/

The `tpl/` folder contains reusable templates.

Templates in here have full access to the context, with the additional `.Args` which is a `map[string]any`.

#### Example template

Let's create `tpl/image.html` that shows an image and caption.

```html
<!-- tpl/image.html -->
<div>
    <img src="{{ `{{.Args.Src}}` }}" />
    <p>{{ `{{.Args.Caption}}` }}</p>
</div>
```

To use the template, we'd do:

```
{{ `{{ tpl "image.html" "Src" "/path/to/my-image.png" "Caption" "My caption" }}` }}
```

