package dal

import (
    "database/sql"
    "errors"
    "github.com/go-sql-driver/mysql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/nihileon/ticktak/log"
    "github.com/nihileon/ticktak/sqlc"
    "reflect"
    "time"
)

const (
    UserTable = "t_user"
)

const (
    MysqlErDupEntry  = 1062
    MysqlErDupKey    = 1022
    MysqlErDupUnique = 1169
)

var mqDB *sql.DB

func InitDB(dsn string) error {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }

    err = db.Ping()
    if err != nil {
        return err
    }
    mqDB = db
    return nil
}

var ErrDupkeyInsert = errors.New("Insert dupkey error")
var ErrAffectNoRows = errors.New("Exec affect 0 row")
var ErrSelectZeroRows = errors.New("Select zero rows")

func FetchSession() *Session {
    return &Session{
        conn: mqDB,
        tx:   nil,
    }
}

type Session struct {
    conn *sql.DB
    tx   *sql.Tx
}

func (s *Session) Begin() error {
    if s.tx != nil {
        return errors.New("Tx is already exist")
    }
    currTime := time.Now()
    tx, err := s.conn.Begin()
    if err != nil {
        return err
    }
    s.tx = tx
    log.GetLogger().Info("[exec BEGIN], [cost time]: %s", time.Since(currTime))
    return nil
}

func (s *Session) Commit() error {
    if s.tx == nil {
        return errors.New("Tx is already ends")
    }
    currTime := time.Now()
    err := s.tx.Commit()
    if err != nil {
        return err
    }
    s.tx = nil
    log.GetLogger().Info("[exec COMMIT], [cost time]: %s", time.Since(currTime))
    return nil
}

func (s *Session) Rollback() error {
    if s.tx == nil {
        return errors.New("Tx is already ends")
    }
    currTime := time.Now()
    err := s.tx.Rollback()
    if err != nil {
        return err
    }
    s.tx = nil
    log.GetLogger().Info("[exec ROLLBACK], [cost time]: %s", time.Since(currTime))
    return nil
}

func (s *Session) Count(c *sqlc.SQLc) (int, error) {
    countStr := c.ToCount()
    var count int
    var err error
    if s.tx != nil {
        err = s.tx.QueryRow(countStr).Scan(&count)
    } else {
        err = s.conn.QueryRow(countStr).Scan(&count)
    }
    return count, err
}

func (s *Session) Select(c *sqlc.SQLc, input interface{}) error {
    value := reflect.ValueOf(input)
    if value.Kind() != reflect.Ptr || value.IsNil() {
        return errors.New("input must be a pointer")
    }
    direct := value.Elem()
    baseType := direct.Type().Elem()
    baseValue := reflect.New(baseType).Interface()
    sqlStr, err := c.ToSelect(baseValue)
    log.GetLogger().Info("Select SQL: %s", sqlStr)
    currTime := time.Now()
    if err != nil {
        return err
    }
    var rows *sql.Rows
    if s.tx != nil {
        rows, err = s.tx.Query(sqlStr)
    } else {
        rows, err = s.conn.Query(sqlStr)
    }
    if err != nil {
        return err
    }
    defer rows.Close()

    cols, err := rows.Columns()
    if err != nil {
        return err
    }

    values := make([]interface{}, len(cols))
    mapper := make(map[string]interface{})

    var val reflect.Value

    for rows.Next() {
        val = reflect.Indirect(reflect.New(baseType))

        for i := 0; i < baseType.NumField(); i++ {
            mapper[baseType.Field(i).Tag.Get("json")] = reflect.Indirect(val).Field(i).Addr().Interface()
        }

        for j := 0; j < len(cols); j++ {
            values[j] = mapper[cols[j]]
        }

        err = rows.Scan(values...)
        if err != nil {
            return err
        }
        direct.Set(reflect.Append(direct, val))

    }
    log.GetLogger().Info("[select cost time]: %s", time.Since(currTime))
    return nil
}

func (s *Session) Insert(c *sqlc.SQLc, input interface{}) (int64, error) {
    sqlStr, err := c.ToInsert(input)
    if err != nil {
        return 0, err
    }
    ret, _, err := s.Exec(sqlStr)
    if err != nil {
        if isDupInsert(err) {
            return 0, ErrDupkeyInsert
        }
        return 0, err

    }
    id, err := ret.LastInsertId()
    return id, err
}

func (s *Session) Update(c *sqlc.SQLc, input interface{}) error {
    sqlStr, err := c.ToUpdate(input)
    if err != nil {
        return err
    }
    _, affectRows, err := s.Exec(sqlStr)
    if err != nil {
        return err
    }
    if affectRows == 0 {
        return ErrAffectNoRows
    }
    return nil
}

func (s *Session) Delete(c *sqlc.SQLc) error {
    sqlStr, err := c.ToDelete()
    if err != nil {
        return err
    }
    _, affectRows, err := s.Exec(sqlStr)
    if err != nil {
        return err
    }
    if affectRows == 0 {
        return ErrAffectNoRows
    }
    return nil
}

func (s *Session) Exec(sqlStr string) (sql.Result, int64, error) {
    log.GetLogger().Info("[exec sql]: %s", sqlStr)
    currTime := time.Now()
    var ret sql.Result
    var err error
    if s.tx == nil {
        ret, err = s.conn.Exec(sqlStr)
    } else {
        ret, err = s.tx.Exec(sqlStr)
    }
    if err != nil {
        return ret, 0, err
    }
    affectRows, err := ret.RowsAffected()
    log.GetLogger().Info("[exec cost time]: %s", time.Since(currTime))
    return ret, affectRows, err
}

func isDupInsert(err error) bool {
    mysqlErr, ok := err.(*mysql.MySQLError)
    if !ok {
        return false
    }
    switch mysqlErr.Number {
    case MysqlErDupEntry, MysqlErDupKey, MysqlErDupUnique:
        return true
    default:
        return false
    }
}
