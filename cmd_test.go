package main

import (
	"reflect"
	"testing"
)

func TestExecuteCommand(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{})
	cmd.Execute()
	if len(argCurrency) != 0 {
		t.Fatalf("expected []string got %v", argCurrency)
	}
	if len(argDate) != 0 {
		t.Fatalf("expected []string got %v", argDate)
	}
	if argSql != "" {
		t.Fatalf("expected \"\" got %v", argSql)
	}

	cmd = newRootCmd()
	cmd.SetArgs([]string{"-c usd,eur"})
	cmd.Execute()
	if !reflect.DeepEqual(argCurrency, []string{"USD", "EUR"}) {
		t.Fatalf("expected []string{\"USD\",\"EUR\"} got %v", argCurrency)
	}

	cmd = newRootCmd()
	cmd.SetArgs([]string{"-d 12.01.2007,1.1.20"})
	cmd.Execute()
	if !reflect.DeepEqual(argDate, []string{"12.01.2007", "1.1.20"}) {
		t.Fatalf("expected []string{\"12.01.2007\",\"1.1.20\"} got %v", argDate)
	}

	cmd = newRootCmd()
	cmd.SetArgs([]string{"-s /usr/d"})
	err := cmd.Execute()
	if err == nil {
		t.Fatalf("expected a cmd error got nil")
	}

	cmd = newRootCmd()
	cmd.SetArgs([]string{"-s f5.sqlite3"})
	cmd.Execute()
	if argSql != "f5.sqlite3" {
		t.Fatalf("expected \"f5.sqlite3\" got %v", argSql)
	}

	cmd = newRootCmd()
	cmd.SetArgs([]string{"ls /"})
	err = cmd.Execute()
	if err == nil {
		t.Fatalf("expected a cmd error got nil")
	}
}

func TestCmdIsDateCorrect(t *testing.T) {
	if isDateCorrect("") {
		t.Fatalf("date must not be an empty string")
	}
	if isDateCorrect("12.mm.2020") {
		t.Fatalf("date must not contain letters")
	}

	if !isDateCorrect("12.11.2020") {
		t.Fatalf("valid date failed validation")
	}
	if !isDateCorrect("1.1.20") {
		t.Fatalf("valid date failed validation")
	}
}

func TestCmdIsFileNameCorrect(t *testing.T) {
	if isFileNameCorrect("") {
		t.Fatalf("file name must not be an empty string")
	}
	if isFileNameCorrect("1") {
		t.Fatalf("file name must not start with a digit")
	}
	if isFileNameCorrect("/f1") {
		t.Fatalf("file name must contain only letters and digits")
	}

	if !isFileNameCorrect("f1") {
		t.Fatalf("valid file name failed validation")
	}
	if !isFileNameCorrect("f1.sqlite3") {
		t.Fatalf("valid file name failed validation")
	}
}
