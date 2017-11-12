package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pborman/uuid"
)

//Model *
type Model struct {
	ID          string     `sql:"index" gorm:"primary_key;column:id;type:varchar(100)" json:"id,omitempty" xml:"id,omitempty"`
	Name        string     `sql:"index" gorm:"column:name;type:varchar(100)" json:"name,omitempty" xml:"name,omitempty"`
	Description string     `gorm:"column:description;type:text" json:"description,omitempty" xml:"description,omitempty"`
	CreatedAt   time.Time  `sql:"index" gorm:"column:created_at;type:timestamp" json:"created_at,omitempty" xml:"created_at,omitempty"`
	UpdatedAt   time.Time  `sql:"index" gorm:"column:updated_at;type:timestamp" json:"updated_at,omitempty" xml:"updated_at,omitempty"`
	DeletedAt   *time.Time `sql:"index" gorm:"column:deleted_at;type:timestamp" json:"deleted_at,omitempty" xml:"deleted_at,omitempty"`
	SortNumber  int        `gorm:"column:sort_number;type:int(11)" json:"sort_number,omitempty" xml:"sort_number,omitempty"`

	MemcachedFlags            int   `gorm:"column:flags;type:int(11)" json:"flags,omitempty" xml:"flags,omitempty"`
	MemcachedCasColumn        int64 `gorm:"column:cas_column;type:bigint(20)" json:"cas_column,omitempty" xml:"cas_column,omitempty"`
	MemcachedExpireTimeColumn int   `gorm:"column:expire_time_column;int(11)" json:"expire_time_column,omitempty" xml:"expire_time_column,omitempty"`
}

//BeforeCreate ID处理
func (d *Model) BeforeCreate(scope *gorm.Scope) error {
	uuidStr := uuid.NewRandom().String()
	if err := scope.SetColumn("ID", uuidStr); nil != err {
		return err
	}
	return nil
}

//PicModel *
type PicModel struct {
	Pic    string `gorm:"column:pic;type:text" json:"pic,omitempty" xml:"pic,omitempty"`
	Pic3rd string `gorm:"column:pic_3rd;type:text" json:"pic_3rd,omitempty" xml:"pic_3rd,omitempty"`
}

//ComputeOffset *
func ComputeOffset(pageNumber, pageSize int) int {
	offset := 0
	if pageNumber > 1 {
		offset = (pageNumber - 1) * pageSize
	}
	return offset
}

//DBCtx *
type DBCtx struct {
	WebContext echo.Context
}

//NewDBCtx *
func NewDBCtx(c echo.Context) *DBCtx {
	dbCtx := &DBCtx{WebContext: c}
	return dbCtx
}

//MySQL *
func (d *DBCtx) MySQL() *gorm.DB {
	return d.WebContext.Get("mysql").(*gorm.DB)
}

//AutoMigrate * AutoMigrate
func AutoMigrate(c echo.Context, schema ...interface{}) error {
	return NewDBCtx(c).MySQL().AutoMigrate(schema...).Error
}

//ParamsToMaps *
func ParamsToMaps(args []string) map[string]string {

	params := make(map[string]string)

	for i := 0; i < len(args); i = i + 2 {
		key := args[i]
		val := args[i+1]
		params[key] = val
	}

	return params
}
