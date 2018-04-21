package handlers

import (
	"encoding/json"
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
const contentTypeTextHTML = "text/html"

const headerContentType = "Content-Type"
const headerAccessControlAllowOrigin = "Access-Control-Allow-Origin"

// return an absolute path string
func mergeURL(pageURL string, URL string) string {
	u, _ := url.Parse(URL)
	base, _ := url.Parse(pageURL)

	return fmt.Sprintf("%s", base.ResolveReference(u))
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Add(headerContentType, contentTypeJSON)
	w.Header().Add(headerAccessControlAllowOrigin, "*")

	URL := r.URL.Query().Get("url")

	fmt.Println("current url :" + URL)
	if URL == "" {
		http.Error(w, "Bad Request, no parameter key 'url' found", http.StatusBadRequest)
		return
	}

	htmlStream, err := fetchHTML(URL)

	if err != nil {
		http.Error(w, fmt.Sprintf("error fetching HTML: %s\n", err), http.StatusBadRequest)
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
		http.Error(w, fmt.Sprintf("error extracting summary: %s", err), http.StatusBadRequest)
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
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, contentTypeTextHTML) {
		return nil, fmt.Errorf("response content type was %s, not text/html", contentType)
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

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			//get the token
			token := tokenizer.Token()

			switch token.Data {

			case "meta":
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

				// filling info for Type,URL,Title,SiteName,Description,Preview Image
				switch prop {
				case "og:type":
					pageSummary.Type = content

				case "og:url":
					pageSummary.URL = content

				case "og:title":
					pageSummary.Title = content

				case "og:site_name":
					pageSummary.SiteName = content

				case "og:description":
					pageSummary.Description = content
				case "og:image", "og:image:url":
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
				case "description":
					if pageSummary.Description == "" {
						pageSummary.Description = content
					}
				case "author":
					pageSummary.Author = content
				case "keywords":
					arr := regexp.MustCompile(",\\s*")
					pageSummary.Keywords = arr.Split(content, -1)
				}

			case "link":
				var rel string
				var href string
				var typ string
				var sizes string
				for _, attr := range token.Attr {
					if attr.Key == "rel" {
						rel = attr.Val
					} else if attr.Key == "href" {
						href = attr.Val
					} else if attr.Key == "type" {
						typ = attr.Val
					} else if attr.Key == "sizes" {
						sizes = attr.Val
					}
				}

				if strings.Contains(rel, "icon") {
					// create new icon instance
					icon := &PreviewImage{
						URL:  mergeURL(pageURL, href),
						Type: typ,
					}

					index := strings.Index(sizes, "x")
					// if there is a valid size statement
					if index != -1 {
						width, _ := strconv.Atoi(sizes[index+1:])
						height, _ := strconv.Atoi(sizes[:index])
						icon.Width = width
						icon.Height = height
					}
					pageSummary.Icon = icon
				}

			case "title":
				if pageSummary.Title == "" {
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
}
