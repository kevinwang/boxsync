package cache

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.engr.illinois.edu/sp-box/boxsync/box"
	"gitlab.engr.illinois.edu/sp-box/boxsync/sync"
)

var defaultLocalRootDirectory = path.Join(os.Getenv("HOME"), "Box Sync")
var defaultRemoteRootDirectory = "Box Sync"
var defaultDBLocation = path.Join(os.Getenv("HOME"), ".boxsync_cache.db")

type SyncCache interface {
	UpdateCache() error
	HardRefresh() error
	RescanLocalTree() error
	//SetEntryInvalid(path string) error
}

type syncCache struct {
	client              box.Client
	db                  *sql.DB
	localRootDirectory  string
	remoteRootDirectory string
	dbLocation          string
}

type FileCacheEntry struct {
	Path       sql.NullString
	ID         sql.NullString
	SHA1       sql.NullString
	Valid      sql.NullBool
	SequenceID sql.NullString
	ParentID   sql.NullString
}

type FolderCacheEntry struct {
	Path       sql.NullString
	ID         sql.NullString
	Valid      sql.NullBool
	SequenceID sql.NullString
	ParentID   sql.NullString
}

func NewCache(client box.Client) (SyncCache, error) {
	if client == nil {
		return nil, errors.New("Client cannot be nil")
	}

	db, err := sql.Open("sqlite3", defaultDBLocation)
	if err != nil {
		return nil, err
	}

	sqlStmt := "drop table if exists files;"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = "drop table if exists folders;"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	sqlStmt = "pragma foreign_keys=ON;"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Print("Failed enabbling foreign keys")
		return nil, err
	}

	sqlStmt = `create table folders
	(Path text not null primary key,
	ID text unique,
	Valid boolean,
	SequenceID text,
	ParentID text,
	FOREIGN KEY(ParentID) REFERENCES folders(ID));
	delete from folders;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Print("Failed creating the folders table")
		return nil, err
	}

	sqlStmt = `create table files
	(Path text not null primary key,
	ID text unique,
	SHA1 text,
	Valid boolean,
	SequenceID text,
	ParentID text,
	FOREIGN KEY(ParentID) REFERENCES folders(ID));
	delete from files;`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	cache := syncCache{
		client:              client,
		db:                  db,
		localRootDirectory:  defaultLocalRootDirectory,
		remoteRootDirectory: defaultRemoteRootDirectory,
		dbLocation:          defaultDBLocation,
	}

	err = cache.HardRefresh()
	if err != nil {
		return nil, err
	}

	err = cache.UpdateCache()
	if err != nil {
		return nil, err
	}

	return &cache, nil
}

func (c *syncCache) SetEntryInvalid(path string) error {
	var sqlStmtText string
	fi, err := os.Stat(path)
	if err != nil {
		return err
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		sqlStmtText = `update folders Valid = ? where Path = ?;`
	case mode.IsRegular():
		sqlStmtText = `update files Valid = ? where Path = ?;`
	}

	sqlStmt, err := c.db.Prepare(sqlStmtText)
	if err != nil {
		return err
	}

	relPath, err := filepath.Rel(c.localRootDirectory, path)
	if err != nil {
		relPath = path
	}

	path = filepath.Join(filepath.Base(c.remoteRootDirectory), relPath)

	_, err = sqlStmt.Exec(false, path)
	if err != nil {
		return nil
	}

	return nil
}

func (c *syncCache) RescanLocalTree() error {
	deletesFolder := map[string]bool{}
	row, err := c.db.Query("SELECT Path FROM folders;")
	if err != nil {
		return err
	}

	var pathName string
	for row.Next() {
		row.Scan(&pathName)
		deletesFolder[pathName] = true
	}
	row.Close()

	deletesFile := map[string]bool{}
	row, err = c.db.Query("SELECT Path FROM files;")
	if err != nil {
		return err
	}

	for row.Next() {
		row.Scan(&pathName)
		deletesFile[pathName] = true
	}
	row.Close()

	filepath.Walk(c.localRootDirectory, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		origFile := filePath
		relPath, err := filepath.Rel(c.localRootDirectory, filePath)
		if err != nil {
			relPath = filePath
		}

		filePath = filepath.Join(filepath.Base(c.remoteRootDirectory), relPath)

		if info.IsDir() {
			_, ok := deletesFolder[filePath]
			if ok {
				delete(deletesFolder, filePath)
			} else {
				c.addFolderToDB(filePath)
			}
		} else {
			_, ok := deletesFile[filePath]
			if ok {
				delete(deletesFile, filePath)
				row, err := c.db.Query(`select ID, SHA1, SequenceID FROM files where Path = "` + filePath + `";`)
				if err != nil {
					return err
				}

				var ID string
				var SHA1 string
				var SequenceID string
				if row.Next() {
					row.Scan(&ID, &SHA1, &SequenceID)
					row.Close()
				} else {
					return errors.New("Did not find a file where we expected one")
				}

				if SHA1 != sync.SHA1(origFile) {
					file, err := c.client.GetFile(ID)
					if err != nil {
						return err
					}

					if file.SequenceID == SequenceID {
						_, err := c.client.UploadFileVersion(ID, origFile)
						if err != nil {
							return err
						}
					} else {
						return errors.New("Conflict Resolution not yet Implemented.")
					}
				}
			} else {
				c.AddFileToDB(origFile)
			}
		}

		return nil
	})

	return nil
}

func (c *syncCache) HardRefresh() error {
	rootFolder, err := sync.GetSyncRootFolder(c.client)
	if err != nil {
		return err
	}

	sqlStmt, err := c.db.Prepare(`insert into folders (Path, ID, Valid, SequenceID, ParentID) values (?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	_, err = sqlStmt.Exec(rootFolder.Name, rootFolder.ID, true, rootFolder.SequenceID, nil)
	if err != nil {
		return err
	}
	sqlStmt.Close()

	err = c.hardCacheAll(rootFolder.ID, c.localRootDirectory, c.remoteRootDirectory)
	if err != nil {
		return err
	}

	return nil
}

func (c *syncCache) UpdateCache() error {
	return nil
}

func (c *syncCache) hardCacheAll(folderID, destPath, remotePath string) error {
	tx, err := c.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	updateFileStmt, err := tx.Prepare(`update files set ID = ?, Valid = ?, SHA1 = ?, SequenceID = ?, ParentID = ? where Path = ?;`)
	if err != nil {
		return err
	}

	insertFileStmt, err := tx.Prepare(`insert into files (Path, ID, SHA1, Valid, SequenceID, ParentID) values (?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}

	contents, err := c.client.GetFolderContents(folderID)
	if err != nil {
		return err
	}

	for _, file := range contents.Files {
		var filePathLoc string
		var fileSHA string
		remoteFilePath := path.Join(remotePath, file.Name)
		remoteRelPath, err := filepath.Rel(c.remoteRootDirectory, remoteFilePath)
		if err != nil {
			return err
		}

		localFilePath := path.Join(destPath, remoteRelPath)

		rows, err := c.db.Query(`SELECT Path, SHA1 FROM files WHERE Path = "` + remoteFilePath + `";`)
		if err != nil {
			return err
		}

		if rows.Next() {
			rows.Scan(&filePathLoc, &fileSHA)
			rows.Close()

			_, err = updateFileStmt.Exec(file.ID, true, file.SHA1, file.SequenceID, folderID, remoteFilePath)
			if err != nil {
				return err
			}

			if strings.Compare(fileSHA, file.SHA1) != 0 {
				err = c.client.DownloadFile(file.ID, localFilePath)
				if err != nil {
					log.Print("Failed to download file")
					return err
				}
			}
		} else {
			rows.Close()

			_, err = insertFileStmt.Exec(remoteFilePath, file.ID, file.SHA1, true, file.SequenceID, folderID)
			err = c.client.DownloadFile(file.ID, localFilePath)
			if err != nil {
				log.Print("Failed to download file")
				return err
			}
		}
	}

	insertFileStmt.Close()
	updateFileStmt.Close()
	tx.Commit()

	for _, folder := range contents.Folders {

		//Build the local and remote paths
		var folderPathLoc string
		remoteFolderPath := path.Join(remotePath, folder.Name)
		remoteRelPath, err := filepath.Rel(c.remoteRootDirectory, remoteFolderPath)
		if err != nil {
			return err
		}

		localFolderPath := path.Join(destPath, remoteRelPath)

		//Grab the desired DB entry
		rows, err := c.db.Query(`SELECT Path FROM folders WHERE Path = "` + remoteFolderPath + `";`)
		if err != nil {
			return err
		}

		if rows.Next() {
			updateFolderStmt, err := c.db.Prepare(`update folders set ID = ?, Valid = ?, SequenceID = ?, ParentID = ? where Path = ?;`)
			if err != nil {
				return err
			}

			rows.Scan(&folderPathLoc)
			rows.Close()

			_, err = updateFolderStmt.Exec(folder.ID, true, folder.SequenceID, folderID, remoteFolderPath)
			if err != nil {
				return err
			}

			updateFolderStmt.Close()
		} else {
			rows.Close()

			insertFolderStmt, err := c.db.Prepare(`insert into folders (Path, ID, Valid, SequenceID, ParentID) values (?, ?, ?, ?, ?)`)
			if err != nil {
				return nil
			}

			_, err = insertFolderStmt.Exec(remoteFolderPath, folder.ID, true, folder.SequenceID, folderID)
			if err != nil {
				return err
			}

			insertFolderStmt.Close()
		}

		if _, err := os.Stat(localFolderPath); os.IsNotExist(err) {
			fmt.Printf("Creating directory %s\n", localFolderPath)
			err := os.MkdirAll(localFolderPath, 0755)
			if err != nil {
				log.Print("Creating directory error.")
				return err
			}
		}

		err = c.hardCacheAll(folder.ID, destPath, remoteFolderPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *syncCache) AddFolderToDB(folderPath string) (string, error) {
	relPath, err := filepath.Rel(c.localRootDirectory, folderPath)
	if err != nil {
		relPath = folderPath
	}

	folderPath = filepath.Join(filepath.Base(c.remoteRootDirectory), relPath)
	return c.addFolderToDB(folderPath)
}

func (c *syncCache) addFolderToDB(folderPath string) (string, error) {
	parent, err := c.db.Query("SELECT ID FROM folders WHERE Path = \"" + folderPath + "\";")
	if err != nil {
		return "", err
	}

	var ID string
	if parent.Next() {
		parent.Scan(&ID)
		parent.Close()
		return ID, nil
	} else {
		parentID, err := c.addFolderToDB(filepath.Dir(folderPath))
		if err != nil {
			return "", err
		}

		folder, err := c.client.CreateFolder(filepath.Base(folderPath), parentID)
		if err != nil {
			return "", err
		}

		stmt, err := c.db.Prepare("INSERT OR IGNORE into folders (Path, ID, Valid, SequenceID, ParentID) values (?, ?, ?, ?, ?)")
		if err != nil {
			return "", err
		}

		_, err = stmt.Exec(folderPath, folder.ID, true, folder.SequenceID, parentID)
		if err != nil {
			return "", err
		}
		stmt.Close()
	}

	return ID, nil
}

func (c *syncCache) AddFileToDB(filePath string) (string, error) {
	relPath, err := filepath.Rel(c.localRootDirectory, filePath)
	if err != nil {
		relPath = filePath
	}

	filePath = filepath.Join(filepath.Base(c.remoteRootDirectory), relPath)

	return c.addFileToDB(filePath)
}

func (c *syncCache) addFileToDB(filePath string) (string, error) {
	parent, err := c.db.Query("SELECT ParentID FROM files WHERE Path = \"" + filePath + "\";")
	if err != nil {
		return "", err
	}

	var parentID string
	if parent.Next() {
		parent.Scan(&parentID)
		parent.Close()
		return parentID, nil
	} else {
		parentID, err := c.addFolderToDB(filepath.Dir(filePath))
		if err != nil {
			return "", err
		}

		file, err := c.client.UploadFile(filePath, parentID)
		if err != nil {
			return "", err
		}

		stmt, err := c.db.Prepare("INSERT OR IGNORE into files (Path, ID, Valid, SequenceID, SHA1, ParentID) values (?, ?, ?, ?, ?, ?)")
		if err != nil {
			return "", err
		}

		_, err = stmt.Exec(filePath, file.ID, true, file.SequenceID, file.SHA1, parentID)
		stmt.Close()
		if err != nil {
			return "", err
		}
	}

	return parentID, nil
}

/*
func LocalCacheChangeUpdate(client box.Client, db *sql.DB, table string) error {
	rows, err := db.Query("SELECT path, id, sequenceID FROM \"" + table + "\" WHERE valid = 0")
	if err != nil {
		return err
	}
	var filePath string
	var fileID string
	var fileSequenceID string
	for rows.Next() {
		rows.Scan(&filePath, &fileID, &fileSequenceID)

		if fileID == "" {
			fi, err := os.Stat(path.Join(table, filePath))
			if err != nil {
				return err
			}
			switch mode := fi.Mode(); {
			case mode.IsDir():
				err := addFolderToDB(client, db, table, filePath)
				if err != nil {
					return err
				}
			case mode.IsRegular():
				err := addFileToDB(client, db, table, filePath)
				if err != nil {
					return err
				}
			}
		}
	}
	rows.Close()

	return nil
}
*/
