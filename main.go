package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

//DONE: Make sure that whatever imagesearchkeyword we get, we encode it into searchable keyword (add + in spaces etc)
//TODO: GET subsequent pages for results
//TODO: incorporate goroutines
//TODO: Some of the images result in format not supported, figure that out.

func main() {

	var imageSearchKeyword, googleSearchKey, searchEngineId, folderPath string
	var numberOfImages, numberOfGoRoutines int

	flag.StringVar(&imageSearchKeyword, "keyword", "cats", "keyword refers to the keyword whose images you want to download from Google.")
	flag.StringVar(&googleSearchKey, "key", "", "Authentication key for Google searches. Please look up here.")
	flag.StringVar(&searchEngineId, "searchId", "", "Search Engine ID")
	flag.StringVar(&folderPath, "folderPath", "C:\\Users\\saraw\\development\\52\\GooglePhotoDownloader\\photos\\", "Folder path to store all the images")
	flag.IntVar(&numberOfImages, "numberOfImages", 5, "Number of images to downlaod")
	flag.IntVar(&numberOfGoRoutines, "goRoutines", 10, "Number of concurrent Go routines to run to fetch the images.")
	flag.Parse()

	fmt.Println(imageSearchKeyword, googleSearchKey, searchEngineId, numberOfGoRoutines)

	gImageRequestDetails := GoogleImageRequestDetails{
		"https://www.googleapis.com/customsearch/v1?",
		imageSearchKeyword,
		googleSearchKey,
		searchEngineId,
	}
	imageChan := make(chan string)
	//defer close(imageChan)
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		for i := 0; i < numberOfImages; i++ {
			wg.Add(1)
			go insertImagelinkForPageInChannel(&gImageRequestDetails, (i*10 + 1), numberOfImages, imageChan, &wg)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < numberOfGoRoutines; i++ {
			wg.Add(1)
			go downloadSaveImageFromUrl(folderPath, imageChan, &wg)
		}
		wg.Done()
	}()

}

func getGoogleImageUri(g *GoogleImageRequestDetails, startIndex int) string {
	return (g.Uri + "q=" + url.QueryEscape(g.ImageSearchKeyword) + "&searchType=image&key=" + g.GoogleSearchKey + "&cx=" + g.SearchEngineId + "&start=" + strconv.Itoa(startIndex))
}

func googleImageRequest(g *GoogleImageRequestDetails, startIndex int) (*GoogleImageSearchResponse, error) {

	requestUri := getGoogleImageUri(g, startIndex)
	fmt.Println(requestUri)
	response, err := http.Get(requestUri)
	if err != nil {
		fmt.Println("Error getting HTTP GET request :: ", err)
		return nil, err
	}
	defer response.Body.Close()
	//respBody, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	fmt.Println("Error converting response body into bytes")
	// }
	//fmt.Println(respBody)
	gResponse := GoogleImageSearchResponse{}
	json.NewDecoder(response.Body).Decode(&gResponse)
	fmt.Println(gResponse)
	return &gResponse, nil
}

//responsibiites of this worker: query API anf get result, decode, prase through each image link and add the link to channel
func insertImagelinkForPageInChannel(g *GoogleImageRequestDetails, startIndex int, numberOfImages int, imageChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	gResponse, _ := googleImageRequest(g, startIndex)
	//respBody, _ := io.ReadAll(response.Body)
	// fmt.Println(response)
	// gResponse := GoogleImageSearchResponse{}
	// json.NewDecoder(response).Decode(&gResponse)
	// fmt.Println(gResponse)
	for _, j := range gResponse.Items {
		fmt.Println("PUTTING: ", j.Link)
		imageChan <- j.Link
	}
	if ((startIndex - 1) / 10) == numberOfImages-1 {
		close(imageChan)
	}

}

func downloadSaveImageFromUrl(folderPath string, imageChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for link := range imageChan {
		fmt.Println(link)
		image, err := http.Get(link)
		if err != nil {
			fmt.Println("Error downloading the image form url")
		}
		file, err := os.Create(folderPath + "\\image" + strconv.Itoa(int(hash(link))) + ".jpeg")
		if err != nil {
			fmt.Println("Error opening image file")
		}
		_, err = io.Copy(file, image.Body)
		if err != nil {
			fmt.Println("error saving image")
		}
		image.Body.Close()
		file.Close()
	}

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

//checks whether the folder path input has trailing slashes or not
func checkFolderPath(folderPAth string) string {
	return ""
}
