package modules

//可执行命令管理
type Kubecmds struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Comname string `gorm:"type:varchar(50);uniqueIndex:cname;comment:'命令名称'" binding:"omitempty,min=3,max=50" json:"comname"`
	Comvers string `gorm:"type:varchar(50);comment:'命令版本'" binding:"omitempty,min=3,max=50" json:"comvers"`
	Ctime   string `gorm:"type:varchar(20);comment:'命令创建时间'" json:"ctime"`
	Compath string `gorm:"type:varchar(200)" binding:"omitempty,min=10,max=200" json:"compath"`
	State   int    `gorm:"type:tinyint;default:0;comment:'命令状态，1表示启用，2表示禁用'" binding:"gte=0,lte=10" json:"state"`
}

func (k *Kubecmds) ComsCreate() error {
	return Db.Model(&k).Create(&k).Error
}

//查询一条命令
func (k *Kubecmds) ComsSelectByName() error {
	return Db.Model(&k).Where("`comname` = ?", k.Comname).Find(&k).Error
}

//查询命令是否为激活状态
func (k *Kubecmds) ComsCommandisActive() bool {
	if Db.Model(&k).Where("`comname` = ?", k.Comname).Find(&k).Error != nil {
		return false
	}
	if k.State == 1 {
		return true
	} else {
		return false
	}
}

//查询命令名称列表
type cmcs struct {
	Comname string
}

//查询kubectl命令列表
func ComsSelectNameList() ([]string, error) {
	var H []cmcs
	var R []string
	err := Db.Model(Kubecmds{}).Select("comname").Where("`state` = 1").Scan(&H).Error
	if H != nil && err == nil {
		for _, v := range H {
			R = append(R, v.Comname)
		}
		return R, err
	} else {
		return nil, err
	}
}

//查询命令列表,0查询无可执行文件的命令，1查询激活的命令(命令文件存在)，2查询禁用的命令，3查询所有命令
func ComsSelectList(i int) ([]Kubecmds, error) {
	var H []Kubecmds
	if i < 3 {
		return H, Db.Raw("select comname,comvers,ctime,state from kubecmds where state = ?", i).Scan(&H).Error
	} else {
		return H, Db.Raw("select comname,comvers,ctime,state from kubecmds").Scan(&H).Error
	}
}

//命令的启用和禁用,i==1表示启用，i==2表示禁用
func (k *Kubecmds) ComsEnableOrDisable(i int) error {
	return Db.Model(&k).Where("`comname` = ?", k.Comname).Update("state", i).Error
}

//删除命令
func (k *Kubecmds) ComsDelete() error {
	return Db.Model(&k).Delete(&k).Error
}
