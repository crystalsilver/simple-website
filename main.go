package main

import (
	"bytes"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"github.com/russross/blackfriday"
)

func getLayoutStart(title string) string {
	return `<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link href="https://fonts.googleapis.com/css?family=Roboto:400,400i,500,500i" rel="stylesheet">
			<link href="https://fonts.googleapis.com/css?family=Roboto+Mono:400,400i,500,500i" rel="stylesheet">
			<title>` + title + `</title>
			<style>
				body {
					font-family: 'Roboto', sans-serif;
					line-height: 1.6;
					text-rendering: optimizeLegibility;
				}

				h1, h2, h3 {
					font-weight: 500;
				}

				.container {
					margin: 2.5em auto;
					max-width: 40.625rem;
					padding: 0 0.5rem;
				}

				nav ul {
					list-style-type: none;
					padding: 0;
				}

				nav li {
					margin-bottom: 0.5rem;
				}

				nav li .date {
					margin-right: 1rem;
				}

				.all-posts {
					font-size: 0.889rem;
					margin: 1.25rem 0 2.5rem 0;
				}

				a {
					text-decoration: none;
				}

				a:hover {
					text-decoration: underline;
				}

				.date {
					font-size: 0.889rem;
				}

				.toggle-theme {
					font-size: 0.889rem;
					margin-top: 3.75rem;
					padding-top: 1rem;
				}

				pre {
					overflow: auto;
					padding: 0.25rem 0.75rem;
				}

				code {
					font-family: 'Roboto Mono', monospace;
				}

				body.light {
					background-color: #fafafa;
					color: rgba(0, 0, 0, 0.87);
				}

				body.light a {
					color: #3949ab;
				}

				body.light pre {
					background-color: #ffffff;
					border-left: 3px solid #8c9eff;
				}

				body.light .toggle-theme {
					border-top: 1px solid #e0e0e0;
				}

				body.light .toggle-theme a {
					color: rgba(0, 0, 0, 0.5);
				}

				body.dark {
					background-color: #212121;
					color: rgba(255, 255, 255, 0.7);
				}

				body.dark a {
					color: #8c9eff;
				}

				body.dark pre {
					background-color: #303030;
					border-left: 3px solid #3949ab;
				}

				body.dark .toggle-theme {
					border-top: 1px solid #303030;
				}

				body.dark .toggle-theme a {
					color: rgba(255, 255, 255, 0.5);
				}
			</style>
			<script>
				document.addEventListener('DOMContentLoaded', function(event) {
					if (localStorage.getItem('theme') === 'dark') {
						setDarkTheme();
					} else {
						setLightTheme();
					}
				});

				function toggleTheme(event) {
					event.preventDefault();

					if (document.body.className === 'dark') {
						setLightTheme();
					} else {
						setDarkTheme();
					}
				}

				function setLightTheme() {
						document.body.className = 'light';
						document.getElementsByClassName('toggle-theme')[0].children[0].innerHTML = 'Dark';
						localStorage.setItem('theme', 'light');
				}

				function setDarkTheme() {
						document.body.className = 'dark';
						document.getElementsByClassName('toggle-theme')[0].children[0].innerHTML = 'Light';
						localStorage.setItem('theme', 'dark');
				}
			</script>
		</head>
		<body>
			<div class="container">`
}

func getLayoutEnd() string {
	return `
					<p class="toggle-theme">
						<a href="#" onclick="toggleTheme(event)">Dark</a>
					</p>
			</div>
		</body>
	</html>`
}

func getFile(f string) []byte {
	b, err := ioutil.ReadFile(f)

	if err != nil {
		panic(err)
	}

	return b
}

func getDir(dir string) []os.FileInfo {
	p, err := ioutil.ReadDir(dir)

	if err != nil {
		panic(err)
	}

	return p
}

func writeFile(fileName string, b bytes.Buffer) {
	err := ioutil.WriteFile(fileName+".html", b.Bytes(), 0644)

	if err != nil {
		panic(err)
	}
}

func getSiteTitle() string {
	return strings.Split(string(getFile("_sections/header.md")), "\n")[0][2:]
}

func getPostMeta(fi os.FileInfo) (string, string, string) {
	id := fi.Name()[:len(fi.Name())-3]
	date := fi.Name()[0:10]
	title := strings.Split(string(getFile("_posts/"+fi.Name())), "\n")[0][2:]

	return id, date, title
}

func getPageMeta(fi os.FileInfo) (string, string) {
	id := fi.Name()[:len(fi.Name())-3]
	title := strings.Split(string(getFile("_pages/"+fi.Name())), "\n")[0][2:]

	return id, title
}

func writeIndex() {
	var b bytes.Buffer
	b.WriteString(getLayoutStart(getSiteTitle()))
	b.Write(blackfriday.MarkdownBasic(getFile("_sections/header.md")))
	writePostsSection(&b)
	writePagesSection(&b)
	b.WriteString(getLayoutEnd())
	writeFile("index", b)
}

func writePostsSection(b *bytes.Buffer) {
	b.WriteString("<h2>Posts</h2><nav class=\"posts\"><ul>")

	posts := getDir("_posts")
	limit := int(math.Max(float64(len(posts))-5, 0))

	for i := len(posts) - 1; i >= limit; i-- {
		fileName, date, title := getPostMeta(posts[i])

		b.WriteString("<li><span class=\"date\">" + date +
			"</span><a href=\"posts/" +
			fileName + ".html\">" +
			title + "</a></li>\n")
	}

	b.WriteString("</ul></nav><p class=\"all-posts\"><a href=\"all-posts.html\">All posts</a></p>")
}

func writePagesSection(b *bytes.Buffer) {
	b.WriteString("<h2>Pages</h2><nav class=\"pages\"><ul>")

	pages := getDir("_pages")

	for i := 0; i < len(pages); i++ {
		id, title := getPageMeta(pages[i])

		b.WriteString("<li><a href=\"pages/" +
			id + ".html\">" +
			title + "</a></li>\n")
	}

	b.WriteString("</ul></nav>")
}

func writePosts() {
	posts := getDir("_posts")

	for i := 0; i < len(posts); i++ {
		id, date, title := getPostMeta(posts[i])

		var b bytes.Buffer
		b.WriteString(getLayoutStart(title + " – " + getSiteTitle()))
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.WriteString("<p class=\"date\">" + date + "</p>")
		b.Write(blackfriday.MarkdownBasic(getFile("_posts/" + posts[i].Name())))
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.WriteString(getLayoutEnd())

		writeFile("posts/"+id, b)
	}
}

func writePostsPage() {
	posts := getDir("_posts")
	var b bytes.Buffer

	b.WriteString(getLayoutStart("All posts – " + getSiteTitle()))
	b.WriteString("<p><a href=\"index.html\">←</a></p>")
	b.WriteString("<h1>All posts</h1>")
	b.WriteString("<nav class=\"posts\"><ul>")

	for i := len(posts) - 1; i >= 0; i-- {
		id, date, title := getPostMeta(posts[i])

		b.WriteString("<li><span class=\"date\">" + date +
			"</span><a href=\"posts/" +
			id + ".html\">" +
			title + "</a></li>\n")
	}

	b.WriteString("</ul></nav><p><a href=\"index.html\">←</a></p>")
	b.WriteString(getLayoutEnd())
	writeFile("all-posts", b)
}

func writePages() {
	pages := getDir("_pages")

	for i := 0; i < len(pages); i++ {
		fileName, title := getPageMeta(pages[i])

		var b bytes.Buffer
		b.WriteString(getLayoutStart(title + " – " + getSiteTitle()))
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.Write(blackfriday.MarkdownBasic(getFile("_pages/" + pages[i].Name())))
		b.WriteString("<p><a href=\"../index.html\">←</a></p>")
		b.WriteString(getLayoutEnd())

		writeFile("pages/"+fileName, b)
	}
}

func createFilesAndDirs() {
	os.MkdirAll("_sections", 0755)
	os.MkdirAll("_posts", 0755)
	os.MkdirAll("_pages", 0755)

	if _, err := os.Stat("_sections/header.md"); os.IsNotExist(err) {
		err := ioutil.WriteFile(
			"_sections/header.md",
			[]byte("# Title\n\nDescription"),
			0644)

		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat("posts"); os.IsNotExist(err) {
		err := ioutil.WriteFile(
			"_posts/"+time.Now().Format("2006-01-02")+"-initial-post.md",
			[]byte("# Initial post\n\nThis is the initial post."),
			0644)

		if err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat("pages"); os.IsNotExist(err) {
		err := ioutil.WriteFile(
			"_pages/about.md",
			[]byte("# About\n\nThis is the about page."),
			0644)

		if err != nil {
			panic(err)
		}
	}

	os.MkdirAll("posts", 0755)
	os.MkdirAll("pages", 0755)
}

func main() {
	createFilesAndDirs()
	writeIndex()
	writePosts()
	writePostsPage()
	writePages()
}
