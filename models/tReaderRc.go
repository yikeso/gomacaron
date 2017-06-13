package models

type readerRc struct {
	id int64
	maxresourceId int64
}
//更新创建电子书txt进度
//maxId电子书创建txt进度，id保存进度的对应记录
func UpdateReaderRcMaxresourceid(maxId int64,id int64)(err error){
	query := "Update T_READER_RC Set maxresourceid = ? Where id = ?"
	if err != nil {
		return
	}
	result,err := resourceDb.Exec(query,maxId,id)
	if err != nil {
		return
	}
	_,err = result.RowsAffected()
	return
}
