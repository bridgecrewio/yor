package reports

type ReportService struct {
}

type Report struct {
	PreviouslyTaggedResources []*ResourceRecord
	UntaggedResources         []*ResourceRecord
	NewlyTaggedResources      []*ResourceRecord
}

func (r *Report) GetTotalResources() int {
	//sum of resources seen in this run
	return -1
}

func (r *Report) GetTotalTaggedResources() int {
	//sum of resources previously tagged by Yor
	return -1
}

func (r *Report) GetTotalUnTaggedResources() int {
	//sum of resources that were never tagged by Yor
	return -1
}

func (r *Report) GetTotalNewlyTaggedResources() int {
	//sum of resources newly tagged in this run by Yor
	return -1
}

func (r *ReportService) CreateReport() interface{} {
	//	TODO - determine structure later
	return nil
}
