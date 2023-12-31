package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "os/exec"
    "sync"
    "golang.org/x/net/html"
    "log"
    "strconv"
    "net/http"

)

type ImageData struct {
    URL   string
    Index int
}
var wg sync.WaitGroup // Waits for a collection of goroutines to finish


// downloadImage downloads the image from the given URL and returns the path to the local temp file
func downloadImage(url string) (string, error) {
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    tmpFile, err := ioutil.TempFile("", "image-*.jpg")
    if err != nil {
        return "", err
    }
    defer tmpFile.Close()

    _, err = io.Copy(tmpFile, resp.Body)
    if err != nil {
        return "", err
    }

    return tmpFile.Name(), nil
}


func processImageTag(node *html.Node, index int, ch chan<- ImageData) {
    // Extract the image URL from the node
    defer wg.Done()
    var imgURL string
    for _, attr := range node.Attr {
        if attr.Key == "src" {
            imgURL = attr.Val
            break
        }
    }
    // Construct the command to run the Python script
    cmd := exec.Command("python", "./translate_Images.py", imgURL, strconv.Itoa(index))

    // Execute the command and capture the combined standard output and standard error
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Failed to execute python script at index %d: %v, with output: %s\n", index, err, string(output))
        return
    }

    // Parse the JSON output
    var result map[string]interface{}
    err = json.Unmarshal(output, &result)
    if err != nil {
        log.Printf("Failed to unmarshal python script output at index %d: %v\n", index, err)
        return
    }

    // Extract the image URL from the response
    imageUrl, ok := result["image"].(string)
    if !ok {
        log.Printf("Error at index %d: Expected a string for the 'image' key\n", index)
        return
    }
    
    fmt.Printf("Translated image URL at index %d: %s\n", index, imageUrl)

    // Send the url to the channel
    ch <- ImageData{
        URL:   imageUrl,
        Index: index,
    }

}

//Delete?
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

func modifyHTML(inputFile string) ([]ImageData, error) {
    content, err := ioutil.ReadFile(inputFile)
    var image_URLS []ImageData
    if err != nil {
        return image_URLS, err
    }

    doc, err := html.Parse(bytes.NewReader(content))
    if err != nil {
        return image_URLS, err
    }

    ch := make(chan ImageData)
    index := 0

    var processNode func(*html.Node)
    //processNode is a recursive function that traverses html document node by node
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
    //started and waits for all images processing goroutines to finsh
    go func() {
        wg.Wait()
        close(ch)
    }()
    
    var images []ImageData

    for imgData := range ch {
        // Append the ImageData object to the images slice
        images = append(images, imgData)

    }
   
    return images, nil
}


