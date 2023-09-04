package scraping

import (
	"fmt"

	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
)

func ScrapingInSite(urlSite string) {
	//fmt.Println(urlSite)
	c := colly.NewCollector(colly.AllowedDomains("www.tibia.com", "tibia.com"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	c.OnHTML("div.TableContainer", func(h *colly.HTMLElement) {
		selection := h.DOM

		childNodes := selection.Children().Nodes

		caption := selection.Find("div.Text").Text()

		//fmt.Printf("childNodes count: \n", len(childNodes))
		switch caption {
		case "Could not find character":
			break
		case "Character Information":
			readAccountInformation(childNodes, selection)
		case "Account Information":
		case "Account Badges":
		case "Account Achievements":
		case "Character Deaths":
		case "Characters":
		}
		//fmt.Println("...End...Child...Node...")
	})

	c.Visit(urlSite)
}

func cleanDesc(s string) string {
	return strings.TrimSpace(s)
}

func readAccountInformation(childNodes []*html.Node, selection *goquery.Selection) {
	//fmt.Printf("childNodes count: \n", len(childNodes))
	for i := range childNodes {
		fmt.Printf("\n")
		fmt.Println(i)

		//description := selection.Find("td.LabelV175").Text()
		value := selection.FindNodes(childNodes[i]).Text()
		//fmt.Println(description)
		fmt.Println(value)
	}
}
