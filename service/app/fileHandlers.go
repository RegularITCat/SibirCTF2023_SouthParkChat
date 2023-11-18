package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//do i really need files in service?

func GetFilesHandler(w http.ResponseWriter, r *http.Request) {
	page, pageSize, err := parseFormPageParams(r)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	files, err := GetFiles(page, pageSize)
	result, err := json.Marshal(files)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	file, err := GetFile(id)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	result, err := json.Marshal(file)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	defer file.Close()
	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	if _, err := io.Copy(dst, file); err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	fid, err := CreateFile(handler.Filename)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	var fileObj File
	fileObj.ID = int(fid)
	fileObj.Name = filepath.Base(handler.Filename)
	fileObj.Path = handler.Filename
	fileObj.UploadTimestamp = time.Now().Unix()
	resultJSON, err := json.Marshal(fileObj)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}

func DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	file, err := GetFile(id)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	fileDescriptor, err := os.Open(file.Path)
	defer fileDescriptor.Close()
	tempBuffer := make([]byte, 512)
	fileDescriptor.Read(tempBuffer)
	fileContentType := http.DetectContentType(tempBuffer)
	fileStat, _ := fileDescriptor.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)
	w.Header().Set("Content-Type", fileContentType+";"+file.Name)
	w.Header().Set("Content-Length", fileSize)
	fileDescriptor.Seek(0, 0)
	io.Copy(w, fileDescriptor)
}

func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	file, err := GetFile(id)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	err = os.Remove(file.Path)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	err = DeleteFile(id)
	if err != nil {
		printError(w, r, err, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
