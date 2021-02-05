package gotainer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type Image struct {
	id   string
	size int64
}

func fileNameWithoutExt(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func (i *Image) GetID() string {
	return i.id
}

func (i *Image) GetSizeMB() float64 {
	return float64(i.size) / 1000000.0
}

func (i *Image) Run(commands []string) error {
	tarPath := filepath.Join(imagesPath, i.id+imageExt)
	tar, err := os.Open(tarPath)

	if err != nil {
		return err
	}

	id := uuid.New()
	containerPath := filepath.Join(containerPath, id.String())
	os.Mkdir(containerPath, 0755)
	err = Untar(containerPath, tar)

	if err != nil {
		return err
	}

	container := NewContainer(id, i)
	defer container.Remove()
	container.Exec(commands)

	return nil

}

func GetImage(id string) (*Image, error) {
	fileName := id + imageExt
	imagePath := filepath.Join(imagesPath, fileName)
	if f, err := os.Stat(imagePath); err == nil {
		return NewImage(id, f.Size()), nil
	}
	return nil, ErrImageNotExist
}

func GetAllImages() ([]*Image, error) {
	dirs, err := ioutil.ReadDir(imagesPath)

	if err != nil {
		return nil, ErrImagesDirUnreadable
	}

	var images []*Image
	for _, d := range dirs {
		name := fileNameWithoutExt(d.Name())
		i := NewImage(name, d.Size())
		images = append(images, i)
	}

	return images, nil
}

func NewImage(id string, size int64) *Image {
	i := &Image{id, size}
	return i
}
