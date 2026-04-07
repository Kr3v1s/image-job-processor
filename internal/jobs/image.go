package jobs

import "github.com/disintegration/imaging"

func CreateThumbnail(inputPath, outputPath string) error {
	img, err := imaging.Open(inputPath)
	if err != nil {
		return err
	}

	thumb := imaging.Thumbnail(img, 200, 200, imaging.Lanczos)

	return imaging.Save(thumb, outputPath)
}
