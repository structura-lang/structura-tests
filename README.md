# Structura Language Test Suite

This repository contains an automated test harness for the **Structura** language implementation. The purpose of the suite is to verify that the interpreter correctly handles the language constructs defined in the official specification.

## Repository Layout

```
structura-tests/
├─ go/                 # Go driver that loads each test case, executes it, and writes a report
├─ tests/              # Individual test case directories (one per feature)
│   ├─ <case‑name>/    # Each case contains:
│   │   ├─ program.json   # Structura program definition
│   │   ├─ data.json      # Optional input data for the program (may be empty)
│   │   └─ output.json    # Expected result: {"out": <value>, "error": <bool>}
│   └─ ...
└─ README.md           # This file
```

## How Tests Are Executed

1. **Preparation** – Ensure the Go toolchain is installed.
2. **Set up interpreter** Clone the https://github.com/structura-lang/structura-go repository next to this repository.
3. **Running the Tests** – From the repository root execute the Go driver:
   ```bash
   go run ./go/main.go
   ```
   The driver walks through every sub‑directory in `tests/`, loads the three JSON files, invokes the interpreter, and records whether the actual output matches the expected `output.json`.
4. **Results** – After execution two report files are produced next to the binary:
   * `report.json` – machine‑readable array of `{name, passed, error_msg, expected, got}` objects.
   * `report.txt` – human‑readable summary with PASS/FAIL markers.

## Adding New Test Cases

To add a test, create a new folder under `tests/` with a unique name and provide the three JSON files described above. The Go driver will automatically discover and run the new case on the next execution.

## Reference Implementation & Specification

The interpreter and language definition are maintained in the main Structura repository:

* **Repository:** https://github.com/structura-lang/structura-lang
* **Specification:** see `spec.md` in that repository for the full language description.

For any questions about language semantics, refer to the official spec rather than this test suite documentation.