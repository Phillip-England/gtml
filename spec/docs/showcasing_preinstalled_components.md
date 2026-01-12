# Showcasing Preinstalled Components

## Purpose
The documentation site includes pages dedicated to showcasing the preinstalled components that come with gtml. This gives users a visual reference for what each component looks like and how to use it.

## Component Showcase Structure

### Gallery Page
Located at `/docs/components`, this page displays all documented components in a grid layout organized by category. Each component card links to its individual showcase page.

### Individual Component Pages
Each documented component has its own page at `/components/{category}/{component-name}`. These pages follow a consistent structure:

1. **Navigation**: The DocsNavbar at the top with "Components" highlighted as active
2. **Sidebar**: The DocsSidebar on the left showing all components, with the current component highlighted
3. **Back Link**: A "Back to Documentation" link at the top of the content area
4. **Title and Description**: The component name as an h1 and a brief description
5. **gtml Code Section**: Shows the gtml markup needed to use the component
6. **Example Section**: A live rendered example of the component in a bordered container
7. **Props Table**: A table listing all props with their name, type, and description

## Currently Documented Components

The following 19 components have documentation pages:

### Alerts (4 components)
- `/components/alerts/alert-success` - AlertSuccess
- `/components/alerts/alert-error` - AlertError
- `/components/alerts/alert-warning` - AlertWarning
- `/components/alerts/alert-info` - AlertInfo

### Badges (5 components)
- `/components/badges/badge-primary` - BadgePrimary
- `/components/badges/badge-success` - BadgeSuccess
- `/components/badges/badge-warning` - BadgeWarning
- `/components/badges/badge-danger` - BadgeDanger
- `/components/badges/skill-badge` - SkillBadge

### Buttons (7 components)
- `/components/buttons/button-primary` - ButtonPrimary
- `/components/buttons/button-secondary` - ButtonSecondary
- `/components/buttons/button-outline` - ButtonOutline
- `/components/buttons/button-danger` - ButtonDanger
- `/components/buttons/button-success` - ButtonSuccess
- `/components/buttons/button-sm` - ButtonSm
- `/components/buttons/button-lg` - ButtonLg

### Cards (1 component)
- `/components/cards/card-basic` - CardBasic

### Display (1 component)
- `/components/display/avatar` - Avatar

### Feedback (1 component)
- `/components/feedback/notification` - Notification

## Adding New Component Documentation

When adding documentation for a new preinstalled component:

1. Create the route file at `docs/routes/components/{category}/{component-name}.html`
2. Follow the existing page structure (navbar, sidebar, code example, live example, props table)
3. Add the component to the DocsSidebar component at `docs/components/DocsSidebar.html`
4. Add the component to the gallery page at `docs/routes/docs/components.html`
5. Update this spec file to include the new component in the list above

## Route File Structure

The category folder structure for routes mirrors the logical grouping of components:

```
docs/routes/components/
  alerts/
    alert-success.html
    alert-error.html
    alert-warning.html
    alert-info.html
  badges/
    badge-primary.html
    badge-success.html
    badge-warning.html
    badge-danger.html
    skill-badge.html
  buttons/
    button-primary.html
    button-secondary.html
    button-outline.html
    button-danger.html
    button-success.html
    button-sm.html
    button-lg.html
  cards/
    card-basic.html
  display/
    avatar.html
  feedback/
    notification.html
```

## Sidebar Categories

The DocsSidebar organizes components into the following categories:
- Alerts
- Badges
- Buttons
- Cards
- Display
- Feedback

Note: The sidebar only includes components that have corresponding route pages. Many preinstalled components exist in `docs/components/` but are not yet documented with showcase pages.
