package main

import (
	"encoding/csv"
	"encoding/json"
	"os"

    "fmt"
	"strings"

	"github.com/gocolly/colly"
)


type Post struct {
    Title string
    URL   string
    Score string
    User  string
}


func main() {
	
	visitedURLs := make(map[string]bool)

	
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 5})

	
	var posts []Post


	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		fmt.Println("User-Agent:", r.Headers.Get("User-Agent"))
	})

	
	c.OnHTML("tr.athing", func(h *colly.HTMLElement) {
		post := Post{}
		post.Title = h.ChildText("td.title > span.titleline > a")
		post.URL = h.ChildAttr("td.title > span.titleline > a", "href")

		
		subtext := h.DOM.Next()
		post.Score = strings.TrimSpace(subtext.Find("span.score").Text())
		post.User = strings.TrimSpace(subtext.Find("a.hnuser").Text())

		
		fmt.Printf("Title: %s\nURL: %s\nScore: %s\nUser: %s\n\n", post.Title, post.URL, post.Score, post.User)

		
		posts = append(posts, post)
	})

	
	c.OnHTML(".morelink", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		if !visitedURLs[nextPage] {
			visitedURLs[nextPage] = true
			c.Visit(nextPage)
		}
	})

	
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	
	c.Visit("https://news.ycombinator.com/")
	
	
	c.Wait()

	
	saveToCSV("results.csv", posts)
	saveToJSON("results.json", posts)
	
}


func saveToCSV(filename string, posts []Post) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	
	err = writer.Write([]string{"Title", "URL", "Score", "User"})
	if err != nil {
		return err
	}

	
	for _, post := range posts {
		err := writer.Write([]string{post.Title, post.URL, post.Score, post.User})
		if err != nil {
			return err
		}
	}

	return nil
}



type JSONPost struct {
    Title string `json:"title"`
    URL   string `json:"url"`
    Score string `json:"score"`
    User  string `json:"user"`
}

func saveToJSON(filename string, posts []Post) error {
    
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(posts)
	if err != nil {
		return err
	}

	return nil
}
