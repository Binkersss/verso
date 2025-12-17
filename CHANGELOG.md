## [v0.0.2] 

### Fixed
- **Site title configuration**: Fixed bug where the `SiteTitle` config option was not being passed to the Builder, causing the `-title` flag to be ignored. Page titles now correctly display as "Page Title - Site Title" when a site title is specified.

### Added
- Added syntax highlighting and copy-paste button generation in `pkg/client/generator.go`. Updated `examples/` (lines 189-246 of `style.css` and lines 10-11 of `index.html`) to match and integrate with that feature. 
- Page metadata display with title, date, and author(s) automatically rendered at the top of each page
- Support for both single `author` and multiple `authors` fields in frontmatter

### Changed
- Reorganized project structure to include `examples/` directory with sample implementations for testing and reference
- **BREAKING**: Page titles are now automatically rendered from frontmatter. Remove duplicate `# Title` headings from your markdown content to avoid double titles.

### Migration Guide
Update your markdown files:

**Before:**
```markdown
---
title: My Post
---

# My Post

Content here...
```

**After:**
```markdown
---
title: My Post
date: December 16th, 2025
author: Your Name
---

Content here...
```

## [0.0.1]

Literally everything lol.


