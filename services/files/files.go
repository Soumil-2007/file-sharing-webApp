package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Soumil-2007/file-sharing-webApp/services/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const UploadDir = "./uploads"

var allowedMIMEs = map[string]struct{}{
	"image/jpeg":      {},
	"image/png":       {},
	"application/pdf": {},
	"text/plain":      {},
}

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

type Handler struct {
	store *Store
}

func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(r *mux.Router, authMiddleware mux.MiddlewareFunc) {
	fileRouter := r.PathPrefix("/files").Subrouter()
	fileRouter.Use(authMiddleware)

	fileRouter.HandleFunc("", h.ListFiles).Methods("GET")
	fileRouter.HandleFunc("/{id}", h.GetFile).Methods("GET")
	fileRouter.HandleFunc("", h.HandleUpload).Methods("POST")
}

// ------------------- UPLOAD -------------------

func (h *Handler) HandleUpload(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") == "" || !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size == 0 {
		http.Error(w, "uploaded file is empty", http.StatusBadRequest)
		return
	}

	filetype, err := detectMIME(file)
	if err != nil {
		http.Error(w, "failed to detect file type", http.StatusInternalServerError)
		return
	}
	if _, ok := allowedMIMEs[filetype]; !ok {
		http.Error(w, fmt.Sprintf("unsupported file type: %s", filetype), http.StatusBadRequest)
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "failed to reset file reader", http.StatusInternalServerError)
		return
	}

	filename := filepath.Base(header.Filename)
	filename = filepath.Clean(filename)
	filename = strings.ReplaceAll(filename, " ", "_")
	newFilename := fmt.Sprintf("%s-%s", uuid.NewString(), filename)
	outPath := filepath.Join(UploadDir, newFilename)

	go func() {
		// Simulate processing (e.g., generate thumbnail, log)
		time.Sleep(1 * time.Second)
		log.Println("File processed:", newFilename)
	}()

	// save file to disk
	outFile, err := os.Create(outPath)
	if err != nil {
		http.Error(w, "unable to create file", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		http.Error(w, "unable to save file", http.StatusInternalServerError)
		return
	}

	// 🔹 TEMPORARY: hardcoded user id until auth is ready
	userID := 1

	// save metadata in DB
	_, err = h.store.DB.Exec(
		"INSERT INTO files (owner_id, original_name, stored_name, mime_type, size_bytes, path) VALUES (?, ?, ?, ?, ?, ?)",
		userID, header.Filename, newFilename, filetype, header.Size, outPath,
	)
	if err != nil {
		http.Error(w, "db insert failed", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"message":       "file uploaded successfully",
		"filename":      newFilename,
		"mime_type":     filetype,
		"size_bytes":    header.Size,
		"original_name": header.Filename,
	}
	json.NewEncoder(w).Encode(resp)
}

func detectMIME(f multipart.File) (string, error) {
	buffer := make([]byte, 512)
	if _, err := f.Read(buffer); err != nil && err != io.EOF {
		return "", err
	}
	return http.DetectContentType(buffer), nil
}

// ------------------- LIST FILES -------------------

func (h *Handler) ListFiles(w http.ResponseWriter, r *http.Request) {
	// 🔹 TEMPORARY: hardcoded user id until auth is ready
	userID := 1
	
	rows, err := h.store.DB.Query(
		"SELECT id, original_name, stored_name, mime_type, size_bytes, created_at FROM files WHERE owner_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		http.Error(w, "query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []map[string]interface{}

	var fileCache = make(map[string][]map[string]interface{})
	key := fmt.Sprintf("files:%d", userID)
	if cached, ok := fileCache[key]; ok {
		json.NewEncoder(w).Encode(cached)
		return
	}
	// Query DB...
	fileCache[key] = out
	json.NewEncoder(w).Encode(out)
	
	for rows.Next() {
		var id int
		var orig, stored, mime string
		var size int64
		var created string
		if err := rows.Scan(&id, &orig, &stored, &mime, &size, &created); err != nil {
			continue
		}
		out = append(out, map[string]interface{}{
			"id":            id,
			"original_name": orig,
			"stored_name":   stored,
			"mime_type":     mime,
			"size":          size,
			"created_at":    created,
		})
	}

	json.NewEncoder(w).Encode(out)
}

// ------------------- GET FILE -------------------

func (h *Handler) GetFile(w http.ResponseWriter, r *http.Request) {
	// 🔹 TEMPORARY: hardcoded user id until auth is ready
	userID := auth.GetUserIDFromContext(r.Context())

	vars := mux.Vars(r)
	id := vars["id"]

	var path, mime string
	err := h.store.DB.QueryRow(
		"SELECT path, mime_type FROM files WHERE id = ? AND owner_id = ?",
		id, userID,
	).Scan(&path, &mime)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	http.ServeFile(w, r, path)
}
