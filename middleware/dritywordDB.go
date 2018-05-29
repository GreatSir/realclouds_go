package middleware

import (
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

//DrityWordDB *
type DrityWordDB struct {
	ID                        string     `sql:"index" gorm:"primary_key;column:id;type:varchar(100)" json:"id,omitempty" xml:"id,omitempty"`
	Name                      string     `sql:"index" gorm:"column:name;type:varchar(100)" json:"name,omitempty" xml:"name,omitempty"`
	Description               string     `gorm:"column:description;type:text" json:"description,omitempty" xml:"description,omitempty"`
	CreatedAt                 time.Time  `sql:"index" gorm:"column:created_at;type:timestamp" json:"created_at,omitempty" xml:"created_at,omitempty"`
	UpdatedAt                 time.Time  `gorm:"column:updated_at;type:timestamp NULL" json:"updated_at,omitempty" xml:"updated_at,omitempty"`
	DeletedAt                 *time.Time `sql:"index" gorm:"column:deleted_at;type:timestamp NULL" json:"deleted_at,omitempty" xml:"deleted_at,omitempty"`
	MemcachedFlags            int        `gorm:"column:flags;type:int(11)" json:"flags,omitempty" xml:"flags,omitempty"`
	MemcachedCasColumn        int64      `gorm:"column:cas_column;type:bigint(20)" json:"cas_column,omitempty" xml:"cas_column,omitempty"`
	MemcachedExpireTimeColumn int        `gorm:"column:expire_time_column;int(11)" json:"expire_time_column,omitempty" xml:"expire_time_column,omitempty"`

	MD5   string `json:"md5,omitempty" xml:"md5,omitempty" gorm:"primary_key;column:md5;type:varchar(100)"`
	Value string `json:"value,omitempty" xml:"value,omitempty" gorm:"column:value;type:text"`
}

//TableName *
func (DrityWordDB) TableName() string {
	return "sys_drityword"
}

//BeforeCreate ID处理
func (d *DrityWordDB) BeforeCreate(scope *gorm.Scope) error {
	uuidStr := uuid.NewRandom().String()
	if err := scope.SetColumn("ID", uuidStr); nil != err {
		return err
	}
	return nil
}

//AddDrityWord *
func AddDrityWord(db *gorm.DB, data *DrityWordDB) (err error) {
	return db.Create(&data).Error
}

//FindDrityWordByID *
func FindDrityWordByID(db *gorm.DB, id string) (data DrityWordDB, boo bool) {
	boo = db.Where(&DrityWordDB{
		ID: strings.TrimSpace(id),
	}).First(&data).RecordNotFound()
	return
}

//FindDrityWordByMD5 *
func FindDrityWordByMD5(db *gorm.DB, md5 string) (data DrityWordDB, boo bool) {
	boo = db.Where(&DrityWordDB{
		MD5: strings.TrimSpace(md5),
	}).First(&data).RecordNotFound()
	return
}

//UpdateDrityWord *
func UpdateDrityWord(db *gorm.DB, data *DrityWordDB) (err error) {
	return db.Model(&DrityWordDB{}).Update(&data).Error
}

//FindDrityWords *
func FindDrityWords(db *gorm.DB, args ...string) (count int, data []DrityWordDB) {
	argMap := paramsToMaps(args)

	ids, _ := argMap["ids"]
	md5s, _ := argMap["md5s"]
	keywords, _ := argMap["keywords"]

	if len(ids) != 0 {
		dwIDs := strings.Split(ids, ",")
		idsLen := len(dwIDs)
		if idsLen > 0 {
			db = db.Where("id in (?)", dwIDs)
		}
	}

	if len(md5s) != 0 {
		dwMD5s := strings.Split(md5s, ",")
		idsLen := len(dwMD5s)
		if idsLen > 0 {
			db = db.Where("md5 in (?)", dwMD5s)
		}
	}

	if len(keywords) != 0 {
		if len(keywords) != 0 {
			db = db.Where("name LIKE ?", "%"+keywords+"%").Or("description LIKE ?", "%"+keywords+"%").Or("value LIKE ?", "%"+keywords+"%").Or("md5 LIKE ?", "%"+keywords+"%")
		}
	}

	db.Model(&DrityWordDB{}).Count(&count).Find(&data)

	return
}

// DeleteDrityWordByID *
func DeleteDrityWordByID(db *gorm.DB, id string) (err error) {
	err = db.Where(&DrityWordDB{
		ID: strings.TrimSpace(id),
	}).Delete(&DrityWordDB{}).Error
	return
}

//paramsToMaps *
func paramsToMaps(args []string) map[string]string {

	params := make(map[string]string)

	for i := 0; i < len(args); i = i + 2 {
		key := args[i]
		val := args[i+1]
		params[key] = val
	}

	return params
}
