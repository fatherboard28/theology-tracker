#!/bin/bash

# 1. Download the four JS libraries into static/js/
#    (do this once — they're vendored so the app works offline)
curl -o static/js/htmx.min.js https://unpkg.com/htmx.org@1.9.12/dist/htmx.min.js
curl -o static/js/alpine.min.js https://unpkg.com/alpinejs@3.14.1/dist/cdn.min.js
curl -o static/js/sortable.min.js https://unpkg.com/sortablejs@1.15.2/Sortable.min.js
curl -o static/js/marked.min.js https://unpkg.com/marked@12.0.0/marked.min.js

# 2. Pull dependencies and generate go.sum
go mod tidy

# 3. Generate the Templ Go files
templ generate

# 4. Run with hot reload
air
