package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
	"github.com/apstndb/adcplus"
	"github.com/apstndb/adcplus/signer"
)

func main() {
	if err := _main(); err != nil {
		log.Fatalln(err)
	}
}

func _main() error {
	method := flag.String("m", "GET", "HTTP method")
	duration := flag.Duration("d", 1*time.Hour, "duration")
	impersonateServiceAccount := flag.String("i", "", "impersonate service account")
	jsonKeyFile := flag.String("k", "", "service account key file")
	useSignatureVersion4 := flag.Bool("4", false, "Use signature version 4")
	flag.Parse()
	ctx := context.Background()
	s, err := signer.SmartSigner(ctx, adcplus.WithTargetPrincipal(*impersonateServiceAccount), adcplus.WithCredentialsFile(*jsonKeyFile))
	if err != nil {
		return err
	}
	var scheme storage.SigningScheme
	if *useSignatureVersion4 {
		scheme = storage.SigningSchemeV4
	}

	re := regexp.MustCompile(`^gs://(?P<bucket>[^/]*)/(?P<name>.*)$`)
	fmt.Println(flag.Args()[0])
	submatch := re.FindStringSubmatch(flag.Args()[0])
	bucket := submatch[subexpIndex(re, "bucket")]
	name := submatch[subexpIndex(re, "name")]
	fmt.Printf("bucket: %s, name: %s\n", bucket, name)
	url, err := storage.SignedURL(bucket, name, &storage.SignedURLOptions{
		GoogleAccessID: s.ServiceAccount(ctx),
		SignBytes:      signer.SignWithoutKeyAdaptor(ctx, s),
		Method:         *method,
		Scheme: 		scheme,
		Expires:        time.Now().Add(*duration),
	})
	if err != nil {
		return err
	}
	fmt.Println(url)
	return nil
}

func subexpIndex(re *regexp.Regexp, name string) int {
	if name != "" {
		for i, s := range re.SubexpNames() {
			if name == s {
				return i
			}
		}
	}
	return -1
}
