// package main

// import (
//     "bytes"
//     "fmt"
//     "io/ioutil"
//     "sync"

//     "golang.org/x/net/html"
// )

// type ImageData struct {
//     ModifiedNode *html.Node
//     Index        int
// }

// func processImageTag(node *html.Node, index int, ch chan<- ImageData) {
//     divNode := &html.Node{
//         Type: html.ElementNode,
//         Data: "div",
//     }
//     divNode.AppendChild(cloneNode(node)) // Clone and wrap the original image node

//     ch <- ImageData{ModifiedNode: divNode, Index: index}
// }

// func cloneNode(n *html.Node) *html.Node {
//     newNode := &html.Node{
//         Type:     n.Type,
//         DataAtom: n.DataAtom,
//         Data:     n.Data,
//         Attr:     make([]html.Attribute, len(n.Attr)),
//     }
//     copy(newNode.Attr, n.Attr)

//     return newNode
// }

// func modifyHTML(inputFile, outputFile string) error {
//     content, err := ioutil.ReadFile(inputFile)
//     if err != nil {
//         return err
//     }

//     doc, err := html.Parse(bytes.NewReader(content))
//     if err != nil {
//         return err
//     }

//     var wg sync.WaitGroup
//     ch := make(chan ImageData)
//     index := 0

//     var processNode func(*html.Node)
//     processNode = func(n *html.Node) {
//         if n.Type == html.ElementNode && n.Data == "img" {
//             wg.Add(1)
//             go processImageTag(n, index, ch)
//             index++
//         }
//         for c := n.FirstChild; c != nil; c = c.NextSibling {
//             processNode(c)
//         }
//     }

//     processNode(doc)

//     go func() {
//         wg.Wait()
//         close(ch)
//     }()

//     var buf bytes.Buffer
//     buf.WriteString("<html><body>\n")

//     images := make([]*html.Node, index)
//     for imgData := range ch {
//         images[imgData.Index] = imgData.ModifiedNode
//         wg.Done()
//     }

//     for _, img := range images {
//         html.Render(&buf, img)
//         buf.WriteString("\n")
//     }

//     buf.WriteString("</body></html>\n")

//     return ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
// }

// func main() {
//     err := modifyHTML("output.html", "output.html")
//     if err != nil {
//         fmt.Println("Error:", err)
//     }
// }
