package webdevdata

import "code.google.com/p/go.net/html"
import "code.google.com/p/cascadia"
import "os"
import "io"
import "fmt"
import "flag"
import "bufio"

func ProcessMatchingTags(file string, cssSel string, run func(*html.Node)) {
  selector := cascadia.MustCompile(cssSel)
  htmlReader := reader(file)
  node, err := html.Parse(htmlReader)
  if err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }
  matchedNodes := selector.MatchAll(node)
  for _, node := range matchedNodes {
    run(node)
  }
}

func ProcessTags(file string, process func(html.Token)) {
  htmlReader := reader(file)
  d := html.NewTokenizer(htmlReader)
  for {
    // token type
    tokenType := d.Next()
    if tokenType == html.ErrorToken {
      err := d.Err()
      if err != io.EOF {
        fmt.Println(err)
        os.Exit(-1)
      } else {
        return
      }
    }
    token := d.Token()
    switch tokenType {
      case html.StartTagToken, html.SelfClosingTagToken: // <tag>
      process(token)
    }
  }
  return
}

func GetAttr(key string, attrs []html.Attribute) (string) {
  for _, attr := range attrs {
    if attr.Key == key {
      return attr.Val
    }
  }
  return ""
}

func reader(file string) io.Reader {
  reader, err := os.Open(file)
  if err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }
  return reader
}

func GetFiles(filesChan chan string, skip int) {
  if flag.NArg() > skip {
    files := flag.Args()
    for i, file := range files {
      if i < skip { continue }
      filesChan <- file
    }
  } else {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
      filesChan <- scanner.Text()
    }
    if err := scanner.Err(); err != nil {
      fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
  }
  close(filesChan)
}

