# Sitemap Builder

This project builds an XML sitemap for any given URL up to a specified depth using BFS/DFS traversal.

## Features

- 🚀 **Fast crawling** with configurable depth limits
- 🔄 **Two traversal methods**: Breadth-First Search (BFS) and Depth-First Search (DFS)
- 🔗 **Smart URL handling**: Automatically resolves relative and absolute URLs
- 🛡️ **Cycle detection**: Prevents infinite loops and duplicate URL processing
- 📄 **Standard XML output**: Generates sitemap.xml following official sitemap protocol
- 🎯 **Domain filtering**: Only crawls links within the same domain
- ⚡ **Generic queue implementation**: Type-safe queue using Go generics

## Installation

### Prerequisites
- Go 1.18 or higher (for generics support)

### Clone and Build
```bash
git clone https://github.com/yourusername/siteMapBuilder.git
cd siteMapBuilder
go mod download
go build
```

## Usage

### Basic Usage
```bash
go run main.go -url="https://example.com" -depth=3
```

### Command-Line Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-url` | `https://gobyexample.com` | The URL to start building the sitemap from |
| `-depth` | `10` | Maximum depth of links to traverse |
| `-traversal` | `bfs` | Traversal method: `bfs` (breadth-first) or `dfs` (depth-first) |

### Examples

**Crawl with BFS (recommended for most cases):**
```bash
go run main.go -url="https://gobyexample.com" -depth=5 -traversal=bfs
```

**Crawl with DFS:**
```bash
go run main.go -url="https://example.com" -depth=3 -traversal=dfs
```

**Quick shallow crawl:**
```bash
go run main.go -url="https://example.com" -depth=1
```

**⚠️ Sites that will NOT work:**
```bash
# These will fail or hang - DO NOT USE
go run main.go -url="https://golang.org" -depth=2      # Parser error: exceeds 512 nodes
go run main.go -url="https://github.com" -depth=2      # Too complex
go run main.go -url="https://wikipedia.org" -depth=2   # Too large
```

## Output
The program generates a `sitemap.xml` file in the current directory with the following format:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
  <url>
    <loc>https://example.com</loc>
  </url>
  <url>
    <loc>https://example.com/about</loc>
  </url>
  <!-- ... more URLs ... -->
</urlset>
```

### Key Components

- **Generic Queue**: Type-safe queue implementation using Go 1.18+ generics
- **URL Resolution**: Proper handling of relative, absolute, and protocol-relative URLs
- **Cycle Prevention**: Visited map ensures each URL is processed only once
- **XML Generation**: Standards-compliant sitemap output

## Project Structure

```
siteMapBuilder/
├── main.go           # Main program with traversal logic
├── linker/           # HTML parsing package
├── queue/            # Generic queue implementation
├── go.mod            # Go module definition
└── README.md         # This file
```

## Limitations
⚠️ **Important Performance & Compatibility Notes:**

- **Not suitable for large/complex websites**: Sites like `golang.org`, Wikipedia, or major commercial sites may cause parsing errors or hang indefinitely
- **HTML parser limitations**: The underlying HTML parser has a limit of 512 nested nodes. Sites with deeply nested HTML structures will fail with: `html: open stack of elements exceeds 512 nodes`
- **No concurrent requests**: Single-threaded crawling makes it slow for sites with many pages
- **Memory usage**: All URLs are kept in memory, which can be problematic for very large sites
- Only crawls same-domain links (no external domains)
- Does not handle JavaScript-rendered content

**Recommended Use Cases:**
- ✅ Small to medium websites (< 500 pages)
- ✅ Simple HTML structure sites
- ✅ Personal blogs, documentation sites, small business websites
- ❌ Large enterprise websites
- ❌ Complex modern web applications
- ❌ Sites with heavy JavaScript rendering

## Future Enhancements

- [ ] Concurrent crawling with goroutines
- [ ] Support for sitemap index files (for large sites)
- [ ] Retry logic for failed requests

## Acknowledgments
- Built as part of the [Gophercises](https://gophercises.com/) coding exercises
- Uses the [Sitemaps Protocol](https://www.sitemaps.org/) standard

---

**Happy crawling! 🕷️**
