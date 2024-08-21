package tpcmon

import (
	"check_ebs_check/tool"
	"database/sql"
	"errors"
	"fmt"

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

type DBInfo struct {
	Id     int    `json:"db_id"`
	User   string `json:"db_user"`
	Passwd string `json:"db_password"`
	Port   int    `json:"db_port"`
	Name   string `json:"db_name"`
	Host   string `json:"db_host"`
}

type AKSKInfo struct {
	Ep string `json:"ep"`
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}

type Databases struct {
	Ascm          DBInfo `json:"ascm"`
	Perf          DBInfo `json:"perf"`
	AsoStm        DBInfo `json:"aso_stm"`
	HouyiRegionDB DBInfo `json:"houyiregiondb"`
	VpcRegionDB   DBInfo `json:"vpcregiondb"`
	Xuanyuan      DBInfo `json:"xuanyuan"`
	Dbaas         DBInfo `json:"dbaas"`
	EcsDriver     DBInfo `json:"ecsdriver"`
	RiverDB       DBInfo `json:"riverdb"`
}

type Conf struct {
	DB   Databases `json:"db"`
	AKSK struct {
		SLS AKSKInfo `json:"sls"`
	} `json:"aksk"`
}

// 获取指定产品的db的连接信息
// （仅供内部调用）
func (c *Conf) getDBConfig(dbName string) (DBInfo, error) {
	switch dbName {
	case "ascm":
		return c.DB.Ascm, nil
	case "perf":
		return c.DB.Perf, nil
	case "aso_stm":
		return c.DB.AsoStm, nil
	case "houyiregiondb":
		return c.DB.HouyiRegionDB, nil
	case "vpcregiondb":
		return c.DB.VpcRegionDB, nil
	case "xuanyuan":
		return c.DB.Xuanyuan, nil
	case "dbaas":
		return c.DB.Dbaas, nil
	case "ecsdriver":
		return c.DB.EcsDriver, nil
	case "riverdb":
		return c.DB.RiverDB, nil
	default:
		return DBInfo{}, errors.New("database not found")
	}
}

// 获取数据库连接实例
func (c *Conf) GetDbConn(dbName string) (*sql.DB, error) {
	// 获取db连接信息
	dbinfo, err := c.getDBConfig(dbName)
	if err != nil {
		return nil, err
	}

	// password：base64 decode
	encoded_pwd := dbinfo.Passwd
	decoded_pwsswd, err := tool.DecodeBase64(encoded_pwd)
	if err != nil {
		fmt.Printf("base64 decoder failed! err:%v\n", err)
		return &sql.DB{}, err
	}

	// make dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbinfo.User,
		decoded_pwsswd,
		dbinfo.Host,
		dbinfo.Port,
		dbinfo.Name,
	)

	// dial db
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("open dns [%s] failed! err:%v\n", dsn, err)
		return &sql.DB{}, err
	}

	// ping db
	if err := db.Ping(); err != nil {
		fmt.Printf("can not ping database, err:%v\n", err)
		return &sql.DB{}, err
	}

	// fmt.Printf("获取数据库[%s]conn成功!\n", dbName)
	return db, nil
}

// 获取ascm_brm库conn
func (c *Conf) GetASCMBrmDbConn() (*sql.DB, error) {
	// 获取db连接信息
	dbinfo, err := c.getDBConfig("ascm")
	if err != nil {
		return nil, err
	}

	// password：base64 decode
	encoded_pwd := dbinfo.Passwd
	decoded_pwsswd, err := tool.DecodeBase64(encoded_pwd)
	if err != nil {
		fmt.Printf("base64 decoder failed! err:%v\n", err)
		return &sql.DB{}, err
	}

	// make dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbinfo.User,
		decoded_pwsswd,
		dbinfo.Host,
		dbinfo.Port,
		"brm",
	)

	// dial db
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("open dns [%s] failed! err:%v\n", dsn, err)
		return &sql.DB{}, err
	}

	// ping db
	if err := db.Ping(); err != nil {
		fmt.Printf("can not ping database, err:%v\n", err)
		return &sql.DB{}, err
	}

	// fmt.Printf("获取数据库[%s]conn成功!\n", dbName)
	return db, nil
}

// 获取指定产品的aksk的连接信息
/*
	product:
		sls
*/
func (c *Conf) GetAkskConfig(product string) (AKSKInfo, error) {
	switch product {
	case "sls":
		return c.AKSK.SLS, nil
	default:
		return AKSKInfo{}, errors.New("product not found in config")
	}
}

// 获取sls客户端
func (c *Conf) GetSlsClinet() sls.ClientInterface {
	return sls.CreateNormalInterface(c.AKSK.SLS.Ep, c.AKSK.SLS.Ak, c.AKSK.SLS.Sk, "")
}
