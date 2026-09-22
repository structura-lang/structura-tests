package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"structura/interpreter"
)

type caseResult struct {
	Name     string `json:"name"`
	Passed   bool   `json:"passed"`
	ErrorMsg string `json:"error_msg,omitempty"`
	Expected any    `json:"expected,omitempty"`
	Got      any    `json:"got,omitempty"`
}

func loadJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func runCase(base, name string) caseResult {
	var prog map[string]any
	var inp map[string]any
	var exp struct {
		Out   any  `json:"out"`
		Error bool `json:"error"`
	}
	var result any
	var execErr error
	passed := true
	var errMsg string

	// Defer recover to handle unexpected panics.
	defer func() {
		if r := recover(); r != nil {
			passed = false
			errMsg = fmt.Sprintf("panic: %v", r)
		}
	}()

	progPath := filepath.Join(base, name, "program.json")
	dataPath := filepath.Join(base, name, "data.json")
	outPath := filepath.Join(base, name, "output.json")

	if err := loadJSON(progPath, &prog); err != nil {
		passed = false
		errMsg = fmt.Sprintf("load program: %v", err)
	}

	if passed {
		if err := loadJSON(dataPath, &inp); err != nil {
			if os.IsNotExist(err) {
				inp = map[string]any{}
			} else {
				passed = false
				errMsg = fmt.Sprintf("load data: %v", err)
			}
		}
	}

	if passed {
		if err := loadJSON(outPath, &exp); err != nil {
			passed = false
			errMsg = fmt.Sprintf("load output: %v", err)
		}
	}

	if passed {
		result, execErr = interpreter.Execute(prog, inp)
		if exp.Error {
			if execErr == nil && passed {
				passed = false
				errMsg = "expected error but got none"
			}
		} else {
			if execErr != nil {
				passed = false
				errMsg = execErr.Error()
			} else if !reflect.DeepEqual(result, exp.Out) {
				passed = false
				errMsg = "output mismatch"
			}
		}
	}

	return caseResult{Name: name, Passed: passed, ErrorMsg: errMsg, Expected: exp.Out, Got: result}
}

func main() {
	base := "../tests"

	entries, err := os.ReadDir(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read cases dir: %v\n", err)
		os.Exit(1)
	}

	var results []caseResult

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		results = append(results, runCase(base, name))
	}

	// Write JSON report next to this binary.
	jsonPath := "report.json"
	jsonData, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(jsonPath, jsonData, 0644)

	// Write human‑readable report.
	txtPath := "report.txt"
	f, _ := os.Create(txtPath)
	defer f.Close()

	for _, r := range results {
		status := "PASS"
		if !r.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(f, "%-50s %s\n", r.Name, status)
		if r.ErrorMsg != "" {
			fmt.Fprintf(f, "    %s\n", r.ErrorMsg)
		}
	}

	total := len(results)
	passedCount := 0

	for _, r := range results {
		if r.Passed {
			passedCount++
		}
	}

	fmt.Fprintf(f, "\nTotal: %d, Passed: %d, Failed: %d\n", total, passedCount, total-passedCount)
}
