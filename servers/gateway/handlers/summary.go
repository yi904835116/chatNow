package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

const contentTypeJSON = "application/json"
const contentTypeText = "text/plain"

const headerContentType = "Content-Type"
const headerAccessControlAllowOrigin = "Access-Control-Allow-Origin"

//PageSummarySlice is a slice of *PageSummary,
//that is, pointers to PageSummary struct
type PageSummarySlice []*PageSummary

const commonPrefix = "og:"

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.

	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	w.Header().Add(headerContentType, contentTypeJSON)
	w.Header().Add(headerAccessControlAllowOrigin, "*")

	URL := r.URL.Query().Get("url")
	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	if URL == "" {
		http.Error(w, "Bad Request, no parameter key 'url' found", http.StatusBadRequest)
		return
	}

	htmlStream, err := fetchHTML(URL)

	if err != nil {
		http.Error(w, "error fetching HTML", http.StatusBadRequest)
		return
	}

	//make sure the response body gets closed
	defer htmlStream.Close()
	//call getPageSummary() passing the requested URL
	//and holding on to the returned openGraphProps map
	page, err := extractSummary(URL, htmlStream)

	//if get back an error, respond to the client
	//with that error and an http.StatusBadRequest code
	if err != nil {
		http.Error(w, fmt.Sprintf("error extracting summary: %v", err), http.StatusBadRequest)
		return
	}

	//otherwise, respond by writing the openGrahProps
	//map as a JSON-encoded object
	encoder := json.NewEncoder(w)
	encoder.Encode(page)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.

	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML

	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/
	resp, err := http.Get(pageURL)

	if err != nil {
		return nil, fmt.Errorf("fetchHTML failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("fetchHTML failed, status code: %v", err)
	}

	//check if the response's Content-Type header
	//starts with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, errors.New("response content type was " + ctype + " not text/htm")
	}

	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {

	pageSummary := &PageSummary{}
	previewImages := []*PreviewImage{}
	previewImage := &PreviewImage{}

	// Create a new tokenizer instance for http response body
	tokenizer := html.NewTokenizer(htmlStream)

	for {
		//get the next token type
		tokenType := tokenizer.Next()

		//if it's an error token, we either reached
		//the end of the file, or the HTML was malformed
		if tokenType == html.ErrorToken {
			return pageSummary, tokenizer.Err()
		}

		// stop tokenizing after you encounter the </head> tag
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				return pageSummary, nil
			}
		}

		//process the token according to the token type...
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			//get the token
			token := tokenizer.Token()

			if token.Data == "meta" {
				var prop string
				var content string
				var name string
				for _, attr := range token.Attr {
					if attr.Key == "property" {
						prop = attr.Val
					} else if attr.Key == "content" {
						content = attr.Val
					} else if attr.Key == "name" {
						name = attr.Val
					}
				}

				fmt.Println("prop" + prop)
				fmt.Println("content" + content)
				fmt.Println("name" + name)
				// filling info for Type,URL,Title,SiteName,Description,Preview Image
				switch prop {
				case "og:type":
					pageSummary.Type = content

				case "og:url":
					pageSummary.URL = content

				case "og:title":
					// twitterTitlePriority = false
					pageSummary.Title = content

				case "og:site_name":
					pageSummary.SiteName = content

				case "og:description":
					if pageSummary.Description == "" {
						pageSummary.Description = content
					}
					// Preview images.
					// og:image or og:iamge:url indicates this is a new preview image.
				case "og:image", "og:image:url":
					// Create a new instance of PreviewImage.
					previewImage = &PreviewImage{}

					previewImage.URL = mergeURL(pageURL, content)
					previewImages = append(previewImages, previewImage)
					pageSummary.Images = previewImages

				case "og:image:secure_url":
					previewImage.SecureURL = mergeURL(pageURL, content)

				case "og:image:type":
					previewImage.Type = content

				case "og:image:width", "og:image:height":
					size, _ := strconv.Atoi(content)

					if prop == "og:image:width" {
						previewImage.Width = size
					} else {
						previewImage.Height = size
					}

				case "og:image:alt":
					previewImage.Alt = content
				}

				// filling info forDescription,Author,Keywords
				switch name {
				case "og:description":
					pageSummary.Description = name
				case "author":
					pageSummary.Author = name
				case "keywords":
					pageSummary.Keywords = regexp.MustCompile(",\\s*").Split(content, -1)
				}

			}

			if token.Data == "link" {
				var ref string
				var href string
				var typ string
				var size string
				for _, attr := range token.Attr {
					if attr.Key == "ref" {
						ref = attr.Val
					} else if attr.Key == "href" {
						href = attr.Val
					} else if attr.Key == "type" {
						typ = attr.Val
					} else if attr.Key == "size" {
						size = attr.Val
					}
				}

				if ref == "icon" {
					icon := &PreviewImage{
						URL:  mergeURL(pageURL, href),
						Type: typ,
					}
					arr := strings.Split(size, "x")
					h, _ := strconv.Atoi(arr[0])
					w, _ := strconv.Atoi(arr[1])

					icon.Height = h
					icon.Width = w

					pageSummary.Icon = icon
				}
			}

			//if the name of the element is "title"
			if token.Data == "title" && pageSummary.Title == "" {
				// //the next token should be the page title
				tokenType = tokenizer.Next()

				// Just make sure it is actually a text token.
				if tokenType == html.TextToken {
					title := tokenizer.Token().Data
					if pageSummary.Title == "" {
						pageSummary.Title = title
					}
				}
			}

		}

	}

}

func mergeURL(pageURL string, URL string) string {
	u, _ := url.Parse(URL)
	base, _ := url.Parse(pageURL)

	return fmt.Sprintf("%s", base.ResolveReference(u))
}
