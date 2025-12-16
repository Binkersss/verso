## [0.0.1]

Literally everything lol.

## [v0.0.2] 

### Fixed
- **Site title configuration**: Fixed bug where the `SiteTitle` config option was not being passed to the Builder, causing the `-title` flag to be ignored. Page titles now correctly display as "Page Title - Site Title" when a site title is specified.

### Changed
- Reorganized project structure to include `examples/` directory with sample implementations for testing and reference
- Added syntax highlighting and copy-paste button generation in `pkg/client/generator.go`. Updated `examples/` (lines 189-246 of `style.css` and lines 10-11 of `index.html`) to match and integrate with that feature. 
