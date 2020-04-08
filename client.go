package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as
	// defined in the function definition
	return
}

// Call shell command wget to download. The reason to use wget is that wget
// supports automatically resume download. So this package only runs on Linux
// systems.
func wget(url, filepath string) error {
	// run shell `wget URL -O filepath`
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Extract all http** links from a given webpage
func crawl(parentUrl string, filePath string) {
	resp, err := http.Get(parentUrl)
	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", parentUrl)
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function completes

	z := html.NewTokenizer(b)
	var urls []string

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the parentUrl begines in http**

			urls = append(urls, parentUrl+url)

			iDir := strings.HasSuffix(url, "/")

			if iDir {
				fmt.Print("recurring ", parentUrl+url, "\n")
				crawl(parentUrl+url, filePath+url)
			} else {
				wget(parentUrl+url, filePath+url)
			}

			for i := 0; i < len(urls); i++ {
				fmt.Println("url ", i, " ", urls[i])
			}

		}
	}
}

func main() {
	seedUrls := os.Args[1:]

	chUrls := make(chan string)

	for _, url := range seedUrls {
		crawl(url, "")
	}

	close(chUrls)
}
