# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

```bash
# Install dependencies (uses pnpm based on lock file)
pnpm install

# Start development server (Expo Go compatible)
pnpm start

# Platform-specific development
pnpm ios      # Start with iOS simulator
pnpm android  # Start with Android emulator
pnpm web      # Start web version

# Linting
pnpm lint
```

## Architecture

This is an Expo SDK 54 React Native app using the new architecture with React 19 and expo-router for file-based routing.

### Routing Structure
- `app/_layout.tsx` - Root layout with theme provider (Stack navigator)
- `app/(tabs)/` - Tab-based navigation group
  - `_layout.tsx` - Tab configuration with bottom tabs
  - `index.tsx` - Home screen
  - `explore.tsx` - Explore screen
- `app/modal.tsx` - Modal screen (presentation: 'modal')

### Key Patterns
- **Path alias**: Use `@/` to import from project root (e.g., `@/components/themed-text`)
- **Platform-specific files**: Use `.ios.tsx` / `.web.ts` suffixes for platform-specific implementations (e.g., `icon-symbol.ios.tsx`)
- **Theming**: Use `useColorScheme` hook from `@/hooks/use-color-scheme` and `Colors` from `@/constants/theme`
- **Typed routes**: Enabled via `experiments.typedRoutes` in app.json

### Component Structure
- `components/` - Shared components (ThemedText, ThemedView, ParallaxScrollView, etc.)
- `components/ui/` - UI primitives (IconSymbol, Collapsible)
- `hooks/` - Custom hooks (useColorScheme, useThemeColor)
- `constants/theme.ts` - Color and font definitions for light/dark themes

### Expo Configuration
- New Architecture enabled (`newArchEnabled: true`)
- React Compiler enabled (`experiments.reactCompiler: true`)
- Uses expo-image, expo-router, expo-haptics, react-native-reanimated
