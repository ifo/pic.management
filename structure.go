package main

import ()

type User struct {
	ID       int64
	Email    string
	Username string
}

type Photo struct {
	ID     int64
	Images []Image           // TODO? not include
	UserID int64             // photo uploader
	Exif   map[string]string // TODO? make this a specific type
}

// TODO? include meta info, size, etc?
type Image struct {
	PhotoID  int64  // parent
	Location string // TODO? make this a specific type
}

// TODO? Make User Groups (Group)
// TODO? Make Photo Groups (Album)
