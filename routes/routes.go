package routes

type Note struct {
	ViewId, Title, PublicDate, Body, ReleaseNumber string
}

type User struct {
	Email, Pass string
}

const (
	time_layout = "2006-01-02"
	time_format = "January 02, 2006"
)
