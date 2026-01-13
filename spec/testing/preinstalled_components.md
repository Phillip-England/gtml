# Testing Preinstalled Components

## Test Purpose
Validate that all preinstalled components included with gtml compile correctly and render expected HTML output. Each preinstalled component should have corresponding tests to ensure functionality is preserved across updates.

## Test Location
All preinstalled component tests should be located in:
```
./tests/preinstalled_components/
```

Tests for each component category should be in separate files:
```
./tests/preinstalled_components/buttons_test.md
./tests/preinstalled_components/forms_test.md
./tests/preinstalled_components/cards_test.md
./tests/preinstalled_components/alerts_test.md
./tests/preinstalled_components/badges_test.md
./tests/preinstalled_components/navigation_test.md
./tests/preinstalled_components/tables_test.md
./tests/preinstalled_components/layout_test.md
./tests/preinstalled_components/data_display_test.md
./tests/preinstalled_components/content_test.md
./tests/preinstalled_components/form_elements_test.md
./tests/preinstalled_components/utility_test.md
./tests/preinstalled_components/pricing_test.md
```

## Test Requirements

### Component Existence Tests
Each preinstalled component must have a test verifying it exists in the preinstalled components directory:

```
./myapp/components/buttons/ButtonPrimary.html
./myapp/components/forms/InputText.html
./myapp/components/cards/CardBasic.html
```

### Component Compilation Tests
Each preinstalled component must compile without errors when used in a route:

```html
<ButtonPrimary text='Click Me' />
```

Should produce valid HTML output with Tailwind CSS classes intact.

### Prop Validation Tests
Each preinstalled component with props must validate that props are correctly rendered:

```html
<ButtonPrimary text='Submit' />
```

Should produce:
```html
<button class="px-4 py-2 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200 cursor-pointer">
  Submit
</button>
```

### Tailwind Class Preservation Tests
All Tailwind CSS classes in preinstalled components must be preserved after compilation:

Input:
```html
<ButtonSuccess text='Success' />
```

Output should contain all original classes:
```html
class="px-4 py-2 bg-green-600 text-white font-medium rounded-lg hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition-colors duration-200 cursor-pointer"
```

### Component Rendering Tests

#### Button Components
Test each button variant renders correctly:
- `ButtonPrimary` - Blue background classes
- `ButtonSecondary` - Gray background classes
- `ButtonOutline` - Border-only classes
- `ButtonDanger` - Red background classes
- `ButtonSuccess` - Green background classes
- `ButtonSm` - Small sizing classes
- `ButtonLg` - Large sizing classes

#### Form Input Components
Test each form input renders with correct attributes:
```html
<InputText label='Name' placeholder='Enter name' name='username' />
```

Should produce:
```html
<div class="mb-4">
  <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
  <input type="text" name="username" placeholder="Enter name" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
</div>
```

#### Card Components
Test card components render with correct structure:
```html
<CardBasic title='Card Title' content='Card content goes here' />
```

#### Alert Components
Test alert components render with correct styling:
```html
<AlertSuccess message='Operation completed successfully' />
```

#### Badge Components
Test badge components render with correct colors:
```html
<BadgePrimary text='Active' />
```

### Route Integration Tests
Each preinstalled component should be tested within a route to ensure full compilation pipeline works:

```html
<Route: components-demo.html>
<ButtonPrimary text='Click Me' />
<AlertSuccess message='All components working' />
```

### Category Integration Tests
Test multiple components from the same category working together:

```html
<div>
  <ButtonPrimary text='Save' />
  <ButtonSecondary text='Cancel' />
  <ButtonDanger text='Delete' />
</div>
```

### Cross-Category Tests
Test components from different categories used together:

```html
<CardBasic title='User Profile' content='Profile details'>
  <Avatar name='John Doe' size='md' />
  <BadgeSuccess text='Verified' />
</CardBasic>
```

### Prop Type Validation Tests
Verify that components enforce correct prop types:

Valid usage:
```html
<ProgressBar value={50} max={100} label='Progress' />
```

Invalid usage (should error):
```html
<ProgressBar value='50' max={100} label='Progress' />
```

### Default Values Tests
Components with optional props should handle missing props gracefully.

### Conditional Rendering Tests
Components using ternary operators should work correctly:

```html
<ToggleSwitch label='Enable Feature' isOn={true} />
```

### Layout Component Tests
Test layout components with nested content:

```html
<SidebarLayout sidebar={SidebarItems} main={MainContent} />
```

### Data Display Component Tests
Test data display components render data correctly:

```html
<StatisticCard label='Total Users' value='1,234' change='+12%' changeType='positive' />
```

### Nested Slot Tests
Test components that use slots render child content correctly:

```html
<CardWithImage title='Product' content='Description' imageUrl='/image.jpg'>
  <ButtonPrimary text='Buy Now' />
</CardWithImage>
```

## Automated Test Execution
Tests for preinstalled components should be executable via the CLI:
```
gtml test ./tests/preinstalled_components/
```

## Test Coverage Requirements
Each preinstalled component must have:
1. Existence test
2. Compilation test
3. Rendering output test
4. Prop validation test

## Regression Prevention
Tests should prevent:
- Tailwind class removal or modification
- Prop type changes without documentation
- HTML structure changes that break existing usage
- Missing or broken subcomponents in component directories

## Documentation Tests
Each component category test file should document:
- List of components in the category
- Required props for each component
- Optional props for each component
- Example usage patterns
- Expected output for each component
