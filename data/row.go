package data

func NewRowList(rows []Row, headings []string, err error) *RowList {
	return &RowList{
		Rows:     rows,
		Headings: headings,
		Err:      err,
	}
}

type RowList struct {
	Rows     []Row
	Headings []string
	Err      error
}

type Row map[string]any

func NewRow(rowMap map[string]any) Row {
	return Row(rowMap)
}

func NewSingleRow(rowMap map[string]any) []Row {
	return []Row{NewRow(rowMap)}
}
