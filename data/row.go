package data

func NewRowList(rows []Row, err error) *RowList {
	return &RowList{
		Rows: rows,
		Err:  err,
	}
}

type RowList struct {
	Rows []Row
	Err  error
}

type Row map[string]any

func NewRow(rowMap map[string]any) Row {
	return Row(rowMap)
}

func NewSingleRow(rowMap map[string]any) []Row {
	return []Row{NewRow(rowMap)}
}
