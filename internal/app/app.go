package app

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/heltirj/image_previewer/internal/image_transformer"
	"image"
	"image/jpeg"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

type Cache interface {
	Save(key string, img image.Image) error
	Get(key string) image.Image
	Load() error
	Clear() error
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	DebugKV(msg string, keysAndValues ...interface{})
	InfoKV(msg string, keysAndValues ...interface{})
	WarnKV(msg string, keysAndValues ...interface{})
	ErrorKV(msg string, keysAndValues ...interface{})
}

type App struct {
	Logger Logger
	Cache  Cache
	client *http.Client
}

func New(logg Logger, cache Cache) *App {
	return &App{
		Logger: logg,
		Cache:  cache,
		client: &http.Client{},
	}
}

var re = regexp.MustCompile(`^/(\d+)/(\d+)/(.*)$`)

func (a *App) GetResizedImage(w http.ResponseWriter, r *http.Request) {
	width, height, imgURL, err := parse(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filename, err := getFileNameByURL(r.URL.RequestURI())
	if err != nil {
		http.Error(w, "invalid url", http.StatusBadRequest)
		return
	}

	if img := a.Cache.Get(filename); img != nil {
		returnImage(w, img)
		return
	}

	response, err := a.doRequest(imgURL, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		for key, value := range response.Header {
			w.Header()[key] = value
		}

		w.WriteHeader(response.StatusCode)
		http.Error(w, "undefined source", response.StatusCode)
		return
	}

	w.Header().Set("origin", response.Request.URL.String())
	srcImg, _, err := image.Decode(response.Body)
	if err != nil {
		if errors.Is(err, image.ErrFormat) {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img, err := image_transformer.Resize(srcImg, width, height)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = a.Cache.Save(filename, img)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	returnImage(w, img)
}

func (a *App) ClearCache(w http.ResponseWriter, _ *http.Request) {
	err := a.Cache.Clear()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("cache has been cleared"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) doRequest(imgURL string, r *http.Request) (*http.Response, error) {
	parsedURL, err := url.Parse(imgURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	sendRequest := func(scheme string) (*http.Response, error) {
		parsedURL.Scheme = scheme
		req, err := http.NewRequest(r.Method, parsedURL.String(), r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}
		req.Header = r.Header
		return a.client.Do(req)
	}

	resp, err := sendRequest("http")
	if err == nil {
		return resp, nil
	}

	resp, err = sendRequest("https")
	if err != nil {
		return nil, fmt.Errorf("failed to do request with https: %w", err)
	}

	return resp, nil
}

func getFileNameByURL(urlPath string) (string, error) {
	parsedURL, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	_, err = hash.Write([]byte(parsedURL.Path))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)) + ".jpg", nil
}

func returnImage(w http.ResponseWriter, img image.Image) {
	w.Header().Set("Content-Type", "image/jpeg")
	err := jpeg.Encode(w, img, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func parse(r *http.Request) (width, height int, imgURL string, err error) {
	query := r.URL.RequestURI()
	matches := re.FindStringSubmatch(query)
	if len(matches) != 4 {
		err = fmt.Errorf("invalid query: %s", query)
		return
	}

	width, err = strconv.Atoi(matches[1])
	if err != nil {
		err = fmt.Errorf("invalid width")
		return
	}

	height, err = strconv.Atoi(matches[2])
	if err != nil {
		err = fmt.Errorf("invalid height")
		return
	}

	imgURL = matches[3]
	return
}
