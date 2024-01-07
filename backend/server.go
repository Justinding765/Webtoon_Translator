package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "sort"
    "sync"
    "github.com/jung-kurt/gofpdf"
    "os"


)
var urlsOnly []string

type TranslateRequest struct {
    URL string `json:"url"`
}
// Helper function to set CORS headers
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") // Allowed methods
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func runPythonScript(scriptPath string, args ...string) {
    cmdArgs := append([]string{scriptPath}, args...)
    cmd := exec.Command("python", cmdArgs...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Failed to execute command: %s, with output: %s", err, string(output))
    }
    //Print the output from the Python script
    fmt.Printf("Python Script Output:\n%s\n", string(output))
}






// Handler for the "/translate_images" route
func translateImagesHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w) // Enable CORS

    // Handle preflight requests for CORS
    if r.Method == "OPTIONS" {
        return
    }

    if r.Method != "POST" {
        http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
        return
    }
	var req TranslateRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        runPythonScript("./web_scrapper.py", req.URL, "./output.html")
    }()

    go func() {
        defer wg.Done()
        runPythonScript("./clean_up.py")
    }()
    wg.Wait()
    fmt.Println("Finished!")

	// Now call modifyHTML
    imageUrls := make([]ImageData, 0)
    imageUrls, err = modifyHTML("./output.html")
    if err != nil {
        fmt.Println("Error:", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    // Sort images by index
    sort.Slice(imageUrls, func(i, j int) bool {
        return imageUrls[i].Index < imageUrls[j].Index
    })
    urlsOnly = make([]string, len(imageUrls))

    // Extract only URLs into a new slice
    for i, imageUrl := range imageUrls {
        urlsOnly[i] = imageUrl.URL
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Explicitly setting the status code to 200
    json.NewEncoder(w).Encode(urlsOnly)
    
    
    

    
}

// downloadPdfHandler handles the PDF generation
func downloadPdfHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w) // Enable CORS

    if r.Method == "OPTIONS" {
        return
    }

    const (
        pageWidth = 210.0 // A4 width in mm
        margin    = 10.0  // Margin in mm
        scale     = 0.75  // Scale factor for the images
    )

    totalHeight := margin // Starting height for first image

    // Temporary PDF instance to calculate total height
    tempPdf := gofpdf.New("P", "mm", "A4", "")
    for _, url := range urlsOnly {
        localPath, err := downloadImage(url)
        if err != nil {
            log.Printf("Failed to download image: %v", err)
            continue
        }

        options := gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}
        tempPdf.RegisterImageOptions(localPath, options)
        info := tempPdf.GetImageInfo(localPath)

        imgWidth := (pageWidth - 2*margin) * scale
        imgHeight := imgWidth * info.Height() / info.Width()
        totalHeight += imgHeight

        // Remove the temp file
        err = os.Remove(localPath)
        if err != nil {
            log.Printf("Failed to remove temp file: %v", err)
        }
    }

    // Create the actual PDF with custom size
    pdf := gofpdf.NewCustom(&gofpdf.InitType{
        UnitStr: "mm",
        Size:    gofpdf.SizeType{Wd: pageWidth, Ht: totalHeight},
    })
    pdf.AddPage()

    // Add images to the single page
    yPosition := margin
    for _, url := range urlsOnly {
        localPath, err := downloadImage(url)
        if err != nil {
            log.Printf("Failed to download image: %v", err)
            continue
        }

        // Define options here within the scope
        options := gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}
        pdf.RegisterImageOptions(localPath, options)
        info := pdf.GetImageInfo(localPath)

        imgWidth := (pageWidth - 2*margin) * scale
        imgHeight := imgWidth * info.Height() / info.Width()

        centerX := (pageWidth - imgWidth) / 2
        pdf.ImageOptions(localPath, centerX, yPosition, imgWidth, imgHeight, false, options, 0, "")

        yPosition += imgHeight

        // Remove the temp file
        err = os.Remove(localPath)
        if err != nil {
            log.Printf("Failed to remove temp file: %v", err)
        }
    }

    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")
    if err := pdf.Output(w); err != nil {
        log.Printf("Failed to generate PDF: %v", err)
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
    }
}

func main() {
    // Define the route and its handler function
    http.HandleFunc("/translate_images", translateImagesHandler)
    http.HandleFunc("/download_pdf", downloadPdfHandler)



    // Start the HTTP server
    fmt.Println("Server is starting on port 5000...")
    if err := http.ListenAndServe(":5000", nil); err != nil {
        log.Fatal(err)
    }
}
