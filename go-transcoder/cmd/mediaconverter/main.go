package main

import (
	"gotranscoder/internal/converter"
	"gotranscoder/internal/database"
)

func main() {
	db, err := database.ConnectPostgres()
	if err != nil {
		panic(err)
	}

	vc := converter.NewVideoConverter(db)

	vc.TaskHandler(
		[]byte(
			`
			{
				"video_id": 1,
				"path": "/tmp/videos/3"
			}
			`,
		))
}
