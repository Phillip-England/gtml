# gtml

An opinionated static site generator that compiles HTML components into static websites.

## What is gtml?

gtml is a templating system and compiler that lets you create reusable HTML components and compile them into static HTML files. Write your UI once as components, use them throughout your site, and compile everything into plain HTML with zero runtime overhead.

## Installation

```bash
go install github.com/phillip-england/gtml@latest
```

Verify the installation:

```bash
gtml --help
```

## Quick Start

### 1. Initialize a new project

```bash
gtml init mysite
```

This creates a new gtml project with the following structure:

```
mysite/
  components/    # Your reusable HTML components
  routes/        # Your page routes
  dist/          # Compiled output
  static/        # Static assets (CSS, JS, images)
```

### 2. Create a component

`components/Greeting.html`
```html
<div props='name string'>
  <h1>Hello, {name}!</h1>
</div>
```

### 3. Use it in a route

`routes/index.html`
```html
<html>
  <body>
    <Greeting name="World" />
  </body>
</html>
```

### 4. Compile

```bash
gtml compile mysite
```

### 5. Output

`dist/index.html`
```html
<html>
  <body>
    <div data-greeting="">
      <h1>Hello, World!</h1>
    </div>
  </body>
</html>
```

## Features

### Props

Define typed props on your components:

```html
<div props='name string, age int, isAdmin boolean'>
  <p>Name: {name}</p>
  <p>Age: {age}</p>
</div>
```

Supported types: `string`, `int`, `boolean`

### Expressions

Use expressions with `{}` for dynamic values:

```html
<Calculator result={10 + 5 * 2} />
<Greeting name={"Bo" + "b"} />
```

Supports arithmetic operations (`+`, `-`, `*`, `/`, `%`) with PEMDAS order.

### Slots

Insert content into specific places within components:

**Define a slot:**
```html
<html props='title string'>
  <head><title>{title}</title></head>
  <body>
    {{ slot: content }}
  </body>
</html>
```

**Fill a slot:**
```html
<Layout title="My Page">
  <slot name='content' tag='main'>
    <h1>Welcome!</h1>
  </slot>
</Layout>
```

### Conditional Rendering

Use ternary operators for conditionals:

```html
<div props='isAdmin boolean'>
  { isAdmin ? (
    <p>Welcome, Admin!</p>
  ) : (
    <p>Welcome, User!</p>
  ) }
</div>
```

Comparison operators: `==`, `!=`, `<`, `>`, `<=`, `>=`
Logical operators: `&&`, `||`

### Prop Drilling

Pass props down through nested components:

```html
<div props='title string'>
  <Inner text={title} />
</div>
```

### Component Styles

Scoped styles within components:

```html
<style>
  h1 { color: red; }
</style>

<div props='heading string'>
  <h1>{heading}</h1>
</div>
```

## CLI Commands

### `gtml init <path> [--force]`

Initialize a new gtml project with starter files and preinstalled components.

- `--force`: Overwrite existing directory

### `gtml compile <path> [--watch]`

Compile all routes to static HTML in the `dist` directory.

- `--watch`: Watch for changes and recompile automatically

## Preinstalled Components

gtml comes with 80+ ready-to-use Tailwind-styled components:

| Category | Components |
|----------|------------|
| **Buttons** | ButtonPrimary, ButtonSecondary, ButtonOutline, ButtonDanger, ButtonSuccess, ButtonSm, ButtonLg |
| **Forms** | InputText, InputEmail, InputPassword, Textarea, Select, Checkbox, RadioGroup, FormLayout, LoginForm |
| **Cards** | CardBasic, CardWithImage, CardClickable, ProfileCard, PricingCard |
| **Alerts** | AlertSuccess, AlertError, AlertWarning, AlertInfo, Notification |
| **Badges** | BadgePrimary, BadgeSuccess, BadgeWarning, BadgeDanger, SkillBadge |
| **Navigation** | Navbar, NavLink, Breadcrumb, Pagination, Tabs, TabItem |
| **Layout** | HeroSection, Footer, Sidebar, Modal, Grid2Columns, Grid3Columns, Section |
| **Data Display** | Avatar, ProgressBar, Accordion, ListGroup, StatisticCard, FeatureCard, TestimonialCard |
| **Tables** | TableBasic, TableHeader, TableRow, TableCell |
| **Utility** | Divider, LoadingSpinner, LoadingSkeleton, StepIndicator |

## Project Structure

```
myproject/
  components/          # Reusable components
    Button.html
    Card.html
    Layout.html
  routes/              # Page routes (become HTML files)
    index.html         # -> dist/index.html
    about.html         # -> dist/about.html
    blog/
      post.html        # -> dist/blog/post.html
  static/              # Static assets (copied to dist/static)
    styles.css
    script.js
  dist/                # Compiled output
```

## Example Component

```html
<div props='title string, description string, imageUrl string'>
  <div class="bg-white rounded-lg shadow-md overflow-hidden">
    <img src={imageUrl} alt={title} class="w-full h-48 object-cover" />
    <div class="p-4">
      <h2 class="text-xl font-bold">{title}</h2>
      <p class="text-gray-600 mt-2">{description}</p>
    </div>
  </div>
</div>
```

## Links

- [Documentation](https://gtml.dev/docs)
- [Component Gallery](https://gtml.dev/docs/components)

## License

MIT
