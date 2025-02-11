package workers

type Job interface {
	Work(ID int)
}

type UploadJob struct {
	ID       int
	FileName string
	URL      string
}
