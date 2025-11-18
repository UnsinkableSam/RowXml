# RowXML Converter

Convert a custom line-based format (P/T/A/F records) into clean XML.


This tool converts input like:

```
P|Carl Gustaf|Bernadotte
T|0768-101801|08-101801
A|Drottningholms slott|Stockholm|10001
F|Victoria|2012
A|Solliden|Öland|10002
T|0702-020202|02-202020
```

Into:

```xml
<people>
  <person>
    <firstname>Carl Gustaf</firstname>
    <lastname>Bernadotte</lastname>
    ...
  </person>
</people>
```

---

## 🚀 Features

* Converts custom P/T/A/F formatted text into structured XML
* Fully validated structure (wrong records produce meaningful errors)
* Supports multiple phone/address/family entries
* Strict ordering rules (T grouped together, A grouped together, etc.)
* Works on **macOS**, **Linux**, **Windows**
* Distributed as binaries via GitHub Releases
* Fully tested (unit + integration tests)

---

## 📦 Installation

### Option 1 — Download Prebuilt Executables

Each push to `main` triggers GitHub Actions which automatically builds releases.

1. Go to **GitHub → Releases**
2. Download the binary for your OS:

| OS              | File                       |
| --------------- | -------------------------- |
| macOS (arm64)   | `rowxml-darwin-arm64`      |
| macOS (amd64)   | `rowxml-darwin-amd64`      |
| Linux (amd64)   | `rowxml-linux-amd64`       |
| Windows (amd64) | `rowxml-windows-amd64.exe` |

### Make it executable (Linux/macOS)

```bash
chmod +x rowxml-*
```

---

### Option 2 — Install From Source

```bash
git clone https://github.com/<yourname>/<repo>.git
cd <repo>
go build -o rowxml ./cmd/rowxml
```

---

## 🏃 Running the CLI

The program accepts input as **a single argument string**.

### Example:

```bash
./rowxml "P|Joe|Biden
A|White House|Washington, D.C|00000"
```

Or with a heredoc:

```bash
./rowxml "$(cat <<EOF
P|Victoria|Bernadotte
T|070-0101010|0459-123456
A|Haga Slott|Stockholm|101
EOF
)"
```

It prints XML to stdout.

---

## 📘 Usage

### Input Specification

| Tag | Meaning |           |                   |                   |
| --- | ------- | --------- | ----------------- | ----------------- |
| `P  | first   | last`     | Start new person  |                   |
| `T  | mobile  | landline` | Add phone entry   |                   |
| `A  | street  | city      | postcode`         | Add address entry |
| `F  | name    | bornYear` | Add family member |                   |

### Structural Rules

* `P` starts a new `<person>`
* `P` may be followed by: `T`, `A`, `F`
* `F` may be followed by: `T`, `A`
* Multiple `T` and `A` entries allowed
* XML tag order is normalized (all phone entries grouped, all address entries grouped, etc.)

---

## 🧪 Running Tests

All tests are standard Go tests.

### Run the entire test suite:

```bash
go test ./...
```

### Run with verbose output:

```bash
go test -v ./...
```

### Run a specific test:

```bash
go test -run TestFullExampleCorrected ./...
```

---

## 🛠 Project Structure

```
.
├── cmd/
│   └── rowxml/
│       └── main.go      # CLI entrypoint
├── internal/
│   ├── model/
│   │   └── model.go     # XML structs
│   ├── parser/
│   │   ├── parser.go    # Parse P/T/A/F to People model
├── go.mod
└── README.md
```

---

## 🧩 Exit Codes

| Code | Meaning                      |
| ---- | ---------------------------- |
| `0`  | OK                           |
| `1`  | Invalid input or parse error |

---

## 🔧 Building Binaries Manually

Build for all major platforms:

```bash
GOOS=linux   GOARCH=amd64 go build -o rowxml-linux-amd64   ./cmd/rowxml
GOOS=darwin  GOARCH=arm64 go build -o rowxml-darwin-arm64  ./cmd/rowxml
GOOS=darwin  GOARCH=amd64 go build -o rowxml-darwin-amd64  ./cmd/rowxml
GOOS=windows GOARCH=amd64 go build -o rowxml-windows-amd64.exe ./cmd/rowxml
```

---

## 🤖 GitHub Actions (Automatic Releases)

Your workflow automatically:

1. Builds binaries for all OS targets
2. Runs all tests
3. Uploads binaries to a GitHub Release

Users can download the latest version without compiling.

---

## 📄 License

MIT — free to use anywhere.

---

## 🎉 Next Steps / Optional Enhancements

* Add interactive mode (read from stdin)
* Add `--input-file` flag
* Add golden XML output tests
* Generate man pages or completions (bash/zsh/fish)
* Add benchmarks and fuzz tests
* CI checks: `go vet`, `golangci-lint`

