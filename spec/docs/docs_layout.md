# Documentation Site Layout and Structure

## Site Pages

The documentation site consists of the following pages:

### Home Page (`/`)
Route: `routes/index.html`

The home page introduces gtml and encourages users to get started. It contains:
- A hero section with the gtml branding and tagline
- Call-to-action buttons linking to the documentation and GitHub
- A "Hello World" example showing component definition and usage
- A features section highlighting why to use gtml (Simple & Fast, Component-Based, Zero Runtime)
- A footer with links to GitHub

### Documentation Page (`/docs`)
Route: `routes/docs.html`

The main documentation page provides comprehensive information about using gtml. It features:
- A sticky sidebar navigation for quick access to different sections
- Sections covering: Installation, Project Structure, Hello World, Components, Props, Expressions, Slots, CLI commands, and Preinstalled Components
- Code examples with syntax highlighting using dark-themed code blocks

The sidebar contains links organized by category:
- Getting Started (Installation, Project Structure, Hello World)
- Components (Basic Components, Props, Prop Types, Expressions, Slots, Prop Drilling, Conditional Rendering)
- CLI (gtml init, gtml compile)
- Extras (Preinstalled Components)

### Components Gallery (`/docs/components`)
Route: `routes/docs/components.html`

A gallery page that lists all documented preinstalled components organized by category. Each component links to its individual showcase page.

### Individual Component Pages (`/components/{category}/{component-name}`)
Routes: `routes/components/{category}/{component-name}.html`

Each documented component has its own page showing:
- The gtml code needed to use the component
- A live example of the component rendered
- A props table describing available props, their types, and descriptions
- A "Back to Documentation" link
- The DocsSidebar for navigating between components

## Reusable Documentation Components

The documentation site uses these custom components located in `docs/components/`:

### DocsLayout
The base HTML layout wrapper used by all pages. Accepts a `title` prop for the page title.
- Sets up the HTML document structure
- Includes Tailwind CSS via CDN
- Provides a content slot for page content

### DocsNavbar
The main navigation bar displayed at the top of every page. Accepts a `currentPage` prop to highlight the active section.
- Sticky positioning (stays visible while scrolling)
- Links: Home, Documentation, Components, GitHub
- Active state styling based on `currentPage` value ('home', 'docs', or 'components')

### DocsSidebar
The sidebar navigation used on individual component showcase pages. Accepts a `currentPage` prop to highlight the active component.
- Fixed sidebar on large screens (hidden on mobile)
- Lists all documented components organized by category
- Active state styling for the current component

## Styling and Theming

### Tailwind CSS
The documentation site uses Tailwind CSS loaded via CDN for all styling. No custom CSS files are used.

### Color Scheme
- Background: White (`bg-white`)
- Text: Gray-900 for primary text, Gray-500/600 for secondary text
- Primary accent: Blue-600 for links, active states, and CTAs
- Code blocks: Gray-900 background with gray-100 text
- Borders: Gray-200 for subtle dividers
- Sidebar: Gray-50 background

### Typography
- Headings use `font-bold` or `font-extrabold` with appropriate size classes
- Body text uses default Tailwind text sizing
- Code uses `font-mono` for monospace display

### Layout Patterns
- Max width container: `max-w-7xl mx-auto`
- Responsive padding: `px-4 sm:px-6 lg:px-8`
- Sticky navigation: `sticky top-0 z-50`
- Sidebar: `w-64` fixed width, `sticky top-16` positioned below navbar
