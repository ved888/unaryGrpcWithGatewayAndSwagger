package dbhelper

import (
	"uploadImage/model"
)

// UploadImage upload image on aws
func (d DAO) UploadImage(user model.Image) (*string, error) {
	// SQL=language
	SQL := `INSERT INTO image(
                      bucket_name,
                      path 
                      )    
                VALUES($1,$2)
                returning id::uuid
                                  `
	arg := []interface{}{
		user.BucketName,
		user.Path,
	}
	var imageId string
	err := d.DB.QueryRowx(SQL, arg...).Scan(&imageId)
	if err != nil {
		return nil, err
	}
	return &imageId, nil
}
