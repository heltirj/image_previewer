package cache

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path"
	"sync"
)

type kvPair struct {
	key   string
	value image.Image
}

type LruImageCache struct {
	capacity int
	queue    List
	items    map[string]*ListItem
	m        *sync.RWMutex
	dirPath  string
}

func NewLruImageCache(capacity int, dirPath string) *LruImageCache {
	return &LruImageCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[string]*ListItem, capacity),
		m:        &sync.RWMutex{},
		dirPath:  dirPath,
	}
}

func (l *LruImageCache) Save(fileName string, img image.Image) error {
	l.m.Lock()
	defer l.m.Unlock()
	item, ok := l.items[fileName]
	if ok {
		item.Value = kvPair{key: fileName, value: img}
		l.items[fileName] = item
		l.queue.MoveToFront(item)
		return nil
	}

	if l.queue.Len() == l.capacity {
		last := l.queue.Back()
		l.queue.Remove(last)
		removeKey := last.Value.(kvPair).key
		err := os.Remove(path.Join(l.dirPath, removeKey))
		if err != nil {
			return fmt.Errorf("failed to remove image from storage: %v", err)
		}

		delete(l.items, last.Value.(kvPair).key)
	}

	file, err := os.Create(path.Join(l.dirPath, fileName))
	if err != nil {
		return fmt.Errorf("failed to save image to storage: %v", err)
	}

	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return fmt.Errorf("failed to save image to storage: %v", err)
	}

	newItem := kvPair{key: fileName, value: img}
	l.items[fileName] = l.queue.PushFront(newItem)

	return nil
}

func (l *LruImageCache) Get(fileName string) image.Image {
	l.m.RLock()
	defer l.m.RUnlock()
	if item, ok := l.items[fileName]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(kvPair).value
	}
	return nil
}

func (l *LruImageCache) Load() error {
	files, err := os.ReadDir(l.dirPath)
	if err != nil {
		log.Fatalf("failed to read dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err = l.loadFileToStorage(file.Name())
		if err != nil {
			log.Printf("failed to load image from storage: filename: %s, err: %s", file.Name(), err)
			continue
		}
	}

	return nil
}

func (l *LruImageCache) loadFileToStorage(filename string) error {
	f, err := os.Open(path.Join(l.dirPath, filename))
	if err != nil {
		return fmt.Errorf("failed to open image file: %v", err)
	}

	if img, err := jpeg.Decode(f); err == nil {
		return l.Save(filename, img)
	}

	return nil
}
