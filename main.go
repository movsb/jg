package main

import (
	"context"
	"log"

	jr "github.com/movsb/jg/runtime"
)

func main() {
	ctx := context.Background()

	rt := jr.MustNewRuntime(
		ctx,
		jr.WithStd(),
	)

	output, err := rt.ExecuteFile(ctx, `main.js`)
	if err != nil {
		log.Fatalln(err)
	}

	if output != nil {
		log.Println(`exited with:`, output)
	}
}
