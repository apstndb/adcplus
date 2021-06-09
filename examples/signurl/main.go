package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"

	"cloud.google.com/go/storage"
	"github.com/apstndb/signer"
)

func main() {
	if err := _main(); err != nil {
		log.Fatalln(err)
	}
}
func _main() error {
	method := flag.String("m", "GET", "")
	duration := flag.Duration("d", 1*time.Hour, "")
	flag.Parse()
	ctx := context.Background()
	s, err := signer.SmartSigner(ctx)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`^gs://(?P<bucket>[^/]*)/(?P<name>.*)$`)
	fmt.Println(flag.Args()[0])
	submatch := re.FindStringSubmatch(flag.Args()[0])
	bucket := submatch[re.SubexpIndex("bucket")]
	name := submatch[re.SubexpIndex("name")]
	fmt.Printf("bucket: %s, name: %s\n", bucket, name)
	url, err := storage.SignedURL(bucket, name, &storage.SignedURLOptions{
		GoogleAccessID: s.ServiceAccount(ctx),
		SignBytes:      s.Signer(ctx),
		Method:         *method,
		Expires:        time.Now().Add(*duration),
	})
	if err != nil {
		return err
	}
	fmt.Println(url)
	return nil
}
