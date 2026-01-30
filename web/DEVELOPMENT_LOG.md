# Development Log

## 2026-01-30

### 🎨 UI/UX Improvements
- **Dark Mode Implementation**: 
  - Integrated `next-themes` for seamless light/dark mode switching.
  - Updated `PublicLayout`, `SiteHeader`, and `HomePage` to support dark mode with adaptive colors (`slate-950` backgrounds, contrasted text).
  - Added a `ThemeToggle` component with Sun/Moon/System options.
- **Navigation Redesign**:
  - Removed the bottom-right `FloatingMenu`.
  - Moved Language and Theme switchers to the `SiteHeader` (top-right) for better accessibility and standard layout compliance.
- **Admin Dashboard Redesign**:
  - Completely rewrote the Admin Dashboard to match the "Dark Terminal" design spec.
  - Implemented a dedicated dark theme for the Admin area (`bg-slate-950`).
  - Added visual elements: Status lights, Gradient cards, Terminal-style breadcrumbs (`root @ kaldalis-cms`), and "Activity Log" timeline.

### 🌐 Internationalization (I18n)
- **Routing Fixes**:
  - Fixed `LanguageSwitcher` to use `next-intl/link` and proper routing logic, resolving 404 errors during locale switching.
  - Ensured `NextIntlClientProvider` receives the correct `locale` prop in `layout.tsx` to prevent stale state.
- **Translation Coverage**:
  - Added `admin` namespace to `en.json` and `zh-CN.json`.
  - Fully translated the Admin Sidebar and Dashboard content (Statistics, Table headers, Charts labels).

### 🛠️ Technical Debt & Cleanup
- Refactored `LanguageSwitcher` to use declarative navigation.
- Cleaned up unused components and redundant routing logic.
