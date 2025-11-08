package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jongio/azd-script/internal/cmd"
)

func main() {
	ctx := context.Background()
	
	if err := cmd.Execute(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
