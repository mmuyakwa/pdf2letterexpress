Hier ist eine Textdatei mit Anweisungen zum Erstellen neuer Git-Tags:

```text
# Git Tag Erstellung - Anleitung

## Was sind Git Tags?
Git Tags sind Referenzen zu bestimmten Punkten in der Git-Historie. Sie werden typischerweise verwendet, um Release-Versionen zu markieren (z.B. v1.0.0, v2.1.3).

## Arten von Tags

### 1. Lightweight Tags (Einfache Tags)
- Nur ein Zeiger auf einen bestimmten Commit
- Keine zusätzlichen Metadaten

### 2. Annotated Tags (Annotierte Tags)
- Vollständige Objekte in der Git-Datenbank
- Enthalten Tagger-Name, E-Mail, Datum und Nachricht
- Empfohlen für Releases

## Befehle zum Erstellen von Tags

### Lightweight Tag erstellen:
```bash
git tag v1.0.0
```

### Annotated Tag erstellen:
```bash
git tag -a v1.0.0 -m "Release Version 1.0.0"
```

### Tag für einen bestimmten Commit erstellen:
```bash
git tag -a v1.0.0 -m "Release Version 1.0.0" <commit-hash>
```

## Tags anzeigen

### Alle Tags auflisten:
```bash
git tag
```

### Tags mit Pattern anzeigen:
```bash
git tag -l "v1.*"
```

### Tag-Details anzeigen:
```bash
git show v1.0.0
```

## Tags hochladen (Push)

### Einzelnen Tag hochladen:
```bash
git push origin v1.0.0
```

### Alle Tags hochladen:
```bash
git push origin --tags
```

### Alle Tags und Commits gleichzeitig hochladen:
```bash
git push --follow-tags
```

## Empfohlener Workflow für PDF2LetterExpress

### 1. Versionsnummer festlegen
Verwende Semantic Versioning (MAJOR.MINOR.PATCH):
- MAJOR: Breaking changes
- MINOR: Neue Features (backward compatible)
- PATCH: Bug fixes

### 2. Tag erstellen
```bash
# Beispiel für Version 1.0.0
git tag -a v1.0.0 -m "Initial release with LetterExpress compatibility

Features:
- Precise 5mm margin addition
- Multiple PDF processing methods
- Cross-platform support
- Robust error handling"
```

### 3. Tag hochladen
```bash
git push origin v1.0.0
```

### 4. Release auf GitHub erstellen
- Gehe zu GitHub → Releases → Create a new release
- Wähle den erstellten Tag aus
- Füge Release Notes hinzu
- Lade Binaries hoch (optional)

## Beispiel-Tags für PDF2LetterExpress

### Erste Veröffentlichung:
```bash
git tag -a v1.0.0 -m "Initial release - PDF2LetterExpress v1.0.0

✨ Features:
- Automatic 5mm margin addition
- LetterExpress compatibility
- Multi-platform support (Windows, macOS, Linux)
- Intelligent content scaling
- Robust error handling

🔧 Technical:
- Go 1.19+ support
- ImageMagick integration
- pdfcpu-based processing
- Comprehensive logging"
```

### Bug Fix Release:
```bash
git tag -a v1.0.1 -m "Bug fix release v1.0.1

🐛 Fixed:
- Corrected page dimension calculations
- Improved error handling for corrupted PDFs
- Fixed memory leak in large file processing"
```

### Feature Release:
```bash
git tag -a v1.1.0 -m "Feature release v1.1.0

✨ New Features:
- Batch processing support
- Custom margin configuration
- Progress indicators
- Enhanced CLI interface

🐛 Bug Fixes:
- Various stability improvements"
```

## Tag-Management

### Tag löschen (lokal):
```bash
git tag -d v1.0.0
```

### Tag löschen (remote):
```bash
git push origin --delete v1.0.0
```

### Tag verschieben:
```bash
# Alten Tag löschen
git tag -d v1.0.0
git push origin --delete v1.0.0

# Neuen Tag erstellen
git tag -a v1.0.0 -m "Neue Nachricht" <neuer-commit>
git push origin v1.0.0
```

## Best Practices

1. **Immer annotierte Tags für Releases verwenden**
2. **Semantic Versioning befolgen**
3. **Aussagekräftige Tag-Nachrichten schreiben**
4. **Tags vor dem Push testen**
5. **Changelog in Tag-Nachricht einbeziehen**
6. **Tags nie für instabile/ungetestete Versionen erstellen**

## Schnell-Referenz

```bash
# Tag erstellen und hochladen (Komplett-Workflow)
git tag -a v1.0.0 -m "Release v1.0.0 - Beschreibung"
git push origin v1.0.0

# Alle Tags anzeigen
git tag

# Tag-Details anzeigen
git show v1.0.0

# Alle Tags hochladen
git push origin --tags
```
```

Diese Datei speichern Sie als `new_tag.txt` in Ihrem Projektverzeichnis. Sie enthält alle wichtigen Informationen zum Erstellen und Verwalten von Git-Tags, speziell angepasst für Ihr PDF2LetterExpress-Projekt.Diese Datei speichern Sie als `new_tag.txt` in Ihrem Projektverzeichnis. Sie enthält alle wichtigen Informationen zum Erstellen und Verwalten von Git-Tags, speziell angepasst für Ihr PDF2LetterExpress-Projekt.