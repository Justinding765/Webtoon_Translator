package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
    "golang.org/x/net/html"
)

type ImageData struct {
    ModifiedNode *html.Node
    Index        int
}

func processImageTag(node *html.Node, index int, ch chan<- ImageData) {
    // Extract the image URL from the node
    var imgURL string
    for _, attr := range node.Attr {
        if attr.Key == "src" {
            imgURL = attr.Val
            break
        }
    }

    // Make a request to the Python API
    apiURL := "http://localhost:5000/translate_image" // Replace with your actual API URL
    jsonData, _ := json.Marshal(map[string]string{"url": imgURL})
    resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Printf("Failed to request image translation at index %d: %v\n", index, err)
        return
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Printf("Failed to read response body at index %d: %v\n", index, err)
        return
    }

    // Unmarshal the JSON response into a map
    var result map[string]interface{}
    err = json.Unmarshal(body, &result)
    if err != nil {
        fmt.Printf("Failed to unmarshal response body: %v\n", err)
        return
    }

    // Extract the image URL from the response
    imageUrl, ok := result["image"].(string)
    if !ok {
        fmt.Println("Error: Expected a string for the 'image' key")
        return
    }


    // Create a new image node with the translated image
    translatedNode := &html.Node{
        Type: html.ElementNode,
        Data: "img",
        Attr: []html.Attribute{
            {Key: "src", Val: imageUrl},
        },
    }

    // Wrap the translated image in a div
    divNode := &html.Node{
        Type: html.ElementNode,
        Data: "div",
    }
    divNode.AppendChild(translatedNode)

    // Send the modified node to the channel
    ch <- ImageData{ModifiedNode: divNode, Index: index}
}

func cloneNode(n *html.Node) *html.Node {
    newNode := &html.Node{
        Type:     n.Type,
        DataAtom: n.DataAtom,
        Data:     n.Data,
        Attr:     make([]html.Attribute, len(n.Attr)),
    }
    copy(newNode.Attr, n.Attr)

    return newNode
}

func modifyHTML(inputFile, outputFile string) error {
    content, err := ioutil.ReadFile(inputFile)
    if err != nil {
        return err
    }

    doc, err := html.Parse(bytes.NewReader(content))
    if err != nil {
        return err
    }

    var wg sync.WaitGroup
    ch := make(chan ImageData)
    index := 0

    var processNode func(*html.Node)
    processNode = func(n *html.Node) {
        if n.Type == html.ElementNode && n.Data == "img" {
            wg.Add(1)
            go processImageTag(n, index, ch)
            index++
        }
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            processNode(c)
        }
    }

    processNode(doc)

    go func() {
        wg.Wait()
        close(ch)
    }()

    var buf bytes.Buffer
    buf.WriteString("<html><body>\n")

    images := make([]*html.Node, index)
    for imgData := range ch {
        images[imgData.Index] = imgData.ModifiedNode
        wg.Done()
    }

    for _, img := range images {
        html.Render(&buf, img)
        buf.WriteString("\n")
    }

    buf.WriteString("</body></html>\n")

    return ioutil.WriteFile(outputFile, buf.Bytes(), 0644)
}

func main() {
    err := modifyHTML("../frontend/src/pages/output.html", "../frontend/src/pages/output.html")
    if err != nil {
        fmt.Println("Error:", err)
    }
}
