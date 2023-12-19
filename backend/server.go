package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "sort"


)

type TranslateRequest struct {
    URL string `json:"url"`
}
// Helper function to set CORS headers
func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") // Allowed methods
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func runPythonScript(scriptPath, url, outputFilePath string) {
    cmd := exec.Command("python", scriptPath, url, outputFilePath)

    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("Failed to execute command: %s, with output: %s", err, string(output))
    }

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
	runPythonScript("./web_scrapper.py", req.URL, "./output.html")

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
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // Explicitly setting the status code to 200
    json.NewEncoder(w).Encode(imageUrls)
    
    
    

    
}

func main() {
    // Define the route and its handler function
    http.HandleFunc("/translate_images", translateImagesHandler)



    // Start the HTTP server
    fmt.Println("Server is starting on port 5000...")
    if err := http.ListenAndServe(":5000", nil); err != nil {
        log.Fatal(err)
    }
}
