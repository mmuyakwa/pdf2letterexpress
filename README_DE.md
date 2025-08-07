# 📄 PDF2LetterExpress

> 🚀 Wandeln Sie Ihre PDFs für LetterExpress-Kompatibilität um, indem Sie präzise 5mm Ränder hinzufügen

[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![Lizenz](https://img.shields.io/badge/Lizenz-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Erfolgreich-brightgreen.svg)]()

## 🎯 Was macht es?

PDF2LetterExpress ist ein spezialisiertes Tool, das automatisch **exakt 5mm Ränder** zu PDF-Dokumenten hinzufügt und sie damit kompatibel für den LetterExpress-Postdienst macht. Das Tool skaliert intelligent Ihre PDF-Inhalte und behält dabei alle Texte, Bilder und Vektorelemente bei.

## ✨ Funktionen

- 🎯 **Präzise 5mm Ränder** auf allen vier Seiten
- 📐 **Intelligente Inhaltsskalierung** mit Beibehaltung der Seitenverhältnisse
- 🔧 **Mehrere Verarbeitungsmethoden** mit automatischen Fallbacks
- 💼 **Plattformübergreifende Unterstützung** (Windows, macOS, Linux)
- 🚀 **Schnelle Verarbeitung** mit optimierten Algorithmen
- 📝 **Detailliertes Logging** für Fehlerbehebung
- 🛡️ **Robuste Fehlerbehandlung** für verschiedene PDF-Typen

## 🚀 Schnellstart

### Installation

Laden Sie die neueste Binary für Ihre Plattform von der [Releases-Seite](releases/) herunter oder kompilieren Sie aus dem Quellcode:

```bash
git clone https://github.com/mmuyakwa/pdf2letterexpress.git
cd pdf2letterexpress
go build -o pdf2letterexpress main.go
```

### Verwendung

```bash
# Grundlegende Verwendung
./pdf2letterexpress eingabe.pdf

# Die Ausgabe wird sein: eingabe - converted.pdf
```

### 📋 Beispiele

```bash
# Konvertieren eines Arztberichts
./pdf2letterexpress "2025-07-28_Befund_Nephrologe.pdf"
# Ausgabe: 2025-07-28_Befund_Nephrologe - converted.pdf

# Konvertieren einer Rechnung
./pdf2letterexpress rechnung_2025.pdf
# Ausgabe: rechnung_2025 - converted.pdf
```

## 🔧 Wie es funktioniert

1. **📖 Analysiert** die Seitendimensionen Ihrer PDF
2. **🧮 Berechnet** den exakten Skalierungsfaktor für 5mm Ränder
3. **🎨 Transformiert** den Inhalt mit mehreren Methoden:
   - Primär: Direkte PDF-Manipulation mit pdfcpu
   - Fallback 1: Ghostscript-basierte Skalierung
   - Fallback 2: ImageMagick-Konvertierung
4. **💾 Gibt** eine neue PDF mit erhaltener Qualität aus

## 📊 Technische Details

- **Randgröße**: Exakt 5mm (≈14,17 Punkte)
- **Seitenformat-Unterstützung**: Alle Standardformate (A4, Letter, Legal, etc.)
- **Inhaltserhaltung**: 100% Treue für Text, Bilder und Vektoren
- **Verarbeitungsgeschwindigkeit**: ~1-2 Sekunden pro Seite
- **Speicherverbrauch**: Optimiert für große Dateien

## 🛠️ Aus Quellcode kompilieren

### Voraussetzungen

- Go 1.19 oder höher
- Git

### Build-Schritte

```bash
# Repository klonen
git clone https://github.com/mmuyakwa/pdf2letterexpress.git
cd pdf2letterexpress

# Abhängigkeiten herunterladen
go mod download

# Für Ihre Plattform kompilieren
go build -o pdf2letterexpress main.go

# Oder für alle Plattformen kompilieren
make build-all
```

### Cross-Kompilierung

```bash
# Für macOS (ARM64) kompilieren
GOOS=darwin GOARCH=arm64 go build -o pdf2letterexpress-darwin-arm64 main.go

# Für Windows kompilieren
GOOS=windows GOARCH=amd64 go build -o pdf2letterexpress-windows-amd64.exe main.go

# Für Linux kompilieren
GOOS=linux GOARCH=amd64 go build -o pdf2letterexpress-linux-amd64 main.go
```

## 🐛 Fehlerbehebung

### Häufige Probleme

**❌ "Fehler beim Verarbeiten der PDF"**
- Prüfen Sie, ob die Eingabedatei eine gültige PDF ist
- Stellen Sie sicher, dass Sie Leseberechtigungen haben
- Versuchen Sie es mit einer anderen PDF, um das Problem zu isolieren

**❌ "Ghostscript nicht gefunden"**
- Installieren Sie Ghostscript: `brew install ghostscript` (macOS) oder laden Sie es von [ghostscript.com](https://www.ghostscript.com/) herunter
- Das Tool verwendet automatisch andere Methoden als Fallback

**❌ Ausgabe-PDF hat falsche Dimensionen**
- Dies wurde in der neuesten Version behoben
- Stellen Sie sicher, dass Sie die neueste Version verwenden

### Debug-Modus

Aktivieren Sie detailliertes Logging:

```bash
export LOG_LEVEL=debug
./pdf2letterexpress eingabe.pdf
```

## 📦 Abhängigkeiten

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF-Verarbeitungsbibliothek
- [logrus](https://github.com/sirupsen/logrus) - Strukturiertes Logging
- [cobra](https://github.com/spf13/cobra) - CLI-Framework

## 🤝 Mitwirken

Wir begrüßen Beiträge! Bitte lesen Sie unseren [Mitwirken-Leitfaden](CONTRIBUTING.md) für Details.

1. Repository forken
2. Feature-Branch erstellen (`git checkout -b feature/tolles-feature`)
3. Änderungen committen (`git commit -m 'Tolles Feature hinzufügen'`)
4. Zum Branch pushen (`git push origin feature/tolles-feature`)
5. Pull Request öffnen

## 📄 Lizenz

Dieses Projekt ist unter der MIT-Lizenz lizenziert - siehe die [LICENSE](LICENSE)-Datei für Details.

## 🙏 Danksagungen

- Mit ❤️ in Go entwickelt
- PDF-Verarbeitung ermöglicht durch [pdfcpu](https://github.com/pdfcpu/pdfcpu)
- Inspiriert durch die Notwendigkeit der LetterExpress-Kompatibilität

## 📞 Support

- 🐛 **Fehlerberichte**: [Issue öffnen](issues/new?template=bug_report.md)
- 💡 **Feature-Wünsche**: [Issue öffnen](issues/new?template=feature_request.md)
- 💬 **Fragen**: [Diskussion starten](discussions/)

---

<div align="center">

**Mit 🇩🇪 in Deutschland entwickelt**

⭐ Geben Sie diesem Repo einen Stern, wenn es Ihnen geholfen hat! ⭐

</div>