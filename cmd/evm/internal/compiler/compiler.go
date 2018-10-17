package compiler

import (
	"errors"
	"fmt"

	"github.com/TreasureChain/go-tstchain/core/asm"
)

func Compile(fn string, src []byte, debug bool) (string, error) {
	compiler := asm.NewCompiler(debug)
	compiler.Feed(asm.Lex(fn, src, debug))

	bin, compileErrors := compiler.Compile()
	if len(compileErrors) > 0 {
		// report errors
		errs := ""
		for _, err := range compileErrors {
			errs += fmt.Sprintf("%s:%v\n", fn, err)
		}
		return "", errors.New(errs + "compiling failed\n")
	}
	return bin, nil
}
