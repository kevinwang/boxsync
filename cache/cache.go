package cache

import (
	"database/sql"
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

type CacheEntry struct {
	Path       sql.NullString
	ID         sql.NullInt64
	SHA1       sql.NullString
	Valid      sql.NullBool
	SequenceID sql.NullInt64
}

func InitCache(client box.Client, root string) *sql.DB {
	db, err := sql.Open("sqlite3", path.Join(os.Getenv("HOME"), ".cache.db"))
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := "drop table if exists \"" + root + "\";"

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal("%q : %s\n", err, sqlStmt)
	}

	sqlStmt = `
	create table "` + root + `" (path text not null primary key, id integer unique, sha1 text, valid boolean, sequenceID text);
	delete from "` + root + "\";"

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q : %s\n", err, sqlStmt)
		log.Fatal(err)
	}

	if client != nil {
		initCacheRemote(client, root, db)
	}
	initCacheLocal(root, db)

	return db
}

func initCacheLocal(root string, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Print("Failed to begin database")
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert or ignore into \"" + root + "\" (path, id, sha1, valid) values (?, ?, ?, ?)")
	if err != nil {
		log.Print("Prepare returned an error")
		log.Fatal(err)
	}
	defer stmt.Close()

	filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			//log.Print("Failed on path: " + filePath)
			return err
		}
		if info.IsDir() {
			return nil
		}

		filePath, _ = filepath.Rel(root, filePath)
		_, err = stmt.Exec(filePath, nil, sync.SHA1(filePath), false)
		if err != nil {
			log.Print("Failed to add " + filePath)
		}
		return nil
	})

	tx.Commit()
}

func initCacheRemote(client box.Client, root string, db *sql.DB) {
	rootFolder, errRoot := sync.GetSyncRootFolder(client)
	if errRoot != nil {
		log.Print("GetSyncRootFolder returned an error")
		log.Fatal(errRoot)
	}

	err := CacheAll(client, rootFolder.ID, root, "", db, root)
	if err != nil {
		log.Fatal(err)
	}
}

func CacheAll(client box.Client, folderID, destPath string, remotePath string, db *sql.DB, table string) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("update \"" + table + "\" set id = ? , valid = ?, sha1 = ? , sequenceID = ? where path = ?")
	if err != nil {
		log.Print("Prepare stmt returned an error")
		log.Fatal(err)
	}

	stmtInsert, err := tx.Prepare("insert into  \"" + table + "\" (path, id, sha1, valid, sequenceID) values (?, ?, ?, ?, ?)")
	if err != nil {
		log.Print("Prepare stmtInsert returned an error")
		log.Fatal(err)
	}

	contents, err := client.GetFolderContents(folderID)
	if err != nil {
		return err
	}

	for _, file := range contents.Files {
		var filePathLoc string
		var fileSHA string
		remoteFilePath := path.Join(remotePath, file.Name)
		localFilePath := path.Join(destPath, remoteFilePath)
		rows, err := db.Query("SELECT path, sha1 FROM \"" + table + "\" WHERE path = \"" + remoteFilePath + "\"")
		if err != nil {
			log.Print("recursive file call problem on table " + table)
			log.Fatal(err)
		}

		if rows.Next() {
			rows.Scan(&filePathLoc, &fileSHA)
			rows.Close()

			if strings.Compare(fileSHA, file.SHA1) == 0 {
				_, err = stmt.Exec(file.ID, true, file.SHA1, remoteFilePath)
				if err != nil {
					return err
				}
			} else {
				_, err = stmt.Exec(file.ID, true, file.SHA1, remoteFilePath)
				if err != nil {
					return err
				}

				err = client.DownloadFile(file.ID, localFilePath)
				if err != nil {
					return err
				}
			}
		} else {
			rows.Close()

			_, err = stmtInsert.Exec(remoteFilePath, file.ID, file.SHA1, true, file.SequenceID)
			err = client.DownloadFile(file.ID, localFilePath)
			if err != nil {
				return err
			}
		}
	}

	tx.Commit()
	stmt.Close()

	for _, folder := range contents.Folders {
		remoteFolderPath := path.Join(remotePath, folder.Name)
		localFolderPath := path.Join(destPath, remoteFolderPath)

		if _, err := os.Stat(localFolderPath); os.IsNotExist(err) {
			fmt.Printf("Creating directory %s\n", localFolderPath)
			err := os.MkdirAll(localFolderPath, 0755)
			if err != nil {
				log.Print("Creating directory error.")
				return err
			}
		}

		err = CacheAll(client, folder.ID, destPath, remoteFolderPath, db, table)
		if err != nil {
			return err
		}
	}

	return nil
}
