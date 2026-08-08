# gomakeagame.com

The website for [Go Make A Game](https://gomakeagame.com), built with Go, [templ](https://templ.guide/), [htmx](https://htmx.org/), and [missing.css](https://missing.style/).

## Setup

```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Install dependencies
go mod tidy

# Generate templates and run
templ generate && go run .
```

The server listens on `PORT` (default `8080`).

## Dev with hot reload

```bash
# Install air
go install github.com/air-verse/air@latest

# From this directory
air
```
