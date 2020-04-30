package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/LUSHDigital/core-lush/currency/internal/cmd/scaffold"
)

const (
	templateFile = "cmd/std.txt" // must include the cmd prefix because this code is called from the Makefile
	outputFile   = "../std.go"   // we want this to output in the top directory

	// For the source on this, please check:
	// - International Organization for Standardization: https://www.iso.org/iso-4217-currency-codes.html
	// - Currency Code Services – ISO 4217 Maintenance Agency: https://www.currency-iso.org
	isoStdDownload = "https://www.currency-iso.org/dam/downloads/lists/list_one.xml"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	err, iso := getLatestISO4217()
	if err != nil {
		log.Fatalf("could not get latest iso: %v", err)
	}

	currencies := buildCurrencyList(iso)
	if len(currencies) == 0 {
		log.Fatalf("could not build currency list")
	}

	for _, gen := range generators {
		gen(currencies)
	}
}
func getLatestISO4217() (err error, iso scaffold.ISO4217) {
	res, err := http.Get(isoStdDownload)
	if err != nil {
		log.Fatal(err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	if err = xml.Unmarshal(b, &iso); err != nil {
		log.Fatal(err)
	}
	return err, iso
}

type currency struct {
	Code   string
	Units  int
	Factor string
	Name   string
}

func buildCurrencyList(iso scaffold.ISO4217) []currency {
	var currencies []currency
	for _, entry := range iso.Table.Entries {
		if entry.Code == "" {
			continue
		}

		var (
			unit int
			err  error
		)
		unit, err = strconv.Atoi(entry.MinorUnits)
		if err != nil {
			// nothing really
			// it's always because of "N.A."
			// But just in case...
			if entry.MinorUnits != "N.A." {
				log.Fatalln(err)
			}
		}
		if func() bool {
			for _, cur := range currencies {
				if cur.Code == entry.Code {
					return true
				}
			}
			return false
		}() {
			continue
		}

		currencies = append(currencies, currency{
			Code:   entry.Code,
			Units:  unit,
			Factor: fmt.Sprintf("1%s", strings.Repeat("0", unit)),
			Name:   entry.Description,
		})
	}
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Code < currencies[j].Code
	})
	return currencies
}

type generatorFunc func(currencies []currency)

var generators = []generatorFunc{
	generateGoPackage,
}

func generateGoPackage(currencies []currency) {
	tpl, err := ioutil.ReadFile(templateFile)
	if err != nil {
		log.Fatalf("cannot open template file: %v", err)
	}

	t := template.Must(template.New("go").Parse(string(tpl)))
	buf := new(bytes.Buffer)
	err = t.Execute(buf, currencies)
	if err != nil {
		log.Fatal(err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	buf = bytes.NewBuffer(formatted)
	to, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = to.Close(); err != nil {
			log.Println(err)
		}
	}()

	_, err = io.Copy(to, buf)
	if err != nil {
		log.Fatal(err)
	}
}
