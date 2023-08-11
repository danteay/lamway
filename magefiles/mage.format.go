package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Vet execute go vet checks over all the code.
func Vet() error {
	command := "go"
	args := []string{"vet", "./..."}

	out, err := sh.Output(command, args...)

	if out != "" {
		println(out)
	}

	return err
}

// Lint Runs revive checks over the code.
func Lint() error {
	mg.Deps(Vet)

	command := "revive"
	args := []string{"-config=revive.toml", "-formatter=friendly", "-exclude=magefiles/...", "./..."}

	out, err := sh.Output(command, args...)

	if out != "" {
		println(out)
	}

	return err
}

// Format Runs gofmt over the code.
func Format() error {
	outImp, err := sh.Output("goimports-reviser", "-format", "./...")
	if err != nil {
		return err
	}

	if outImp != "" {
		println(outImp)
	}

	return nil
}
