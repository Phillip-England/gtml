# gtml

An opinionated static site generator with a built-in component system that compiles HTML components into static websites. gtml features reactive props, client-side fetch with suspense/fallback, scoped CSS, and 80+ preinstalled Tailwind-styled components.

## What is gtml?

gtml is three things:

1. **A Static Site Generator** - Compiles routes and components into static HTML files
2. **A Templating System** - Create reusable HTML components with props, slots, and expressions
3. **A Compiler** - Transforms component syntax into plain HTML with zero runtime overhead

## Quick Start

### Installation

```bash
go install github.com/phillip-england/gtml@latest
```

Verify installation:

```bash
gtml --help
```

### Initialize a Project

```bash
gtml init mysite
```

This creates:

```
mysite/
  components/    # Reusable HTML components
  routes/        # Page routes (become HTML files)
  static/        # Static assets (CSS, JS, images)
  dist/          # Compiled output
```

### Create a Component

`components/Greeting.html`:
```html
<div props='name string'>
  <h1>Hello, {name}!</h1>
</div>
```

### Use in a Route

`routes/index.html`:
```html
<html>
  <body>
    <Greeting name="World" />
  </body>
</html>
```

### Compile

```bash
gtml compile mysite
```

### Output

`dist/index.html`:
```html
<html>
  <body>
    <div data-greeting="">
      <h1>Hello, World!</h1>
    </div>
  </body>
</html>
```

## Project Structure

```
myproject/
  components/          # Reusable components (PascalCase names)
    Button.html
    Layout.html
    ui/                # Subdirectories allowed
      Card.html
  routes/              # Page routes (kebab-case names only)
    index.html         # -> dist/index.html
    about.html         # -> dist/about.html
    blog/
      post.html        # -> dist/blog/post.html
  static/              # Static assets (copied to dist/static)
    styles.css
    script.js
  dist/                # Compiled output (generated)
```

### Components Directory Rules

- **Naming**: PascalCase only (e.g., `MyComponent.html`, NOT `myComponent.html`)
- **Uniqueness**: Each component name must be unique across all subdirectories
- **One root element**: Each file must contain a single root element
- **Subdirectories**: Allowed (e.g., `components/ui/Button.html`)

### Routes Directory Rules

- **Naming**: kebab-case only (e.g., `my-route.html`, NOT `MyRoute.html`)
- **Subdirectories**: Allowed (e.g., `routes/blog/post.html`)
- **Naming collisions**: Allowed in different subdirectories

## Component System

### Basic Component

```html
<style>
  h1 { color: red; }
</style>

<div props='heading string, subheading string'>
  <h1>{heading}</h1>
  <p>{subheading}</p>
</div>
```

### Props

Define typed props on components:

```html
<div props='name string, age int, isAdmin boolean'>
  <p>Name: {name}</p>
  <p>Age: {age}</p>
  { isAdmin ? (
    <p>Admin Access</p>
  ) : (
    <p>User Access</p>
  ) }
</div>
```

### Prop Types

| Type | Declaration | Usage | Notes |
|------|------------|-------|-------|
| `string` | `name string` | `name='Bob'` or `name={"Bob"}` | Raw strings allowed |
| `int` | `count int` | `count={5}` or `count={2+3}` | Must use `{}` |
| `boolean` | `enabled boolean` | `enabled={true}` or `enabled={2>1}` | Must use `{}` |

### Prop Drilling

Props pass down through nested components:

```html
<div props='title string'>
  <ChildComponent text={title} />
</div>
```

### Slots

Insert content into specific places within components:

**Define a slot:**
```html
<html props='title string'>
  <head><title>{title}</title></head>
  <body>
    <slot name='content' />
  </body>
</html>
```

**Fill a slot:**
```html
<Layout title="My Page">
  <slot name='content' tag='main' class='container'>
    <h1>Welcome!</h1>
  </slot>
</Layout>
```

The `tag` attribute wraps content in the specified HTML tag. Any other attributes are preserved.

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

#### Comparison Operators

- `==` (equals)
- `!=` (not equals)
- `<` (less than)
- `>` (greater than)
- `<=` (less than or equal)
- `>=` (greater than or equal)

#### Logical Operators

- `&&` (and) - both conditions must be true
- `||` (or) - at least one condition must be true

#### Nested Ternaries

```html
<div props='age int'>
  { age < 21 ? (
    <p>You cannot drink</p>
  ) : (
    { age > 21 && age < 25 ? (
      <p>You can drink but shouldn't</p>
    ) : (
      <p>You can drink</p>
    ) }
  ) }
</div>
```

## Expressions

Use `{}` for dynamic values:

```html
<Calculator result={10 + 5 * 2} />
<Greeting name={"Bo" + "b"} />
<User age={currentAge + 1} />
```

### Supported Operations

| Operator | Types | Example | Result |
|----------|-------|---------|--------|
| `+` | string + string | `"A" + "B"` | `"AB"` |
| `+` | int + int | `2 + 3` | `5` |
| `-` | int - int | `5 - 2` | `3` |
| `*` | int * int | `2 * 3` | `6` |
| `/` | int / int | `6 / 2` | `3` |
| `%` | int % int | `5 % 2` | `1` |

### Evaluation Order

Expressions follow PEMDAS order (parentheses, exponents, multiplication/division/modulo, addition/subtraction).

### Raw Strings vs Expressions

```html
<!-- String - raw value preferred -->
<Component name='Bob' />

<!-- String - expression also works -->
<Component name={"Bob"} />

<!-- Non-string - must use expression -->
<Component count={42} />
<Component enabled={true} />
```

## Styling

### Component Styles

Add `<style>` blocks at the top of components:

```html
<style>
  h1 { color: red; }
  .card { padding: 1rem; }
</style>

<div props='heading string'>
  <h1>{heading}</h1>
  <div class="card">Content</div>
</div>
```

### Scoped CSS

Styles are automatically scoped to prevent conflicts:

- Each component gets a unique `data-gtml-scope` attribute
- CSS selectors are prefixed to target only the component
- All component styles are aggregated into `dist/static/styles.css`

### Static Assets

Place global assets in the `static/` directory:

```
static/
  styles.css    # Global styles
  script.js     # Global JavaScript
  images/       # Images
```

These are copied to `dist/static/` during compilation.

## Interactivity

### Signals

Props can be treated as reactive signals for client-side interactivity:

```html
<div props='name string'>
  <p>My name is {name}!</p>
  <button id='btn'>Change Name</button>
</div>

<script type='gtml'>
  #btn.onclick(() => {
    name = "John"
  })
</script>
```

### Element Selection

| Selector | Description |
|----------|-------------|
| `#id` | Select element by ID |
| `.class` | Select element by class |
| `#id*` | Select all matching elements |

### Event Handlers

```html
<button id='submit'>Submit</button>
<div id='output'></div>

<script type='gtml'>
  #submit.onclick(() => {
    #output.innerHTML = "Clicked!"
  })
</script>
```

### Inline Events

```html
<div props='count int'>
  <p>Count: {count}</p>
  <button onclick={() => { count = count + 1 }}>
    Increment
  </button>
</div>
```

### Signal Library

gtml includes a built-in signal library for reactivity. When signals are modified, the DOM updates automatically.

## Fetch System

### Basic Fetch

Fetch data from APIs with automatic client-side rendering:

```html
<div fetch='GET /api/users' as='users'>
  <ul>
    <li for='user in users'>{user.name}</li>
  </ul>
</div>
```

### Fetch Syntax

```
fetch='METHOD URL' as='variableName'
```

- `METHOD`: HTTP method (GET, POST, PUT, DELETE, etc.)
- `URL`: Endpoint URL
- `as`: Variable name for the response data

### Suspense

Show loading state while fetching:

```html
<div fetch='GET /api/users' as='users'>
  <div suspense>
    <p>Loading users...</p>
  </div>
  <ul>
    <li for='user in users'>{user.name}</li>
  </ul>
</div>
```

### Fallback

Show error state if fetch fails:

```html
<div fetch='GET /api/users' as='users'>
  <div fallback>
    <p>Failed to load users</p>
  </div>
  <ul>
    <li for='user in users'>{user.name}</li>
  </ul>
</div>
```

### For Loops

Iterate over arrays with the `for` attribute:

```html
<li for='user in users'>{user.name}
  <ul>
    <li for='color in user.colors'>{color.name}</li>
  </ul>
</li>
```

Syntax: `for='item in array'` or `for='item in parent.nested'`

### Prop Values in Fetch URLs

Fetch URLs can use prop expressions:

```html
<UserList apiUrl='localhost:8080/api/users'>
  <div fetch='GET {apiUrl}' as='users'>
    <ul><li for='user in users'>{user.name}</li></ul>
  </div>
</UserList>
```

## CLI Commands

### `gtml init <PATH> [--force]`

Initialize a new gtml project with starter files and preinstalled components.

- `--force`: Overwrite existing directory

### `gtml compile <PATH> [--watch]`

Compile all routes to static HTML in the `dist` directory.

- `--watch`: Watch for changes and recompile automatically

### `gtml test [PATH]`

Run component tests. Tests are `-test.html` files that compile successfully.

## Compilation Process

1. **Read components directory**: Load all components into a registry
2. **Process routes recursively**: Find and compile nested components
3. **Generate static HTML**: Output to `dist/` directory
4. **Copy static assets**: Copy `static/` to `dist/static/`
5. **Aggregate styles**: Combine component styles into `dist/static/styles.css`
6. **Inject interactivity**: Add signal library and event handlers

### Watch Mode

With `--watch`, gtml monitors all files in the project directory and recompiles on any change.

## Preinstalled Components

gtml comes with 80+ ready-to-use Tailwind-styled components.

### Buttons

| Component | Description |
|-----------|-------------|
| `ButtonPrimary` | Standard blue primary button |
| `ButtonSecondary` | Gray secondary button |
| `ButtonOutline` | Border-only outline button |
| `ButtonDanger` | Red destructive action button |
| `ButtonSuccess` | Green positive action button |
| `ButtonSm` | Small sized button |
| `ButtonLg` | Large sized button |

### Forms

| Component | Description |
|-----------|-------------|
| `InputText` | Standard text input with label |
| `InputEmail` | Email input field |
| `InputPassword` | Password input field |
| `Textarea` | Multi-line text area |
| `Select` | Dropdown select component |
| `Checkbox` | Checkbox with label |
| `RadioGroup` | Radio button group |
| `FormLayout` | Form wrapper with title/submit |
| `FormField` | Generic form field wrapper |
| `LoginForm` | Pre-built login form |

### Cards

| Component | Description |
|-----------|-------------|
| `CardBasic` | Simple card with title/content |
| `CardWithImage` | Card with image header |
| `CardClickable` | Hoverable clickable card |
| `ProfileCard` | User profile display |
| `PricingCard` | Pricing tier display |

### Alerts

| Component | Description |
|-----------|-------------|
| `AlertSuccess` | Green success alert |
| `AlertError` | Red error alert |
| `AlertWarning` | Yellow warning alert |
| `AlertInfo` | Blue information alert |
| `Notification` | Multi-type notification |

### Badges

| Component | Description |
|-----------|-------------|
| `BadgePrimary` | Blue status badge |
| `BadgeSuccess` | Green status badge |
| `BadgeWarning` | Yellow status badge |
| `BadgeDanger` | Red status badge |
| `SkillBadge` | Skill level indicator |

### Navigation

| Component | Description |
|-----------|-------------|
| `Navbar` | Top nav bar with brand/links |
| `NavLink` | Individual nav link |
| `Breadcrumb` | Breadcrumb container |
| `BreadcrumbItem` | Individual breadcrumb |
| `Pagination` | Page navigation controls |
| `Tabs` | Tab navigation container |
| `TabItem` | Individual tab |

### Layout

| Component | Description |
|-----------|-------------|
| `HeroSection` | Hero/landing with CTA |
| `Footer` | Site footer |
| `FooterLink` | Footer link item |
| `Sidebar` | Sidebar navigation |
| `SidebarItem` | Sidebar nav item |
| `Modal` | Modal dialog |
| `SidebarLayout` | Layout with sidebar |
| `DashboardLayout` | Dashboard with header |
| `TwoColumnLayout` | Two column flex |
| `ThreeColumnLayout` | Three column flex |
| `Grid2Columns` | Two column grid |
| `Grid3Columns` | Three column grid |
| `PageHeader` | Page title with breadcrumbs |
| `Section` | Generic section wrapper |

### Data Display

| Component | Description |
|-----------|-------------|
| `Avatar` | User avatar image/initials |
| `ProgressBar` | Progress indicator |
| `Accordion` | Collapsible accordion |
| `AccordionItem` | Individual accordion item |
| `ListGroup` | List group container |
| `ListItem` | List group item |
| `StatisticCard` | Statistics display |
| `StatsGrid` | Stats grid container |
| `StatItem` | Individual statistic |
| `FeatureCard` | Feature highlight |
| `TestimonialCard` | Testimonial display |
| `TeamMemberCard` | Team member profile |
| `ArticleCard` | Blog/article preview |
| `Comment` | Comment display |

### Tables

| Component | Description |
|-----------|-------------|
| `TableBasic` | Full table wrapper |
| `TableHeader` | Table header cell |
| `TableRow` | Table row container |
| `TableCell` | Table data cell |

### Utility

| Component | Description |
|-----------|-------------|
| `Divider` | Horizontal divider |
| `LoadingSpinner` | Loading animation |
| `LoadingSkeleton` | Skeleton placeholder |
| `StepIndicator` | Step progress |
| `StepItem` | Individual step |

### Content

| Component | Description |
|-----------|-------------|
| `SocialLinks` | Social media icons |
| `Tag` | Clickable tag |
| `TagList` | Multiple tags container |
| `CodeBlock` | Code snippet display |
| `CallToAction` | CTA section |
| `EmptyState` | Zero data display |

### Form Elements

| Component | Description |
|-----------|-------------|
| `Dropdown` | Dropdown menu |
| `DropdownItem` | Dropdown menu item |
| `ToggleSwitch` | Toggle control |
| `Tooltip` | Hover tooltip |
| `SearchInput` | Search with icon |

### Pricing

| Component | Description |
|-----------|-------------|
| `PricingFeature` | Pricing feature item |

## Testing

Create test files with the `-test.html` suffix:

`routes/component-test.html`:
```html
<MyComponent prop='value' />
```

Run tests:

```bash
gtml test
# or
make test
```

Tests compile the component and verify no errors occur.

## Development

### Build

```bash
go build -o gtml .
```

### Run Tests

```bash
make test
# or
go test ./...
```

### Project Philosophy

gtml uses a "Spec-First" approach where the `spec/` directory contains natural language documentation that describes how the system should work. The spec serves as an intermediate representation (IR) that can be used to generate the full implementation.

## Links

- [Documentation](https://gtml.dev/docs)
- [Component Gallery](https://gtml.dev/docs/components)

## License

MIT
