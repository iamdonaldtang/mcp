# Design TODO — Dark Theme Restyle

> Active design tasks for `design/pencil-new.pen`
> Previous work archived in `design/progress-archive.md`
> Last updated: 2026-03-06

---

## Phase: Dark Theme Restyle — COMPLETE

> **Objective**: Apply the production dark theme (from `design/legacy/theme/`) to all 81 pages + 19 modals.
> **Status**: DONE — all 15 batches executed + second-pass color fix + visual QA

### Execution Summary

| Step | What | Status |
|------|------|--------|
| Color mapping | Built light→dark mapping from production screenshots | Done |
| Pass 1 (Batch 1-15) | `replace_all_matching_properties` on all page groups | Done |
| Pass 2 (fix) | Caught 11 additional colors: `#f7fafc`, `#e6f2ff`, `#fff5f5`, `#fef2f2`, `#f3f0ff`, `#e9e0ff`, `#f1f5f9`(fill), `#bfdbfe`, `#93c5fd`(fill), `#fffbeb`, `#fee2e2` | Done |
| Visual QA | Spot-checked 12+ pages across all types + layout checks | Done, 0 issues |

### Pages Restyled (100 nodes total)

- **Components**: 327GX, kXg2k, 43K4T, ClmXH
- **B-End (50+ pages)**: Dashboard, Quest, Community, WL, Boost, Analytics, Settings, Dev Kit
- **C-End (9 pages)**: Community Home through Activity Feed
- **Marketing (14 pages)**: M02-M14 (M01 was already dark)
- **Modals (19)**: D01-D19
- **Skipped**: M01 (already dark), CXtOH/EDoSn/PWQV6 (deprecated)

---

## Color Mapping Reference (Light → Dark)

### Fill Colors

| Light (From) | Dark (To) | Usage |
|-------------|-----------|-------|
| `#f8fafc` | `#0A0F1A` | Page background |
| `#ffffff` | `#111B27` | Card/panel background |
| `#fafafa` | `#0D1520` | Subtle bg variant |
| `#f7fafc` | `#0A0F1A` | Near-white bg variant |
| `#f1f5f9` | `#1E293B` | Light gray bg / status |
| `#e2e8f0` | `#1E293B` | Dividers, separators |
| `#d1d5db` | `#2D3748` | Gray fill elements |
| `#eff6ff` | `#0F1A2E` | Blue tint bg |
| `#dbeafe` | `#1A2744` | Blue bg |
| `#ebf4ff` | `#0F1A2E` | Blue tint bg |
| `#e6f2ff` | `#0F1A2E` | Blue tint bg |
| `#bfdbfe` | `#1A2744` | Light blue |
| `#93c5fd` | `#1E3A5E` | Medium blue (as fill) |
| `#f0fdf4` | `#0A1F1A` | Green tint bg |
| `#ecfdf5` | `#0A1F1A` | Green tint bg |
| `#f0fff4` | `#0A1F1A` | Green tint bg |
| `#dcfce7` | `#0A2E1A` | Green badge/status bg |
| `#faf5ff` | `#1A1033` | Purple tint bg |
| `#f3e8ff` | `#1A1033` | Purple tint bg |
| `#ede9fe` | `#1A1033` | Purple tint bg |
| `#f5f3ff` | `#1A1033` | Purple tint bg |
| `#f3f0ff` | `#1A1033` | Purple tint bg |
| `#e9e0ff` | `#1A1033` | Purple tint bg |
| `#fff5eb` | `#1F1508` | Orange tint bg |
| `#ffedd5` | `#1F1508` | Orange tint bg |
| `#fff7ed` | `#1F1508` | Orange tint bg |
| `#fef3c7` | `#1F1A08` | Amber/yellow tint bg |
| `#fffbeb` | `#1F1A08` | Amber light bg |
| `#fff5f5` | `#1F0D0D` | Red tint bg |
| `#fef2f2` | `#1F0D0D` | Red tint bg |
| `#fee2e2` | `#2D1515` | Red status bg |

### Text Colors

| Light (From) | Dark (To) | Usage |
|-------------|-----------|-------|
| `#1e293b` | `#F1F5F9` | Primary text |
| `#2d3748` | `#E2E8F0` | Strong text |
| `#1a202c` | `#FFFFFF` | Strongest text |
| `#4a5568` | `#CBD5E1` | Medium text |
| `#475569` | `#CBD5E1` | Medium text |
| `#718096` | `#94A3B8` | Muted text |
| `#a0aec0` | `#94A3B8` | Muted text |
| `#000000` | `#FFFFFF` | Black → white |
| `#cbd5e1` | `#475569` | Placeholder/disabled |
| `#e2e8f0` | `#334155` | Decorative text |

### Stroke Colors

| Light (From) | Dark (To) | Usage |
|-------------|-----------|-------|
| `#e2e8f0` | `#1E293B` | Default border |
| `#93c5fd` | `#1E3A5E` | Blue accent border |
| `#c6f6d5` | `#1A3D2A` | Green accent border |

### Unchanged Colors
- Product brands: Quest `#5D7EF1`, Community `#48BB78`, WL `#9B7EE0`, Boost `#ED8936`
- Status: Active `#16A34A`, Draft `#D97706`, semantic blues/purples/oranges
- Social: Discord `#5865F2`, Telegram `#229ED9`, YouTube `#FF0000`, LinkedIn `#0A66C2`
