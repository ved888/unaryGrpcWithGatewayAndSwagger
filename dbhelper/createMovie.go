package dbhelper

import (
	"fmt"
	"gRPC2/model"
)

// MovieCreate create movie by given payload
func (d DAO) MovieCreate(movie *model.Movie) (*string, error) {
	//SQL==language
	SQL := `INSERT INTO movie (
                   title,
                   genre
                   ) 
            VALUES ($1, $2)
            RETURNING id::uuid ;`

	arg := []interface{}{
		movie.Title,
		movie.Genre,
	}
	var movieId string

	err := d.DB.QueryRowx(SQL, arg...).Scan(&movieId)
	if err != nil {
		fmt.Println("error in movieId:", err.Error())
		return nil, err
	}

	return &movieId, nil
}

// GetMovieById get movie by given id
func (d DAO) GetMovieById(movieId string) (*model.Movie, error) {

	// SQL=language
	SQL := `SELECT 
                title,
                genre
          FROM movie
          WHERE id=$1::uuid
                     `

	var movie model.Movie
	err := d.DB.Get(&movie, SQL, &movieId)
	if err != nil {
		fmt.Println("error in movieId:", err.Error())
		return nil, err
	}
	return &movie, nil
}

// GetAllMovies get all movies
func (d DAO) GetAllMovies() ([]model.Movie, error) {
	// SQL=language
	SQL := `SELECT 
                id,
                title,
                genre
          FROM movie
                    `

	var movie []model.Movie

	err := d.DB.Select(&movie, SQL)
	if err != nil {
		return nil, err
	}
	return movie, nil

}

// UpdateMovie update movie by id with given fields
func (d DAO) UpdateMovie(movie model.Movie, movieId string) error {
	// SQL query
	SQL := `
		UPDATE movie
		SET title = COALESCE($1, title),
			genre = COALESCE($2, genre),
			updated_at = NOW()
		WHERE id = $3::uuid
	`
	args := []interface{}{
		movie.Title,
		movie.Genre,
		movieId,
	}
	_, err := d.DB.Exec(SQL, args...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMovie delete movie by id
func (d DAO) DeleteMovie(movieId string) error {
	// Sql=language
	SQL := `
         UPDATE movie
		SET 
			archived_at = NOW()
		WHERE id = $1::uuid
          `

	_, err := d.DB.Exec(SQL, movieId)
	if err != nil {
		return err
	}

	return nil
}
