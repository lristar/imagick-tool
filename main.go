package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/gographics/imagick.v2/imagick"
)

var mw *imagick.MagickWand

func main() {

	fmt.Println("Imagick setup")

	// Setup
	imagick.Initialize()
	defer imagick.Terminate()

	mw = imagick.NewMagickWand()
	defer mw.Destroy()

	mux := http.NewServeMux()

	mux.HandleFunc("/convert", ConvertToJPG)

	fmt.Println("Start http server")
	fmt.Println(http.ListenAndServe(":900", mux))

}

func ConvertToJPG(w http.ResponseWriter, r *http.Request) {

	// Check request type
	if r.Method != "POST" {
		fmt.Println("Recivied non Post request")
		return
	}

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)

	// FormFile returns the first file for the given key `file`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("file")

	// check error
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Read file to byte
	f, e := ioutil.ReadAll(file)
	if e != nil {
		fmt.Println(e)
		return
	}

	// Must be *before* ReadImageFile
	// Make sure our image is high quality
	if err := mw.SetResolution(300, 300); err != nil {
		fmt.Println(err)
		return
	}

	// Load the byte into imagick
	if err := mw.ReadImageBlob(f); err != nil {
		fmt.Println(err)
		return
	}

	// Must be *after* ReadImageFile
	// Flatten image and remove alpha channel, to prevent alpha turning black in jpg
	if mw.GetImageAlphaChannel() {
		if err := mw.SetImageAlphaChannel(imagick.ALPHA_CHANNEL_REMOVE); err != nil {
			fmt.Println(err)
			return
		}
	}

	// Set any compression (100 = max quality)
	if err := mw.SetCompressionQuality(95); err != nil {
		fmt.Println(err)
		return
	}

	// Select only first page of pdf
	mw.SetIteratorIndex(0)

	// Convert into JPG
	if err := mw.SetFormat("jpg"); err != nil {
		fmt.Println(err)
		return
	}

	Filename := "image.jpg"
	// Save File
	if err := mw.WriteImage(Filename); err != nil {
		fmt.Println(err)
		return
	}

	mw.Clear()

	//Check if file exists and open
	Openfile, err := os.Open(Filename)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client

	fmt.Println("Done")
	return
}
