package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
)

func main() {
	f := flag.String("img", "../test_data/test2.jpg", "input image_file")
	flag.Parse()
	filepath := *f
	rawExif, err := exif.SearchFileAndExtractExif(filepath)
	if err != nil {
		log.Fatalf("1: %s", err)
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		log.Fatalf("2: %s", err)
	}

	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	if err != nil {
		log.Fatalf("3: %s", err)
	}
	// fmt.Print(index.RootIfd.String())
	tagName := "FileName"

	results, err := index.Tree[2].FindTagWithName(tagName)
	if err != nil {
		log.Fatalf("4: %s", err)
	}
	// fmt.Print(results)

	ite := results[0]
	valueRaw, err := ite.Value()
	if err != nil {
		log.Fatalf("5: %s", err)
	}
	fmt.Print(valueRaw)

	// value := valueRaw.(string)
	// fmt.Println(value)
	// gi, err := ifd.GpsInfo()
	// if err != nil {
	// 	log.Fatalf("5: %s", err)
	// }
	// fmt.Printf("%s\n", gi)
	// fmt.Printf("%f\n", gi.Latitude.Decimal())
}
