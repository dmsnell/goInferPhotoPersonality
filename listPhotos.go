package main

import (
	"path/filepath"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var serverAddress = "127.0.0.1:2346"
var photoHome = "/Users/dmsnell/Pictures/"
var photoList []string;
var photoSubmitter = make( chan string )

func addPhotoToList() {
	for {
		photoList = append( photoList, <- photoSubmitter )
	}
}

func dispatchList( path string, info os.FileInfo, err error ) error {
	if ( ".jpg" == filepath.Ext( path ) ) {
		photoSubmitter <- path
	}

	return nil
}

func listPhotos( writer http.ResponseWriter, request *http.Request ) {
	photoList = make([]string, 0)

	go addPhotoToList()
	filepath.Walk( photoHome, dispatchList )

	fmt.Fprint( writer, fmt.Sprintf( "Found %d JPEG images!", len(photoList) ) )
}

func init() {
	flag.StringVar( &serverAddress, "address", serverAddress, "Server listening address and port" )
}

func main() {
	flag.Parse()
	http.HandleFunc( "/listPhotos", listPhotos )
	log.Fatal( http.ListenAndServe( serverAddress, nil ) )
}