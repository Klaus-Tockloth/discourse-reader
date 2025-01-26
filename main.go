/*
Purpose:
- Discourse Reader

Description:
- Retrieves data (e.g. site, category, topic) from Discourse forum.

Releases:
- v1.0.0 - 2022/11/18: initial release
- v1.0.1 - 2025/01/24: compiled with go v1.23.5

Author:
- Klaus Tockloth

Copyright:
- Copyright (c) 2022-2025 Klaus Tockloth

Contact:
- klaus.tockloth@googlemail.com

Remarks:
- Lint: golangci-lint run --no-config --enable gocritic
- Vulnerability detection: govulncheck ./...

ToDo:
- NN

Links:
- https://docs.discourse.org/
- https://meta.discourse.org/t/available-settings-for-global-rate-limits-and-throttling/78612
- https://meta.discourse.org/t/api-can-pull-only-20-posts/163406/5
*/

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// general program info
var (
	progName    = filepath.Base(os.Args[0])
	progVersion = "v1.0.1"
	progDate    = "2025/01/24"
	progPurpose = "Discourse Reader"
	progInfo    = "Retrieves data (e.g. site, category, topic) from Discourse forum."
	userAgent   = progName + "/" + progVersion
)

// httpClient represents HTTP client for communication with Discourse
var httpClient *http.Client

// command line parameters
var (
	forum      *string
	category   *int
	topic      *int
	pages      *int
	query      *string
	output     *string
	userapikey *string
	sleeptime  *int
)

// (optional) environment variables
const (
	evUSERAPIKEY = "USER_API_KEY"
	evHTTPSPROXY = "HTTPS_PROXY"
)

/*
init initializes this program.
*/
func init() {
	// initialize standard logger
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

/*
main starts this program.
*/
func main() {
	fmt.Printf("Program:\n")
	fmt.Printf("  Name    : %s\n", progName)
	fmt.Printf("  Release : %s - %s\n", progVersion, progDate)
	fmt.Printf("  Purpose : %s\n", progPurpose)
	fmt.Printf("  Info    : %s\n\n", progInfo)

	forum = flag.String("forum", "", "Discourse forum URL")
	category = flag.Int("category", -1, "retrieve data (list of topics) for category with identifier")
	topic = flag.Int("topic", -1, "retrieve data (list of posts) for topic with identifier")
	pages = flag.Int("pages", 19, "pages of data to retrieve")
	query = flag.String("query", "", "general data retrieve query (full URL)")
	output = flag.String("output", "", "name of JSON output file")
	userapikey = flag.String("userapikey", "", fmt.Sprintf("personal user API key (can also be set as environment var '%s')", evUSERAPIKEY))
	sleeptime = flag.Int("sleeptime", 2, "sleep time in seconds before retrieving the next page (avoids user rate limiting)")

	flag.Usage = printUsage
	flag.Parse()
	if flag.NFlag() == 0 {
		printUsage()
	}
	if *output == "" {
		log.Fatalf("Error: Option '-output=string' required.\n")
	}
	if *userapikey == "" {
		*userapikey = os.Getenv(evUSERAPIKEY)
		if *userapikey == "" {
			log.Fatalf("Error: User-API-Key not found (neither as option '-userapikey=string' nor as environment variable '%s').\n", evUSERAPIKEY)
		}
	}
	if *sleeptime < 0 {
		log.Fatalf("Error: Option '-sleeptime=int' must be >= 0.\n")
	}

	// create HTTP transport object
	httpTransport := &http.Transport{}
	httpTransport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}

	// get internet proxy from environment
	internetProxy := os.Getenv(evHTTPSPROXY)
	if internetProxy != "" {
		internetProxyURL, err := url.Parse(internetProxy)
		if err != nil {
			log.Fatalf("Error: Could not parse internet proxy URL from environment, error=[%v], proxy=[%v].", err, internetProxy)
		}
		httpTransport.Proxy = http.ProxyURL(internetProxyURL)
	}

	// create HTTPS client
	httpClient = &http.Client{
		Transport: httpTransport,
		Timeout:   time.Second * time.Duration(30),
	}

	switch {
	case *query != "":
		*query = "https://" + *query
		fmt.Printf("Requesting data for query %s ...\n", *query)
		queryData, err := getQueryData()
		if err != nil {
			_ = os.WriteFile(*output, queryData, 0666)
			log.Fatalf("Error: Error requesting query data, error=[%s].", err)
		}
		fmt.Printf("Writing data to file %v ...\n", *output)
		err = os.WriteFile(*output, queryData, 0666)
		if err != nil {
			log.Fatalf("Error: Could not write output file, error=[%v].", err)
		}

	case *category > 0:
		if *forum == "" {
			log.Fatalf("Error: Option '-forum=string' required.\n")
		}
		fmt.Printf("Requesting data (list of topics) for category %d ...\n", *category)
		categoryData, err := getCategoryData()
		if err != nil {
			_ = os.WriteFile(*output, []byte(categoryData), 0666)
			log.Fatalf("Error: Error requesting category data, error=[%s].", err)
		}
		fmt.Printf("Writing data to file %v ...\n", *output)
		err = os.WriteFile(*output, []byte(categoryData), 0666)
		if err != nil {
			log.Fatalf("Error: Could not write output file, error=[%v].", err)
		}

	case *topic > 0:
		if *forum == "" {
			log.Fatalf("Error: Option '-forum=string' required.\n")
		}
		fmt.Printf("Request data (list of posts) for topic %d ...\n", *topic)
		topicData, err := getTopicData()
		if err != nil {
			_ = os.WriteFile(*output, []byte(topicData), 0666)
			log.Fatalf("Error: Error requesting topic data, error=[%s].", err)
		}
		fmt.Printf("Writing data to file %v ...\n", *output)
		err = os.WriteFile(*output, []byte(topicData), 0666)
		if err != nil {
			log.Fatalf("Error: Could not write output file, error=[%v].", err)
		}

	default:
		log.Fatalf("Error: Something went wrong, nothing to do.")
	}

	fmt.Printf("Done.\n")
}

/*
printUsage prints the usage of this program.
*/
func printUsage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s -forum=string -query=string -category=int -topic=int -pages=int -output=string -userapikey -sleeptime=int\n", os.Args[0])

	fmt.Printf("\nExamples for general query:\n")
	fmt.Printf("  %s\n", os.Args[0])
	fmt.Printf("  %s -query=community.openstreetmap.org/site.json -output=community.openstreetmap.org.json\n", os.Args[0])
	fmt.Printf("  %s -query=community.openstreetmap.org/site.json -output=community.openstreetmap.org.json -userapikey=bd38603815e3f2562c3eb3988c69eb77\n", os.Args[0])
	fmt.Printf("  %s -query=meta.discourse.org/site.json -output=meta.discourse.org.json\n", os.Args[0])
	fmt.Printf("  %s -query=meta.discourse.org/session/current.json -output=session-current.json\n", os.Args[0])

	fmt.Printf("\nExamples for category:\n")
	fmt.Printf("  %s -forum=community.openstreetmap.org -category=56 -output=category-56.json\n", os.Args[0])
	fmt.Printf("  %s -forum=community.openstreetmap.org -category=56 -output=category-56.json -userapikey=bd38603815e3f2562c3eb3988c69eb77\n", os.Args[0])
	fmt.Printf("  %s -forum=meta.discourse.org -category=67 -pages=99 -sleeptime=6 -output=category-67.json\n", os.Args[0])

	fmt.Printf("\nExamples for topic:\n")
	fmt.Printf("  %s -forum=community.openstreetmap.org -topic=4120 -output=topic-4120.json\n", os.Args[0])
	fmt.Printf("  %s -forum=community.openstreetmap.org -topic=4120 -pages=99 -sleeptime=6 -output=topic-4120.json\n", os.Args[0])
	fmt.Printf("  %s -forum=community.openstreetmap.org -topic=4120 --output=topic-4120.json -userapikey=bd38603815e3f2562c3eb3988c69eb77\n", os.Args[0])
	fmt.Printf("  %s -forum=meta.discourse.org -topic=112837 -output=topic-112837.json\n", os.Args[0])

	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()

	fmt.Printf("\nRemarks:\n")
	fmt.Printf("  - User API key can be set as environment variable [%s].\n", evUSERAPIKEY)
	fmt.Printf("  - Internet proxy can be set as environment variable [%s].\n", evHTTPSPROXY)
	fmt.Printf("  - Examples for Linux:\n")
	fmt.Printf("    export %s=bd38603815e3f2562c3eb3988c69eb77\n", evUSERAPIKEY)
	fmt.Printf("    export %s=http://user:password@194.114.63.23:8080\n", evHTTPSPROXY)
	fmt.Printf("  - Examples for Windows:\n")
	fmt.Printf("    set %s=bd38603815e3f2562c3eb3988c69eb77\n", evUSERAPIKEY)
	fmt.Printf("    set %s=http://user:password@194.114.63.23:8080\n", evHTTPSPROXY)

	fmt.Printf("\nRate limiting by forum service:\n")
	fmt.Printf("  - This program does functionally no different than a user via a browser. However, the\n" +
		"    data is retrieved somewhat faster. This can lead to rejections (rate limiting) by the\n" +
		"    service. To prevent this, the program can pause between fetching pages. The pause time\n" +
		"    can be specified with the option '-sleeptime=int'.\n")

	fmt.Printf("  - Typical user rate limit settings are:\n" +
		"    - requests per minute : 20\n" +
		"    - requests per day    : 2880\n")

	fmt.Printf("\n")
	os.Exit(1)
}
