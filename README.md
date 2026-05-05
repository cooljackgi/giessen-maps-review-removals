# Gießen Google-Bewertungen: Diffamierungs-Löschbanner

Lokaler Go-Workflow, um öffentlich sichtbare Google-Maps-Ortsdaten für Gießen zu sammeln, Hinweise auf entfernte Bewertungen zu erkennen und daraus Auswertungen sowie ein Dashboard zu erzeugen.

## Status

- Die Instanz ist auf Gießen als Standardstadt vorbereitet.
- Standard-PLZ: `35390, 35392, 35394, 35396, 35398`
- Output-Dateien und Modulname sind auf `giessen-*` umgestellt.
- Statistische Bezirks- bzw. Stadtteil-Geometrien sind aktuell **nicht** eingebunden. Dashboard und Scraper laufen trotzdem; die Bezirks-/Overlay-Funktion ist bewusst deaktiviert.

## Voraussetzungen

- Go 1.25+
- Chrome oder Chromium
- Optional: ImageMagick `magick` für PNG-Export

## Schnellstart

```bash
make setup
make scrape ARGS="--postcodes all --headless=false"
make charts
make dashboard
```

## Wichtige Kommandos

```bash
make scrape ARGS="--postcodes 35390,35392 --queries restaurant,cafe,imbiss"
make backfill
make validate
go run ./cmd/validate --strict-city
make charts ARGS="--png"
make dashboard
```

## Outputs

- `output/discovery.json`
- `output/places.json`
- `output/places.csv`
- `output/metadata.json`
- `output/charts/giessen_dashboard.html`
- `output/charts/giessen_overall_summary.svg`
- `output/charts/giessen_<PLZ>_summary.svg`
- `output/charts/giessen_most_removed.csv`
- `output/charts/giessen_most_removed.md`
- `output/charts/giessen_most_removed.html`

## Offener Punkt

Wenn du die Stadtteil-/Gebietsflächen im Dashboard wiederhaben willst, muss als Nächstes eine belastbare Gießener Geometriequelle eingebunden werden. Die Codebasis ist dafür bereits so vorbereitet, dass diese Logik separat ergänzt werden kann, ohne die restliche Gießen-Instanz erneut umzubauen.
