package storage

var (
	RequestSaveFile    = "INSERT INTO files(name, creation_date, update_date) VALUES($1,$2,$3)"
	RequestUpdateFile  = "UPDATE files SET update_date=$1 WHERE name=$2"
	RequestDeleteFile  = "DELETE FROM files WHERE name=$1"
	RequestGetFullData = "SELECT name, creation_date, update_date FROM files"
)
