package cache

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
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
			return fmt.Errorf("failed to remove image from storage: %w", err)
		}

		delete(l.items, last.Value.(kvPair).key)
	}

	file, err := os.Create(path.Join(l.dirPath, fileName))
	if err != nil {
		return fmt.Errorf("failed to save image to storage: %w", err)
	}

	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return fmt.Errorf("failed to save image to storage: %w", err)
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
	_, err := os.Stat(l.dirPath)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(l.dirPath, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	files, err := os.ReadDir(l.dirPath)
	if err != nil {
		return fmt.Errorf("failed to read dir: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err = l.loadFileToStorage(file.Name())
		if err != nil {
			continue
		}
	}

	return nil
}

func (l *LruImageCache) Clear() error {
	l.m.Lock()
	defer l.m.Unlock()

	l.items = make(map[string]*ListItem, l.capacity)
	l.queue = NewList()

	entries, err := os.ReadDir(l.dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(l.dirPath, entry.Name())
		if entry.IsDir() {
			err := os.RemoveAll(fullPath)
			if err != nil {
				return fmt.Errorf("error removing directory %s: %w", entry.Name(), err)
			}
		} else {
			err := os.Remove(fullPath)
			if err != nil {
				return fmt.Errorf("error removing file %s: %w", entry.Name(), err)
			}
		}
	}
	return nil
}

func (l *LruImageCache) loadFileToStorage(filename string) error {
	f, err := os.Open(path.Join(l.dirPath, filename))
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}

	if img, err := jpeg.Decode(f); err == nil {
		return l.Save(filename, img)
	}

	return nil
}
