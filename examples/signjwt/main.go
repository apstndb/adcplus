package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/apstndb/signer"
)

func main() {
	if err := _main(); err != nil {
		log.Fatalln(err)
	}
}

func _main() error {
	input := flag.String("input", "", "")
	output := flag.String("output","" , "")
	_ = output
	flag.Parse()

	var reader io.Reader
	if *input != "" {
		file, err := os.Open(*input)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}
	b, err := ioutil.ReadAll(reader)
	ctx := context.Background()
	s, err := signer.SmartSigner(ctx)
	if err != nil {
		return err
	}

	jwt, err := s.SignJwt(ctx, string(b))
	if err != nil {
		return err
	}
	fmt.Println(jwt)
	return nil
}
