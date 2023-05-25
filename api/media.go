package api

import (
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/google/uuid"
	"github.com/jasontconnell/sitecore/data"
	"github.com/jasontconnell/sqlhelp"
)

func LoadBlob(connstr string, id uuid.UUID) (data.Blob, error) {
	query := fmt.Sprintf(`select Data from Blobs where BlobId = '%s' order by [Index]`, id)
	conn, cerr := sql.Open("mssql", connstr)
	if cerr != nil {
		return nil, cerr
	}
	defer conn.Close()

	records, rerr := sqlhelp.GetResultSet(conn, query)
	if rerr != nil {
		return nil, rerr
	}

	var blob data.Blob
	var bytes []byte

	for _, row := range records {
		bytes = append(bytes, row["Data"].([]byte)...)
	}

	blob = data.NewBlob(id, bytes)
	return blob, nil
}
