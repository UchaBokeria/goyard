# Simple Example

A basic example demonstrating the integration of Tailwind CSS with DaisyUI in a Go project using templ templates.

## Features

- Tailwind CSS for styling
- DaisyUI components for enhanced UI
- templ templates for server-side rendering
- Hot reload during development

## Setup

1. Install dependencies:
   ```bash
   npm install
   ```

2. Run the development server:
   ```bash
   go run .
   ```

3. Start Tailwind CSS watcher:
   ```bash
   npx tailwindcss -i ./input.css -o ./public/output.css --watch
   ```

## Project Structure

- `views/` - templ template files
- `public/` - static assets and compiled CSS
- `tailwind.config.js` - Tailwind configuration with DaisyUI plugin

## Technologies

- Go
- Templ
- Htmx
- Alpine
- Tailwind CSS
- DaisyUI
