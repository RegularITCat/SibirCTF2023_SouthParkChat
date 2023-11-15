package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//do i really need files in service?

func GetFiles(w http.ResponseWriter, r *http.Request) {
	files := make([]File, 0)
	rows, err := db.Query("SELECT id,name,path,upload_timestamp FROM files;")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for rows.Next() {
		var file File
		err = rows.Scan(&file.ID, &file.Name, &file.Path, &file.UploadTimestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}
	result, err := json.Marshal(files)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func GetFileByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT id,name,path,upload_timestamp FROM files WHERE id=%v;", id))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var file File
	for rows.Next() {
		err = rows.Scan(&file.ID, &file.Name, &file.Path, &file.UploadTimestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	result, err := json.Marshal(file)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result, err := db.Exec(fmt.Sprintf("INSERT INTO files (name, path, upload_timestamp) VALUES ('%v', '%v', %v);", filepath.Base(handler.Filename), handler.Filename, time.Now().Unix()))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fid, _ := result.LastInsertId()
	var fileObj File
	fileObj.ID = int(fid)
	fileObj.Name = filepath.Base(handler.Filename)
	fileObj.Path = handler.Filename
	fileObj.UploadTimestamp = time.Now().Unix()
	resultJSON, err := json.Marshal(fileObj)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resultJSON)
}

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT id,name,path,upload_timestamp FROM files WHERE id=%v;", id))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var file File
	for rows.Next() {
		err = rows.Scan(&file.ID, &file.Name, &file.Path, &file.UploadTimestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	rows, err := db.Query(fmt.Sprintf("SELECT id,name,path,upload_timestamp FROM files WHERE id=%v;", id))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var file File
	for rows.Next() {
		err = rows.Scan(&file.ID, &file.Name, &file.Path, &file.UploadTimestamp)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	err = os.Remove(file.Path)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(fmt.Sprintf("DELETE FROM files WHERE id = %v;", id))
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}