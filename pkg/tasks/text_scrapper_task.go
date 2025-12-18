package tasks

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"

	hatchet "github.com/hatchet-dev/hatchet/sdks/go"
)

type ScrapperInput struct {
	URL     string `json:"url"`
	Element string `json:"element"` // selector ou xpath
	By      string `json:"by"`      // id | class | tag | css | xpath
}

type ScrapperOutput struct {
	Value string `json:"value"`
}

func TextScrapperTask(
	ctx hatchet.Context,
	input ScrapperInput,
) (ScrapperOutput, error) {

	log.Println("Starting web text scrapper workflow...")
	log.Printf(
		"URL: %s | Element: %s | By: %s",
		input.URL,
		input.Element,
		input.By,
	)

	// HTTP Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, input.URL, nil)
	if err != nil {
		return ScrapperOutput{}, err
	}

	req.Header.Set("User-Agent", "Hatchet-Scrapper/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ScrapperOutput{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return ScrapperOutput{}, fmt.Errorf("http error: %d", resp.StatusCode)
	}

	// By XPath
	if strings.ToLower(input.By) == "xpath" {
		doc, err := htmlquery.Parse(resp.Body)
		if err != nil {
			return ScrapperOutput{}, err
		}

		node := htmlquery.FindOne(doc, input.Element)
		if node == nil {
			return ScrapperOutput{}, fmt.Errorf("xpath not found")
		}

		value := strings.TrimSpace(htmlquery.InnerText(node))
		log.Printf("Scrapped text value (xpath): %s", value)

		return ScrapperOutput{Value: value}, nil
	}

	// CSS Selector
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ScrapperOutput{}, err
	}

	var selection *goquery.Selection

	switch strings.ToLower(input.By) {
	case "id":
		selection = doc.Find("#" + input.Element)
	case "class":
		selection = doc.Find("." + input.Element)
	case "tag", "css":
		selection = doc.Find(input.Element)
	default:
		return ScrapperOutput{}, fmt.Errorf("invalid selector type: %s", input.By)
	}

	if selection.Length() == 0 {
		return ScrapperOutput{}, fmt.Errorf("element not found")
	}

	value := strings.TrimSpace(selection.First().Text())
	log.Printf("Scrapped text value (css): %s", value)

	return ScrapperOutput{Value: value}, nil
}
