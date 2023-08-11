package main

import (
	"github.com/magefile/mage/sh"
)

// PreCommit Install pre-commit hooks
func PreCommit() error {
	preCommit := sh.OutCmd("pre-commit")

	out, err := preCommit("install", "--hook-type", "commit-msg")
	if out != "" {
		println(out)
	}

	if err != nil {
		return err
	}

	out, err = preCommit("install")
	if out != "" {
		println(out)
	}

	if err != nil {
		return err
	}

	return nil
}

// Install install dependencies
func Install() error {
	return sh.Run("go", "mod", "download")
}
