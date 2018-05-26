package models

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/shibingli/realclouds_go/models"
)

//DrityWord *
type DrityWord struct {
	models.Model
	MD5   string `json:"md5,omitempty" xml:"md5,omitempty" gorm:"primary_key;column:md5;type:varchar(100)"`
	Value string `json:"value,omitempty" xml:"value,omitempty" gorm:"column:value;type:text"`
}

//TableName *
func (DrityWord) TableName() string {
	return "sys_drityword"
}

//AddDrityWord *
func AddDrityWord(c echo.Context, data *DrityWord) (err error) {
	return models.NewDBCtx(c).MySQL().Create(data).Error
}

//FindDrityWordByID *
func FindDrityWordByID(c echo.Context, id string) (data DrityWord, boo bool) {
	boo = models.NewDBCtx(c).MySQL().Where(&DrityWord{
		Model: models.Model{
			ID: strings.TrimSpace(id),
		},
	}).First(&data).RecordNotFound()
	return
}

//FindDrityWordByMD5 *
func FindDrityWordByMD5(c echo.Context, md5 string) (data DrityWord, boo bool) {
	boo = models.NewDBCtx(c).MySQL().Where(&DrityWord{
		MD5: strings.TrimSpace(md5),
	}).First(&data).RecordNotFound()
	return
}

//UpdateDrityWord *
func UpdateDrityWord(c echo.Context, data *DrityWord) (err error) {
	return models.NewDBCtx(c).MySQL().Model(&DrityWord{}).Update(&data).Error
}

//FindDrityWords *
func FindDrityWords(c echo.Context, args ...string) (count int, data []DrityWord) {
	argMap := models.ParamsToMaps(args)

	ids, _ := argMap["ids"]
	md5s, _ := argMap["md5s"]

	viewData := models.NewDBCtx(c).MySQL()

	if len(ids) != 0 {
		dwIDs := strings.Split(ids, ",")
		c.Logger().Debugf("Drity words IDs: %v", dwIDs)
		idsLen := len(dwIDs)
		if idsLen > 0 {
			viewData = viewData.Where("id in (?)", dwIDs)
		}
	}

	if len(md5s) != 0 {
		dwMD5s := strings.Split(md5s, ",")
		c.Logger().Debugf("Drity words MD5s: %v", dwMD5s)
		idsLen := len(dwMD5s)
		if idsLen > 0 {
			viewData = viewData.Where("md5 in (?)", dwMD5s)
		}
	}

	viewData.Model(&DrityWord{}).Count(&count).Find(&data)

	return
}

// DeleteDrityWordByID *
func DeleteDrityWordByID(c echo.Context, id string) (err error) {
	err = models.NewDBCtx(c).MySQL().Where(&DrityWord{
		Model: models.Model{
			ID: strings.TrimSpace(id),
		},
	}).Delete(&DrityWord{}).Error
	return
}
