package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type ArticleMetadata struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Uri    string `json:"uri"`
	Date   string `json:"date"`
}

const (
	EXT_HTML          string = ".html"
	EXT_MD            string = ".md"
	DirPrefix         string = "../home/public/blog_contents"
	BookPath          string = "blog_contents/book_think"
	DevPath           string = "blog_contents/dev_record"
	PathBookThinkBlog string = "../home/src/components/BookThinkBlog.vue"
	PathDevRecordBlog string = "../home/src/components/DevRecordBlog.vue"
)

var bookArticles []string
var devArticles []string

func execRmOldHtml(path string) {
	if (len(path) > len(EXT_HTML)) &&
		path[len(path)-len(EXT_HTML):] == EXT_HTML {

		cmd := exec.Command("rm", path)
		// cmd := exec.Command("pwd")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("pandoc output: %q\n", out.String())
		log.Printf("deleted html document: %q\n", path)
	}
}

func execPandocMdToHtml(path string) {
	// Length of a file have to be logner than extension
	if (len(path) > len(EXT_MD)) &&
		path[len(path)-len(EXT_MD):] == EXT_MD {

		result := path[0:len(path)-len(EXT_MD)] + EXT_HTML

		cmd := exec.Command(
			"pandoc",
			path,
			"-o",
			result,
			"--to",
			"html5",
			"--no-highlight",
		)
		// cmd := exec.Command("pwd")
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
		//log.Printf("pandoc output: %q\n", out.String())

		if strings.Contains(path, BookPath) {
			bookArticles = append(bookArticles, path[0:len(path)-len(EXT_MD)]+EXT_HTML)
			log.Printf("generated document: %q\n", result)
		} else if strings.Contains(path, DevPath) {
			devArticles = append(devArticles, path[0:len(path)-len(EXT_MD)]+EXT_HTML)
			log.Printf("generated document: %q\n", result)
		}
	}
}

func iterate(path string, dirInfo []os.FileInfo) {
	for _, info := range dirInfo {
		if info.Name()[0:2] == "aA" {
			continue
		}

		absolutePathName := path + "/" + info.Name()

		if info.IsDir() {
			subDirInfo, err := ioutil.ReadDir(absolutePathName)
			if err != nil {
				log.Println("(iterate) ", err)
				log.Fatal(err)
			}
			iterate(absolutePathName, subDirInfo)
		} else {
			execRmOldHtml(absolutePathName)
			execPandocMdToHtml(absolutePathName)
		}
		fmt.Println(info.Name())
	}
}

func GetArticleMetadata(htmlPaths []string) []ArticleMetadata {
	r := regexp.MustCompile("(19|20)[0-9][0-9][- /.](0[1-9]|1[012])[- /.](0[1-9]|[12][0-9]|3[01])")

	var metadata []ArticleMetadata
	var uri string

	for _, htmlPath := range htmlPaths {
		htmlSource, err := ioutil.ReadFile(htmlPath)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Get title
		titleFrom := bytes.Index(htmlSource, []byte(">"))
		titleFrom += len(">")
		titleTo := bytes.Index(htmlSource, []byte("</h1>"))

		// Get author, get first ">" after title
		authorFromTemp := titleTo + len("</h1>")
		authorFrom := bytes.Index(htmlSource[authorFromTemp:], []byte(">"))
		authorFrom += len(">")
		authorFrom += authorFromTemp
		authorTo := bytes.Index(htmlSource, []byte("</h2>"))
		fmt.Println(authorTo == -1)
		// Get Uri
		if strings.Contains(htmlPath, BookPath) {
			uri = strings.Replace(htmlPath, "../home/public/blog_contents/book_think/", "/bookThinkBlog/", 1)
			uri = strings.Replace(uri, ".html", "/", 1)
		} else if strings.Contains(htmlPath, DevPath) {
			uri = strings.Replace(htmlPath, "../home/public/blog_contents/dev_record/", "/devRecordBlog/", 1)
			uri = strings.Replace(uri, ".html", "/", 1)
		}

		// Get date
		date := r.FindString(htmlPath)

		// if authorTo is -1, it's devRecordBlog
		if authorTo == -1 {
			metadata = append(metadata, ArticleMetadata{
				Title:  string(htmlSource[titleFrom:titleTo]),
				Author: "myself",
				Uri:    uri,
				Date:   date,
			})
		} else {
			metadata = append(metadata, ArticleMetadata{
				Title:  string(htmlSource[titleFrom:titleTo]),
				Author: string(htmlSource[authorFrom:authorTo]),
				Uri:    uri,
				Date:   date,
			})
		}
	}
	return metadata
}

func insertJsonToBlogComponent(json string, PathBlogComponent string) {
	componentSourceBytes, err := ioutil.ReadFile(PathBlogComponent)
	if err != nil {
		log.Fatal(err)
	}
	componentSource := string(componentSourceBytes)

	i := strings.Index(componentSource, "__INSERTION_POSITION__")

	from := strings.Index(componentSource[i:], "[")
	from += i

	k := strings.Index(componentSource, "__INSERTION_POSITION_END__")
	to := strings.LastIndex(componentSource[i:k+2], "]")
	to += i

	fmt.Println(" -- old json -- ")
	fmt.Println(string(componentSource[from : to+1]))
	fmt.Println("-- old json end -- ")

	componentSource = componentSource[0:from] + json + componentSource[to+1:]

	ioutil.WriteFile(PathBlogComponent, []byte(componentSource), 0644)
}

func main() {
	dirInfo, err := ioutil.ReadDir(DirPrefix)
	if err != nil {
		log.Println("(main) ", err)
		log.Fatal(err)
	}
	iterate(DirPrefix, dirInfo)

	fmt.Println("\n\n bookArticles \n\n")
	for _, bookArticle := range bookArticles {
		fmt.Println(bookArticle)
	}

	fmt.Println("\n\n devArticles \n\n")
	for _, devArticle := range devArticles {
		fmt.Println(devArticle)
	}

	// Reverse articles' order
	i := 0
	k := len(bookArticles) - 1
	for {
		bookArticles[i], bookArticles[k] = bookArticles[k], bookArticles[i]
		i++
		k--
		if i >= k {
			break
		}
	}

	i = 0
	k = len(devArticles) - 1
	for {
		devArticles[i], devArticles[k] = devArticles[k], devArticles[i]
		i++
		k--
		if i >= k {
			break
		}
	}

	bookMetadata := GetArticleMetadata(bookArticles)
	devMetadata := GetArticleMetadata(devArticles)

	bb, err := json.Marshal(bookMetadata)
	if err != nil {
		log.Fatal(err)
	}

	db, err := json.Marshal(devMetadata)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\n book json \n\n")
	fmt.Println(string(bb))

	fmt.Println("\n\n dev json \n\n")
	fmt.Println(string(db))
	insertJsonToBlogComponent(string(bb), string(PathBookThinkBlog))
	insertJsonToBlogComponent(string(db), string(PathDevRecordBlog))
}
