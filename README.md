# Verso

An awesome static site generator and server with built-in SPA routing. Write markdown, get a fast single-page site.

## Features

- **Markdown with frontmatter** - YAML metadata support
- **Built-in SPA routing** - No JavaScript needed from you
- **Live reloading** - Dev Server auto-rebuilds on file changes
- **No Code, No Config** - Sensible defaults, works out of the box
- **Go-powered** - Fast builds, simple deployment

## Quick Start

### Installation

```bash
go get github.com/Binkersss/verso
```

### Project Structure

```
my-site/
├── content/
│   ├── about.md
│   ├── projects.md
│   └── blog/
│       └── post-1.md
├── templates/
│   └── index.html
├── static/
│   ├── style.css
│   └── favicon.ico
└── main.go
```

### Create `main.go`

```go
package main

import (
    "log"
    "github.com/Binkersss/verso"
)

func main() {
    site := verso.New(verso.Config{
        ContentDir:   "content",
        TemplateDir:  "templates",
        StaticDir:    "static",
        OutputDir:    "dist",
    })
    
    if err := site.Serve(":3000"); err != nil {
        log.Fatal(err)
    }
}
```

### Create Your Template

**templates/index.html:**
```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Site</title>
    <link rel="stylesheet" href="/style.css">
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <div class="sidebar-header">MY SITE</div>
            <nav id="nav"></nav>
        </aside>
        <main id="content"></main>
    </div>
    <script src="/app.js"></script>
</body>
</html>
```

### Write Content

**content/about.md:**
```markdown
---
title: About Me
date: 2024-12-13
---

# About Me

I build things with Go.
```

### Run

```bash
go run main.go
```

Visit `http://localhost:3000`

## How It Works

1. **Markdown parsing** - All `.md` files in `content/` are converted to HTML
2. **Frontmatter** - YAML metadata at the top of markdown files
3. **Routing** - File paths become routes (`content/blog/post.md` → `/blog/post`)
4. **Navigation** - Top-level pages auto-populate the nav (nested pages are hidden)
5. **SPA magic** - Client-side routing generated automatically

## Frontmatter

Add metadata to any markdown file:

```markdown
---
title: My Post
date: 2024-12-13
author: Your Name
tags: [go, web]
---

# Content starts here
```

The `title` field is used for:
- Document title (browser tab)
- Navigation text (if not nested)

## Navigation

Only **top-level pages** appear in the nav sidebar:

```
content/
├── about.md          ✓ Shows in nav
├── projects.md       ✓ Shows in nav  
└── blog/
    ├── post-1.md     ✗ Hidden from nav
    └── post-2.md     ✗ Hidden from nav
```

Link to nested pages from your content:

```markdown
# Blog

Recent posts:
- [Building Verso](/blog/post-1)
- [Go Tips](/blog/post-2)
```

## Building for Production

```go
// Build only (no server)
if err := site.Build(); err != nil {
    log.Fatal(err)
}
```

This creates a `dist/` directory with:
- `index.html` - Your template
- `content.json` - Rendered pages
- `app.js` - Generated SPA router
- `style.css` - Your styles
- Other static assets

Deploy the `dist/` folder to any static host, or serve it with the built-in server.

## Deployment

### Static Hosting
Deploy `dist/` to:
- Cloudflare Pages
- Netlify
- Vercel
- GitHub Pages

### Self-Hosted
Run the Go server:
- Fly.io
- Railway
- Your own server + Cloudflare Tunnel
- Coolify

## Styling

Verso doesn't impose any styling. The CSS selector requirements:

- `#nav` - Container for navigation items
- `#content` - Container for page content
- `.nav-item` - Navigation items (auto-generated)
- `.nav-item.active` - Current page
- `.fade` - Optional transition class for content area

## Configuration

```go
verso.Config{
    ContentDir:   "content",   // Where markdown lives
    TemplateDir:  "templates",  // Where index.html lives
    StaticDir:    "static",     // CSS, images, etc.
    OutputDir:    "dist",       // Build output
}
```
