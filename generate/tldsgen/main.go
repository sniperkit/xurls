/* Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc> */
/* See LICENSE for licensing information */

package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"text/template"
)

const path = "tlds.go"

var tldsTmpl = template.Must(template.New("tlds").Parse(`// Generated by tldsgen

package xurls

// TLDs is a sorted list of all public top-level domains
//
// Sources:
//  * https://data.iana.org/TLD/tlds-alpha-by-domain.txt
//  * https://publicsuffix.org/list/effective_tld_names.dat
var TLDs = []string{
{{range $_, $value := .TLDs}}` + "\t`" + `{{$value}}` + "`" + `,
{{end}}}
`))

func cleanTld(tld string) string {
	tld = strings.ToLower(tld)
	if strings.HasPrefix(tld, "xn--") {
		return ""
	}
	return tld
}

func fromURL(url, pat string) {
	log.Printf("Fetching %s", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetcihng %s: %v", url, err)
		return
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	re := regexp.MustCompile(pat)
	for scanner.Scan() {
		line := scanner.Text()
		tld := re.FindString(line)
		tld = cleanTld(tld)
		if tld == "" {
			continue
		}
		tldChan <- tld
	}
	wg.Done()
}

var (
	wg      sync.WaitGroup
	tldChan chan string
)

func tldList() ([]string, error) {

	tldChan = make(chan string)
	wg.Add(2)

	go fromURL("https://data.iana.org/TLD/tlds-alpha-by-domain.txt", `^[^#]+$`)
	go fromURL("https://publicsuffix.org/list/effective_tld_names.dat", `^[^/.]+$`)

	tlds := make(map[string]struct{})
	go func() {
		for tld := range tldChan {
			tlds[tld] = struct{}{}
		}
	}()
	wg.Wait()

	list := make([]string, 0, len(tlds))
	for tld := range tlds {
		list = append(list, tld)
	}

	sort.Strings(list)
	return list, nil
}

func writeTlds(tlds []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return tldsTmpl.Execute(f, struct {
		TLDs []string
	}{
		TLDs: tlds,
	})
}

func main() {
	tlds, err := tldList()
	if err != nil {
		log.Fatalf("Could not get TLD list: %s", err)
	}
	log.Printf("Generating %s...", path)
	if err := writeTlds(tlds); err != nil {
		log.Fatalf("Could not write path: %s", err)
	}
}
