package sqlc

type SQLc struct {
	fields Fields
	table  string
	conds  []Cond
	extras []Extra
}

func NewSQLc(table string) *SQLc {
	sqlc := &SQLc{
		table: table,
	}
	return sqlc
}

func (s *SQLc) And(c Cond) *SQLc {
	s.conds = append(s.conds, c)
	return s
}

func (s *SQLc) Ext(e Extra) *SQLc {
	s.extras = append(s.extras, e)
	return s
}

func (s *SQLc) condToStr() string {
	if len(s.conds) == 0 {
		return ""
	}
	condStr := " WHERE "
	var sep string
	for _, cond := range s.conds {
		condStr += sep + cond.ToString()
		sep = " AND "
	}
	return condStr
}

func (s *SQLc) extraToStr() string {
	if len(s.extras) == 0 {
		return ""
	}
	extraStr := " "
	var sep string
	for _, ext := range s.extras {
		extraStr += sep + ext.ToString()
		sep = " "
	}
	return extraStr
}

func (s *SQLc) extraExceptLimitToStr() string {
	if len(s.extras) == 0 {
		return ""
	}
	extraStr := " "
	var sep string
	for _, ext := range s.extras {
		_, ok := ext.(*LimitEx)
		if !ok {
			extraStr += sep + ext.ToString()
			sep = " "
		}
	}
	if extraStr == " " {
		return ""
	}
	return extraStr
}

// input -> struct
func (s *SQLc) ToSelect(input interface{}) (string, error) {
	fields, err := NewFieldSlice(input)
	if err != nil {
		return "", err
	}

	var fieldSQL, sep string
	for _, v := range fields.fieldSlice {
		fieldSQL += sep + v
		sep = ", "
	}

	sql := "SELECT " + fieldSQL + " FROM " + s.table + s.condToStr() + s.extraToStr()
	return sql, nil
}

func (s *SQLc) ToCount() string {
	sql := "SELECT COUNT(*) FROM " + s.table + s.condToStr() + s.extraExceptLimitToStr()
	return sql
}

// input -> struct
func (s *SQLc) ToInsert(input interface{}) (string, error) {
	fields, err := NewFieldMap(input)
	if err != nil {
		return "", err
	}

	var fieldSQL, valueSQL, sep string
	for k, v := range fields.fieldMap {
		fieldSQL += sep + k
		valueSQL += sep + "'" + MysqlRealEscapeString(v) + "'"
		sep = ", "
	}

	fieldSQL = "(" + fieldSQL + ")"
	valueSQL = "(" + valueSQL + ")"

	sql := "INSERT INTO " + s.table + " " + fieldSQL + " VALUES " + valueSQL
	return sql, nil
}

// input -> struct
func (s *SQLc) ToUpdate(input interface{}) (string, error) {
	fields, err := NewFieldMap(input)
	if err != nil {
		return "", err
	}
	var setSQL, sep string
	for k, v := range fields.fieldMap {
		setSQL += sep + k + " = " + "'" + MysqlRealEscapeString(v) + "'"
		sep = ", "
	}
	sql := "UPDATE " + s.table + " SET " + setSQL + s.condToStr()
	return sql, nil
}

func (s *SQLc) ToDelete() (string, error) {
	sql := "DELETE FROM " + s.table + s.condToStr()
	return sql, nil
}
