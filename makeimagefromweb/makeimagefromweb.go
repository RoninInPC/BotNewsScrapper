package makeimagefromweb

import "image"

type MakeImageFromWeb interface {
	Get(url string) (string, image.Image, error)
}
