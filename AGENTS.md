## Architecture

The stack is Go with the `verb` framework, HTMX, Bootstrap 5, and GORM against PostgreSQL. The `verb` framework organizes the app into four primitives: `v.Page` for full-page renders, `v.Component` for HTMX partials (auto-registered at `/htmx/<template-name>`), `v.Action` for plain HTTP handlers, and `v.Bridge` (via `verb.Map("Key", fn)`) for injecting data into templates under a named key like `.Search` or `.Products`.

A critical GORM constraint: never combine `Group()` and `Preload()` in PostgreSQL, because GROUP BY conflicts with the extra columns association queries emit. The workaround is a two-step pattern — first query DISTINCT IDs via JOIN, then run a separate `Preload` with `WHERE id IN ?`.

The most important HTMX architectural rule is that any component renderable in multiple container contexts must receive its own results container selector, passed as a `results_target` query param, mapped through the bridge into the template, and forwarded into JavaScript as a function parameter. This prevents hardcoded DOM selectors from breaking when the same component is mounted in different pages.

## CSS Design Spec

Every component shares a single visual language. Surfaces use `#f8f9fa` as the background with `#ffffff` for inset elements like inputs and pills. Borders are `1px solid #dee2e6` throughout. Border radii follow a three-tier scale: `12px` for cards and wrappers, `8px` for inputs, and `20px` for pills. Shadows are `0 2px 8px rgba(0,0,0,0.04)` at rest and `0 4px 16px rgba(0,0,0,0.08)` on hover or focus. Primary text is `#495057`, muted text and icons are `#adb5bd`, and the focus accent is `#0d6efd` (Bootstrap primary blue). The font stack is `'Segoe UI', system-ui, -apple-system, sans-serif` and all transitions run `0.2–0.3s ease`.

Structurally, cards are `flex-direction: column` with `gap: 0.65rem`. Metric stat grids use `grid-template-columns: 1fr 1fr` and require `min-width: 0` on children to allow text truncation via ellipsis. Scrollable pill lists cap at `max-height: 6rem` with a 3px custom scrollbar. When multiple components share a page, CSS is scoped with a component prefix (e.g. `sr-*` for product search results) to prevent class collisions.
