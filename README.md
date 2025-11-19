# RowXml -- Line-Based Format → XML Converter

A small and robust CLI tool written in Go that converts a custom Swedish
line-based format into structured XML.

The tool supports: - Reading input from **file**, **raw string**, or
**stdin** - Multiple `P` (person), `F` (family), `A` (address), `T`
(phone) entries - Validated and ordered XML output - Cross-platform
binaries (Linux / macOS / Windows) - Automated releases created via
GitHub Actions

## 📦 Installation

### Download prebuilt binaries (Linux/macOS/Windows)

Browse the latest release here:

➡️ **GitHub Releases**\
https://github.com/UnsinkableSam/RowXml/releases

Make it executable:

## 🐧 Linux
```bash
wget https://github.com/UnsinkableSam/RowXml/releases/latest/download/rowxml-linux-amd64
chmod +x rowxml-linux-amd64
sudo mv rowxml-linux-amd64 /usr/local/bin/rowxml
```


## 🍎 macOS Installation (Intel)
```bash
wget https://github.com/UnsinkableSam/RowXml/releases/latest/download/rowxml-darwin-amd64
chmod +x rowxml-darwin-amd64
sudo mv rowxml-darwin-amd64 /usr/local/bin/rowxml
```



## 🍎 macOS Installation
```bash
wget https://github.com/UnsinkableSam/RowXml/releases/latest/download/rowxml-darwin-amd64
chmod +x rowxml-darwin-amd64
sudo mv rowxml-darwin-amd64 /usr/local/bin/rowxml
```

## 🪟 Windows Installation
```powershell
Invoke-WebRequest -Uri "https://github.com/UnsinkableSam/RowXml/releases/latest/download/rowxml-windows-amd64.exe" -OutFile rowxml-windows-amd64.exe
```

## 🚀 Usage

### 1️⃣ Using a File

``` sh
rowxml --file input.txt
```

### 2️⃣ Passing Raw Data as a String

``` sh
rowxml --string "P|Joe|Biden
A|White House|Washington, D.C|00000"
```

### 3️⃣ Using Stdin

``` sh
cat input.txt | rowxml
```

## 🧪 Running Tests

``` sh
go test ./...
```

## 🏗 Project Structure

    RowXml/
      cmd/rowxml/
      internal/
        parser/
        validator/
        xml/
        model/
      testdata/
      .github/workflows/

## 🔄 GitHub CI/CD

Tag a release:

``` sh
git tag v1.0.0
git push origin v1.0.0
```

## 📖 Specification

P\|förnamn\|efternamn\
T\|mobilnummer\|fastnätsnummer\
A\|gata\|stad\|postnummer\
F\|namn\|födelseår

This tool fully implements the specification.
