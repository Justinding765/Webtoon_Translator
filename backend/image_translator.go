package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "sync"
    "golang.org/x/net/html"
    "log"
    "net/http"
    "os/exec"
    "strconv"
    "runtime"
    

)

type ImageData struct {
    URL   string
    Index int
}

type TranslateRequest struct {
    URL string `json:"url"`
    SessionID  string `json:"sessionId"`

}
var wg sync.WaitGroup // Waits for a collection of goroutines to finish

func runPythonScript(scriptPath string, args ...string) ([]byte, error) {
    cmdArgs := append([]string{scriptPath}, args...)
    cmd := exec.Command("python", cmdArgs...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Printf("Failed to execute command: %s, with output: %s", err, string(output))
        return nil, err
    }
    // Print the output from the Python script
    fmt.Printf("Python Script Output:\n%s\n", string(output))
    return output, nil
}


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


func processImageTag(node *html.Node, index int, ch chan<- ImageData, req TranslateRequest) {
    fmt.Println("Number of goroutines:", runtime.NumGoroutine())
    // Extract the image URL from the node
    defer wg.Done()
    var imgURL string
    for _, attr := range node.Attr {
        if attr.Key == "src" {
            imgURL = attr.Val
            break
        }
    }
    output, err := runPythonScript("./translate_Images.py", imgURL, strconv.Itoa(index), req.SessionID)
    if err != nil {
        // handle the error
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


func modifyHTML(inputFile string, req TranslateRequest) ([]ImageData, error) {
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
            go processImageTag(n, index, ch, req)
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


